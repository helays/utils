package streamdecode_json

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/helays/utils/v2/dataType/customWriter"
	"github.com/helays/utils/v2/excelTools"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/tools"
	"github.com/xuri/excelize/v2"
)

// 参数验证
func (i *Import) valid() error {

	if (i.FieldRow == 0 && (i.fields == nil || len(i.fields) == 0)) || i.FieldRow >= i.DataRow {
		return fmt.Errorf("无有效字段、数据所在的行数")
	}
	return nil
}

func (i *Import) importExcelWithHandlerStream(handler JSONHandler) (int64, error) {
	start := time.Now()
	ulogs.Debugf("开始处理excel上传数据")
	if err := i.valid(); err != nil {
		return 0, err
	}

	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(i.rd, counter)
	excel, err := excelize.OpenReader(teeReader)
	defer excelTools.CloseExcel(excel)
	if err != nil {
		return counter.TotalSize, fmt.Errorf("excel文件打开失败：%s", err.Error())
	}
	ulogs.Debugf("excel文件打开完成 耗时 %.2f秒", time.Since(start).Seconds())
	rows, err := excel.Rows(excel.GetSheetName(0))
	if err != nil {
		return counter.TotalSize, fmt.Errorf("sheet读取失败：%s", err.Error())
	}
	ulogs.Debugf("sheet读取完成 耗时 %.2f秒", time.Since(start).Seconds())
	var fieldRowMap = i.fields
	idx := 0 // 为了方便，就直接从1开始计数
	for rows.Next() {
		idx++
		if i.FieldRow > 0 && idx == i.FieldRow {
			r, e := rows.Columns()
			if e != nil {
				return counter.TotalSize, fmt.Errorf("sheet读取字段行失败：%s", e.Error())
			}
			fieldRowMap = r
			continue
		}
		if idx >= i.DataRow && idx != i.FieldRow {
			r, e := rows.Columns()
			if e != nil {
				return counter.TotalSize, fmt.Errorf("sheet读取数据行失败：%s", e.Error())
			}
			err = handler(i.ctx, tools.Slice2MapWithHeader(r, fieldRowMap))
			if err != nil {
				return counter.TotalSize, err
			}
		}
	}
	ulogs.Debugf("excel文件处理完成 耗时 %.2f秒", time.Since(start).Seconds())
	return counter.TotalSize, nil
}

func (i *Import) importExcelWithHandler(handler JSONHandler) (int64, error) {
	if err := i.valid(); err != nil {
		return 0, err
	}
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(i.rd, counter)
	excel, err := excelize.OpenReader(teeReader)
	defer excelTools.CloseExcel(excel)
	if err != nil {
		return counter.TotalSize, fmt.Errorf("excel文件打开失败：%s", err.Error())
	}
	rows, err := excel.GetRows(excel.GetSheetName(0))
	if err != nil {
		return counter.TotalSize, fmt.Errorf("sheet读取失败：%s", err.Error())
	}
	if len(rows) < i.DataRow {
		return counter.TotalSize, fmt.Errorf("未读取到有效数据")
	}
	var (
		dataRow     = i.DataRow - 1
		fieldRowMap = i.fields
	)
	if i.FieldRow > 0 {
		fieldRowMap = rows[i.FieldRow-1]
	}
	for idx, row := range rows {
		if idx < dataRow {
			continue
		}
		err = handler(i.ctx, tools.Slice2MapWithHeader(row, fieldRowMap))
		if err != nil {
			return counter.TotalSize, err
		}
	}
	return counter.TotalSize, nil
}

func (i *Import) importCsvWithHandler(handler JSONHandler) (int64, error) {
	if err := i.valid(); err != nil {
		return 0, err
	}
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(i.rd, counter)
	var (
		scanner   = bufio.NewScanner(teeReader)
		idx       int
		fieldRows []string
	)
	i.Sep = tools.Ternary(i.Sep == "", ",", i.Sep)
	for scanner.Scan() {
		idx++
		line := scanner.Text()
		lineRows := strings.Split(line, i.Sep)
		if idx == i.FieldRow {
			fieldRows = lineRows
			continue
		}
		if idx < i.DataRow {
			continue
		}
		if err := handler(i.ctx, tools.Slice2MapWithHeader(lineRows, fieldRows)); err != nil {
			return 0, err
		}
	}
	return counter.TotalSize, nil
}
