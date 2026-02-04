package sqliteconnect

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"helay.net/go/utils/v3/db"
)

func InitDB(c *db.Dbbase) (*gorm.DB, error) {
	dialector := sqlite.Open(c.Dsn())
	return c.Connect(&dialector)
}
