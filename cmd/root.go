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
	From Data `json:"from"`
	To   Data `json:"to"`
}

func checkFileExists(path string) bool {
	_, error := os.Stat(path)
	return !errors.Is(error, os.ErrNotExist)
}

func checkIsJson(path string) bool {
	fileExtension := filepath.Ext(path)
	return fileExtension == ".json"
}

func generateSecretSantaMatches(data []Data) []MatchPair {
	shuffledData := make([]Data, len(data))
	copy(shuffledData, data)
	rand.Shuffle(len(shuffledData), func(i, j int) {
		shuffledData[i], shuffledData[j] = shuffledData[j], shuffledData[i]
	})

	return circlularMatchingAlgorithm(shuffledData)
}

func circlularMatchingAlgorithm(data []Data) []MatchPair {
	var matches []MatchPair
	for i := 0; i < len(data); i++ {
		from := data[i]
		to := data[(i+1)%len(data)]
		matches = append(matches, MatchPair{From: from, To: to})
	}
	return matches
}

func sendEmail(to, subject, body string) error {
	host := os.Getenv("EMAIL_HOST")
	strPort := os.Getenv("EMAIL_PORT")
	user := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")

	port, err := strconv.Atoi(strPort)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(host, port, user, password)
	s, err := d.Dial()
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return gomail.Send(s, m)
}

func sendEmails(matches []MatchPair) {
	subject := "Your Secret Santa Match!"
	for _, match := range matches {
		emailBody := "Hello " + match.From.Name + ",<br><br>You are the secret Santa for " + match.To.Name + "!<br><br>Best regards,<br>Secret Santa Match Generator"
		if err := sendEmail(match.From.Email, subject, emailBody); err != nil {
			log.Printf("Error sending email to %s: %v", match.From.Email, err)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:        "run [path]",
	Short:      "A cli tool that generates secret santa matches and notifies the participants by email",
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

		matches := generateSecretSantaMatches(payload)

		sendEmails(matches)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
