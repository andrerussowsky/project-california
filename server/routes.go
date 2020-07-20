package server

import (
	"github.com/gorilla/pat"
	"html/template"
	"net/http"
	"project-california/components"
	"project-california/controllers"
)

func Config(c *components.Components) *pat.Router {

	p := pat.New()

	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		controllers.Callback(c, res, req)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		controllers.Authenticate(c, res, req)
	})

	p.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		controllers.Logout(c, res, req)
	})

	p.Get("/user/sign-in", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserSignIn(c, res, req)
	})

	p.Post("/user/sign-in", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserSignInPost(c, res, req)
	})

	p.Get("/user/sign-up-complete", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserSignUpComplete(c, res, req)
	})

	p.Post("/user/sign-up-complete", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserSignUpCompletePost(c, res, req)
	})

	p.Get("/user/sign-up", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserSignUp(c, res, req)
	})

	p.Post("/user/sign-up", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserSignUpPost(c, res, req)
	})

	p.Get("/user/profile-edit", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserProfileEdit(c, res, req)
	})

	p.Post("/user/profile-edit", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserProfileEditPost(c, res, req)
	})

	p.Get("/user/profile", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserProfile(c, res, req)
	})

	p.Get("/user/forgot-password-reset", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserForgotPasswordReset(c, res, req)
	})

	p.Post("/user/forgot-password-reset", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserForgotPasswordResetPost(c, res, req)
	})

	p.Get("/user/forgot-password", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserForgotPassword(c, res, req)
	})

	p.Post("/user/forgot-password", func(res http.ResponseWriter, req *http.Request) {
		controllers.UserForgotPasswordPost(c, res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {

		sessionModel := controllers.LoadErrorSession(res, req)

		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, sessionModel)
	})

	return p
}

