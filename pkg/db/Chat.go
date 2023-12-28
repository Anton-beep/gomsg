package db

import (
	"database/sql"
	"github.com/lib/pq"
	"gomsg/pkg/models"
	"time"
)

func (d *APIDB) GetChatsByUserID(id int) ([]models.Chat, error) {
	rows, err := d.db.Query("SELECT * FROM chats WHERE $1 = ANY(usersids)", id)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		var usersIDs []sql.NullInt64
		err = rows.Scan(&chat.ChatID, &chat.ChatName, pq.Array(&usersIDs), &chat.Timestamp)
		if err != nil {
			return nil, err
		}
		for _, userID := range usersIDs {
			if userID.Valid {
				chat.UsersIDs = append(chat.UsersIDs, int(userID.Int64))
			}
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (d *APIDB) GetChatByID(id int) (*models.Chat, error) {
	rows, err := d.db.Query("SELECT * FROM chats WHERE chatid = $1", id)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	var chat models.Chat
	var usersIDs []sql.NullInt64
	if !rows.Next() {
		return nil, nil
	}
	err = rows.Scan(&chat.ChatID, &chat.ChatName, pq.Array(&usersIDs), &chat.Timestamp)
	if err != nil {
		return nil, err
	}
	for _, userID := range usersIDs {
		if userID.Valid {
			chat.UsersIDs = append(chat.UsersIDs, int(userID.Int64))
		}
	}
	return &chat, nil
}

func (d *APIDB) CreateNewChat(chat models.Chat) (int, error) {
	var newID int
	timestamp := int(time.Now().Unix())
	err := d.db.QueryRow("INSERT INTO chats (chatname, usersids, timestamp) VALUES ($1, $2, $3) RETURNING chatid",
		chat.ChatName, pq.Array(chat.UsersIDs), timestamp).Scan(&newID)
	return newID, err
}

func (d *APIDB) DeleteChat(chatID int) (bool, error) {
	result, err := d.db.Exec("DELETE FROM chats WHERE chatid = $1", chatID)
	if err != nil {
		return false, nil
	}
	return handleResultAfterEdit(result)
}

func (d *APIDB) EditChatNameByChatID(chatID int, newChatName string) (bool, error) {
	timestamp := int(time.Now().Unix())
	result, err := d.db.Exec("UPDATE chats SET chatname = $1, timestamp = $2 WHERE chatid = $3", newChatName, timestamp, chatID)
	if err != nil {
		return false, nil
	}
	return handleResultAfterEdit(result)
}

func (d *APIDB) AddUserToChat(chatID int, userID int) (bool, error) {
	result, err := d.db.Exec("UPDATE chats SET usersids = array_append(usersids, $1) WHERE chatid = $2", userID, chatID)
	if err != nil {
		return false, nil
	}
	return handleResultAfterEdit(result)
}

func (d *APIDB) RemoveUserFromChat(chatID int, userID int) error {
	_, err := d.db.Exec("UPDATE chats SET usersids = array_remove(usersids, $1) WHERE chatid = $2", userID, chatID)
	return err
}
