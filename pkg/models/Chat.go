package models

type Chat struct {
	ChatID    int    `db:"chatid" json:"chatID"`
	ChatName  string `db:"chatname" json:"chatName"`
	UsersIDs  []int  `db:"usersids" json:"usersIDs"`
	Timestamp int    `db:"timestamp" json:"timestamp"`
}
