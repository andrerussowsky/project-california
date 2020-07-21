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
	"project-california/components"
	"project-california/db"
	"project-california/models"
	"strings"
)

func UserForgotPassword(c *components.Components, res http.ResponseWriter, req *http.Request) {

	sessionModel := LoadErrorSession(res, req)

	t, _ := template.ParseFiles("templates/forgot-password.html")
	t.Execute(res, sessionModel)
}

func UserForgotPasswordPost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	email := strings.TrimSpace(req.Form.Get("email"))
	if len(email) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid email", "/user/forgot-password")
		return
	}

	dbUser, err := db.GetUserWithEmail(c, email)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Email not found", "/user/forgot-password")
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
		RedirectWithErrorMessage(res, req, "Unexpected error", "/user/forgot-password")
		return
	}

	tmpl, _ := template.ParseFiles("templates/forgot-password-email.html")

	data := struct {
		Url string
	}{
		fmt.Sprintf(`%s?uuid=%s`, c.Settings.GetString("email.forgotPasswordUrl"), uuid),
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

func UserForgotPasswordReset(c *components.Components, res http.ResponseWriter, req *http.Request) {

	sessionModel := LoadErrorSession(res, req)

	keys, ok := req.URL.Query()["uuid"]
	if !ok || len(keys[0]) < 1 {
		RedirectWithErrorMessage(res, req, "Invalid forgot password", "/")
		return
	}

	fp, err := db.GetForgotPassword(c, keys[0])
	if err != nil {
		RedirectWithErrorMessage(res, req, "Link expired", "/")
		return
	}

	SetForgotPasswordSession(res, req, keys[0], fp.UserID)

	t, _ := template.ParseFiles("templates/forgot-password-reset.html")
	t.Execute(res, sessionModel)
}

func UserForgotPasswordResetPost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	fpSession, _ := gothic.Store.Get(req, "forgot-password-session")
	if fpSession.Values["user_id"] == nil {
		RedirectWithErrorMessage(res, req, "Unexpected error", "/")
		return
	}

	id := fmt.Sprintf("%v", fpSession.Values["user_id"])
	keyID := fmt.Sprintf("%v", fpSession.Values["key_id"])

	req.ParseForm()

	password := strings.TrimSpace(req.Form.Get("password"))
	if len(password) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid password", "/user/forgot-password-reset?uuid="+keyID)
		return
	}

	confirmPassword := strings.TrimSpace(req.Form.Get("confirm-password"))
	if len(confirmPassword) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid password confirmation", "/user/forgot-password-reset?uuid="+keyID)
		return
	}

	if password != confirmPassword {
		RedirectWithErrorMessage(res, req, "Password doesnt match", "/user/forgot-password-reset?uuid="+keyID)
		return
	}

	dbUser, err := db.GetUser(c, id)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Unexpected error", "/user/forgot-password-reset?uuid="+keyID)
		return
	}

	hash := md5.Sum([]byte(password))
	passwordMD5 := hex.EncodeToString(hash[:])
	dbUser.Password = passwordMD5

	err = db.UpdateUser(c, dbUser)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Update error, please try again later", "/user/forgot-password-reset?uuid="+keyID)
		return
	}

	db.DeleteForgotPassword(c, dbUser.ID)
	RemoveForgotPasswordSession(res, req)

	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusFound)
	return
}