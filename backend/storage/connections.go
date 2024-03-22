package storage

import "sync"

type ConnectionsStorage struct {
	storage *localStorage
	mutex   sync.Mutex
}

func NewConnections() *ConnectionsStorage {
	return &ConnectionsStorage{
		storage: NewLocalStore("connections.yaml"),
	}
}
