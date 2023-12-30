package models

type User struct {
	UserID   int    `db:"userid" json:"userID"`
	Username string `db:"username" json:"username"`
	Token    string `db:"token" json:"token"`
	Picture  string `db:"picture" json:"picture"`
	Status   string `db:"status" json:"status"`
}
