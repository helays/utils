package gormkit

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/db/exportkit"
	"github.com/helays/utils/v2/excelTools"
	"github.com/helays/utils/v2/tools"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type Export struct {
	db             *gorm.DB
	header         map[string]string
	config         exportkit.ExportConfig
	excel          *excelize.File
	sw             *excelize.StreamWriter
	cw             *csv.Writer
	onInitColumn   exportkit.Callback
	beforeFinalize exportkit.Callback
}

// New 创建一个导出对象
func New(db *gorm.DB, config ...exportkit.ExportConfig) *Export {
	e := &Export{
		db: db,
		config: exportkit.ExportConfig{
			IterationType: exportkit.IterationBatch,
			BatchSize:     1000,
			SheetName:     "Sheet1",
		},
	}
	if len(config) > 0 {
		e.config = config[0]
	}
	return e
}

func (e *Export) SetHeader(header map[string]string) {
	e.header = header
}

// OnInitColumn 设置初始化的回调函数
func (e *Export) OnInitColumn(callback exportkit.Callback) {
	e.onInitColumn = callback
}

// BeforeFinalize 设置最终处理回调函数
func (e *Export) BeforeFinalize(callback exportkit.Callback) {
	e.beforeFinalize = callback
}

func (e *Export) close() {
	excelTools.CloseExcel(e.excel)
}

// Execute 执行导出
func (e *Export) Execute(et exportkit.ExportType, f io.Writer) error {
	// 第一步，查询表头字段
	columns, err := e.getColumns()
	if err != nil {
		return fmt.Errorf("获取表头字段失败:%v", err)
	}
	if e.onInitColumn != nil {
		if err = e.onInitColumn(et); err != nil {
			return err
		}
	}
	defer e.close()
	if err = e.setFile(et, f, columns); err != nil {
		return err
	}
	// 根据数据获取方式进行导出
	switch e.config.IterationType {
	case exportkit.IterationCursor: // 游标方式导出
		err = e.cursor(et, columns)
	case exportkit.IterationBatch: // 批量查询方式
		err = e.batch(et, columns)
	case exportkit.IterationAll: // 直接全部查询出来
		err = e.all(et, columns)
	default: // 批量查询
		err = e.batch(et, columns)
	}
	if err != nil {
		return err
	}
	if e.beforeFinalize != nil {
		if err = e.beforeFinalize(et); err != nil {
			return err
		}
	}
	// 输出文件
	switch et {
	case exportkit.ExportTypeExcel:
		_ = e.sw.Flush()
		return e.excel.Write(f) // 写入文件
	case exportkit.ExportTypeCsv:
		e.cw.Flush() // csv就再同步一次文件
	}
	return nil
}

// 获取表头字段
func (e *Export) getColumns() ([]string, error) {
	rows, err := e.db.Rows()
	if err != nil {
		return nil, err
	}
	defer vclose.Close(rows)
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("查询字段失败：%s", err.Error())
	}
	// 表头处理替换
	for idx, columnName := range columnNames {
		if name, ok := e.header[columnName]; ok {
			columnNames[idx] = name
		}
	}
	return columnNames, nil
}

func (e *Export) setFile(et exportkit.ExportType, f io.Writer, columns []string) error {
	switch et {
	case exportkit.ExportTypeCsv:
		_, _ = f.Write([]byte("\xef\xbb\xbf"))
		e.cw = csv.NewWriter(f)
	case exportkit.ExportTypeExcel:
		e.excel = excelize.NewFile()
		sw, err := e.excel.NewStreamWriter(e.config.SheetName)
		if err != nil {
			return fmt.Errorf("创建流写入器失败：%s", err.Error())
		}
		e.sw = sw
	default:
		return exportkit.ErrExportTypeNotSupport
	}
	return e.setFileHeader(et, columns)
}

// 设置文件header
func (e *Export) setFileHeader(et exportkit.ExportType, columns []string) error {
	switch et {
	case exportkit.ExportTypeCsv:
		return e.cw.Write(columns)
	case exportkit.ExportTypeExcel:
		return e.sw.SetRow("A1", tools.StrSlice2AnySlice(columns))
	default:
		return exportkit.ErrExportTypeNotSupport
	}
}
