package db

import (
	"database/sql"
	"fmt"
	"log"
	"project-california/components"
	"project-california/models"
)

func GetUser(c *components.Components, id string) (models.User, error) {

	mysql := Connect(c)
	defer mysql.Close()

	var user models.User
	results, err := mysql.Query(fmt.Sprintf("select * from users where id='%s'", id))
	if err == nil {
		if results.Next() {
			err = results.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone, &user.Address, &user.Status)
		} else {
			err = sql.ErrNoRows
		}
	}

	log.Println("GetUser", user, err, id)

	return user, err
}

func GetUserWithEmail(c *components.Components, email string) (models.User, error) {

	mysql := Connect(c)
	defer mysql.Close()

	var user models.User
	results, err := mysql.Query(fmt.Sprintf("select * from users where email='%s'", email))
	if err == nil {
		if results.Next() {
			err = results.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone, &user.Address, &user.Status)
		} else {
			err = sql.ErrNoRows
		}
	}

	log.Println("GetUserWithEmail", email, user, err)

	return user, err
}

func InsertUser(c *components.Components, user models.User) error {

	mysql := Connect(c)
	defer mysql.Close()

	insForm, err := mysql.Prepare("INSERT INTO users(email,password,name,phone,address) VALUES(?,?,?,?,?)")
	if err == nil {
		insForm.Exec(user.Email, user.Password, user.Name, user.Phone, user.Address)
	}

	log.Println("InsertUser", user, err)

	return err
}

func UpdateUser(c *components.Components, user models.User) error {

	mysql := Connect(c)
	defer mysql.Close()

	insForm, err := mysql.Prepare("UPDATE users set email=?,password=?,name=?,phone=?,address=?,status=? WHERE id=?")
	if err == nil {
		insForm.Exec(user.Email, user.Password, user.Name, user.Phone, user.Address, user.Status, user.ID)
	}

	log.Println("UpdateUser", user, err)

	return err
}
