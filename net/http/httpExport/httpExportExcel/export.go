package httpExportExcel

import (
	"net/http"

	"github.com/helays/utils/v2/db/exportkit"
	"github.com/helays/utils/v2/db/exportkit/gormkit"
	"github.com/helays/utils/v2/net/http/httpkit"
	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm"
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
	DB           *gorm.DB
	FileType     exportkit.ExportType
	FileName     string
	ExportHeader map[string]string
}

func (e *RowsExport) Response(w http.ResponseWriter) error {
	export := gormkit.New(e.DB, exportkit.ExportConfig{
		IterationType: exportkit.IterationCursor, // 这里还是适合用游标方式默认导出
		BatchSize:     1000,
		SheetName:     "Sheet1",
	})
	export.SetHeader(e.ExportHeader)
	if e.FileType == "" {
		e.FileType = exportkit.ExportTypeCsv // 默认为csv
	}
	e.FileName = tools.Ternary(e.FileName == "", "export", e.FileName)
	w.Header().Del("Accept-Ranges")
	// 设置csv的响应头
	export.OnInitColumn(func(et exportkit.ExportType) error {
		if et == exportkit.ExportTypeCsv {
			w.Header().Set("Content-Type", "text/csv")
			httpkit.SetDisposition(w, e.FileName+".csv")
		}
		return nil
	})
	// 设置excel的响应头
	export.BeforeFinalize(func(et exportkit.ExportType) error {
		if et == exportkit.ExportTypeExcel {
			w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			httpkit.SetDisposition(w, e.FileName+".xlsx")
		}
		return nil
	})

	return export.Execute(e.FileType, w)
}
