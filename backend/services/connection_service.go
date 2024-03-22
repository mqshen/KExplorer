package services

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	. "kafkaexplorer/backend/storage"
	"kafkaexplorer/backend/types"
	"sync"
)

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
	conn, err := kafka.Dial("tcp", fmt.Sprint("%s:%s", config.Addr, config.Port))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *connectionService) Start(ctx context.Context) {
	c.ctx = ctx
}
