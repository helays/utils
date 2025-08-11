package tableRotate

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/db/userDb"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/tools/retention"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const dateFormat = "20060102150405"

var crond *cron.Cron

func init() {
	crond = cron.New()
	crond.Start()
}

func Close() error {
	crond.Stop()
	return nil
}

// TableRotate db自动轮转配置
type TableRotate struct {
	Enable                  bool          `json:"enable" yaml:"enable" ini:"enable"` // 是否启用自动轮转
	Duration                time.Duration `json:"duration" yaml:"duration" ini:"duration"`
	Crontab                 string        `json:"crontab" yaml:"crontab" ini:"crontab"`                                                          // crontab 表达式 ,定时器和crontab二选一
	SplitTable              bool          `json:"split_table" yaml:"split_table" ini:"split_table"`                                              // 是否开启按天切分日志 ，开启后，自动回收数据 只会看表的保留数量，不开启，就看数据保留时长
	MaxTableRetention       int           `json:"max_table_retention" yaml:"max_table_retention" ini:"max_table_retention"`                      // 分表后，最大保留表的数量 -1 不限制
	TimeFormat              string        `json:"time_format" yaml:"time_format" ini:"time_format"`                                              // 时间格式
	SeqFields               []string      `json:"seq_fields" yaml:"seq_fields" ini:"seq_fields"`                                                 // 序列字段
	DataRetentionPeriod     int           `json:"data_retention_period" yaml:"data_retention_period" ini:"data_retention_period"`                // 数据保留时长 -1 不限制
	DataRetentionPeriodUnit string        `json:"data_retention_period_unit" yaml:"data_retention_period_unit" ini:"data_retention_period_unit"` // 数据保留时间单位 支持 second minute hour day month year
	FilterField             string        `json:"filter_field" yaml:"filter_field" ini:"filter_field"`                                           // 过滤字段 默认create_time
	tx                      *gorm.DB
	tableName               string
}

// AddTask 添加自动轮转任务
func (r TableRotate) AddTask(ctx context.Context, tx *gorm.DB, tableName string) {
	if !r.Enable {
		return
	}
	if r.Duration <= 0 && r.Crontab == "" {
		return
	}
	ulogs.Log("【表自动轮转配置】", "数据库", tx.Dialector.Name(), tx.Migrator().CurrentDatabase(), tableName)
	ulogs.Log("【表自动轮转配置】", "周期策略", r.Crontab, r.Duration)
	if r.SplitTable {
		ulogs.Log("【表自动轮转配置】", "回收策略：", "分表", "最大保留数量", r.MaxTableRetention)
	} else {
		ulogs.Log("【表自动轮转配置】", "回收策略：", "数据", "数据保留时长", r.DataRetentionPeriod, r.DataRetentionPeriodUnit)
	}
	r.tx = tx
	r.tableName = tableName
	if r.TimeFormat == "" {
		r.TimeFormat = dateFormat
	}
	if r.Crontab != "" {
		go r.toCrontab(ctx)
		return
	}
	go r.toTicker(ctx)
}

// 通过 crontab方式运行
func (r *TableRotate) toCrontab(ctx context.Context) {
	eid, err := crond.AddFunc(r.Crontab, r.run)
	if err != nil {
		ulogs.Error("添加自动轮转任务失败", "表", r.tableName, "定时", r.Crontab)
		return
	}
	go func() {
		<-ctx.Done()      // 等待上下文取消
		crond.Remove(eid) // 移除任务
		ulogs.Log("【表自动轮转配置终止】", "crontab", "数据库", r.tx.Dialector.Name(), r.tx.Migrator().CurrentDatabase(), r.tableName)
	}()
}

// 通过 定时器方式运行
func (r *TableRotate) toTicker(ctx context.Context) {
	tck := time.NewTicker(r.Duration)
	defer tck.Stop()
	for {
		select {
		case <-ctx.Done():
			ulogs.Log("【表自动轮转配置终止】", "定时器", "数据库", r.tx.Dialector.Name(), r.tx.Migrator().CurrentDatabase(), r.tableName)
			return
		case <-tck.C:
			r.run()
		}
	}
}

func (r *TableRotate) run() {
	if r.SplitTable {
		r.runSplitTable()
		return
	}
	r.runRotateTableData()
}

