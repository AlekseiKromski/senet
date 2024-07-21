package storage

type Message struct {
	Id        string  `json:"id"`
	ChatId    string  `json:"chat_id"`
	SenderId  string  `json:"sender_id"`
	Message   string  `json:"message"`
	CreateAt  string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type StorageMessage interface {
	CreateMessage(cid, sid, message string) (*Message, error)
	GetMessagesByChatId(cid string, offset, limit int) ([]*Message, error)
}
