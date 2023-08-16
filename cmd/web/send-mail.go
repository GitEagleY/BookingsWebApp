package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/GitEagleY/BookingsWebApp/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {

	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()

}

// sendMsg sends an email using the provided mail data.
func sendMsg(m models.MailData) {
	// Create a new SMTP client.
	server := mail.NewSMTPClient()

	// Set the SMTP server host and port.
	server.Host = "localhost"
	server.Port = 1025

	// Configure connection settings.
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// Uncomment and provide credentials if required.
	// server.Username = "your_username"
	// server.Password = "your_password"

	// Uncomment and set encryption type if needed.
	// server.Encryption = mail.EncryptionTLS

	// Connect to the SMTP server.
	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	// Create a new email message.
	email := mail.NewMSG()

	// Set the sender, recipient, and subject of the email.
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)

	if m.Template == "" {
		// Set the body of the email as HTML text.
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}

		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)

	}

	// Send the email using the SMTP client.
	err = email.Send(client)
	if err != nil {
		// Log any errors that occur during email sending.
		log.Println(err)
	} else {
		// Log a message indicating successful email sending.
		log.Println("Email sent successfully")
	}
}
