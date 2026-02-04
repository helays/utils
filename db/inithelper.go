package db

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"helay.net/go/utils/v3/close/gormClose"
	"helay.net/go/utils/v3/logger/zaploger"
	"helay.net/go/utils/v3/tools"
	"strings"
	"time"
)

func (c *Dbbase) Connect(dialector *gorm.Dialector) (*gorm.DB, error) {
	var err error
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
	_db, err := gorm.Open(*dialector, &cfg)
	if err != nil {
		return nil, err
	}
	_sqlDb, err := _db.DB()
	if err != nil {
		gormClose.Close(_db)
		return nil, err
	}
	_sqlDb.SetMaxIdleConns(tools.Ternary(c.MaxIdleConns < 1, 2, c.MaxIdleConns))      // 设置连接池中空闲连接的最大数量
	_sqlDb.SetMaxOpenConns(tools.Ternary(c.MaxOpenConns < 1, 5, c.MaxOpenConns))      // 设置打开数据库连接的最大数量
	_sqlDb.SetConnMaxLifetime(tools.AutoTimeDuration(c.MaxConnLifetime, time.Second)) // 连接的总生存时间，从连接创建开始计时
	_sqlDb.SetConnMaxIdleTime(tools.AutoTimeDuration(c.MaxConnIdleTime, time.Second)) // 连接的空闲时间，从连接变为空闲开始计时
	return _db, nil
}
