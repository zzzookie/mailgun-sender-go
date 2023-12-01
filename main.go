package main

import (
	"fmt"
	"os"

	"mailgun-sender-go/database"
	"mailgun-sender-go/models"
	"mailgun-sender-go/utils"

	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	requiredEnvVars := []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD", "MAILGUN_DOMAIN", "MAILGUN_API_KEY"}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			panic(fmt.Sprintf("Environment variable %s is not set", envVar))
		}
	}

	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	fmt.Println(mg)

	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  utils.GetEnvWithDefault("DB_SSLMODE", "disable"),
	}

	pg, err := database.NewConnection(config)

	if err != nil {
		fmt.Println("Error:", err)
		panic("Could not load the database")
	}

	err = models.Migrate(pg)
	if err != nil {
		fmt.Println("Error:", err)
		panic("Could not migrate DB")
	}

	db := &database.Database{DB: pg}

	fileName, campaignName := utils.ReadClArgs()
	userList, err := utils.ParseMaillist(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	userLangs := utils.GetUniqueValues(userList, "Lang")
	fmt.Println("🚀 ~ file: main.go ~ line 61 ~ funcmain ~ userLangs : ", userLangs)

	campaign, err := db.GetCampaignByName(campaignName, userLangs)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("🚀 ~ file: main.go ~ line 62 ~ funcmain ~ campaign : ", campaign)
}
