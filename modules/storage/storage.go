package storage

type Storage interface {
	StorageUser        // StorageUser - crud interface for working with users
	StoragePermissions // StoragePermissions - all commands related to permissions (roles / endpoints)
	StorageEndpoints   // StorageEndpoints - all commands related to endpoints
	StorageRoles       // StorageRoles - all commands related to roles
}
