package db

import (
	"database/sql"
	"fmt"
	"project-california/components"
	"project-california/models"
)

func GetForgotPassword(c *components.Components, uuid string) (models.ForgotPassword, error) {

	var fp models.ForgotPassword
	results, err := c.DB.Query(fmt.Sprintf("select * from forgot_password where uuid='%s'", uuid))
	if err == nil {
		if results.Next() {
			err = results.Scan(&fp.ID, &fp.UserID, &fp.UUID)
		} else {
			err = sql.ErrNoRows
		}
	}

	return fp, err
}

func InsertForgotPassword(c *components.Components, fp models.ForgotPassword) error {

	insForm, err := c.DB.Prepare("INSERT INTO forgot_password(user_id,uuid) VALUES(?,?)")
	if err == nil {
		insForm.Exec(fp.UserID, fp.UUID)
	}

	return err
}

func DeleteForgotPassword(c *components.Components, userID string) error {

	insForm, err := c.DB.Prepare("DELETE FROM forgot_password where user_id=?")
	if err == nil {
		insForm.Exec(userID)
	}

	return err
}
