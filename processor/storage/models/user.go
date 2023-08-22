package models

type User struct {
	ID         string  `json:"id"`
	Username   string  `json:"username"`
	Password   *string `json:"password,omitempty"` //null for safety
	LastOnline string  `json:"lastOnline"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
	DeletedAt  *string `json:"deletedAt"` // null
}
