package models

type User struct {
	Username   string  `json:"username"`
	LastOnline string  `json:"lastOnline"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
	DeletedAt  *string `json:"deletedAt"` // null
}
