package mysqlconnect

import (
	"github.com/helays/utils/v2/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(c *db.Dbbase) (*gorm.DB, error) {
	//dsn = c.User + ":" + c.Pwd + "@tcp(" + strings.Join(c.Host, ",") + ")/" + c.Dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	//fmt.Println(dsn)
	dialector := mysql.New(mysql.Config{
		DSN:                       c.Dsn(),
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	})

	return c.Connect(&dialector)
}
