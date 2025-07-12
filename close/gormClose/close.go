package gormClose

import (
	"github.com/helays/utils/v2/close/vclose"
	"gorm.io/gorm"
)

func Close(db *gorm.DB) {
	if db == nil {
		return
	}
	if sqlDb, err := db.DB(); err == nil {
		vclose.Close(sqlDb)
	}
}
