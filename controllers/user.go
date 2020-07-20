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

	var sessionModel  = models.Session{}
	errorSession, _ := gothic.Store.Get(req, "error-session")
	if errorSession.Values["error"] != nil {
		sessionModel.Error=  fmt.Sprintf("%v", errorSession.Values["error"])
		RemoveErrorSession(res,req)
	}

	t, _ := template.ParseFiles("templates/sign-in.html")
	t.Execute(res, sessionModel)
}

func UserSignInPost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	email := strings.TrimSpace(req.Form.Get("email"))
	if len(email) == 0 {
		SetErrorSession(res, req, "Invalid email")

		res.Header().Set("Location", "/user/sign-in")
		res.WriteHeader(http.StatusFound)
		return
	}
	password := strings.TrimSpace(req.Form.Get("password"))
	if len(password) == 0 {
		SetErrorSession(res, req, "Invalid password")

		res.Header().Set("Location", "/user/sign-in")
		res.WriteHeader(http.StatusFound)
		return
	}

	dbUser, err := db.GetUserWithEmail(c, email)
	if err != nil {
		if err == sql.ErrNoRows {
			SetErrorSession(res, req, "User not found")

			res.Header().Set("Location", "/user/sign-in")
			res.WriteHeader(http.StatusFound)
			return
		}

		SetErrorSession(res, req, "Unexpected error")

		res.Header().Set("Location", "/user/sign-in")
		res.WriteHeader(http.StatusFound)
		return
	}

	hash := md5.Sum([]byte(password))
	passwordMD5 := hex.EncodeToString(hash[:])
	if dbUser.Password != passwordMD5 {
		SetErrorSession(res, req, "Wrong password")

		res.Header().Set("Location", "/user/sign-in")
		res.WriteHeader(http.StatusFound)
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

	var sessionModel = models.Session{}
	errorSession, _ := gothic.Store.Get(req, "error-session")
	if errorSession.Values["error"] != nil {
		sessionModel.Error = fmt.Sprintf("%v", errorSession.Values["error"])
		RemoveErrorSession(res,req)
	}

	t, _ := template.ParseFiles("templates/sign-up.html")
	t.Execute(res, sessionModel)
}

func UserSignUpPost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	email := strings.TrimSpace(req.Form.Get("email"))
	if len(email) == 0 {
		SetErrorSession(res, req, "Invalid email")

		res.Header().Set("Location", "/user/sign-up")
		res.WriteHeader(http.StatusFound)
		return
	}

	password := strings.TrimSpace(req.Form.Get("password"))
	if len(password) == 0 {
		SetErrorSession(res, req, "Invalid password")

		res.Header().Set("Location", "/user/sign-up")
		res.WriteHeader(http.StatusFound)
		return
	}

	confirmPassword := strings.TrimSpace(req.Form.Get("confirm-password"))
	if len(confirmPassword) == 0 {
		SetErrorSession(res, req, "Invalid password confirmation")

		res.Header().Set("Location", "/user/sign-up")
		res.WriteHeader(http.StatusFound)
		return
	}

	if password != confirmPassword {
		SetErrorSession(res, req, "Password doesnt match")

		res.Header().Set("Location", "/user/sign-up")
		res.WriteHeader(http.StatusFound)
		return
	}

	_, err := db.GetUserWithEmail(c, email)
	if err == nil {
		SetErrorSession(res, req, "Email is already taken")

		res.Header().Set("Location", "/user/sign-up")
		res.WriteHeader(http.StatusFound)
		return
	}

	hash := md5.Sum([]byte(password))
	passwordMD5 := hex.EncodeToString(hash[:])

	err = db.InsertUser(c, models.User{Email: email, Password: passwordMD5})
	if err != nil {
		SetErrorSession(res, req, "Insert error, please try again later")

		res.Header().Set("Location", "/user/sign-up")
		res.WriteHeader(http.StatusFound)
		return
	}

	dbUser, _ := db.GetUserWithEmail(c, email)

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

	var sessionModel = models.Session{}
	errorSession, _ := gothic.Store.Get(req, "error-session")
	if errorSession.Values["error"] != nil {
		sessionModel.Error = fmt.Sprintf("%v", errorSession.Values["error"])
		RemoveErrorSession(res,req)
	}

	userSession, _ := gothic.Store.Get(req, "user-session")
	if userSession.Values["id"] == nil {
		SetErrorSession(res, req, "Session expired")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	t, _ := template.ParseFiles("templates/sign-up-complete.html")
	t.Execute(res, sessionModel)
}

