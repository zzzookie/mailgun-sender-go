package models

import (
	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	MgTemplate  string `gorm:"not null"`
	DefaultLang string `gorm:"default:en;not null"`
}

type Translation struct {
	gorm.Model
	CampID   uint     `gorm:"not null"`
	Campaign Campaign `gorm:"foreignKey: CampID"`
	Lang     string   `gorm:"not null"`
	From     string   `gorm:"not null"`
	Subject  string   `gorm:"not null"`
}

type SendStat struct {
	gorm.Model
	Ts       int64    `gorm:"autoCreateTime"`
	CampID   uint     `gorm:"not null"`
	Campaign Campaign `gorm:"foreignKey: CampID"`
	Lang     string   `gorm:"not null"`
	Email    string   `gorm:"not null"`
	ExtID    string   `gorm:"not null"`
	Success  bool     `gorm:"not null"`
	ErrorMsg string
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Campaign{}, &Translation{}, &SendStat{})
	return err
}

func Rollback(db *gorm.DB) {
	db.Migrator().DropTable(&Campaign{}, &Translation{}, &SendStat{})
	db.Migrator().DropTable("Campaings", "Translations", "Send_stats")
	db.Migrator().DropConstraint("Translations", "unique_camp_lang")
	db.Migrator().DropConstraint("Send_stats", "unique_camp_lang")
}
