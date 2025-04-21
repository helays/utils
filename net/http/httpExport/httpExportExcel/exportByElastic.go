package httpExportExcel

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/helays/utils/close/esClose"
	"github.com/helays/utils/db/elastic/elasticModel"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/helays/utils/config"
	"github.com/helays/utils/excelTools"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/net/http/httpTools"
	"github.com/helays/utils/tools"
	"github.com/xuri/excelize/v2"
)

// ElasticsearchExport 用于导出Elasticsearch数据的结构体
type ElasticsearchExport struct {
	ESClient      *elasticsearch.Client // Elasticsearch客户端
	SearchRequest *esapi.SearchRequest  // 搜索请求
	FileType      string                // 文件类型 (excel/csv)
	FileName      string                // 文件名
	ExportHeader  map[string]string     // 字段名映射
	BatchSize     int                   // 每批处理的数据量
	ScrollTime    time.Duration         // 滚动查询保持时间
	QueryBody     io.Reader             // 查询体
	SortFields    []string              // 排序字段
}

// Response 导出Elasticsearch数据到HTTP响应
func (e *ElasticsearchExport) Response(w http.ResponseWriter) error {
	// 设置默认值
	e.FileType = strings.ToLower(tools.Ternary(e.FileType == "", config.ExportFileTypeExcel, e.FileType))
	e.FileName = tools.Ternary(e.FileName == "", "export", e.FileName)
	e.BatchSize = tools.Ternary(e.BatchSize == 0, 1000, e.BatchSize)
	e.ScrollTime = tools.Ternary(e.ScrollTime == 0, 10*time.Minute, e.ScrollTime)

	// 初始化导出文件
	var (
		f            *excelize.File
		streamWriter *excelize.StreamWriter
		cw           *csv.Writer
		sheetName    = "Sheet1"
		err          error
	)

	if e.FileType == config.ExportFileTypeExcel {
		f = excelize.NewFile()
		defer excelTools.CloseExcel(f)
		if streamWriter, err = f.NewStreamWriter(sheetName); err != nil {
			return fmt.Errorf("创建sheet失败：%w", err)
		}
	} else if e.FileType == config.ExportFileTypeCsv {
		w.Header().Set("Content-Type", "text/csv")
		httpTools.SetDisposition(w, e.FileName+".csv")
		cw = csv.NewWriter(w)
		defer cw.Flush()
		// 写入UTF-8 BOM头
		_, _ = w.Write([]byte("\xef\xbb\xbf"))
	} else {
		return fmt.Errorf("不支持的导出类型：%s", e.FileType)
	}

	// 构建SearchRequest
	if e.SearchRequest == nil {
		// 如果没有提供SearchRequest，则创建一个默认的
		e.SearchRequest = &esapi.SearchRequest{
			Body:    e.QueryBody,
			Size:    &e.BatchSize,
			Scroll:  e.ScrollTime,
			Sort:    e.SortFields,
			Pretty:  false,
			Human:   false,
			Timeout: 30 * time.Second,
		}
	} else {
		// 确保SearchRequest有必要的参数
		if e.SearchRequest.Size == nil {
			e.SearchRequest.Size = &e.BatchSize
		}
		if e.SearchRequest.Scroll == 0 {
			e.SearchRequest.Scroll = e.ScrollTime
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 执行初始搜索请求
	res, err := e.SearchRequest.Do(ctx, e.ESClient)
	if err != nil {
		return fmt.Errorf("执行搜索请求失败: %w", err)
	}
	defer esClose.CloseResp(res)

	if res.IsError() {
		return fmt.Errorf("搜索请求错误: %s", res.String())
	}

	// 解析初始响应
	var initialResponse elasticModel.ESSearchResponse

	if err = json.NewDecoder(res.Body).Decode(&initialResponse); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	hits := initialResponse.GetHits()
	hitLen := len(hits)

	// 获取字段名
	var columnNames []string
	if hitLen > 0 {
		columnNames = e.getFieldNames(hits[0].Source)
	} else {
		return fmt.Errorf("没有找到可导出的数据")
	}

	// 写入表头
	if err = e.writeHeader(columnNames, streamWriter, cw); err != nil {
		return err
	}

	// 处理第一批结果
	if err = e.processHits(hits, columnNames, 2, streamWriter, cw); err != nil {
		return err
	}

	// 使用滚动查询获取剩余结果
	// 使用滚动查询获取剩余结果
	scrollID := initialResponse.ScrollID
	rowCount := hitLen + 2 // 从第二行开始写入数据(第一行是表头)

	for {
		if scrollID == "" {
			break
		}

		scrollReq := esapi.ScrollRequest{
			ScrollID: scrollID,
			Scroll:   e.ScrollTime,
		}

		res, err = scrollReq.Do(context.Background(), e.ESClient)
		if err != nil {
			return fmt.Errorf("滚动查询失败: %w", err)
		}

		if res.IsError() {
			esClose.CloseResp(res)
			return fmt.Errorf("滚动查询错误: %s", res.String())
		}

		var scrollResponse elasticModel.ESSearchResponse
		if err = json.NewDecoder(res.Body).Decode(&scrollResponse); err != nil {
			esClose.CloseResp(res)
			return fmt.Errorf("解析滚动响应失败: %w", err)
		}
		hits = scrollResponse.GetHits()
		if len(hits) == 0 {
			esClose.CloseResp(res)
			break
		}
		// 处理当前批次结果
		if err = e.processHits(hits, columnNames, rowCount, streamWriter, cw); err != nil {
			esClose.CloseResp(res)
			return err
		}

		rowCount += len(hits)
		scrollID = scrollResponse.ScrollID
		esClose.CloseResp(res) // 立即关闭响应体
	}
	// 清除滚动上下文
	clearScrollReq := esapi.ClearScrollRequest{
		ScrollID: []string{scrollID},
	}
	_, err = clearScrollReq.Do(context.Background(), e.ESClient)
	if err != nil {
		ulogs.Warn("清除滚动上下文失败: %v", err)
	}

	// 完成Excel导出
	if e.FileType == config.ExportFileTypeExcel {
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		httpTools.SetDisposition(w, e.FileName+".xlsx")
		_ = streamWriter.Flush()
		if err = f.Write(w); err != nil {
			return fmt.Errorf("导出excel失败：%w", err)
		}
	}

	return nil
}

// writeHeader 写入表头
func (e *ElasticsearchExport) writeHeader(columnNames []string, streamWriter *excelize.StreamWriter, cw *csv.Writer) error {
	// 应用字段名映射
	for idx, field := range columnNames {
		if header, ok := e.ExportHeader[field]; ok {
			columnNames[idx] = header
		}
	}

	if e.FileType == config.ExportFileTypeExcel {
		if err := streamWriter.SetRow("A1", tools.StrSlice2AnySlice(columnNames)); err != nil {
			return fmt.Errorf("导出excel失败，表头写入失败: %w", err)
		}
	} else if e.FileType == config.ExportFileTypeCsv {
		if err := cw.Write(columnNames); err != nil {
			return fmt.Errorf("写入csv头失败：%w", err)
		}
		cw.Flush()
	}
	return nil
}

// processHits 处理一批查询结果
func (e *ElasticsearchExport) processHits(hits []*elasticModel.Hit, columnNames []string, startRow int, streamWriter *excelize.StreamWriter, cw *csv.Writer) error {
	for i, hit := range hits {
		row := make([]string, len(columnNames))
		for j, col := range columnNames {
			if val, ok := hit.Source[col]; ok {
				row[j] = tools.Any2string(val)
			}
		}

		if e.FileType == config.ExportFileTypeExcel {
			// 转换为any类型切片
			rowAny := make([]any, len(row))
			for k, v := range row {
				rowAny[k] = v
			}
			// 计算单元格位置 (A2, A3, ...)
			rowNum := startRow + i
			cell, _ := excelize.CoordinatesToCellName(1, rowNum)
			if err := streamWriter.SetRow(cell, rowAny); err != nil {
				return fmt.Errorf("写入Excel行失败: %w", err)
			}
		} else {
			if err := cw.Write(row); err != nil {
				return fmt.Errorf("写入CSV行失败: %w", err)
			}
			cw.Flush() // 手动刷新，确保数据及时发送到客户端
		}
	}
	return nil
}

// getFieldNames 从文档中获取字段名
func (e *ElasticsearchExport) getFieldNames(doc map[string]interface{}) []string {
	var fields []string
	for field := range doc {
		fields = append(fields, field)
	}
	return fields
}
