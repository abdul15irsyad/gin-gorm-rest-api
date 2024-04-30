package utils

import (
	"os"
	"strconv"

	gomail "github.com/go-mail/mail"
)

type SendMailOptions struct {
	To            string
	Subject       string
	Message       string
	AttachmentURL *string
}

func SendMail(options SendMailOptions) error {
	APP_NAME := os.Getenv("APP_NAME")
	SMTP_HOST := os.Getenv("SMTP_HOST")
	SMTP_PORT, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	SMTP_USERNAME := os.Getenv("SMTP_USERNAME")
	SMTP_PASSWORD := os.Getenv("SMTP_PASSWORD")
	SMTP_FROM := os.Getenv("SMTP_FROM")

	message := gomail.NewMessage()
	message.SetHeader("From", APP_NAME+" <"+SMTP_FROM+">")
	message.SetHeader("To", options.To)
	message.SetHeader("Subject", options.Subject)
	message.SetBody("text/html", options.Message)
	if options.AttachmentURL != nil {
		message.Attach(*options.AttachmentURL)
	}

	dial := gomail.NewDialer(SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD)
	err := dial.DialAndSend(message)
	if err != nil {
		return err
	}

	return nil
}
