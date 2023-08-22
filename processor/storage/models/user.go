package models

type User struct {
	ID         string  `json:"ID"`
	Username   string  `json:"username"`
	LastOnline string  `json:"lastOnline"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
	DeletedAt  *string `json:"deletedAt"` // null
}
