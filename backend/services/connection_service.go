package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-zookeeper/zk"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	. "kafkaexplorer/backend/storage"
	"kafkaexplorer/backend/types"
	"strings"
	"sync"
	"time"
)

type cmdHistoryItem struct {
	Timestamp int64  `json:"timestamp"`
	Server    string `json:"server"`
	Cmd       string `json:"cmd"`
	Cost      int64  `json:"cost"`
}

type connectionService struct {
	ctx          context.Context
	conns        *ConnectionsStorage
	topicService *topicService
}

var connection *connectionService
var onceConnection sync.Once

func Connection() *connectionService {
	if connection == nil {
		onceConnection.Do(func() {
			connection = &connectionService{
				conns:        NewConnections(),
				topicService: Topic(),
			}
		})
	}
	return connection
}

func (c *connectionService) TestConnection(config types.ConnectionConfig) (resp types.JSResp) {
	conn, err := c.createZkClient(config)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	defer conn.Close()

	//if _, err = conn.ReadPartitions(); err != nil {
	//	resp.Msg = err.Error()
	//} else {
	resp.Success = true
	//}
	return
}

func (c *connectionService) createKafkaClient(config types.ConnectionConfig) (*kafka.AdminClient, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers":        config.Bootstrap,
		"allow.auto.create.topics": "false"}
	adminClient, err := kafka.NewAdminClient(&conf)
	if err != nil {
		runtime.LogErrorf(c.ctx, "Failed to create AdminClient: %s\n", err)
		return nil, err
	}

	return adminClient, nil
}

func (c *connectionService) createZkClient(config types.ConnectionConfig) (*zk.Conn, error) {
	server := fmt.Sprintf("%s:%d", config.Addr, config.Port)
	conn, _, err := zk.Connect([]string{server}, time.Second*20)
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

// GetConnection get connection profile by name
func (c *connectionService) GetConnection(name string) (resp types.JSResp) {
	conn := c.getConnection(name)
	resp.Success = conn != nil
	resp.Data = conn
	return
}

// DeleteConnection remove connection by name
func (c *connectionService) DeleteConnection(name string) (resp types.JSResp) {
	err := c.conns.DeleteConnection(name)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	c.topicService.DeleteByServer(name)
	resp.Success = true
	return
}
