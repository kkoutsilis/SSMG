package cmd

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"gopkg.in/gomail.v2"
)

var version string = "0.0.1"

type Data struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	WishList []string `json:"wishlist omitempty"`
}

type MatchPair struct {
	Person1 Data `json:"person1"`
	Person2 Data `json:"person2"`
}

func checkFileExists(path string) bool {
	_, error := os.Stat(path)
	return !errors.Is(error, os.ErrNotExist)
}

func checkIsJson(path string) bool {
	fileExtension := filepath.Ext(path)
	return fileExtension == ".json"
}

func generateSecretSantaMatches(data []Data) ([]MatchPair, error) {
	if len(data)%2 != 0 {
		return nil, errors.New("number of participants is odd, please provide an even number of participants for the Secret Santa matches")
	}

	shuffledData := make([]Data, len(data))
	copy(shuffledData, data)
	rand.Shuffle(len(shuffledData), func(i, j int) {
		shuffledData[i], shuffledData[j] = shuffledData[j], shuffledData[i]
	})

	var matches []MatchPair
	for i := 0; i < len(shuffledData); i += 2 {
		person1 := shuffledData[i]
		person2 := shuffledData[(i+1)%len(shuffledData)]
		matches = append(matches, MatchPair{Person1: person1, Person2: person2})
	}

	return matches, nil
}

func sendEmails(matches []MatchPair) {
	host := os.Getenv("EMAIL_HOST")
	strPort := os.Getenv("EMAIL_PORT")
	user := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")

	port, err := strconv.Atoi(strPort)
	if err != nil {
		panic(err)
	}

	d := gomail.NewDialer(host, port, user, password)
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()
	for _, match := range matches {
		from := os.Getenv("EMAIL_FROM")
		subject := "Your Secret Santa Match!"

		m.SetHeader("From", from)
		m.SetHeader("To", match.Person1.Email)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", "Hello "+match.Person1.Name+",<br><br>You are the secret santa for "+match.Person2.Name+"!<br><br>Best regards,<br>Secret Santa Match Generator")

		if err := gomail.Send(s, m); err != nil {
			log.Fatalf("Could not send email: %v", err)
		}
		m.Reset()

		m.SetHeader("From", from)
		m.SetHeader("To", match.Person2.Email)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", "Hello "+match.Person2.Name+",<br><br>You are the secret santa for "+match.Person1.Name+"!<br><br>Best regards,<br>Secret Santa Generator")

		if err := gomail.Send(s, m); err != nil {
			log.Fatalf("Could not send email: %v", err)
		}
		m.Reset()
	}
}

var rootCmd = &cobra.Command{
	Use:        "run [path]",
	Short:      "A cli tool that generates secret santa matches",
	ArgAliases: []string{"path"},
	Version:    version,
	Run: func(cmd *cobra.Command, args []string) {
		var filePath string
		if len(args) == 0 {
			filePath = "data.json"
		} else {
			filePath = args[0]
		}
		if filePath == "" {
			filePath = "data.json"
		}

		if !checkFileExists(filePath) {
			log.Fatal("File ", filePath, " does not exist")
		}

		if !checkIsJson(filePath) {
			log.Fatal("File ", filePath, " is not a json file")
		}

		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}

		var payload []Data
		err = json.Unmarshal(file, &payload)
		if err != nil {
			log.Fatal("Error reading file content: ", err)
		}

		if len(payload) == 0 {
			log.Fatal("File is empty")
		}

		matches, err := generateSecretSantaMatches(payload)
		if err != nil {
			log.Fatal(err)
		}

		sendEmails(matches)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
