package gormkit

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"helay.net/go/utils/v3/db/exportkit"
	"helay.net/go/utils/v3/tools"
)

func (e *Export) batch(et exportkit.ExportType, columns []string) error {
	var records []map[string]any
	rowIdx := 0
	return e.db.FindInBatches(&records, e.config.BatchSize, func(tx *gorm.DB, batch int) error {
		for _, record := range records {
			rowIdx++
			switch et {
			case exportkit.ExportTypeExcel:
				row := make([]any, len(columns))
				for i, column := range columns {
					row[i] = tools.Any2string(record[column])
				}
				cell, _ := excelize.CoordinatesToCellName(1, rowIdx+2)
				if err := e.sw.SetRow(cell, row); err != nil {
					return fmt.Errorf("写入数据失败：%v", err)
				}
			case exportkit.ExportTypeCsv:
				row := make([]string, len(columns))
				for i, column := range columns {
					row[i] = tools.Any2string(record[column])
				}
				if err := e.cw.Write(row); err != nil {
					return fmt.Errorf("写入数据失败：%v", err)
				}
			default:
				return exportkit.ErrExportTypeNotSupport
			}
		}
		return nil
	}).Error
}
