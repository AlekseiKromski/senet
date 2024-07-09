package storage

type Chat struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	ChatType      string  `json:"chat_type"`
	SecurityLevel string  `json:"security_level"`
	CreateAt      string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     *string `json:"deleted_at"`
	Users         []*User `json:"users,omitempty"`
}

type StorageChat interface {
	CreateChat(name, chatType, securityLevel string) (*Chat, error)
	AddUserToChat(uid, cid string) error
	IsChatBetweenUsersExists(u1id, u2id string) (string, error)
	GetChats(uid string) ([]*Chat, error)
}
