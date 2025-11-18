package exportkit

import (
	"errors"
)

type ExportType string

func (t ExportType) String() string {
	return string(t)
}

const (
	ExportTypeExcel = "excel"
	ExportTypeCsv   = "csv"
)

// IterationType 迭代方式枚举
type IterationType int

const (
	IterationCursor IterationType = iota // 游标迭代（推荐大数据量）
	IterationBatch                       // 批量查询
	IterationAll                         // 一次性查询
)

var ErrExportTypeNotSupport = errors.New("不支持的导出类型")

type ExportConfig struct {
	IterationType IterationType
	BatchSize     int // 批量大小
	SheetName     string
}

type Callback func(et ExportType) error
