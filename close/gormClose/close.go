package gormClose

import (
	"gorm.io/gorm"
	"helay.net/go/utils/v3/close/vclose"
)

func Close(db *gorm.DB) {
	if db == nil {
		return
	}
	if sqlDb, err := db.DB(); err == nil {
		vclose.Close(sqlDb)
	}
}
