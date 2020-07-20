package controllers

import (
	"database/sql"
	"github.com/markbates/goth/gothic"
	"net/http"
	"project-california/components"
	"project-california/db"
	"project-california/models"
)

func Authenticate(c *components.Components, res http.ResponseWriter, req *http.Request) {
	gothic.BeginAuthHandler(res, req)
}

func Callback(c *components.Components, res http.ResponseWriter, req *http.Request) {

	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		RedirectWithErrorMessage(res, req, "Authenticate error, please try again later", "/")
		return
	}

	dbUser, err := db.GetUserWithEmail(c, user.Email)
	if err != nil && err == sql.ErrNoRows {

		err = db.InsertUser(c, models.User{Email: user.Email, Name: user.Name})
		if err != nil {
			RedirectWithErrorMessage(res, req, "Insert error, please try again later", "/")
			return
		}

		dbUser, _ = db.GetUserWithEmail(c, user.Email)
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

func Logout(c *components.Components, res http.ResponseWriter, req *http.Request) {

	RemoveUserSession(res, req)
	gothic.Logout(res, req)

	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusFound)
}

