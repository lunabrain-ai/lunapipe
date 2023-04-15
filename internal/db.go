package internal

import (
	"github.com/glebarez/sqlite"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"gorm.io/gorm"
)

func NewDatabase(cache cache.Cache) (*gorm.DB, error) {
	dbPath, err := cache.GetFile("db.sqlite")
	if err != nil {
		return nil, err
	}

	openedDb := sqlite.Open(dbPath + "?cache=shared&mode=rwc&_journal_mode=WAL")

	db, err := gorm.Open(openedDb, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, InitDB(db)
}
