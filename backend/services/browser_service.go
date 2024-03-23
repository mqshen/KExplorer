package services

import (
	"context"
	"github.com/segmentio/kafka-go"
	"kafkaexplorer/backend/types"
	"sync"
)

type connectionItem struct {
	client     *kafka.Conn
	ctx        context.Context
	cancelFunc context.CancelFunc
	cursor     map[int]uint64 // current cursor of databases
	stepSize   int64
	db         int // current database index
}

type browserService struct {
	ctx        context.Context
	connMap    map[string]*connectionItem
	cmdHistory []cmdHistoryItem
	mutex      sync.Mutex
}

var browser *browserService
var onceBrowser sync.Once

func Browser() *browserService {
	if browser == nil {
		onceBrowser.Do(func() {
			browser = &browserService{
				connMap: map[string]*connectionItem{},
			}
		})
	}
	return browser
}

func (b *browserService) Start(ctx context.Context) {
	b.ctx = ctx
}

func (b *browserService) Stop() {
	for _, item := range b.connMap {
		if item.client != nil {
			if item.cancelFunc != nil {
				item.cancelFunc()
			}
			item.client.Close()
		}
	}
	b.connMap = map[string]*connectionItem{}
}

// OpenConnection open redis server connection
func (b *browserService) OpenConnection(name string) (resp types.JSResp) {
	// get connection config
	selConn := Connection().getConnection(name)

	resp.Success = true
	resp.Data = map[string]any{
		"name": selConn.Name,
	}
	return
}

// CloseConnection close redis server connection
func (b *browserService) CloseConnection(name string) (resp types.JSResp) {
	item, ok := b.connMap[name]
	if ok {
		delete(b.connMap, name)
		if item.client != nil {
			if item.cancelFunc != nil {
				item.cancelFunc()
			}
			item.client.Close()
		}
	}
	resp.Success = true
	return
}
