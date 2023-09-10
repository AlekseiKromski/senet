package chat

import (
	"database/sql"
	"fmt"
	"senet/processor/storage/models"
)

type PrivateChatCreator struct {
	ChatID      string
	RequestorID string
	ReceiverID  string
	chatType    string
}

func NewPrivateChatCreator(chatID, requestorID, receiverID string) *PrivateChatCreator {
	return &PrivateChatCreator{
		ChatID:      chatID,
		RequestorID: requestorID,
		ReceiverID:  receiverID,
		chatType:    "private",
	}
}

func (pcc *PrivateChatCreator) Create(tx *sql.Tx) (models.Chat, error) {

	//Check availability
	isAvailable, err := pcc.isAvailable(tx)
	if err != nil {
		return models.Chat{}, fmt.Errorf("cannot check isAvailable: %v", err)
	}
	if !isAvailable {
		return models.Chat{}, fmt.Errorf("users are not available: %v", err)
	}

	//Create chat
	if err := pcc.createChat(tx); err != nil {
		return models.Chat{}, fmt.Errorf("cannot create chat: %v", err)
	}

	//Get chat
	chat, err := pcc.getChat(tx)
	if err != nil {
		return models.Chat{}, fmt.Errorf("cannot get created chat: %v", err)
	}

	return chat, nil
}

func (pcc *PrivateChatCreator) GetMembers() []string {
	return []string{pcc.RequestorID, pcc.ReceiverID}
}

func (pcc *PrivateChatCreator) isAvailable(tx *sql.Tx) (bool, error) {
	//Check chats
	checkQuery := `
			 SELECT COUNT(id) FROM chats_users WHERE chat_id in (
				SELECT chat_id FROM chats_users WHERE user_id = ?
			) AND user_id = ?;
		`
	rows, err := tx.Query(checkQuery, pcc.RequestorID, pcc.ReceiverID)

	if err != nil {
		return false, fmt.Errorf("cannot get chats: %v", err)
	}

	defer rows.Close()
	var count int

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return false, fmt.Errorf("cannot get chats: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		return false, fmt.Errorf("cannot get chats: %v", err)
	}

	if count != 0 {
		return false, fmt.Errorf("cannot create chat, already exists: %v", err)
	}

	return true, nil
}

func (pcc *PrivateChatCreator) createChat(tx *sql.Tx) error {
	query := `INSERT INTO chats (id, security_level) VALUES (?, ?)`
	if _, err := tx.Exec(query, pcc.ChatID, pcc.chatType); err != nil {
		return fmt.Errorf("cannot create new chat: %v", err)
	}

	for _, user := range []string{pcc.RequestorID, pcc.ReceiverID} {
		query := `INSERT INTO chats_users (chat_id, user_id) VALUES (?, ?)`
		if _, err := tx.Exec(query, pcc.ChatID, user); err != nil {
			return fmt.Errorf("cannot create new chat: %v", err)
		}
	}

	return nil
}

func (pcc *PrivateChatCreator) getChat(tx *sql.Tx) (models.Chat, error) {
	getChat := `
			 SELECT * FROM chats WHERE id = ?;
		`
	rows, err := tx.Query(getChat, pcc.ChatID)

	if err != nil {
		return models.Chat{}, fmt.Errorf("cannot get chats: %v", err)
	}

	defer rows.Close()
	chat := models.Chat{}

	for rows.Next() {
		err := rows.Scan(&chat.ID, &chat.SecurityLevel, &chat.ServerPublicKey, &chat.ServerPrivateKey, &chat.Created, &chat.Updated, &chat.Deleted)
		if err != nil {
			return models.Chat{}, fmt.Errorf("cannot get chats: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		return models.Chat{}, fmt.Errorf("cannot get chats: %v", err)
	}

	//get users

	chat.Messages = []models.Message{}

	requestorUser, err := pcc.getUser(tx, pcc.RequestorID)
	if err != nil {
		return models.Chat{}, fmt.Errorf("cannot get requestor user: %v", err)
	}

	receiverUser, err := pcc.getUser(tx, pcc.ReceiverID)
	if err != nil {
		return models.Chat{}, fmt.Errorf("cannot get receiver user: %v", err)
	}

	chat.Users = []models.User{}
	chat.Users = []models.User{requestorUser, receiverUser}
	return chat, nil
}

func (pcc *PrivateChatCreator) getUser(tx *sql.Tx, userID string) (models.User, error) {
	getUser := `
			 SELECT * FROM users WHERE id = ?;
		`
	rows, err := tx.Query(getUser, userID)

	if err != nil {
		return models.User{}, fmt.Errorf("cannot get user by id (%s): %v", userID, err)
	}

	defer rows.Close()
	user := models.User{}

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Image, &user.LastOnline, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return models.User{}, fmt.Errorf("cannot read all users from rows: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		return models.User{}, fmt.Errorf("cannot get user by id (%s): %v", userID, err)
	}

	user.Password = nil

	return user, nil
}
