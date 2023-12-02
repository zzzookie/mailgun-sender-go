package utils

import (
	"flag"
	"fmt"
	"mailgun-sender-go/structs"
	"os"
	"reflect"
	"strings"
)

func ReadClArgs() (string, string) {
	var maillistFile string
	var campaignName string

	const (
		maillistDesc = "Specify a file with a list of mailing addresses in 'file' or 'file.csv' format"
		campaignDesc = "Specify the campaign name (for example, 'registration_without_activation')"
	)

	flag.StringVar(&maillistFile, "maillist", GetEnvWithDefault("DEFAULT_FILE_NAME", "maillist"), maillistDesc)
	flag.StringVar(&maillistFile, "ml", GetEnvWithDefault("DEFAULT_FILE_NAME", "maillist"), maillistDesc+" (shorthand)")
	flag.StringVar(&campaignName, "campaign", "", campaignDesc)
	flag.StringVar(&campaignName, "camp", "", campaignDesc+" (shorthand)")

	flag.Parse()

	if campaignName == "" {
		fmt.Println("Error: Please specify the Mailgun campaign name with \"-camp\"")
		os.Exit(1)
	}

	if !strings.HasSuffix(maillistFile, ".csv") {
		maillistFile += ".csv"
	}

	return maillistFile, strings.TrimSpace(campaignName)
}

func ParseMaillist(fileName string) ([]structs.Recipient, error) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Error: File \"%s\" not found", fileName)
	} else if err != nil {
		return nil, fmt.Errorf("Error: Unable to access file - %s", err.Error())
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Error: Parsing file error %s", err.Error())
	}

	usersStrArr := strings.Split(strings.TrimSpace(string(content)), "\n")
	if strings.Contains(usersStrArr[0], "name,email") {
		usersStrArr = usersStrArr[1:]
	}

	var users []structs.Recipient
	for _, userStr := range usersStrArr {
		fields := strings.Split(userStr, ",")
		name := strings.TrimSpace(fields[0])
		email := strings.TrimSpace(fields[1])
		lang := strings.TrimSpace(fields[2])
		extID := strings.TrimSpace(fields[3])
		users = append(users, structs.Recipient{Name: name, Email: email, Lang: lang, ExtID: extID})
	}

	return users, nil
}

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetUniqueValues(arr interface{}, fieldName string) []string {
	values := make(map[string]bool)

	elementType := reflect.TypeOf(arr).Elem()
	if elementType.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < reflect.ValueOf(arr).Len(); i++ {
		elemValue := reflect.ValueOf(arr).Index(i)
		fieldValue := elemValue.FieldByName(fieldName).String()
		values[fieldValue] = true
	}

	result := make([]string, 0, len(values))
	for key := range values {
		result = append(result, key)
	}

	return result
}
