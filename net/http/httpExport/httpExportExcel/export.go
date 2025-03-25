package httpExportExcel

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/helays/utils/config"
	"github.com/helays/utils/excelTools"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/net/http/httpTools"
	"github.com/helays/utils/tools"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strings"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2024/11/23 15:42
//

type RowsExport struct {
	Rows         *sql.Rows
	FileType     string
	FileName     string
	ExportHeader map[string]string
}

func (this *RowsExport) Response(w http.ResponseWriter) error {
	var (
		ii           int
		f            *excelize.File
		streamWriter *excelize.StreamWriter
		cw           *csv.Writer
		sheetName    = "Sheet1"
	)
	this.FileType = strings.ToLower(tools.Ternary(this.FileType == "", config.ExportFileTypeExcel, this.FileType))
	this.FileName = tools.Ternary(this.FileName == "", "export", this.FileName)
	columnNames, err := this.Rows.Columns()
	if err != nil {
		return fmt.Errorf("获取表头字段失败:%w", err)
	}
	for idx, field := range columnNames {
		// 获取字段名作为表头，并设置到对应的单元格
		if header, ok := this.ExportHeader[field]; ok {
			columnNames[idx] = header
		}
	}
	w.Header().Del("Accept-Ranges")
	if this.FileType == config.ExportFileTypeExcel {
		f = excelize.NewFile()
		defer excelTools.CloseExcel(f)
		if streamWriter, err = f.NewStreamWriter(sheetName); err != nil {
			return fmt.Errorf("创建sheet失败：%w", err)
		}
		ulogs.Checkerr(streamWriter.SetRow("A1", tools.StrSlice2AnySlice(columnNames)), "导出excel失败，表头写入失败")
	} else if this.FileType == config.ExportFileTypeCsv {
		w.Header().Set("Content-Type", "text/csv")
		httpTools.SetDisposition(w, this.FileName+".csv")
		cw = csv.NewWriter(w)
		defer cw.Flush()
		// 写入CSV头
		_, _ = w.Write([]byte("\xef\xbb\xbf"))
		if err = cw.Write(columnNames); err != nil {
			return fmt.Errorf("写入csv头失败：%w", err)
		}
	} else {
		return fmt.Errorf("不支持的导出类型：%s", this.FileType)
	}

	var values = make([]sql.RawBytes, len(columnNames))
	scanArgs := make([]any, len(columnNames))
	for i := range scanArgs {
		scanArgs[i] = &values[i]
	}
	for this.Rows.Next() {
		if err = this.Rows.Scan(scanArgs...); err != nil {
			ulogs.Error(err, "导出数据报错")
			continue
		}
		// 将[]sql.RawBytes转换为[]string
		ii++
		if this.FileType == config.ExportFileTypeExcel {
			row := make([]any, len(values))
			for i, value := range values {
				row[i] = string(value)
			}
			// 计算单元格位置 (A2, A3, ...)
			rowId := ii + 1
			cell, _ := excelize.CoordinatesToCellName(1, rowId)
			_ = streamWriter.SetRow(cell, row)
		} else {
			row := make([]string, len(values))
			for i, value := range values {
				row[i] = string(value)
			}
			_ = cw.Write(row) // 写入当前行到CSV
			cw.Flush()        // 手动刷新，确保数据及时发送到客户端
		}
	}

	if this.FileType == config.ExportFileTypeExcel {
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		httpTools.SetDisposition(w, this.FileName+".xlsx")
		_ = streamWriter.Flush()
		if err := f.Write(w); err != nil {
			return fmt.Errorf("导出excel失败：%w", err)
		}
	}
	return nil
}
