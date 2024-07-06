package storage

import "time"

type Endpoint struct {
	Id          string     `json:"id"`
	Urn         string     `json:"urn"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type StorageEndpoints interface {
	GetAllEndpoints() ([]*Endpoint, error)
	CreateEndpoint(urn, description string) error
	UpdateEndpoint(id, urn, description string) error
	DeleteEndpoint(id string) error
}
