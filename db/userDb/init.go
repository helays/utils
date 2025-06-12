package userDb

import (
	"errors"
	"fmt"
	"github.com/helays/utils/config"
	"github.com/helays/utils/db"
	"github.com/helays/utils/logger/zaploger"
	"github.com/helays/utils/tools"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

// InitDb 连接数据库
func InitDb(c db.Dbbase) (*gorm.DB, error) {
	var (
		dsn       = c.Dsn()
		dialector gorm.Dialector
		err       error
	)
	switch c.DbType {
	case config.DbTypePostgres, config.DbTypePg:
		//postgres://user:password@host1:port1/database?target_session_attrs=read-write&TimeZone=Asia/Shanghai
		//dsn = "postgres://" + c.User + ":" + c.Pwd + "@" + strings.Join(c.Host, ",") + "/" + c.Dbname + "?TimeZone=Asia/Shanghai"
		dialector = postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		})
	case config.DbTypeMysql:
		//dsn = c.User + ":" + c.Pwd + "@tcp(" + strings.Join(c.Host, ",") + ")/" + c.Dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
		//fmt.Println(dsn)
		dialector = mysql.New(mysql.Config{
			DSN:                       dsn,
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		})
	default:
		return nil, errors.New("不支持的数据库")
	}
	namingStrategy := schema.NamingStrategy{}
	if c.TablePrefix != "" {
		namingStrategy.TablePrefix = c.TablePrefix
	}
	namingStrategy.SingularTable = c.SingularTable == 1
	lger := logger.Default.LogMode(logger.Silent)
	if c.Logger.LogLevelConfigs != nil {
		_logger := zaploger.Config{
			ConsoleSeparator: c.Logger.ConsoleSeparator,
			LogFormat:        c.Logger.LogFormat,
			LogLevel:         c.Logger.LogLevel,
			LogLevelConfigs:  make(map[string]zaploger.LogConfig),
		}
		for k, v := range c.Logger.LogLevelConfigs {
			_logger.LogLevelConfigs[k] = v
		}
		for level, cfg := range _logger.LogLevelConfigs {
			if cfg.FileName == "" {
				cfg.FileName = fmt.Sprintf("%s_%s", strings.ReplaceAll(c.Host[0], ":", "_"), c.Dbname)
				if c.Schema != "" {
					cfg.FileName += "_" + c.Schema
				}
				_logger.LogLevelConfigs[level] = cfg
			}
		}
		lger, err = zaploger.New(&_logger)
		if err != nil {
			return nil, fmt.Errorf("日志初始化失败:%s", err)
		}
	}
	cfg := gorm.Config{
		SkipDefaultTransaction:                   true,
		Logger:                                   lger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy:                           namingStrategy,
	}
	_db, err := gorm.Open(dialector, &cfg)
	if err != nil {
		return nil, err
	}
	_sqlDb, err := _db.DB()
	if err != nil {
		CloseGormDb(_db)
		return nil, err
	}
	_sqlDb.SetMaxIdleConns(tools.Ternary(c.MaxIdleConns < 1, 2, c.MaxIdleConns))      // 设置连接池中空闲连接的最大数量
	_sqlDb.SetMaxOpenConns(tools.Ternary(c.MaxOpenConns < 1, 5, c.MaxOpenConns))      // 设置打开数据库连接的最大数量
	_sqlDb.SetConnMaxLifetime(tools.AutoTimeDuration(c.MaxConnLifetime, time.Second)) // 连接的总生存时间，从连接创建开始计时
	_sqlDb.SetConnMaxIdleTime(tools.AutoTimeDuration(c.MaxConnIdleTime, time.Second)) // 连接的空闲时间，从连接变为空闲开始计时
	return _db, nil
}

func CloseGormDb(_db *gorm.DB) {
	if _db == nil {
		return
	}
	if sqlDb, err := _db.DB(); err == nil {
		CloseDb(sqlDb)
	}
}
