package models

type Message struct {
	MessageID  int    `db:"messageid" json:"messageID"`
	ChatID     int    `db:"chatid" json:"chatID"`
	Text       string `db:"text" json:"text"`
	SenderID   int    `db:"senderid" json:"senderID"`
	SenderName string `db:"sendername" json:"senderName"`
	Timestamp  int    `db:"timestamp" json:"timestamp"`
}
