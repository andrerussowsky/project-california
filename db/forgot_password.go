package db

import (
	"database/sql"
	"fmt"
	"log"
	"project-california/components"
	"project-california/models"
)

func GetForgotPassword(c *components.Components, uuid string) (models.ForgotPassword, error) {

	mysql := Connect(c)
	defer mysql.Close()

	var fp models.ForgotPassword
	results, err := mysql.Query(fmt.Sprintf("select * from forgot_password where uuid='%s'", uuid))
	if err == nil {
		if results.Next() {
			err = results.Scan(&fp.ID, &fp.UserID, &fp.UUID)
		} else {
			err = sql.ErrNoRows
		}
	}

	log.Println("GetForgotPassword", uuid, err)

	return fp, err
}

func InsertForgotPassword(c *components.Components, fp models.ForgotPassword) error {

	mysql := Connect(c)
	defer mysql.Close()

	insForm, err := mysql.Prepare("INSERT INTO forgot_password(user_id,uuid) VALUES(?,?)")
	if err == nil {
		insForm.Exec(fp.UserID, fp.UUID)
	}

	log.Println("InsertForgotPassword", fp, err)

	return err
}

func DeleteForgotPassword(c *components.Components, userID string) error {

	mysql := Connect(c)
	defer mysql.Close()

	insForm, err := mysql.Prepare("DELETE FROM forgot_password where user_id=?")
	if err == nil {
		insForm.Exec(userID)
	}

	log.Println("DeleteForgotPassword", userID, err)

	return err
}
