package utils

import (
	"fmt"
	"net/smtp"
)

func SendEmail(email, url string) error {
	password := ReadENV("password")
	from := ReadENV("from")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Subject: " + "Email confirmation from Network" + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"Confirm your email by clicking this link:" + url)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	to := []string{email}

	fmt.Println(to)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	if err != nil {
		return err
	}

	return nil
}
