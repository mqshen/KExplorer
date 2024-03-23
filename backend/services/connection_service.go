package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	. "kafkaexplorer/backend/storage"
	"kafkaexplorer/backend/types"
	"strings"
	"sync"
)

type cmdHistoryItem struct {
	Timestamp int64  `json:"timestamp"`
	Server    string `json:"server"`
	Cmd       string `json:"cmd"`
	Cost      int64  `json:"cost"`
}

type connectionService struct {
	ctx   context.Context
	conns *ConnectionsStorage
}

var connection *connectionService
var onceConnection sync.Once

func Connection() *connectionService {
	if connection == nil {
		onceConnection.Do(func() {
			connection = &connectionService{
				conns: NewConnections(),
			}
		})
	}
	return connection
}

func (c *connectionService) TestConnection(config types.ConnectionConfig) (resp types.JSResp) {
	conn, err := c.createKafkaClient(config)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	defer func(conn *kafka.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	if _, err = conn.ReadPartitions(); err != nil {
		resp.Msg = err.Error()
	} else {
		resp.Success = true
	}
	return
}

func (c *connectionService) createKafkaClient(config types.ConnectionConfig) (*kafka.Conn, error) {
	conn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", config.Addr, config.Port))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *connectionService) Start(ctx context.Context) {
	c.ctx = ctx
}

// SaveConnection save connection config to local profile
func (c *connectionService) SaveConnection(name string, param types.ConnectionConfig) (resp types.JSResp) {
	var err error
	if strings.ContainsAny(param.Name, "/") {
		err = errors.New("connection name contains illegal characters")
	} else {
		if len(name) > 0 {
			// update connection
			err = c.conns.UpdateConnection(name, param)
		} else {
			err = c.conns.CreateConnection(param)
		}
	}
	if err != nil {
		resp.Msg = err.Error()
	} else {
		resp.Success = true
	}
	return
}

// ListConnection list all saved connection in local profile
func (c *connectionService) ListConnection() (resp types.JSResp) {
	resp.Success = true
	resp.Data = c.conns.GetConnections()
	return
}

func (c *connectionService) getConnection(name string) *types.Connection {
	return c.conns.GetConnection(name)
}
