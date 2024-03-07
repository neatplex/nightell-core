package mailer

import (
	"fmt"
	"net/smtp"
)

func Send() {
	// Sender's email address and password
	senderEmail := "i@miladrahimi.com"
	senderPassword := "Zo@22&1a"

	// Recipient's email address
	recipientEmail := "realmiladrahimi@gmail.com"

	// SMTP server address and port
	smtpServer := "smtp.zoho.com"
	smtpPort := 587

	// Message to be sent
	message := []byte("Subject: Hello!\r\n" +
		"From: " + senderEmail + "\r\n" +
		"To: " + recipientEmail + "\r\n" +
		"\r\n" +
		"This is a test email sent using Go!")

	// Authentication credentials
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	// Send email
	err := smtp.SendMail(smtpServer+":"+fmt.Sprintf("%d", smtpPort), auth, senderEmail, []string{recipientEmail}, message)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
