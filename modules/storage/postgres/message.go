package postgres

import (
	"alekseikromski.com/senet/modules/storage"
	"fmt"
)

func (p *Postgres) CreateMessage(cid, sid, message string) (*storage.Message, error) {
	query := "INSERT INTO messages (chatid, senderid, message) VALUES ($1, $2, $3) RETURNING *"

	rows, err := p.db.Query(query, cid, sid, message)
	if err != nil {
		return nil, fmt.Errorf("cannot create message: %v", err)
	}

	m := &storage.Message{}
	for rows.Next() {
		err := rows.Scan(&m.Id, &m.ChatId, &m.SenderId, &m.Message, &m.CreateAt, &m.UpdatedAt, &m.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}
	}

	return m, nil
}

func (p *Postgres) GetMessagesByChatId(cid string, offset, limit int) ([]*storage.Message, error) {
	query := "SELECT * FROM messages WHERE chatid = $1 ORDER BY created_at DESC OFFSET $2 LIMIT $3"

	rows, err := p.db.Query(query, cid, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("cannot get messages: %v", err)
	}

	ms := []*storage.Message{}
	for rows.Next() {
		m := &storage.Message{}
		err := rows.Scan(&m.Id, &m.ChatId, &m.SenderId, &m.Message, &m.CreateAt, &m.UpdatedAt, &m.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}
		ms = append(ms, m)
	}

	return ms, nil
}
