package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"senet/processor/storage/creators/chat"
	"senet/processor/storage/models"
)

func (db *DbStorage) CreateChat(creator chat.ChatCreator) (models.Chat, error) {
	ctx := context.Background()
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return models.Chat{}, fmt.Errorf("cannot start transaction: %v", err)
	}

	ch, err := creator.Create(tx)
	if err != nil {
		if re := tx.Rollback(); re != nil && re != sql.ErrTxDone {
			return models.Chat{}, fmt.Errorf("cannot make rollback transaction: %v", err)
		}

		return models.Chat{}, fmt.Errorf("cannot create chat: %v", err)
	}

	//if err = tx.Commit(); err != nil { TODO
	//	return models.Chat{}, fmt.Errorf("cannot commit transaction: %v", err)
	//}

	return ch, nil
}

func (db *DbStorage) GetChats(userID string) ([]models.Chat, error) {
	chatsQuery := `
		SELECT c.id,
			   c.security_level,
			   c.created,
			   c.updated,
			   c.deleted,
			   u.id,
			   u.username,
			   u.image,
			   u.lastonline,
			   u.created,
			   u.updated,
			   u.deleted
		FROM   chats_users AS cu
			   INNER JOIN chats AS c
					   ON cu.chat_id = c.id
			   INNER JOIN users AS u
					   ON cu.user_id = u.id
		WHERE  chat_id IN (SELECT chat_id
						   FROM   chats_users
						   WHERE  user_id = ?); 
	`

	rows, err := db.Conn.Query(chatsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get chats: %v", err)
	}

	defer rows.Close()
	var ch models.Chat
	var chats []models.Chat

	for rows.Next() {
		tchat := models.Chat{}
		tuser := models.User{}
		err := rows.Scan(
			&tchat.ID,
			&tchat.SecurityLevel,
			&tchat.Created,
			&tchat.Updated,
			&tchat.Deleted,
			&tuser.ID,
			&tuser.Username,
			&tuser.Image,
			&tuser.LastOnline,
			&tuser.CreatedAt,
			&tuser.UpdatedAt,
			&tuser.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot get chat: %v", err)
		}

		if ch.ID == tchat.ID {
			ch.Users = append(ch.Users, tuser)
		} else {
			if len(ch.ID) != 0 {
				//save old temp chat
				chats = append(chats, ch)
			}

			ch = tchat
			ch.Users = append(ch.Users, tuser)
		}
	}

	chats = append(chats, ch)

	messageQuery := `
		SELECT m.*, u.id, u.username,u.Image, u.lastOnline FROM messages as m
		INNER JOIN users as u ON m.user_id = u.id
		 WHERE chat_id = ? AND user_id = ? AND m.deleted IS NULL LIMIT 25 
	`
	for _, chat := range chats {
		rows, err := db.Conn.Query(messageQuery, chat.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("cannot get messages: %v", err)
		}

		defer rows.Close()
		chat.Messages = []models.Message{}
		user := models.User{}

		for rows.Next() {
			message := models.Message{}
			err := rows.Scan(
				&message.ID,
				&message.ChatID,
				&message.UserID,
				&message.Message,
				&message.Created,
				&message.Updated,
				&message.Deleted,
				&user.ID,
				&user.Username,
				&user.Image,
				&user.LastOnline,
			)
			if err != nil {
				return nil, fmt.Errorf("cannot get chat: %v", err)
			}

			message.User = user
			chat.Messages = append(chat.Messages, message)
		}
	}

	return chats, nil
}
