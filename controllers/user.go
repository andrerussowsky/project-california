package controllers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/markbates/goth/gothic"
	"html/template"
	"net/http"
	"project-california/components"
	"project-california/db"
	"project-california/models"
	"strings"
)

func UserSignIn(c *components.Components, res http.ResponseWriter, req *http.Request) {

	sessionModel := LoadErrorSession(res, req)

	t, _ := template.ParseFiles("templates/sign-in.html")
	t.Execute(res, sessionModel)
}

func UserSignInPost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	email := strings.TrimSpace(req.Form.Get("email"))
	if len(email) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid email", "/user/sign-in")
		return
	}
	password := strings.TrimSpace(req.Form.Get("password"))
	if len(password) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid password", "/user/sign-in")
		return
	}

	dbUser, err := db.GetUserWithEmail(c, email)
	if err != nil {
		if err == sql.ErrNoRows {
			RedirectWithErrorMessage(res, req, "User not found", "/user/sign-in")
			return
		}

		RedirectWithErrorMessage(res, req, "Unexpected error", "/user/sign-in")
		return
	}

	hash := md5.Sum([]byte(password))
	passwordMD5 := hex.EncodeToString(hash[:])
	if dbUser.Password != passwordMD5 {
		RedirectWithErrorMessage(res, req, "Wrong password", "/user/sign-in")
		return
	}

	SetUserSession(res, req, dbUser)

	if dbUser.Status == models.UserStatusComplete {
		res.Header().Set("Location", "/user/profile")
	} else {
		res.Header().Set("Location", "/user/sign-up-complete")
	}
	res.WriteHeader(http.StatusFound)
	return
}

func UserSignUp(c *components.Components, res http.ResponseWriter, req *http.Request) {

	sessionModel := LoadErrorSession(res, req)

	t, _ := template.ParseFiles("templates/sign-up.html")
	t.Execute(res, sessionModel)
}

func UserSignUpPost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	email := strings.TrimSpace(req.Form.Get("email"))
	if len(email) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid email", "/user/sign-up")
		return
	}

	password := strings.TrimSpace(req.Form.Get("password"))
	if len(password) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid password", "/user/sign-up")
		return
	}

	confirmPassword := strings.TrimSpace(req.Form.Get("confirm-password"))
	if len(confirmPassword) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid password confirmation", "/user/sign-up")
		return
	}

	if password != confirmPassword {
		RedirectWithErrorMessage(res, req, "Password doesnt match", "/user/sign-up")
		return
	}

	_, err := db.GetUserWithEmail(c, email)
	if err == nil {
		RedirectWithErrorMessage(res, req, "Email is already taken", "/user/sign-up")
		return
	}

	hash := md5.Sum([]byte(password))
	passwordMD5 := hex.EncodeToString(hash[:])

	err = db.InsertUser(c, models.User{Email: email, Password: passwordMD5})
	if err != nil {
		RedirectWithErrorMessage(res, req, "Insert error, please try again later", "/user/sign-up")
		return
	}

	dbUser, err := db.GetUserWithEmail(c, email)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Load user error, please try again later", "/user/sign-up")
		return
	}

	SetUserSession(res, req, dbUser)

	if dbUser.Status == models.UserStatusComplete {
		res.Header().Set("Location", "/user/profile")
	} else {
		res.Header().Set("Location", "/user/sign-up-complete")
	}
	res.WriteHeader(http.StatusFound)
	return
}

func UserSignUpComplete(c *components.Components, res http.ResponseWriter, req *http.Request) {

	sessionModel := LoadErrorSession(res, req)

	userSession, _ := gothic.Store.Get(req, "user-session")
	if userSession.Values["id"] == nil {
		RedirectWithErrorMessage(res, req, "Session expired", "/")
		return
	}

	t, _ := template.ParseFiles("templates/sign-up-complete.html")
	t.Execute(res, sessionModel)
}

