package db

import (
	"go.uber.org/zap"
	"gomsg/pkg/models"
)

func (d *APIDB) GetUserByUsername(name string) (*models.User, error) {
	rows, err := d.db.Query("SELECT * FROM users WHERE username = $1", name)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	var user models.User
	if !rows.Next() {
		return nil, nil
	}
	err = rows.Scan(&user.UserID, &user.Username, &user.Token, &user.Picture, &user.Status)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *APIDB) GetUserByID(id int) (*models.User, error) {
	rows, err := d.db.Query("SELECT * FROM users WHERE userid = $1", id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			zap.L().Error(err.Error())
		}
	}()

	var user models.User
	if !rows.Next() {
		return nil, nil
	}
	err = rows.Scan(&user.UserID, &user.Username, &user.Token, &user.Picture, &user.Status)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *APIDB) CreateNewUser(newUser models.User) (int, error) {
	var newID int
	err := d.db.QueryRow("INSERT INTO users (username, token, picture, status) VALUES ($1, $2, $3, $4) RETURNING userid",
		newUser.Username, newUser.Token, newUser.Picture, newUser.Status).Scan(&newID)
	return newID, err
}

// DeleteUserByUsername returns isDeleted and error
func (d *APIDB) DeleteUserByUsername(username string) (bool, error) {
	result, err := d.db.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		return false, err
	}

	return handleResultAfterEdit(result)
}

// DeleteUserByID returns isDeleted and error
func (d *APIDB) DeleteUserByID(userID int) (bool, error) {
	result, err := d.db.Exec("DELETE FROM users WHERE userid = $1", userID)
	if err != nil {
		return false, err
	}

	return handleResultAfterEdit(result)
}

func (d *APIDB) EditStatus(userID int, newStatus string) (bool, error) {
	result, err := d.db.Exec("UPDATE users SET status = $1 WHERE userid = $2", newStatus, userID)
	if err != nil {
		return false, err
	}

	return handleResultAfterEdit(result)
}

func (d *APIDB) EditPicture(userID int, newPicture string) (bool, error) {
	result, err := d.db.Exec("UPDATE users SET picture = $1 WHERE userid = $2", newPicture, userID)
	if err != nil {
		return false, err
	}

	return handleResultAfterEdit(result)
}
