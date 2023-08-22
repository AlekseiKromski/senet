package storage

import "senet/processor/storage/models"

type Storage interface {
	GetUsers() ([]*models.User, error)
}
