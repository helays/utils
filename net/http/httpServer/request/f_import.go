package request

import (
	"bufio"
	"context"
	"fmt"
	"github.com/helays/utils/v2/dataType/customWriter"
	"github.com/helays/utils/v2/excelTools"
	"github.com/helays/utils/v2/tools"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"strings"
)

type Import struct {
	FileType string `json:"file_type"` // 文件类型 excel、csv
	FieldRow int    `json:"field_row"` // 字段所在行
	DataRow  int    `json:"data_row"`  // 数据开始行
	Sep      string `json:"sep"`       // csv 分割符
}

type JSONHandler func(ctx context.Context, obj map[string]interface{}) error

func (i *Import) Import(r *http.Request) ([]any, int64, error) {
	switch i.FileType {
	case "excel":
		return i.ImportExcel(r)
	case "csv":
		return i.ImportCsv(r)
	default:
		return nil, 0, fmt.Errorf("不支持的文件类型：%s", i.FileType)
	}
}

func (i *Import) ImportWithHandler(r *http.Request, handler JSONHandler) (int64, error) {
	switch i.FileType {
	case "excel":
		return i.ImportExcelWithHandler(r, handler)
	case "csv":
		return i.ImportCsvWithHandler(r, handler)
	default:
		return 0, fmt.Errorf("不支持的文件类型：%s", i.FileType)
	}
}

func (i *Import) ImportExcelWithHandler(r *http.Request, handler JSONHandler) (int64, error) {
	if err := i.valid(); err != nil {
		return 0, err
	}
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(r.Body, counter)

	excel, err := excelize.OpenReader(teeReader)
	defer excelTools.CloseExcel(excel)
	if err != nil {
		return counter.TotalSize, fmt.Errorf("excel文件打开失败：%s", err.Error())
	}
	rows, err := excel.GetRows(excel.GetSheetName(0))
	if err != nil {
		return counter.TotalSize, fmt.Errorf("sheet读取失败：%s", err.Error())
	}
	if len(rows) < i.DataRow {
		return counter.TotalSize, fmt.Errorf("未读取到有效数据")
	}
	var (
		dataRow     = i.DataRow - 1
		fieldRowMap = rows[i.FieldRow-1]
	)

	for idx, row := range rows {
		if idx < dataRow {
			continue
		}
		err = handler(r.Context(), tools.Slice2MapWithHeader(row, fieldRowMap))
		if err != nil {
			return counter.TotalSize, err
		}
	}
	return counter.TotalSize, nil
}

func (i *Import) ImportCsvWithHandler(r *http.Request, handler JSONHandler) (int64, error) {
	if err := i.valid(); err != nil {
		return 0, err
	}
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(r.Body, counter)
	var (
		scanner   = bufio.NewScanner(teeReader)
		idx       int
		fieldRows []string
	)
	i.Sep = tools.Ternary(i.Sep == "", ",", i.Sep)
	for scanner.Scan() {
		idx++
		line := scanner.Text()
		lineRows := strings.Split(line, i.Sep)
		if idx == i.FieldRow {
			fieldRows = lineRows
			continue
		}
		if idx < i.DataRow {
			continue
		}
		if err := handler(r.Context(), tools.Slice2MapWithHeader(lineRows, fieldRows)); err != nil {
			return 0, err
		}
	}
	return counter.TotalSize, nil
}

// ImportExcel 获取excel内容
func (i *Import) ImportExcel(r *http.Request) ([]any, int64, error) {
	var data []any
	flow, err := i.ImportExcelWithHandler(r, func(ctx context.Context, obj map[string]interface{}) error {
		data = append(data, obj)
		return nil
	})
	return data, flow, err

}

// ImportCsv 获取csv内容
func (i *Import) ImportCsv(r *http.Request) ([]any, int64, error) {
	var data []any
	flow, err := i.ImportCsvWithHandler(r, func(ctx context.Context, obj map[string]interface{}) error {
		data = append(data, obj)
		return nil
	})
	return data, flow, err
}

// 参数验证
func (i *Import) valid() error {
	if i.FieldRow == 0 || i.FieldRow >= i.DataRow {
		return fmt.Errorf("无有效字段、数据所在的行数")
	}
	return nil
}
