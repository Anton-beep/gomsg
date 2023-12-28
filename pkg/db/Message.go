package db

import (
	"gomsg/pkg/models"
	"time"
)

func (d *APIDB) GetMessagesByChatID(id, quantity, timestamp int) ([]models.Message, error) {
	rows, err := d.db.Query("SELECT * FROM messages WHERE chatid = $1 AND timestamp <= $2 ORDER BY timestamp DESC LIMIT $3",
		id, timestamp, quantity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.MessageID, &message.ChatID, &message.Text, &message.SenderID, &message.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (d *APIDB) GetMessageByID(id int) (*models.Message, error) {
	rows, err := d.db.Query("SELECT * FROM messages WHERE messageid = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var message models.Message
	if !rows.Next() {
		return nil, nil
	}
	err = rows.Scan(&message.MessageID, &message.ChatID, &message.Text, &message.SenderID, &message.Timestamp)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (d *APIDB) CreateNewMessage(newMessage models.Message) (int, error) {
	var newID int
	timestamp := int(time.Now().Unix())
	err := d.db.QueryRow("INSERT INTO messages (chatid, text, senderid, timestamp) VALUES ($1, $2, $3, $4) RETURNING messageid",
		newMessage.ChatID, newMessage.Text, newMessage.SenderID, timestamp).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (d *APIDB) DeleteMessage(id int) (bool, error) {
	result, err := d.db.Exec("DELETE FROM messages WHERE messageid = $1", id)
	if err != nil {
		return false, err
	}

	return handleResultAfterEdit(result)
}

func (d *APIDB) EditMessage(id int, newText string) (bool, error) {
	result, err := d.db.Exec("UPDATE messages SET text = $1 WHERE messageid = $2", newText, id)
	if err != nil {
		return false, err
	}

	return handleResultAfterEdit(result)
}

func (d *APIDB) GetMessageUpdates(userID, timestamp int) ([]models.Message, error) {
	rows, err := d.db.Query("SELECT * FROM messages WHERE senderid = $1 AND timestamp >= $2 ORDER BY timestamp DESC",
		userID, timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.MessageID, &message.ChatID, &message.Text, &message.SenderID, &message.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
