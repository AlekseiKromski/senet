package storage

import "time"

type Role struct {
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	Endpoints []*Endpoint `json:"endpoints,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt *time.Time  `json:"deleted_at"`
}

type StorageRoles interface {
	GetAllRoles() ([]*Role, error)
	CreateRole(name string) (*string, error)
	UpdateRole(id, name string) error
	DeleteRole(id string) error
}
