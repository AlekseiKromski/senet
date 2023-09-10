package storage

import (
	"github.com/google/uuid"
	"senet/processor/storage/creators/chat"
	"senet/processor/storage/models"
)

type Storage interface {
	GetUser(username string, likeMode bool) ([]models.User, error)
	CreateUser(id uuid.UUID, username, password string) error
	CreateChat(creator chat.ChatCreator) (models.Chat, error)
	GetChats(userID string) ([]models.Chat, error)
}
