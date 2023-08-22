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

func (db *DbStorage) GetUser(username string) (*models.User, error) {
	rows, err := db.Conn.Query("SELECT * FROM users WHERE username = ? LIMIT 1", username)
	if err != nil {
		return nil, fmt.Errorf("cannot get users: %v", err)
	}

	defer rows.Close()

	user := &models.User{}

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.LastOnline, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot read all users from rows: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("problem with rows: %v", err)
	}

	return user, nil
}

func (db *DbStorage) CreateUser(id uuid.UUID, username, password string) error {
	query := `INSERT INTO users (id, username, password, lastOnline) VALUES (?, ?, ?, ?)`
	if _, err := db.Conn.Exec(query, id, username, password, time.Now()); err != nil {
		return fmt.Errorf("cannot create user: %v", err)
	}

	return nil
}
