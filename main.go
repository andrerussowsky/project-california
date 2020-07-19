package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"log"
	"net/http"
	"os"
	"project-california/components"
	"project-california/db"
	"project-california/server"
	"project-california/settings"
)


func main() {

	c := &components.Components{
		Settings: settings.Config(),
	}
	c.DB = db.Config(c)
	defer c.DB.Close()

	configSession(c)

	service := server.Config(c)
	log.Fatal(http.ListenAndServe(getPort(), service))
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":3000"
}

func configSession(c *components.Components) {

	store := sessions.NewCookieStore([]byte("project-california-session"))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false
	gothic.Store = store

	goth.UseProviders(
		google.New(
			c.Settings.GetString("google.clientKey"),
			c.Settings.GetString("google.secret"),
			c.Settings.GetString("google.callback"),
			"email",
			"profile"),
	)
}