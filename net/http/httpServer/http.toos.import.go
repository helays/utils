package httpServer

import (
	"bufio"
	"fmt"
	"github.com/helays/utils/excelTools"
	"github.com/helays/utils/tools"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strings"
)

type Import struct {
	FileType string `json:"file_type"` // 文件类型 excel、csv
	FieldRow int    `json:"field_row"` // 字段所在行
	DataRow  int    `json:"data_row"`  // 数据开始行
	Sep      string `json:"sep"`       // csv 分割符
}

func (this Import) Import(r *http.Request) ([]any, error) {
	switch this.FileType {
	case "excel":
		return this.ImportExcel(r)
	case "csv":
		return this.ImportCsv(r)
	default:
		return nil, fmt.Errorf("不支持的文件类型：%s", this.FileType)
	}
}

// ImportExcel 获取excel内容
func (this Import) ImportExcel(r *http.Request) ([]any, error) {
	if err := this.valid(); err != nil {
		return nil, err
	}
	excel, err := excelize.OpenReader(r.Body)
	defer excelTools.CloseExcel(excel)
	if err != nil {
		return nil, fmt.Errorf("excel文件打开失败：%s", err.Error())
	}
	rows, err := excel.GetRows(excel.GetSheetName(0))
	if err != nil {
		return nil, fmt.Errorf("sheet读取失败：%s", err.Error())
	}
	if len(rows) < this.DataRow {
		return nil, fmt.Errorf("未读取到有效数据")
	}
	var (
		data        []any
		dataRow     = this.DataRow - 1
		fieldRowMap = rows[this.FieldRow-1]
	)

	for idx, row := range rows {
		if idx < dataRow {
			continue
		}
		data = append(data, tools.Slice2MapWithHeader(row, fieldRowMap))
	}
	return data, nil
}

func (this Import) ImportCsv(r *http.Request) ([]any, error) {
	if err := this.valid(); err != nil {
		return nil, err
	}
	var (
		scanner   = bufio.NewScanner(r.Body)
		idx       int
		fieldRows []string
		data      []any
	)
	this.Sep = tools.Ternary(this.Sep == "", ",", this.Sep)
	for scanner.Scan() {
		idx++
		line := scanner.Text()
		lineRows := strings.Split(line, this.Sep)
		if idx == this.FieldRow {
			fieldRows = lineRows
			continue
		}
		if idx < this.DataRow {
			continue
		}
		data = append(data, tools.Slice2MapWithHeader(lineRows, fieldRows))
	}
	return data, nil
}

// 参数验证
func (this *Import) valid() error {
	if this.FieldRow == 0 || this.FieldRow >= this.DataRow {
		return fmt.Errorf("无有效字段、数据所在的行数")
	}
	return nil
}
