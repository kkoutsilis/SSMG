package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
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
		return fmt.Errorf("error reading email server port %w", err)
	}

	d := gomail.NewDialer(host, port, user, password)
	s, err := d.Dial()
	if err != nil {
		return fmt.Errorf("error dialing email server %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return gomail.Send(s, m)
}

func sendEmails(matches []MatchPair) ([]string, error) {
	// TODO: add retry logic for failed emails (maybe use a queue)
	// instead of sending each email individually,
	// generate all the emails and send them in bulk to avoid individual errors
	emailIssues := []string{}
	subject := "Your Secret Santa Match!"

	templateFile := "./templates/email_template.html"
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return emailIssues, fmt.Errorf("failed to parse email template: %w", err)
	}

	for _, match := range matches {
		var emailBodyBuffer = &strings.Builder{}
		err := tmpl.Execute(emailBodyBuffer, match)
		if err != nil {
			emailIssues = append(emailIssues, fmt.Sprintf("failed to execute template for %s: %v", match.From.Email, err))

		}
		emailBody := emailBodyBuffer.String()

		if err := sendEmail(match.From.Email, subject, emailBody); err != nil {
			emailIssues = append(emailIssues, fmt.Sprintf("failed to send email to %s: %v", match.From.Email, err))
		}
	}
	return emailIssues, nil
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
			return fmt.Errorf("error when opening file: %w", err)
		}

		var payload []Data
		err = json.Unmarshal(file, &payload)
		if err != nil {
			return fmt.Errorf("error reading file content: %w", err)
		}

		if len(payload) == 0 {
			return errors.New("file is empty")
		}

		matches := generateSecretSantaMatches(payload)
		emailIssues, err := sendEmails(matches)
		if err != nil {
			return fmt.Errorf("error sending emails: %w", err)
		}
		if len(emailIssues) > 0 {
			for _, issue := range emailIssues {
				fmt.Println(issue)
			}
			fmt.Println("Some emails failed to send, please check the logs for more details")
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
