package internal

import (
	"github.com/glebarez/sqlite"
	"github.com/lunabrain-ai/lunabrain/pkg/store"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	cache, err := store.NewFolderCache()
	if err != nil {
		return nil, err
	}

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
