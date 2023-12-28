package models

type Chat struct {
	ChatID    int    `db:"chatid"`
	ChatName  string `db:"chatname"`
	UsersIDs  []int  `db:"usersids"`
	Timestamp int    `db:"timestamp"`
}
