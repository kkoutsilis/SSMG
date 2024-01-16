package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

var version string = "0.1.0"

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
	var matches = make([]MatchPair, 0, len(data))
	for i := 0; i < len(data); i++ {
		from := data[i]
		to := data[(i+1)%len(data)]
		matches = append(matches, MatchPair{From: from, To: to})
	}
	return matches
}

func generateEmailMessage(to, body string) *gomail.Message {
	subject := "Your Secret Santa Match!"

	message := gomail.NewMessage()

	message.SetHeader("From", os.Getenv("EMAIL_FROM"))
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	return message

}

func generateEmailBody(match MatchPair) (string, error) {
	templateFile := "./templates/email_template.html"
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse email template")
	}

	var emailBodyBuffer = &strings.Builder{}
	err = tmpl.Execute(emailBodyBuffer, match)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template for email body")
	}
	emailBody := emailBodyBuffer.String()

	return emailBody, nil
}

func sendEmails(matches []MatchPair) error {
	host := os.Getenv("EMAIL_HOST")
	strPort := os.Getenv("EMAIL_PORT")
	user := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")

	port, err := strconv.Atoi(strPort)
	if err != nil {
		return errors.Wrap(err, "error parsing email server port")
	}

	d := gomail.NewDialer(host, port, user, password)
	s, err := d.Dial()
	if err != nil {
		return errors.Wrap(err, "error dialing email server")
	}
	defer s.Close()

	emailMessages := make([]*gomail.Message, 0, len(matches))
	for _, match := range matches {
		emailBody, err := generateEmailBody(match)
		if err != nil {
			return errors.Wrap(err, "failed to generate email body")
		}
		emailMessages = append(emailMessages, generateEmailMessage(match.From.Email, emailBody))
	}

	err = gomail.Send(s, emailMessages...)
	if err != nil {
		return err
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:          "run [path]",
	Short:        "A cli tool that generates secret santa matches and notifies the participants by email",
	ArgAliases:   []string{"path"},
	Version:      version,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
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
			return fmt.Errorf("file %s does not exist", filePath)

		}

		if !checkIsJson(filePath) {
			return fmt.Errorf("file %s is not a json file", filePath)
		}

		file, err := os.ReadFile(filePath)
		if err != nil {
			return errors.Wrap(err, "error when opening file")
		}

		var payload []Data
		err = json.Unmarshal(file, &payload)
		if err != nil {
			return errors.Wrap(err, "error reading file content")
		}

		if len(payload) == 0 {
			return errors.New("file is empty")
		}

		matches := generateSecretSantaMatches(payload)
		err = sendEmails(matches)
		if err != nil {
			return errors.Wrap(err, "error sending emails")
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
