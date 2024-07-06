package postgres

import (
	"alekseikromski.com/senet/modules/storage"
	"fmt"
	"time"
)

func (p *Postgres) GetPermissions() (map[string][]*storage.Endpoint, error) {
	rows, err := p.db.Query("SELECT roles.id, endpoints.urn FROM roles_endpoints INNER JOIN roles ON roles.ID = roles_endpoints.roleuuid INNER JOIN endpoints ON endpoints.ID = roles_endpoints.endpointuuid")
	if err != nil {
		return nil, fmt.Errorf("cannot get roles / endpoint permissions: %v", err)
	}
	defer rows.Close()

	permissions := map[string][]*storage.Endpoint{}
	for rows.Next() {
		role_id := ""
		urn := ""
		err := rows.Scan(&role_id, &urn)
		if err != nil {
			return nil, fmt.Errorf("cannot read response from database: %v", err)
		}

		permissions[role_id] = append(permissions[role_id], &storage.Endpoint{
			Urn: urn,
		})
	}

	return permissions, nil
}

func (p *Postgres) CreatePermission(roleId, endpointId string) error {
	query := "INSERT INTO roles_endpoints (roleuuid, endpointuuid) VALUES ($1, $2)"
	if _, err := p.db.Exec(query, roleId, endpointId); err != nil {
		return fmt.Errorf("cannot save permission: %v", err)
	}

	return nil
}

func (p *Postgres) GetEndpointIdsByRoleId(roleId string) ([]string, error) {
	endpoints, err := p.GetEndpointByRoleId(roleId)
	if err != nil {
		return nil, fmt.Errorf("cannot get endpoints: %v", err)
	}

	endpointids := []string{}
	for _, endpoint := range endpoints {
		endpointids = append(endpointids, endpoint.Id)
	}

	return endpointids, nil
}

func (p *Postgres) GetEndpointByRoleId(roleId string) ([]*storage.Endpoint, error) {
	rows, err := p.db.Query(`
SELECT endpoints.id, endpoints.urn, endpoints.description, endpoints.created_at, endpoints.updated_at, endpoints.deleted_at
FROM roles_endpoints
         INNER JOIN roles ON roles.ID = roles_endpoints.roleuuid
         INNER JOIN endpoints ON endpoints.ID = roles_endpoints.endpointuuid
        WHERE roles.ID = $1 AND endpoints.deleted_at IS NULL
`, roleId)
	if err != nil {
		return nil, fmt.Errorf("cannot get endpoint: %v", err)
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

func (p *Postgres) DeletePermission(roleId, endpointId string) error {
	query := "UPDATE roles_endpoints SET deleted_at = $1, updated_at = $2 WHERE endpointuuid = $3 AND roleuuid = $4"
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := p.db.Exec(query, now, now, endpointId, roleId); err != nil {
		return fmt.Errorf("cannot delete permission: %v", err)
	}

	return nil
}
