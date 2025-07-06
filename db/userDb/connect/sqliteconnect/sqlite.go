package sqliteconnect

import (
	"github.com/helays/utils/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(c *db.Dbbase) (*gorm.DB, error) {
	dialector := sqlite.Open(c.Dsn())
	return c.Connect(&dialector)
}
