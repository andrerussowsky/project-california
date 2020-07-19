package db

import (
	"database/sql"
	"fmt"
	"project-california/components"
)

func Config(c *components.Components) *sql.DB {

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s",
			c.Settings.GetString("mysql.username"),
			c.Settings.GetString("mysql.password"),
			c.Settings.GetString("mysql.host"),
			c.Settings.GetString("mysql.database")))
	if err != nil {
		panic(err.Error())
	}

	return db
}
