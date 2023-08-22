package models

type ChatKey struct {
	ID                        string  `json:"id"`
	ChatID                    string  `json:"chat_id"`
	UserID                    string  `json:"user_id"`
	UserPublicKey             string  `json:"user_public_key"`
	ServerEncryptedPrivateKey string  `json:"server_encrypted_private_key"`
	Created                   string  `json:"created"`
	Updated                   string  `json:"updated"`
	Deleted                   *string `json:"deleted,omitempty"`
}
