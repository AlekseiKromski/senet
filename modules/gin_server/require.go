package gin_server

import (
	"fmt"

	server_key_storage "alekseikromski.com/senet/modules/server-key-storage"
	"alekseikromski.com/senet/modules/storage"
	"github.com/AlekseiKromski/server-core/core"
)

func (s *Server) Require() []string {
	return []string{
		"storage",
		"server-key-storage",
	}
}

func (s *Server) getStorageFromRequirement(requirements map[string]core.Module) (storage.Storage, error) {
	storage, ok := requirements["storage"].(storage.Storage)
	if !ok {
		s.Log("All requirements list")
		for k, v := range requirements {
			s.Log("Requirement", k, v.Signature())
		}
		return nil, fmt.Errorf("requiremnt list has wrong storage requirement")
	}

	return storage, nil
}

func (s *Server) getServerKeyStorageRequirement(requirements map[string]core.Module) (server_key_storage.ServerKeyStorage, error) {
	serverKeyStorage, ok := requirements["server-key-storage"].(server_key_storage.ServerKeyStorage)
	if !ok {
		s.Log("All requirements list")
		for k, v := range requirements {
			s.Log("Requirement", k, v.Signature())
		}
		return nil, fmt.Errorf("requiremnt list has wrong server-key-storage requirement")
	}

	return serverKeyStorage, nil
}
