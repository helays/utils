package sqliteconnect

import (
	"github.com/helays/utils/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite" // 必须导入以注册驱动
)

func InitDB(c *db.Dbbase) (*gorm.DB, error) {
	dialector := sqlite.Open(c.Dsn())
	return c.Connect(&dialector)
}
