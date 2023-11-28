package models

import (
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	MgTemplate  string `gorm:"not null"`
	DefaultLang string `gorm:"type:varchar(2);default:en;not null"`
}

type Translation struct {
	ID       uint     `gorm:"primaryKey"`
	CampID   uint     `gorm:"uniqueIndex:unique_campID_lang;not null,constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Campaign Campaign `gorm:"foreignKey: CampID"`
	Lang     string   `gorm:"uniqueIndex:unique_campID_lang;type:varchar(2);not null"`
	From     string   `gorm:"not null"`
	Subject  string   `gorm:"not null"`
}

type SendStat struct {
	ID       uint      `gorm:"primaryKey"`
	Ts       time.Time `gorm:"default:current_timestamp"`
	CampID   uint      `gorm:"not null, constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Campaign Campaign  `gorm:"foreignKey: CampID"`
	Lang     string    `gorm:"type:varchar(2);not null"`
	Email    string    `gorm:"not null"`
	ExtID    string    `gorm:"not null"`
	Success  bool      `gorm:"not null"`
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
