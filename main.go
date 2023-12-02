package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"mailgun-sender-go/database"
	"mailgun-sender-go/models"
	"mailgun-sender-go/structs"
	"mailgun-sender-go/utils"

	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v3"
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
	}

	userLangs := utils.GetUniqueValues(userList, "Lang")

	campaign, err := db.GetCampaignByName(campaignName, userLangs)
	if err != nil {
		fmt.Println(err)
		return
	}

	translations := make([]structs.Translation, len(campaign.Translations))
	for i, translation := range campaign.Translations {
		translations[i] = structs.Translation{
			Lang:       translation.Lang,
			From:       translation.From,
			Subject:    translation.Subject,
			Recipients: make([]structs.Recipient, 0),
		}
	}

	defaultLang := campaign.DefaultLang
	if defaultLang == "" {
		defaultLang = "en"
	}

	for _, user := range userList {
		for i := range translations {
			if translations[i].Lang == user.Lang {
				translations[i].Recipients = append(translations[i].Recipients, user)
				break
			}
		}
	}

	for _, translation := range translations {
		// statsData := make([]structs.SendStat, 0)

		emails := make([]string, 0) // список емэйлов к одному языку
		for _, recipient := range translation.Recipients {
			emails = append(emails, recipient.Email)
		}

		// duplicatedRecipients, err := db.GetDuplicatedRecipients(campaign.ID, emails)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		// if len(duplicatedRecipients) > 0 {
		// 	duplicatedRecipientsEmails := make(map[string]bool)
		// 	for _, r := range duplicatedRecipients {
		// 		duplicatedRecipientsEmails[r.Email] = true
		// 	}

		// 	removedIndexes := make([]int, 0)
		// 	excludedRecipients := make([]structs.RecipientExcluded, 0)
		// 	for i, recipient := range translation.Recipients {
		// 		if duplicatedRecipientsEmails[recipient.Email] {
		// 			var sendStatsID uint
		// 			for _, r := range duplicatedRecipients {
		// 				if recipient.Email == r.Email && translation.Lang == r.Lang {
		// 					sendStatsID = r.ID
		// 					break
		// 				}
		// 			}
		// 			excludedRecipients = append(excludedRecipients, structs.RecipientExcluded{
		// 				Recipient:   recipient,
		// 				SendStatsID: sendStatsID,
		// 			})
		// 			removedIndexes = append(removedIndexes, i)
		// 		}
		// 	}

		// 	filteredRecipients := make([]structs.Recipient, 0)
		// 	for i, recipient := range translation.Recipients {
		// 		if !slices.Contains(removedIndexes, i) {
		// 			filteredRecipients = append(filteredRecipients, recipient)
		// 		}
		// 	}
		// 	translation.Recipients = filteredRecipients

		// 	excludedStatsData := make([]structs.SendStat, 0)
		// 	for _, r := range excludedRecipients {
		// 		excludedStatsData = append(excludedStatsData, structs.SendStat{
		// 			Lang:     translation.Lang,
		// 			Email:    r.Recipient.Email,
		// 			ExtID:    r.Recipient.ExtID,
		// 			Success:  false,
		// 			ErrorMsg: fmt.Sprintf("already sent: <%s>", r.SendStatsID),
		// 		})
		// 	}

		// 	_, err = db.SendStats(excludedStatsData)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return
		// 	}
		// }

		if len(translation.Recipients) == 0 { // Проверка, есть ли кому посылать
			fmt.Printf("Nothing has been sent: all the recipients in the list has received \"%s\" campaign email before.\n", campaignName)
			fmt.Println("Information has been saved in the statistics.")
			return
		}
		if len(translation.Recipients) > 1000 { // Заглушка, Mailgun не может посылать больше 1000 за раз
			fmt.Println("Can't send more than to 1000 recipients at once.")
			return
		}

		message := mg.NewMessage(
			translation.From,
			translation.Subject,
			"")

		message.SetTemplate(campaign.MgTemplate)
		message.AddTag(campaign.MgTemplate)
		message.AddTag(translation.Lang)
		message.AddTemplateVariable("version", translation.Lang)
		for _, recipient := range translation.Recipients {
			message.AddRecipient(fmt.Sprintf("%s <%s>", recipient.Name, recipient.Email))
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		_, _, err := mg.Send(ctx, message)

		var statsData []structs.SendStat
		success := false
		if err != nil {
			for _, recipient := range translation.Recipients {
				statsData = append(statsData, structs.SendStat{
					CampID:   campaign.ID,
					Lang:     recipient.Lang,
					Email:    recipient.Email,
					ExtID:    recipient.ExtID,
					Success:  false,
					ErrorMsg: err.Error(),
				})
			}
		} else {
			for _, recipient := range translation.Recipients {
				statsData = append(statsData, structs.SendStat{
					CampID:   campaign.ID,
					Lang:     recipient.Lang,
					Email:    recipient.Email,
					ExtID:    recipient.ExtID,
					Success:  true,
					ErrorMsg: "",
				})
			}
			success = true
		}

		_, err = db.SendStats(statsData)
		if err != nil {
			fmt.Println(err)
			return
		}

		if success {
			fmt.Printf("%s mails has been successfully sent.\n", strconv.Itoa(len(statsData)))
		} else {
			fmt.Printf("Failed to send to: %s", strings.Join(emails, ", "))
			fmt.Println(". Information has been saved to stats_data")
		}
	}
}
