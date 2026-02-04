package gormkit

import (
	"database/sql"
	"fmt"

	"github.com/xuri/excelize/v2"
	"helay.net/go/utils/v3/close/vclose"
	"helay.net/go/utils/v3/db/exportkit"
)

func (e *Export) cursor(et exportkit.ExportType, columns []string) error {
	rows, err := e.db.Rows()
	if err != nil {
		return err
	}
	defer vclose.Close(rows)
	var values = make([]sql.RawBytes, len(columns))
	scanVal := make([]any, len(columns))
	for i := range scanVal {
		scanVal[i] = &values[i]
	}
	var rowIdx int
	for rows.Next() {
		if err = rows.Scan(scanVal...); err != nil {
			return fmt.Errorf("扫描数据失败：%v", err)
		}
		rowIdx++
		switch et {
		case exportkit.ExportTypeExcel:
			row := make([]any, len(values))
			for i, value := range values {
				row[i] = string(value)
			}
			cell, _ := excelize.CoordinatesToCellName(1, rowIdx+1)
			if err = e.sw.SetRow(cell, row); err != nil {
				return fmt.Errorf("写入数据失败：%v", err)
			}
		case exportkit.ExportTypeCsv:
			row := make([]string, len(values))
			for i, value := range values {
				row[i] = string(value)
			}
			if err = e.cw.Write(row); err != nil {
				return fmt.Errorf("写入数据失败：%v", err)
			}
		default:
			return exportkit.ErrExportTypeNotSupport
		}
	}
	return nil
}
