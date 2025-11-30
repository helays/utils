package httpExportExcel

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/helays/utils/v2/db/elastic/elasticModel"
	"net/http"
	"strings"
	"time"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/excelTools"
	"github.com/helays/utils/v2/net/http/httpkit"
	"github.com/helays/utils/v2/tools"
	"github.com/xuri/excelize/v2"
)

// ElasticsearchExport 用于导出Elasticsearch数据的结构体
type ElasticsearchExport struct {
	*elasticModel.EsScroll
	FileType     string            // 文件类型 (excel/csv)
	FileName     string            // 文件名
	ExportHeader map[string]string // 字段名映射
}

// Response 导出Elasticsearch数据到HTTP响应
func (e *ElasticsearchExport) Response(ctx context.Context, w http.ResponseWriter) error {
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
		rowCount     = 2 // 从第二行开始写入数据(第一行是表头)
		columnNames  []string
	)

	if e.FileType == config.ExportFileTypeExcel {
		f = excelize.NewFile()
		defer excelTools.CloseExcel(f)
		if streamWriter, err = f.NewStreamWriter(sheetName); err != nil {
			return fmt.Errorf("创建sheet失败：%w", err)
		}
	} else if e.FileType == config.ExportFileTypeCsv {
		w.Header().Set("Content-Type", "text/csv")
		httpkit.SetDisposition(w, e.FileName+".csv")
		cw = csv.NewWriter(w)
		defer cw.Flush()
		// 写入UTF-8 BOM头
		_, _ = w.Write([]byte("\xef\xbb\xbf"))
	} else {
		return fmt.Errorf("不支持的导出类型：%s", e.FileType)
	}

	firstBatch := true
	err = e.EsScroll.DoESSearchWithScroll(ctx, func(hits []*elasticModel.Hit) error {
		if firstBatch {
			// 获取字段名
			if len(hits) > 0 {
				columnNames = e.getFieldNames(hits[0].Source)
			} else {
				return fmt.Errorf("没有找到可导出的数据")
			}
			// 写入表头
			if err = e.writeHeader(columnNames, streamWriter, cw); err != nil {
				return err
			}
			firstBatch = false
		}
		// 处理当前批次结果
		if err = e.processHits(hits, columnNames, rowCount, streamWriter, cw); err != nil {
			return err
		}
		rowCount += len(hits)
		return nil
	})
	if err != nil {
		return err
	}
	// 完成Excel导出
	if e.FileType == config.ExportFileTypeExcel {
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		httpkit.SetDisposition(w, e.FileName+".xlsx")
		_ = streamWriter.Flush()
		_ = f.Write(w)
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
