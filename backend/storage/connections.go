package storage

import (
	"errors"
	"gopkg.in/yaml.v3"
	"kafkaexplorer/backend/types"
	"sync"
)

type ConnectionsStorage struct {
	storage *localStorage
	mutex   sync.Mutex
}

func NewConnections() *ConnectionsStorage {
	return &ConnectionsStorage{
		storage: NewLocalStore("connections.yaml"),
	}
}

func (c *ConnectionsStorage) saveConnections(conns types.Connections) error {
	b, err := yaml.Marshal(&conns)
	if err != nil {
		return err
	}
	if err = c.storage.Store(b); err != nil {
		return err
	}
	return nil
}

// CreateConnection create new connection
func (c *ConnectionsStorage) CreateConnection(param types.ConnectionConfig) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conn := c.GetConnection(param.Name)
	if conn != nil {
		return errors.New("duplicated connection name")
	}
	conns := c.getConnections()
	var group *types.Connection
	if len(param.Group) > 0 {
		for i, conn := range conns {
			if conn.Type == "group" && conn.Name == param.Group {
				group = &conns[i]
				break
			}
		}
	}
	if group != nil {
		group.Connections = append(group.Connections, types.Connection{
			ConnectionConfig: param,
		})
	} else {
		if len(param.Group) > 0 {
			// no group matched, create new group
			conns = append(conns, types.Connection{
				Type: "group",
				Connections: types.Connections{
					types.Connection{
						ConnectionConfig: param,
					},
				},
			})
		} else {
			conns = append(conns, types.Connection{
				ConnectionConfig: param,
			})
		}
	}

	return c.saveConnections(conns)
}

// UpdateConnection update existing connection by name
func (c *ConnectionsStorage) UpdateConnection(name string, param types.ConnectionConfig) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.getConnections()
	var updated bool
	var retrieve func(types.Connections, string, types.ConnectionConfig) error
	retrieve = func(conns types.Connections, name string, param types.ConnectionConfig) error {
		for i, conn := range conns {
			if conn.Type != "group" {
				if name != param.Name && conn.Name == param.Name {
					return errors.New("duplicated connection name")
				} else if conn.Name == name && !updated {
					conns[i] = types.Connection{
						ConnectionConfig: param,
					}
					updated = true
				}
			} else {
				if err := retrieve(conn.Connections, name, param); err != nil {
					return err
				}
			}
		}
		return nil
	}

	err := retrieve(conns, name, param)
	if err != nil {
		return err
	}
	if !updated {
		return errors.New("connection not found")
	}

	return c.saveConnections(conns)
}

// GetConnection get connection by name
func (c *ConnectionsStorage) GetConnection(name string) *types.Connection {
	conns := c.getConnections()

	var findConn func(string, string, types.Connections) *types.Connection
	findConn = func(name, groupName string, conns types.Connections) *types.Connection {
		for i, conn := range conns {
			if conn.Type != "group" {
				if conn.Name == name {
					conns[i].Group = groupName
					return &conns[i]
				}
			} else {
				if ret := findConn(name, conn.Name, conn.Connections); ret != nil {
					return ret
				}
			}
		}
		return nil
	}

	return findConn(name, "", conns)
}

func (c *ConnectionsStorage) getConnections() (ret types.Connections) {
	b, err := c.storage.Load()
	ret = c.defaultConnections()
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(b, &ret); err != nil {
		ret = c.defaultConnections()
		return
	}
	if len(ret) <= 0 {
		ret = c.defaultConnections()
	}
	//if !sliceutil.AnyMatch(ret, func(i int) bool {
	//	return ret[i].GroupName == ""
	//}) {
	//	ret = append(ret, c.defaultConnections()...)
	//}
	return
}

func (c *ConnectionsStorage) defaultConnections() types.Connections {
	return types.Connections{}
}

// GetConnections get all store connections from local
func (c *ConnectionsStorage) GetConnections() (ret types.Connections) {
	return c.getConnections()
}
