package models

type Message struct {
	ID      string  `json:"id"`
	ChatID  string  `json:"chat_id"`
	UserID  string  `json:"user_id"`
	Message string  `json:"message"`
	Created string  `json:"created"`
	Updated string  `json:"updated"`
	Deleted *string `json:"deleted,omitempty"`
	User    User    `json:"user,omitempty"` //null
}
