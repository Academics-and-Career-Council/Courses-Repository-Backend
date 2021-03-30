package services

import (
	"bytes"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/mhale/smtpd"
	"github.com/spf13/viper"
)

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
	mail, err := parsemail.Parse(strings.NewReader(string(data)))
	if err != nil {
		log.Println(err)
		return err
	}
	SendMail(mail.TextBody, to[0])
	return nil
}

// SMTPListenAndServe : Starts the SMTP server
func SMTPListenAndServe() error {
	log.Print("Starting SMTP server at :2525")
	err := smtpd.ListenAndServe("0.0.0.0:2525", mailHandler, "anc-mailer", "localhost")
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func SendMail(message string, to string) error {

	msg := "From: " + viper.GetString("mail.from") + "\n" +
		"To: " + to + "\n" +
		"Subject: Account Recovery\n\n" +
		message

	err := smtp.SendMail(viper.GetString("mail.host"),
		smtp.PlainAuth("", viper.GetString("mail.id"), viper.GetString("mail.pwd"), viper.GetString("mail.host")),
		viper.GetString("mail.from"), []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	return nil
}
