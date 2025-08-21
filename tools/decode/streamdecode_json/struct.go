package streamdecode_json

import (
	"context"
	"fmt"
	"io"

	"github.com/helays/utils/v2/dataType"
)

type JSONHandler func(ctx context.Context, obj map[string]interface{}) error

type Import struct {
	fileType dataType.ContentType
	ctx      context.Context
	rd       io.Reader
	bigLine  bool   // 是否启用大行模式
	FieldRow int    `json:"field_row"` // 字段所在行
	DataRow  int    `json:"data_row"`  // 数据开始行
	Sep      string `json:"sep"`       // csv 分割符
}

func New(ctx context.Context, ft dataType.ContentType, rd io.Reader) *Import {
	return &Import{
		ctx:      ctx,
		fileType: ft,
		rd:       rd,
	}
}

func (i *Import) SetFieldRow(fieldRow int) {
	i.FieldRow = fieldRow
}

func (i *Import) SetDataRow(dataRow int) {
	i.DataRow = dataRow
}

func (i *Import) SetSep(sep string) {
	i.Sep = sep
}

func (i *Import) SetBigLine(bigLine bool) {
	i.bigLine = bigLine
}

func (i *Import) Import() ([]any, int64, error) {
	var data []any
	flow, err := i.ImportWithHandler(func(ctx context.Context, obj map[string]interface{}) error {
		data = append(data, obj)
		return nil
	})
	return data, flow, err
}

func (i *Import) ImportWithHandler(handler JSONHandler) (int64, error) {
	switch i.fileType {
	case dataType.PlainText:
		if i.bigLine {
			return i.processLineEnhanced(handler)
		}
		return i.processLine(handler)
	case dataType.JSON:
		return i.processJSON(handler)
	case dataType.Excel:
		return i.importExcelWithHandler(handler)
	case dataType.CSV:
		return i.importCsvWithHandler(handler)
	default:
		return 0, fmt.Errorf("不支持的文件类型：%s", i.fileType)
	}
}
