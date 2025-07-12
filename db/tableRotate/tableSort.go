package tableRotate

import (
	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

type tableSplit struct {
	tableName  string
	createTime time.Time
}

// 定义切片类型
type byCreateTime []tableSplit

// Len 实现 sort.Interface 接口的 Len 方法
func (b byCreateTime) Len() int {
	return len(b)
}

// Less 实现 sort.Interface 接口的 Less 方法
func (b byCreateTime) Less(i, j int) bool {
	return b[i].createTime.After(b[j].createTime) // 按 createTime 降序排序
}

// Swap 实现 sort.Interface 接口的 Swap 方法
func (b byCreateTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// MaxTableRetention 表自动分割后，保留指定数量的表
type MaxTableRetention struct {
	Tx         *gorm.DB
	MaxNum     int
	TableName  string
	Step       string
	Tables     []string
	DateFormat string
}

// Run 执行删除多余的表
func (this *MaxTableRetention) Run() error {
	if this.MaxNum < 1 {
		return nil
	}
	this.Step = tools.Ternary(this.Step == "", "_", this.Step)
	this.DateFormat = tools.Ternary(this.DateFormat == "", dateFormat, this.DateFormat)
	var ts byCreateTime
	for _, tableName := range this.Tables {
		// 根据去除 tableName中的前缀
		createTime := ""
		currentTable := this.TableName + "_"
		if strings.HasPrefix(tableName, currentTable) {
			createTime = tableName[len(currentTable):]
		}
		if t, err := time.Parse(this.DateFormat, createTime); err == nil {
			ts = append(ts, tableSplit{
				tableName:  tableName,
				createTime: t,
			})
		}
	}
	if len(ts) < this.MaxNum {
		return nil
	}
	// 对 列表按时间降序排序
	// 由于保留的表是最新时间，所以指定位置往后面删除
	sort.Sort(ts)
	for _, item := range ts[this.MaxNum:] {
		err := this.Tx.Debug().Migrator().DropTable(item.tableName) // 删除表，一定要打印日志
		if err != nil {
			return err
		}
	}
	return nil
}
