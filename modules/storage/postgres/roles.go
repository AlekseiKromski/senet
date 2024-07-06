package postgres

import (
	"alekseikromski.com/senet/modules/storage"
	"fmt"
	"sort"
	"time"
)

func (p *Postgres) GetAllRoles() ([]*storage.Role, error) {
	rows, err := p.db.Query(`
	SELECT
		roles.id,
		roles.name,
		roles.created_at,
		roles.updated_at,
		endpoints.id,
		endpoints.urn,
		endpoints.description,
		endpoints.created_at,
		endpoints.updated_at
	FROM roles_endpoints
	INNER JOIN endpoints on endpoints.id = roles_endpoints.endpointuuid
	INNER JOIN roles on roles.id = roles_endpoints.roleuuid
	WHERE roles.deleted_at IS NULL AND endpoints.deleted_at IS NULL AND roles_endpoints.deleted_at IS NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("cannot get roles: %v", err)
	}
	defer rows.Close()

	roles := map[string]*storage.Role{}
	for rows.Next() {
		role := &storage.Role{}
		endpoint := &storage.Endpoint{}
		err := rows.Scan(&role.Id, &role.Name, &role.CreatedAt, &role.UpdatedAt, &endpoint.Id, &endpoint.Urn, &endpoint.Description, &endpoint.CreatedAt, &endpoint.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}

		if roles[role.Id] != nil {
			roles[role.Id].Endpoints = append(roles[role.Id].Endpoints, endpoint)
		} else {
			role.Endpoints = append(role.Endpoints, endpoint)
			roles[role.Id] = role
		}
	}

	processedRoles := []*storage.Role{}
	for _, v := range roles {
		processedRoles = append(processedRoles, v)
	}

	sort.Slice(processedRoles, func(i, j int) bool {
		return processedRoles[i].CreatedAt.After(processedRoles[j].CreatedAt)
	})

	return processedRoles, nil
}

func (p *Postgres) CreateRole(name string) (*string, error) {
	query := "INSERT INTO roles (name) VALUES ($1) RETURNING id"

	rows, err := p.db.Query(query, name)
	if err != nil {
		return nil, fmt.Errorf("cannot save role: %v", err)
	}

	roleId := ""
	for rows.Next() {
		err := rows.Scan(&roleId)
		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}
	}

	if len(roleId) == 0 {
		return nil, fmt.Errorf("role id is null")
	}
	return &roleId, nil
}

func (p *Postgres) UpdateRole(id, name string) error {
	query := "UPDATE roles SET name = $1, updated_at = $2 WHERE id = $3"
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := p.db.Exec(query, name, now, id); err != nil {
		return fmt.Errorf("cannot update role: %v", err)
	}

	return nil
}

func (p *Postgres) DeleteRole(id string) error {
	query := "UPDATE roles SET deleted_at = $1, updated_at = $2 WHERE id = $3"
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := p.db.Exec(query, now, now, id); err != nil {
		return fmt.Errorf("cannot delete role: %v", err)
	}

	return nil
}
