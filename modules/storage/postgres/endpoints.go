package postgres

import (
	"alekseikromski.com/senet/modules/storage"
	"fmt"
	"time"
)

func (p *Postgres) GetAllEndpoints() ([]*storage.Endpoint, error) {
	rows, err := p.db.Query(`
	SELECT id,
       urn,
       description,
       created_at,
       updated_at,
       deleted_at
	FROM endpoints
	WHERE deleted_at IS NULL
	ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("cannot get endpoints: %v", err)
	}
	defer rows.Close()

	endpoints := []*storage.Endpoint{}
	for rows.Next() {
		endpoint := &storage.Endpoint{}
		err := rows.Scan(&endpoint.Id, &endpoint.Urn, &endpoint.Description, &endpoint.CreatedAt, &endpoint.UpdatedAt, &endpoint.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}

		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}

func (p *Postgres) CreateEndpoint(urn, description string) error {
	query := "INSERT INTO endpoints (urn, description) VALUES ($1, $2)"
	if _, err := p.db.Exec(query, urn, description); err != nil {
		return fmt.Errorf("cannot save endpoint: %v", err)
	}

	return nil
}

func (p *Postgres) UpdateEndpoint(id, urn, description string) error {
	query := "UPDATE endpoints SET urn = $1, description = $2, updated_at = $3 WHERE id = $4"
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := p.db.Exec(query, urn, description, now, id); err != nil {
		return fmt.Errorf("cannot update endpoint: %v", err)
	}

	return nil
}

func (p *Postgres) DeleteEndpoint(id string) error {
	query := "UPDATE endpoints SET deleted_at = $1, updated_at = $2 WHERE id = $3"
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := p.db.Exec(query, now, now, id); err != nil {
		return fmt.Errorf("cannot delete endpoint: %v", err)
	}

	return nil
}