func UserSignUpCompletePost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	userSession, _ := gothic.Store.Get(req, "user-session")
	if userSession.Values["id"] == nil {
		RedirectWithErrorMessage(res, req, "Session expired", "/")
		return
	}

	id := fmt.Sprintf("%v", userSession.Values["id"])

	dbUser, err := db.GetUser(c, id)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Unexpected error", "/user/sign-up-complete")
		return
	}

	name := strings.TrimSpace(req.Form.Get("name"))
	if len(name) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid name", "/user/sign-up-complete")
		return
	}
	phone := strings.TrimSpace(req.Form.Get("phone"))
	if len(phone) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid phone", "/user/sign-up-complete")
		return
	}
	address := strings.TrimSpace(req.Form.Get("address"))
	if len(address) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid address", "/user/sign-up-complete")
		return
	}

	dbUser.Name = name
	dbUser.Phone = phone
	dbUser.Address = address
	dbUser.Status = models.UserStatusComplete

	err = db.UpdateUser(c, dbUser)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Update error, please try again later", "/user/sign-up-complete")
		return
	}

	SetUserSession(res, req, dbUser)

	res.Header().Set("Location", "/user/profile")
	res.WriteHeader(http.StatusFound)
	return
}

func UserProfile(c *components.Components, res http.ResponseWriter, req *http.Request) {

	var sessionModel  = models.Session{}

	userSession, _ := gothic.Store.Get(req, "user-session")
	if userSession.Values["id"] == nil {
		RedirectWithErrorMessage(res, req, "Session expired", "/")
		return
	}

	id := fmt.Sprintf("%v", userSession.Values["id"])

	dbUser, err := db.GetUser(c, id)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Unexpected error", "/")
		return
	}

	if dbUser.Status != models.UserStatusComplete {
		RedirectWithErrorMessage(res, req, "Forbidden", "/")
		return
	}

	sessionModel.ID = dbUser.ID
	sessionModel.Name = dbUser.Name
	sessionModel.Email = dbUser.Email
	sessionModel.Phone = dbUser.Phone
	sessionModel.Address = dbUser.Address

	t, _ := template.ParseFiles("templates/profile.html")
	t.Execute(res, sessionModel)
}

func UserProfileEdit(c *components.Components, res http.ResponseWriter, req *http.Request) {

	sessionModel := LoadErrorSession(res, req)

	userSession, _ := gothic.Store.Get(req, "user-session")
	if userSession.Values["id"] == nil {
		RedirectWithErrorMessage(res, req, "Session expired", "/")
		return
	}

	id := fmt.Sprintf("%v", userSession.Values["id"])

	dbUser, err := db.GetUser(c, id)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Unexpected error", "/")
		return
	}

	if dbUser.Status != models.UserStatusComplete {
		RedirectWithErrorMessage(res, req, "Forbidden", "/")
		return
	}

	sessionModel.Name = dbUser.Name
	sessionModel.Email = dbUser.Email
	sessionModel.Phone = dbUser.Phone
	sessionModel.Address = dbUser.Address

	t, _ := template.ParseFiles("templates/profile-edit.html")
	t.Execute(res, sessionModel)
}

func UserProfileEditPost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	userSession, _ := gothic.Store.Get(req, "user-session")
	if userSession.Values["id"] == nil {
		RedirectWithErrorMessage(res, req, "Session expired", "/")
		return
	}

	id := fmt.Sprintf("%v", userSession.Values["id"])

	dbUser, err := db.GetUser(c, id)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Unexpected error", "/")
		return
	}

	if dbUser.Status != models.UserStatusComplete {
		RedirectWithErrorMessage(res, req, "Forbidden", "/")
		return
	}

	email := strings.TrimSpace(req.Form.Get("email"))
	if len(email) == 0 {
		RedirectWithErrorMessage(res, req, "Invalid email", "/user/profile-edit")
		return
	}

	if dbUser.Email != email {
		_, err := db.GetUserWithEmail(c, email)
		if err == nil {
			RedirectWithErrorMessage(res, req, "Email is already taken", "/user/profile-edit")
			return
		}
	}

	name := strings.TrimSpace(req.Form.Get("name"))
	if len(name) > 0 {
		dbUser.Name = name
	}

	password := strings.TrimSpace(req.Form.Get("password"))
	if len(password) > 0 {
		hash := md5.Sum([]byte(password))
		passwordMD5 := hex.EncodeToString(hash[:])
		dbUser.Password = passwordMD5
	}

	phone := strings.TrimSpace(req.Form.Get("phone"))
	if len(phone) > 0 {
		dbUser.Phone = phone
	}

	address := strings.TrimSpace(req.Form.Get("address"))
	if len(address) > 0 {
		dbUser.Address = address
	}

	err = db.UpdateUser(c, dbUser)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Update error, please try again later", "/user/profile-edit")
		return
	}

	SetUserSession(res, req, dbUser)

	res.Header().Set("Location", "/user/profile")
	res.WriteHeader(http.StatusFound)
	return
}
