package main

import (
	"log"
	"time"

	"github.com/hippo-an/tiny-go-challenges/back_06/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {

	go func() {
		for {
			m := <-app.MailChan
			sendMessage(m)
		}
	}()

}

func sendMessage(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()

	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)

	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)

	if err != nil {
		log.Println(err)
		return
	}
}
