package chat

import (
	"database/sql"
	"senet/processor/storage/models"
)

type ChatCreator interface {
	Create(tx *sql.Tx) (models.Chat, error)
	GetMembers() []string
}
