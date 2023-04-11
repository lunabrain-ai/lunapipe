package internal

import "gorm.io/gorm"

type Log struct {
	gorm.Model

	Prompt string
	Answer string
}

func InitDB(db *gorm.DB) error {
	return db.AutoMigrate(&Log{})
}
