package httpExportExcel

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/helays/utils/excelTools"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/net/http/httpTools"
	"github.com/xuri/excelize/v2"
	"net/http"
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

// RespExcelOrCsv 将 rows导出成excel 或者 csv
func RespExcelOrCsv(w http.ResponseWriter, rows *sql.Rows, fileType string, args ...map[string]string) error {
	if fileType == "" {
		fileType = "excel"
	} else if fileType != "excel" && fileType != "csv" {
		return fmt.Errorf("不支持的导出类型：%s", fileType)
	}
	var (
		ii           int
		exportHeader map[string]string
		f            *excelize.File
		cw           *csv.Writer
		sheetName    = "Sheet1"
	)
	w.Header().Del("Accept-Ranges")
	if fileType == "excel" {
		f = excelize.NewFile()
		defer excelTools.CloseExcel(f)
		if _, err := f.NewSheet(sheetName); err != nil {
			return fmt.Errorf("创建sheet失败：%w", err)
		}
	} else {
		cw = csv.NewWriter(w)
		defer cw.Flush()
		w.Header().Set("Content-Type", "text/csv")
		httpTools.SetDisposition(w, "export.csv")
	}
	if len(args) > 0 {
		exportHeader = args[0]
	}
	columnNames, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("获取表头字段失败:%w", err)
	}
	if fileType == "excel" {
		for i, k := range columnNames {
			colIndex, _ := excelize.ColumnNumberToName(i + 1)
			// 获取字段名作为表头，并设置到对应的单元格
			header, ok := exportHeader[k]
			if !ok {
				header = k
			}
			ulogs.Checkerr(f.SetCellValue(sheetName, fmt.Sprintf("%s1", colIndex), header), "导出excel失败，表头写入失败")

		}
	} else {
		// 写入CSV头
		_, _ = w.Write([]byte("\xef\xbb\xbf"))
		if err = cw.Write(columnNames); err != nil {
			return fmt.Errorf("写入csv头失败：%w", err)
		}
	}
	var values = make([]sql.RawBytes, len(columnNames))
	scanArgs := make([]any, len(columnNames))
	for i := range scanArgs {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			ulogs.Error(err, "导出数据报错")
			continue
		}
		// 将[]sql.RawBytes转换为[]string
		row := make([]string, len(values))
		for i, value := range values {
			row[i] = string(value)
		}
		ii++
		if fileType == "excel" {
			for i, v := range row {
				colIndex, _ := excelize.ColumnNumberToName(i + 1)
				ulogs.Checkerr(f.SetCellValue(sheetName, fmt.Sprintf("%s%d", colIndex, ii+1), v), "导出excel失败，写入数据行失败")
			}
		} else {
			// 写入当前行到CSV
			ulogs.Checkerr(cw.Write(row), "导出csv失败，写入数据行失败")
			// 手动刷新，确保数据及时发送到客户端
			w.(http.Flusher).Flush()
		}
	}

	if fileType == "excel" {
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		httpTools.SetDisposition(w, "export.xlsx")
		if err := f.Write(w); err != nil {
			return fmt.Errorf("导出excel失败：%w", err)
		}
	}
	return nil
}
