package db

import (
	"database/sql"
	"fmt"
	"project-california/components"
	"project-california/models"
)

func GetUser(c *components.Components, id string) (models.User, error) {

	var user models.User
	results, err := c.DB.Query(fmt.Sprintf("select * from users where id='%s'", id))
	if err == nil {
		if results.Next() {
			err = results.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone, &user.Address, &user.Status)
		} else {
			err = sql.ErrNoRows
		}
	}

	return user, err
}

func GetUserWithEmail(c *components.Components, email string) (models.User, error) {

	var user models.User
	results, err := c.DB.Query(fmt.Sprintf("select * from users where email='%s'", email))
	if err == nil {
		if results.Next() {
			err = results.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone, &user.Address, &user.Status)
		} else {
			err = sql.ErrNoRows
		}
	}

	return user, err
}

func InsertUser(c *components.Components, user models.User) error {

	insForm, err := c.DB.Prepare("INSERT INTO users(email,password,name,phone,address) VALUES(?,?,?,?,?)")
	if err == nil {
		insForm.Exec(user.Email, user.Password, user.Name, user.Phone, user.Address)
	}

	return err
}

func UpdateUser(c *components.Components, user models.User) error {

	insForm, err := c.DB.Prepare("UPDATE users set email=?,password=?,name=?,phone=?,address=?,status=? WHERE id=?")
	if err == nil {
		insForm.Exec(user.Email, user.Password, user.Name, user.Phone, user.Address, user.Status, user.ID)
	}

	return err
}
