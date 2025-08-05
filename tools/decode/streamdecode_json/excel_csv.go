package streamdecode_json

import (
	"bufio"
	"fmt"
	"github.com/helays/utils/v2/dataType/customWriter"
	"github.com/helays/utils/v2/excelTools"
	"github.com/helays/utils/v2/tools"
	"github.com/xuri/excelize/v2"
	"io"
	"strings"
)

// 参数验证
func (i *Import) valid() error {
	if i.FieldRow == 0 || i.FieldRow >= i.DataRow {
		return fmt.Errorf("无有效字段、数据所在的行数")
	}
	return nil
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
		fieldRowMap = rows[i.FieldRow-1]
	)

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
