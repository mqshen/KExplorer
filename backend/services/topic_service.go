package services

import (
	"context"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	. "kafkaexplorer/backend/storage"
	"kafkaexplorer/backend/types"
	"sync"
)

type topicService struct {
	ctx    context.Context
	topics *TopicStorage
}

var topic *topicService
var onceTopic sync.Once

func Topic() *topicService {
	if topic == nil {
		onceTopic.Do(func() {
			topic = &topicService{
				topics: NewTopics(),
			}
		})
	}
	return topic
}

func (c *topicService) Start(ctx context.Context) {
	c.ctx = ctx
}

func (c *topicService) SaveTopic(server string, topic string, param types.TopicConfig) (resp types.JSResp) {
	var err error
	if len(server) == 0 || len(topic) == 0 {
		err = errors.New("server or topic can not be empty")
	} else {
		runtime.LogInfof(c.ctx, "param %v", param)
		// update topic config
		err = c.topics.UpsetTopic(server, topic, param)
	}
	if err != nil {
		resp.Msg = err.Error()
	} else {
		resp.Success = true
	}
	return
}

// ListConnection list all saved connection in local profile
func (c *topicService) ListTopic() (resp types.JSResp) {
	resp.Success = true
	resp.Data = c.topics.GetTopics()
	return
}

func (c *topicService) GetTopic(server, topic string) *types.TopicConfig {
	return c.topics.GetTopic(server, topic)
}

func (c *topicService) DeleteByServer(name string) {
	err := c.topics.DeleteByServer(name)
	if err != nil {
		runtime.LogErrorf(c.ctx, "delete topic by server failed %s", err)
		return
	}

}
