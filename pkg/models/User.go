package models

type User struct {
	UserID   int    `db:"userid"`
	Username string `db:"username"`
	Token    string `db:"token"`
	Picture  string `db:"picture"`
	Status   string `db:"status"`
}
