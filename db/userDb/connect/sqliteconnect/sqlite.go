package sqliteconnect

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"helay.net/go/utils/v3/db"
)

func InitDB(c *db.Dbbase) (*gorm.DB, error) {
	return c.Connect(new(sqlite.Open(c.Dsn())))
}