func UserSignUpCompletePost(c *components.Components, res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	userSession, _ := gothic.Store.Get(req, "user-session")
	if userSession.Values["id"] == nil {
		SetErrorSession(res, req, "Session expired")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	id := fmt.Sprintf("%v", userSession.Values["id"])

	dbUser, err := db.GetUser(c, id)
	if err != nil {
		SetErrorSession(res, req, "Unexpected error")

		res.Header().Set("Location", "/user/sign-up-complete")
		res.WriteHeader(http.StatusFound)
		return
	}

	name := strings.TrimSpace(req.Form.Get("name"))
	if len(name) > 0 {
		dbUser.Name = name
	}
	phone := strings.TrimSpace(req.Form.Get("phone"))
	if len(phone) > 0 {
		dbUser.Phone = phone
	}
	address := strings.TrimSpace(req.Form.Get("address"))
	if len(address) > 0 {
		dbUser.Address = address
	}

	dbUser.Status = models.UserStatusComplete

	err = db.UpdateUser(c, dbUser)
	if err != nil {
		SetErrorSession(res, req, "Update error, please try again later")

		res.Header().Set("Location", "/user/sign-up-complete")
		res.WriteHeader(http.StatusFound)
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
		SetErrorSession(res, req, "Session expired")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	id := fmt.Sprintf("%v", userSession.Values["id"])

	dbUser, err := db.GetUser(c, id)
	if err != nil {
		SetErrorSession(res, req, "Unexpected error")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	if dbUser.Status != models.UserStatusComplete {
		SetErrorSession(res, req, "Forbidden")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
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

	var sessionModel  = models.Session{}
	errorSession, _ := gothic.Store.Get(req, "error-session")
	if errorSession.Values["error"] != nil {
		sessionModel.Error = fmt.Sprintf("%v", errorSession.Values["error"])
		RemoveErrorSession(res,req)
	}

	userSession, _ := gothic.Store.Get(req, "user-session")
	if userSession.Values["id"] == nil {
		SetErrorSession(res, req, "Session expired")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	id := fmt.Sprintf("%v", userSession.Values["id"])

	dbUser, err := db.GetUser(c, id)
	if err != nil {
		SetErrorSession(res, req, "Unexpected error")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	if dbUser.Status != models.UserStatusComplete {
		SetErrorSession(res, req, "Forbidden")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
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
		SetErrorSession(res, req, "Session expired")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	id := fmt.Sprintf("%v", userSession.Values["id"])

	dbUser, err := db.GetUser(c, id)
	if err != nil {
		SetErrorSession(res, req, "Unexpected error")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	if dbUser.Status != models.UserStatusComplete {
		SetErrorSession(res, req, "Forbidden")

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusFound)
		return
	}

	email := strings.TrimSpace(req.Form.Get("email"))
	if len(email) == 0 {
		SetErrorSession(res, req, "Invalid email")

		res.Header().Set("Location", "/user/profile-edit")
		res.WriteHeader(http.StatusFound)
		return
	}

	if dbUser.Email != email {
		_, err := db.GetUserWithEmail(c, email)
		if err == nil {
			SetErrorSession(res, req, "Email is already taken")

			res.Header().Set("Location", "/user/profile-edit")
			res.WriteHeader(http.StatusFound)
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
		SetErrorSession(res, req, "Update error, please try again later")

		res.Header().Set("Location", "/user/profile-edit")
		res.WriteHeader(http.StatusFound)
		return
	}

	SetUserSession(res, req, dbUser)

	res.Header().Set("Location", "/user/profile")
	res.WriteHeader(http.StatusFound)
	return
}
