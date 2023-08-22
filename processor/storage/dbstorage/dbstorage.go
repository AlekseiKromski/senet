package dbstorage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"senet/config"
	"senet/processor/storage/models"
	"time"
)

type DbStorage struct {
	Conn *sql.DB
}

func NewDbStorage(config *config.DbConfig) (*DbStorage, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4,utf8", config.Username, config.Password, config.Host, config.Database))
	if err != nil {
		return nil, fmt.Errorf("cannot start database: %v", err)
	}

	return &DbStorage{
		Conn: db,
	}, nil
}

func (db *DbStorage) GetUsers() ([]*models.User, error) {
	rows, err := db.Conn.Query("SELECT username, lastOnline, created, updated, deleted FROM users")
	if err != nil {
		return nil, fmt.Errorf("cannot get users: %v", err)
	}

	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Username, &user.LastOnline, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot read all users from rows: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (db *DbStorage) CreateUser(id uuid.UUID, username, password string) error {
	query := `INSERT INTO users (id, username, password, lastOnline) VALUES (?, ?, ?, ?)`
	if _, err := db.Conn.Exec(query, id, username, password, time.Now()); err != nil {
		return fmt.Errorf("cannot create user: %v", err)
	}

	return nil
}
