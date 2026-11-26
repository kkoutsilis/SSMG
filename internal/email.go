package internal

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

func EmailConfigFromEnv() (EmailConfig, error) {
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		return EmailConfig{}, fmt.Errorf("parse EMAIL_PORT: %w", err)
	}
	return EmailConfig{
		Host: os.Getenv("EMAIL_HOST"), Port: port, User: os.Getenv("EMAIL_USER"),
		Password: os.Getenv("EMAIL_PASSWORD"), From: os.Getenv("EMAIL_FROM"),
	}, nil
}

func Send(matches []Match, templatePath string, config EmailConfig) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("load email template: %w", err)
	}

	messages := make([]*gomail.Message, 0, len(matches))
	for _, match := range matches {
		body, err := renderEmailBody(match, tmpl)
		if err != nil {
			return err
		}
		message := gomail.NewMessage()
		message.SetHeader("From", config.From)
		message.SetHeader("To", match.From.Email)
		message.SetHeader("Subject", "Your Secret Santa Match!")
		message.SetBody("text/html", body)
		messages = append(messages, message)
	}

	dialer := gomail.NewDialer(config.Host, config.Port, config.User, config.Password)
	sender, err := dialer.Dial()
	if err != nil {
		return fmt.Errorf("connect to email server: %w", err)
	}
	defer sender.Close()
	if err := gomail.Send(sender, messages...); err != nil {
		return fmt.Errorf("send emails: %w", err)
	}
	return nil
}

func renderEmailBody(match Match, tmpl *template.Template) (string, error) {
	var body strings.Builder
	if err := tmpl.Execute(&body, match); err != nil {
		return "", fmt.Errorf("render email template: %w", err)
	}
	return body.String(), nil
}
