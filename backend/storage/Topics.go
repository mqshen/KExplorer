package storage

import (
	"errors"
	"gopkg.in/yaml.v3"
	"kafkaexplorer/backend/types"
	"sync"
)

type TopicStorage struct {
	storage *localStorage
	mutex   sync.Mutex
}

func NewTopics() *TopicStorage {
	return &TopicStorage{
		storage: NewLocalStore("topics.yaml"),
	}
}

func (c *TopicStorage) UpsetTopic(server string, topic string, param types.TopicConfig) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	topics := c.getTopics()
	var updated bool
	var retrieve func(types.Topics, string, string, types.TopicConfig) error
	retrieve = func(topics types.Topics, server string, topicName string, param types.TopicConfig) error {
		for i, topic := range topics {
			if server != param.Server && topic.Server == param.Server {
				return errors.New("duplicated topic name")
			} else if topic.Server == server && topic.Topic == topicName && !updated {
				topics[i] = types.TopicConfig{
					Server:          server,
					Topic:           topicName,
					KeySerializer:   param.KeySerializer,
					ValueSerializer: param.ValueSerializer,
				}
				updated = true
			}
		}
		return nil
	}

	err := retrieve(topics, server, topic, param)
	if err != nil {
		return err
	}
	if !updated {
		topics = append(topics, types.TopicConfig{
			Server:          server,
			Topic:           topic,
			KeySerializer:   param.KeySerializer,
			ValueSerializer: param.ValueSerializer,
		})
	}

	return c.saveTopics(topics)
}

// GetConnections get all store connections from local
func (c *TopicStorage) GetTopics() (ret types.Topics) {
	return c.getTopics()
}

func (c *TopicStorage) saveTopics(topics types.Topics) error {
	b, err := yaml.Marshal(&topics)
	if err != nil {
		return err
	}
	if err = c.storage.Store(b); err != nil {
		return err
	}
	return nil
}

func (c *TopicStorage) getTopics() (ret types.Topics) {
	b, err := c.storage.Load()
	ret = c.defaultTopics()
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(b, &ret); err != nil {
		ret = c.defaultTopics()
		return
	}
	if len(ret) <= 0 {
		ret = c.defaultTopics()
	}
	return
}

func (c *TopicStorage) defaultTopics() types.Topics {
	return types.Topics{}
}

// GetConnection get connection by name
func (c *TopicStorage) GetTopic(server, topic string) *types.TopicConfig {
	topics := c.GetTopics()

	var findConn func(string, string, types.Topics) *types.TopicConfig
	findConn = func(server, topic string, conns types.Topics) *types.TopicConfig {
		for i, conn := range conns {
			if conn.Server == server && conn.Topic == topic {
				return &conns[i]
			}
		}
		return nil
	}

	return findConn(server, topic, topics)
}

// DeleteConnection remove special connection
func (c *TopicStorage) DeleteByServer(name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	topics := c.getTopics()
	var updated bool
	for i, topic := range topics {
		if topic.Server == name {
			topics = append(topics[:i], topics[i+1:]...)
			updated = true
			break
		}
		if updated {
			break
		}
	}
	if !updated {
		return errors.New("no match topic server")
	}
	return c.saveTopics(topics)
}