// 分表
func (r *TableRotate) runSplitTable() {
	newTableName := r.tableName + "_" + time.Now().Format(r.TimeFormat)
	err := r.tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Migrator().RenameTable(r.tableName, newTableName)
		if err != nil {
			return fmt.Errorf("修改表名失败 %s to %s :%s", r.tableName, newTableName, err.Error())
		}
		switch tx.Dialector.Name() {
		case config.DbTypePostgres:
			// 创建新表
			err = tx.Debug().Exec("CREATE TABLE ? (LIKE ? INCLUDING ALL)", clause.Table{Name: r.tableName}, clause.Table{Name: newTableName}).Error
			if err != nil {
				return fmt.Errorf("创建表失败 %s :%s", r.tableName, err.Error())
			}
			// 清理新表的序列字段的默认值，不清理在删表的时候会失败
			if err = userDb.ClearSequenceFieldDefaultValue(tx, newTableName, r.SeqFields); err != nil {
				return err
			}
			// 创建表后，需要将改表后的序列清除掉
		case config.DbTypeMysql:
			err = tx.Debug().Exec("CREATE TABLE ? LIKE ?", clause.Table{Name: r.tableName}, clause.Table{Name: newTableName}).Error
			if err != nil {
				return fmt.Errorf("创建表失败 %s :%s", r.tableName, err.Error())
			}
		}

		return nil
	})
	if err != nil {
		ulogs.Error("自动轮转表，修改表名失败", r.tableName, "新表名", newTableName, err)
	} else {
		ulogs.Log("自动轮转表", r.tableName, "修改表名成功", "新表名", newTableName)
	}
	// 如果最大保留数量为0，就不会清理表
	if r.MaxTableRetention <= 0 {
		return
	}
	// 查询以 this.tableName开头的表名
	tables, _err := userDb.FindTableWithPrefix(r.tx, r.tableName)
	if _err != nil {
		ulogs.Error("自动轮转表：%s", _err.Error())
		return
	}

	tableManager := retention.New(retention.Config{
		MaxRetain:  r.MaxTableRetention,
		Prefix:     r.tableName,
		Delimiter:  "_",
		TimeFormat: r.TimeFormat,
		Order:      false, // 按时间降序排序
		Criteria:   retention.ByTime,
	})
	tableManager.AddItems(tables)
	err = tableManager.Run(func(name string) error {
		return r.tx.Migrator().DropTable(name)
	})
	if err != nil {
		ulogs.Error("自动轮转表，删除表失败", err)
	}
}

// 回收表数据
func (r *TableRotate) runRotateTableData() {
	if r.DataRetentionPeriod <= 0 {
		return
	}
	var queryVal clause.Expr
	unit := strings.ToUpper(r.DataRetentionPeriodUnit)
	retentionPeriod := strconv.Itoa(r.DataRetentionPeriod)
	switch r.tx.Dialector.Name() {
	case config.DbTypeMysql:
		queryVal = clause.Expr{
			SQL:                "NOW() - INTERVAL '? ?'",
			Vars:               []any{clause.Column{Name: retentionPeriod, Raw: true}, clause.Column{Name: unit, Raw: true}},
			WithoutParentheses: false,
		}
	case config.DbTypePostgres:
		queryVal = clause.Expr{
			SQL:                "NOW() - INTERVAL '? ?'",
			Vars:               []any{clause.Column{Name: retentionPeriod, Raw: true}, clause.Column{Name: unit, Raw: true}},
			WithoutParentheses: false,
		}
	}
	err := r.tx.Table(r.tableName).Where(r.FilterField+" < ?", queryVal).Delete(nil).Error
	if err != nil {
		switch _err := err.(type) {
		case *pgconn.PgError:
			if _err.Code == "42P01" {
				return
			}
		case *mysql.MySQLError:

		default:
			fmt.Println("fadsf", _err)
		}
		ulogs.Error("自动轮转表，回收表数据失败", r.tableName, "过滤字段", r.FilterField, "条件", retentionPeriod, unit, err)
	} else {
		ulogs.Log("自动轮转表", r.tableName, "回收表数据成功", "过滤字段", r.FilterField, "条件", retentionPeriod, unit)
	}

}
