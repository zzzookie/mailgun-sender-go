package database

import (
	"fmt"
	"mailgun-sender-go/models"
	"mailgun-sender-go/structs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

type Config struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.Host, config.Port, config.DBName, config.User, config.Password, config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return db, err
	} else {
		return db, nil
	}
}

func (db *Database) GetCampaignByName(campaignName string, userLangs []string) (*models.Campaign, error) {
	var campaign models.Campaign
	err := db.DB.Preload("Translations", "lang IN (?)", userLangs).
		Where("name = ?", campaignName).
		First(&campaign).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Error: A campaign with \"%s\" name was not found in the database", campaignName)
		}
		return nil, fmt.Errorf("Error while using getCampaignByName method: %s", err.Error())
	}

	if len(campaign.Translations) == 0 {
		return nil, fmt.Errorf("Error: A campaign with \"%s\" name does not have any of the requested languages (%v)", campaignName, userLangs)
	}

	return &campaign, nil
}

func (db *Database) GetDuplicatedRecipients(campID uint, emails []string) ([]models.SendStat, error) {
	var duplicatedRecipients []models.SendStat
	err := db.DB.
		Where("camp_id = ? AND email IN ? AND success = true", campID, emails).
		Find(&duplicatedRecipients).Error

	if err != nil {
		return nil, fmt.Errorf("Error while using getDuplicatedRecipients method: %s", err.Error())
	}

	return duplicatedRecipients, nil
}

func (db *Database) SendStats(statsData []structs.SendStat) ([]structs.SendStat, error) {
	err := db.DB.Create(&statsData).Error
	if err != nil {
		return nil, fmt.Errorf("Error while using sendStats method: %s", err.Error())
	}
	return statsData, nil
}
