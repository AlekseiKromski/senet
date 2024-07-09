package postgres

import (
	"alekseikromski.com/senet/modules/storage"
	"fmt"
)

func (p *Postgres) CreateChat(name, chatType, securityLevel string) (*storage.Chat, error) {
	query := "INSERT INTO chats (name, chat_type, security_level) VALUES ($1, $2, $3) RETURNING id, name, chat_type, security_level, created_at, updated_at, deleted_at"
	rows, err := p.db.Query(query, name, chatType, securityLevel)
	if err != nil {
		return nil, fmt.Errorf("cannot save chat: %v", err)
	}

	chat := &storage.Chat{}
	for rows.Next() {
		err := rows.Scan(&chat.Id, &chat.Name, &chat.ChatType, &chat.SecurityLevel, &chat.CreateAt, &chat.UpdatedAt, &chat.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}
	}

	return chat, nil
}

func (p *Postgres) IsChatBetweenUsersExists(u1id, u2id string) (string, error) {
	query := "SELECT chatid FROM chats_users where chatid IN (SELECT chatid FROM chats_users WHERE userid = $1) and userid = $2"
	rows, err := p.db.Query(query, u1id, u2id)
	if err != nil {
		return "", fmt.Errorf("cannot save chat: %v", err)
	}

	chatId := ""
	for rows.Next() {
		err := rows.Scan(&chatId)
		if err != nil {
			return "", fmt.Errorf("cannot read response from database: %v", err)
		}
	}

	if len(chatId) == 0 {
		return "", nil
	}

	return chatId, nil
}

func (p *Postgres) AddUserToChat(uid, cid string) error {
	query := "INSERT INTO chats_users (userid, chatid) VALUES ($1, $2)"
	if _, err := p.db.Exec(query, uid, cid); err != nil {
		return fmt.Errorf("cannot add user to chat: %v", err)
	}

	return nil
}

func (p *Postgres) GetChats(uid string) ([]*storage.Chat, error) {
	query := `
		SELECT chats.id,
			   chats.name,
			   chats.chat_type,
			   chats.security_level,
			   chats.created_at,
			   chats.updated_at,
			   chats.deleted_at,
			   u.id,
			   u.username,
			   u.first_name,
			   u.second_name,
			   u.image,
			   u.created_at,
			   u.updated_at,
			   u.deleted_at
		FROM chats
				 INNER JOIN public.chats_users cu on chats.id = cu.chatid
				 INNER JOIN public.users u on cu.userid = u.id
		WHERE chatid IN (SELECT chatid FROM chats_users WHERE userid = $1)
	`

	rows, err := p.db.Query(query, uid)
	if err != nil {
		return nil, fmt.Errorf("cannot get chats: %v", err)
	}

	chatsLocalMap := map[string]*storage.Chat{}
	for rows.Next() {
		chat := &storage.Chat{}
		user := &storage.User{}
		err := rows.Scan(
			&chat.Id,
			&chat.Name,
			&chat.ChatType,
			&chat.SecurityLevel,
			&chat.CreateAt,
			&chat.UpdatedAt,
			&chat.DeletedAt,
			&user.Id,
			&user.Username,
			&user.First_name,
			&user.Second_name,
			&user.Image,
			&user.CreateAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}

		chat.Users = append(chat.Users, user)
		chatFromMap := chatsLocalMap[chat.Id]

		if chatFromMap == nil {
			chatsLocalMap[chat.Id] = chat
			continue
		}

		// Push new user, no need to change any chat data
		chatFromMap.Users = append(chatFromMap.Users, user)
	}

	chats := []*storage.Chat{}
	for _, chat := range chatsLocalMap {
		chats = append(chats, chat)
	}

	return chats, nil
}
