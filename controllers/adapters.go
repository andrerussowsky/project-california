package controllers

import (
	"fmt"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
	"net/smtp"
	"project-california/components"
	"project-california/models"
)

func LoadErrorSession(res http.ResponseWriter, req *http.Request) models.Session {
	var sessionModel = models.Session{}
	errorSession, _ := gothic.Store.Get(req, "error-session")
	if errorSession.Values["error"] != nil {
		sessionModel.Error = fmt.Sprintf("%v", errorSession.Values["error"])
		RemoveErrorSession(res,req)
	}
	return sessionModel
}

func RedirectWithErrorMessage(res http.ResponseWriter, req *http.Request, err, url string) {
	SetErrorSession(res, req, err)
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusFound)
}

func SendMail(c *components.Components, email string, body []byte) {

	smtpServer := struct {
		host string
		port string
	} {
		"smtp.gmail.com",
		"587",
	}

	from := c.Settings.GetString("email.username")
	password := c.Settings.GetString("email.password")

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err := smtp.SendMail(smtpServer.host + ":" + smtpServer.port, auth, from, []string{email}, body)
	if err != nil {
		log.Println("Email error", err)
		return
	}
	log.Println("Email Sent!")
}