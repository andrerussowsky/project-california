package controllers

import (
	"github.com/markbates/goth/gothic"
	"net/http"
	"project-california/models"
)

func RemoveUserSession(res http.ResponseWriter, req *http.Request) {
	session, _ := gothic.Store.Get(req, "user-session")
	delete(session.Values, "id")
	session.Save(req, res)
}

func SetUserSession(res http.ResponseWriter, req *http.Request, user models.User) {
	session, _ := gothic.Store.Get(req, "user-session")
	session.Values["id"] = user.ID
	session.Save(req, res)
}

func RemoveErrorSession(res http.ResponseWriter, req *http.Request) {
	session, _ := gothic.Store.Get(req, "error-session")
	delete(session.Values, "error")
	session.Save(req, res)
}

func SetErrorSession(res http.ResponseWriter, req *http.Request, error string) {
	session, _ := gothic.Store.Get(req, "error-session")
	session.Values["error"] = error
	session.Save(req, res)
}

func RemoveForgotPasswordSession(res http.ResponseWriter, req *http.Request) {
	session, _ := gothic.Store.Get(req, "forgot-password-session")
	delete(session.Values, "user_id")
	delete(session.Values, "key_id")
	session.Save(req, res)
}

func SetForgotPasswordSession(res http.ResponseWriter, req *http.Request, keyID, userID string) {
	session, _ := gothic.Store.Get(req, "forgot-password-session")
	session.Values["user_id"] = userID
	session.Values["key_id"] = keyID
	session.Save(req, res)
}
