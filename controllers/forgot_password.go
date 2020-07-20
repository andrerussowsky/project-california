package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/markbates/goth/gothic"
	"github.com/satori/go.uuid"
	"html/template"
	"net/http"
	"net/smtp"
	"project-california/components"
	"project-california/db"
	"project-california/models"
	"strings"
)

func UserForgotPassword(c *components.Components, res http.ResponseWriter, req *http.Request) {

	var sessionModel = models.Session{}
	errorSession, _ := gothic.Store.Get(req, "error-session")
	if errorSession.Values["error"] != nil {
		sessionModel.Error = fmt.Sprintf("%v", errorSession.Values["error"])
		RemoveErrorSession(res,req)
	}

	t, _ := template.ParseFiles("templates/forgot-password.html")
	t.Execute(res, sessionModel)
}

func UserForgotPasswordPost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	email := strings.TrimSpace(req.Form.Get("email"))
	if len(email) == 0 {
		SetErrorSession(res, req, "Invalid email")

		res.Header().Set("Location", "/user/forgot-password")
		res.WriteHeader(http.StatusFound)
		return
	}

	dbUser, err := db.GetUserWithEmail(c, email)
	if err != nil {
		SetErrorSession(res, req, "Email not found")

		res.Header().Set("Location", "/user/forgot-password")
		res.WriteHeader(http.StatusFound)
		return
	}

	db.DeleteForgotPassword(c, dbUser.ID)
	RemoveForgotPasswordSession(res, req)

	uuid := uuid.NewV4().String()
	err = db.InsertForgotPassword(c, models.ForgotPassword{
		UserID: dbUser.ID,
		UUID:   uuid,
	})
	if err != nil {
		SetErrorSession(res, req, "Unexpected error")

		res.Header().Set("Location", "/user/forgot-password")
		res.WriteHeader(http.StatusFound)
		return
	}

	tmpl, _ := template.ParseFiles("templates/forgot-password-email.html")

	forgotPasswordUrl := c.Settings.GetString("email.forgotPasswordUrl")
	data := struct {
		Url string
	}{
		fmt.Sprintf(`%s?uuid=%s`, forgotPasswordUrl, uuid),
	}

	var tpl bytes.Buffer
	tmpl.Execute(&tpl, data)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Forgot password\n"
	body := []byte(subject + mime + tpl.String())

	SendMail(c, email, body)

	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusFound)
	return
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
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

func UserForgotPasswordReset(c *components.Components, res http.ResponseWriter, req *http.Request) {

	var sessionModel = models.Session{}
	errorSession, _ := gothic.Store.Get(req, "error-session")
	if errorSession.Values["error"] != nil {
		sessionModel.Error = fmt.Sprintf("%v", errorSession.Values["error"])
		RemoveErrorSession(res,req)
	}

	keys, ok := req.URL.Query()["uuid"]
	if !ok || len(keys[0]) < 1 {
		SetErrorSession(res, req, "Invalid forgot password")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	fp, err := db.GetForgotPassword(c, keys[0])
	if err != nil {
		SetErrorSession(res, req, "Unexpected error")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	SetForgotPasswordSession(res, req, keys[0], fp.UserID)

	t, _ := template.ParseFiles("templates/forgot-password-reset.html")
	t.Execute(res, sessionModel)
}

func UserForgotPasswordResetPost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	fpSession, _ := gothic.Store.Get(req, "forgot-password-session")
	if fpSession.Values["user_id"] == nil {
		SetErrorSession(res, req, "Unexpected error")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	id := fmt.Sprintf("%v", fpSession.Values["user_id"])
	keyID := fmt.Sprintf("%v", fpSession.Values["key_id"])

	req.ParseForm()

	password := strings.TrimSpace(req.Form.Get("password"))
	if len(password) == 0 {
		SetErrorSession(res, req, "Invalid password")

		res.Header().Set("Location", "/user/forgot-password-reset?uuid="+keyID)
		res.WriteHeader(http.StatusFound)
		return
	}

	confirmPassword := strings.TrimSpace(req.Form.Get("confirm-password"))
	if len(confirmPassword) == 0 {
		SetErrorSession(res, req, "Invalid password confirmation")

		res.Header().Set("Location", "/user/forgot-password-reset?uuid="+keyID)
		res.WriteHeader(http.StatusFound)
		return
	}

	if password != confirmPassword {
		SetErrorSession(res, req, "Password doesnt match")

		res.Header().Set("Location", "/user/forgot-password-reset?uuid="+keyID)
		res.WriteHeader(http.StatusFound)
		return
	}

	dbUser, err := db.GetUser(c, id)
	if err != nil {
		SetErrorSession(res, req, "Unexpected error")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	hash := md5.Sum([]byte(password))
	passwordMD5 := hex.EncodeToString(hash[:])
	dbUser.Password = passwordMD5

	err = db.UpdateUser(c, dbUser)
	if err != nil {
		SetErrorSession(res, req, "Update error, please try again later")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	db.DeleteForgotPassword(c, dbUser.ID)
	RemoveForgotPasswordSession(res, req)

	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusFound)
	return
}