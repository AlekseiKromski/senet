package postgres

import (
	"alekseikromski.com/senet/modules/storage"
	"fmt"
	"time"
)

func (p *Postgres) CreateUser(username, email, first_name, second_name, password, role string) error {
	query := "INSERT INTO users (username, email, first_name, second_name, image, password, role) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if _, err := p.db.Exec(query, username, email, first_name, second_name, "", password, role); err != nil {
		return fmt.Errorf("cannot save user: %v", err)
	}

	return nil
}

func (p *Postgres) UpdateUser(id, username, email, first_name, second_name, password, role string) error {
	query := "UPDATE users SET username = $1, email = $2, first_name = $3, second_name = $4, password = $5, role = $6, updated_at = $7 WHERE id = $8"
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := p.db.Exec(query, username, email, first_name, second_name, password, role, now, id); err != nil {
		return fmt.Errorf("cannot update user: %v", err)
	}

	return nil
}

func (p *Postgres) DeleteUser(id string) error {
	query := "UPDATE users SET deleted_at = $1, updated_at = $2 WHERE id = $3"
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := p.db.Exec(query, now, now, id); err != nil {
		return fmt.Errorf("cannot delete user: %v", err)
	}

	return nil
}

func (p *Postgres) GetUserById(id string) (*storage.User, error) {
	rows, err := p.db.Query("SELECT id, username, first_name, second_name, image, email, password, role FROM users WHERE id = $1 AND deleted_at IS NULL LIMIT 1", id)
	if err != nil {
		return nil, fmt.Errorf("cannot send request to check migrations tables: %v", err)
	}
	defer rows.Close()

	user := &storage.User{}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.First_name, &user.Second_name, &user.Image, &user.Email, &user.Password, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}
	}

	return user, nil
}

func (p *Postgres) GetUserByUsername(username string) (*storage.User, error) {
	rows, err := p.db.Query("SELECT id, username, first_name, second_name, image, email, password, role FROM users WHERE username = $1 AND deleted_at IS NULL LIMIT 1", username)
	if err != nil {
		return nil, fmt.Errorf("cannot send request to check migrations tables: %v", err)
	}
	defer rows.Close()

	user := &storage.User{}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.First_name, &user.Second_name, &user.Image, &user.Email, &user.Password, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}
	}

	return user, nil
}

func (p *Postgres) FindUsersByUsername(username string) ([]*storage.User, error) {
	likePattern := "%" + username + "%"
	rows, err := p.db.Query("SELECT id, username, image FROM users WHERE username LIKE $1 AND deleted_at IS NULL", likePattern)
	if err != nil {
		return nil, fmt.Errorf("cannot send request to check migrations tables: %v", err)
	}
	defer rows.Close()

	users := []*storage.User{}
	for rows.Next() {
		user := &storage.User{}
		err := rows.Scan(&user.Id, &user.Username, &user.Image)
		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}

		users = append(users, user)
	}

	return users, nil
}
