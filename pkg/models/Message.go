package models

type Message struct {
	MessageID int    `db:"messageid"`
	ChatID    int    `db:"chatid"`
	Text      string `db:"text"`
	SenderID  int    `db:"senderid"`
	Timestamp int    `db:"timestamp"`
}
