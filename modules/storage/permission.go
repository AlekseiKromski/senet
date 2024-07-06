package storage

type StoragePermissions interface {
	GetPermissions() (map[string][]*Endpoint, error)
	CreatePermission(roleId, endpointId string) error
	GetEndpointIdsByRoleId(roleId string) ([]string, error)
	GetEndpointByRoleId(roleId string) ([]*Endpoint, error)
	DeletePermission(roleId, endpointId string) error
}
