package email

import (
	"crypto/tls"
	"fmt"
	"github.com/ctaperts/messages/src"
	"log"
	"net/smtp"
	"strconv"
	"strings"
)

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type smtpServer struct {
	host string
	port string
}

var (
	LocalConfig configuration.Config
)

func (s *smtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func Send(LocalConfig configuration.Config, emailFrom, emailSubject, emailBody string, emailTo []string) bool {

	mail := Mail{}
	mail.senderId = LocalConfig.SMTP.Username
	mail.toIds = emailTo
	mail.subject = emailSubject
	mail.body = emailBody

	messageBody := mail.BuildMessage()

	smtpServer := smtpServer{host: LocalConfig.SMTP.Server, port: strconv.Itoa(LocalConfig.SMTP.Port)}

	//build an auth
	auth := smtp.PlainAuth("", mail.senderId, LocalConfig.SMTP.Password, smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		log.Println(err)
		return false
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Println(err)
		return false
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Println(err)
		return false
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderId); err != nil {
		log.Println(err)
		return false
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Println(err)
			return false
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Println(err)
		return false
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
		return false
	}

	client.Quit()

	return true

}
