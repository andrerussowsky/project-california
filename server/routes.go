package server

import (
	"fmt"
	"github.com/gorilla/pat"
	"github.com/markbates/goth/gothic"
	"html/template"
	"log"
	"net/http"
	"project-california/components"
	"project-california/models"
)

func Config(c *components.Components) *pat.Router {

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)
	})

	p.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {

		results, _ := c.DB.Query("select * from users")
		for results.Next() {
			var users models.Users
			// for each row, scan the result into our tag composite object
			err := results.Scan(&users.ID, &users.Email)
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}
			// and then print out the tag's Name attribute
			log.Printf(users.Email)
		}

		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})

	return p
}
