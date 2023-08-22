package storage

import (
	"github.com/google/uuid"
	"senet/processor/storage/models"
)

type Storage interface {
	GetUser(username string) (*models.User, error)
	CreateUser(id uuid.UUID, username, password string) error
}
