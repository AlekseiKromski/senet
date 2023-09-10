package models

type User struct {
	ID         string  `json:"id"`
	Username   string  `json:"username"`
	Password   *string `json:"password,omitempty"` //null for safety
	Image      *string `json:"image"`              // null
	LastOnline string  `json:"lastOnline"`
	CreatedAt  string  `json:"createdAt,omitempty"` //null
	UpdatedAt  string  `json:"updatedAt,omitempty"` //null
	DeletedAt  *string `json:"deletedAt,omitempty"` // null
}
