package models

type Chat struct {
	ID               string    `json:"id"`
	SecurityLevel    string    `json:"security_level"`
	ServerPublicKey  *string   `json:"server_public_key,omitempty"`
	ServerPrivateKey *string   `json:"server_private_key,omitempty"`
	Messages         []Message `json:"messages"`
	Users            []User    `json:"users"`
	Created          string    `json:"created"`
	Updated          string    `json:"updated"`
	Deleted          *string   `json:"deleted"`
}
