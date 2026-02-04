package gormkit

import (
	"fmt"

	"github.com/helays/utils/v2/db/exportkit"
	"github.com/helays/utils/v2/tools"
	"github.com/xuri/excelize/v2"
)

func (e *Export) all(et exportkit.ExportType, columns []string) error {
	var records []map[string]any
	if err := e.db.Find(&records).Error; err != nil {
		return fmt.Errorf("查询数据失败：%v", err)
	}
	// 遍历数据
	for rowIdx, record := range records {
		switch et {
		case exportkit.ExportTypeExcel:
			row := make([]any, len(columns))
			for colIdx, columnName := range columns {
				row[colIdx] = tools.Any2string(record[columnName])
			}
			cell, _ := excelize.CoordinatesToCellName(1, rowIdx+2)
			if err := e.sw.SetRow(cell, row); err != nil {
				return fmt.Errorf("写入数据失败：%v", err)
			}
		case exportkit.ExportTypeCsv:
			row := make([]string, len(columns))
			for colIdx, columnName := range columns {
				row[colIdx] = tools.Any2string(record[columnName])
			}
			if err := e.cw.Write(row); err != nil {
				return fmt.Errorf("写入数据失败：%v", err)
			}
		default:
			return exportkit.ErrExportTypeNotSupport
		}
	}
	return nil
}
