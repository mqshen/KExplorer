package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-zookeeper/zk"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"kafkaexplorer/backend/consts"
	"kafkaexplorer/backend/types"
	"math"
	"sort"
	"sync"
)

const KafkaTimeOut = 5000

type connectionItem struct {
	server     string
	client     *zk.Conn
	ctx        context.Context
	cancelFunc context.CancelFunc
	stepSize   int64
	root       string // current database index
	brokers    []types.Broker
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

// OpenDatabase open select database, and list all keys
// @param path contain connection name and db name
func (b *browserService) GetKafkaMetaData(server string) (resp types.JSResp) {

	item, err := b.getZkClient(server)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	brokers := b.getBrokers(item)
	if len(brokers) == 0 {
		resp.Msg = "brokers is empty"
		return
	}

	kafkaClient, err := b.getKafkaClient(server, brokers[0])
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	defer kafkaClient.Close()
	topics := b.getTopics(kafkaClient)
	consumers := b.getConsumers(kafkaClient)

	var brokerServers []string
	for _, broker := range brokers {
		brokerServers = append(brokerServers, fmt.Sprintf("%s:%d", broker.Host, broker.Port))
	}
	resp.Success = true
	resp.Data = map[string]any{
		"brokers":   brokerServers,
		"topics":    topics,
		"consumers": consumers,
	}
	return
}

func (b *browserService) getBrokers(item *connectionItem) []types.Broker {
	path := fmt.Sprintf("%s%s", item.root, "brokers/ids")
	runtime.LogInfo(b.ctx, fmt.Sprintf("path %v", path))
	client := item.client
	ids, _, err := client.Children(path)
	if err != nil {
		return nil
	}
	var brokers []types.Broker
	for _, id := range ids {
		idPath := fmt.Sprintf("%s/%s", path, id)

		runtime.LogInfo(b.ctx, fmt.Sprintf("path %v", idPath))
		jsonData, _, err := client.Get(idPath)
		if err != nil {
			runtime.LogInfo(b.ctx, fmt.Sprintf("path %v", err))
			continue
		}
		var broker types.Broker
		err = json.Unmarshal(jsonData, &broker)
		if err != nil {
			runtime.LogInfo(b.ctx, fmt.Sprintf("path 1 %v", err))
			continue
		}
		brokers = append(brokers, broker)
	}
	item.brokers = brokers
	return brokers
}

func (b *browserService) getTopics(client *kafka.AdminClient) []string {
	metadata, err := client.GetMetadata(nil, true, 5000)
	if err != nil {
		runtime.LogErrorf(b.ctx, fmt.Sprintf("get topic error %v", err))
		return nil
	}
	var topics []string
	for topic := range metadata.Topics {
		topics = append(topics, topic)
	}
	return topics
}

func (b *browserService) getConsumers(conn *kafka.AdminClient) []string {
	result, err := conn.ListConsumerGroups(b.ctx)
	if err != nil {
		return nil
	}

	var consumers []string
	for _, consumer := range result.Valid {
		consumers = append(consumers, consumer.GroupID)
	}

	return consumers
}

// get a redis client from local cache or create a new open
// if db >= 0, will also switch to db index
func (b *browserService) getZkClient(server string) (item *connectionItem, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var ok bool
	var client *zk.Conn
	if item, ok = b.connMap[server]; ok {

		// close previous connection if database is not the same
		if item.cancelFunc != nil {
			item.cancelFunc()
		}
		item.client.Close()
		delete(b.connMap, server)
	}

	// recreate new connection after switch database
	selConn := Connection().getConnection(server)
	if selConn == nil {
		err = fmt.Errorf("no match connection \"%s\"", server)
		return
	}
	var connConfig = selConn.ConnectionConfig
	client, err = b.createZkClient(connConfig)
	if err != nil {
		err = fmt.Errorf("can not connect to server \"%s\", %s", server, err.Error())
		return
	}
	ctx, cancelFunc := context.WithCancel(b.ctx)
	item = &connectionItem{
		server:     server,
		client:     client,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		root:       "/", //connConfig.Root,
	}
	if item.stepSize <= 0 {
		item.stepSize = consts.DEFAULT_LOAD_SIZE
	}
	b.connMap[server] = item
	return
}

// get a redis client from local cache or create a new open
// if db >= 0, will also switch to db index
func (b *browserService) getKafkaClient(server string, broker types.Broker) (client *kafka.AdminClient, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// recreate new connection after switch database
	selConn := Connection().getConnection(server)
	if selConn == nil {
		err = fmt.Errorf("no match connection \"%s\"", server)
		return
	}
	var connConfig = selConn.ConnectionConfig

	connConfig.Bootstrap = fmt.Sprintf("%s:%d", broker.Host, broker.Port)
	runtime.LogInfo(b.ctx, fmt.Sprintf("connect to kafka %s", connConfig.Bootstrap))
	client, err = b.createKafkaClient(connConfig)
	if err != nil {
		return
	}
	return
}

func (b *browserService) createZkClient(selConn types.ConnectionConfig) (client *zk.Conn, err error) {
	client, err = Connection().createZkClient(selConn)
	if err != nil {
		err = fmt.Errorf("create conenction error: %s", err.Error())
		return
	}
	return
}

func (b *browserService) createKafkaClient(selConn types.ConnectionConfig) (client *kafka.AdminClient, err error) {
	client, err = Connection().createKafkaClient(selConn)
	if err != nil {
		err = fmt.Errorf("create conenction error: %s", err.Error())
		return
	}
	return
}

type KafkaMessage struct {
	Partition int32        `json:"partition" yaml:"partition"`
	Offset    kafka.Offset `json:"offset" yaml:"offset"`
	Key       string       `json:"key" yaml:"key"`
	Value     string       `json:"value" yaml:"value"`
	Timestamp string       `json:"timestamp" yaml:"timestamp"`
}

func (b *browserService) fetchOldestMessages(item *connectionItem, topic string, param types.FetchRequest) []KafkaMessage {
	broker := item.brokers[0]
	bootstrap := fmt.Sprintf("%s:%d", broker.Host, broker.Port)
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrap})
	if err != nil {
		runtime.LogErrorf(b.ctx, fmt.Sprintf("Failed to create Admin client: %s\n", err))
		return nil
	}
	defer adminClient.Close()

	metadata, err := adminClient.GetMetadata(&topic, false, KafkaTimeOut)
	if err != nil {
		runtime.LogErrorf(b.ctx, fmt.Sprintf("get topic metadata error %v", err))
		return nil
	}
	totalSize := param.Size * len(metadata.Topics)

	return b.fetchMessages(item, topic, nil, false, totalSize, param)
}

func (b *browserService) fetchMessages(item *connectionItem, topic string, positions []kafka.TopicPartition, needSeek bool, totalSize int, param types.FetchRequest) []KafkaMessage {

	broker := item.brokers[0]
	bootstrap := fmt.Sprintf("%s:%d", broker.Host, broker.Port)
	runtime.LogInfo(b.ctx, fmt.Sprintf("start connect to broker %s position: %v", bootstrap, positions))

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        bootstrap,
		"broker.address.family":    "v4",
		"group.id":                 "kafka_explorer_main",
		"session.timeout.ms":       6000,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
	})
	if err != nil {
		runtime.LogInfo(b.ctx, fmt.Sprintf("connect to broker failed %s", err))
		return nil
	}
	defer c.Close()
	if needSeek {
		_, err = c.SeekPartitions(positions)
		if err != nil {
			runtime.LogInfo(b.ctx, fmt.Sprintf("seek partition failed %s", err))
			return nil
		}
	}

	runtime.LogInfo(b.ctx, fmt.Sprintf("subscribe to topic %s", topic))
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		runtime.LogInfo(b.ctx, fmt.Sprintf("subscribe to topic failed %s", err))
		return nil
	}

	var messages []KafkaMessage
	var times = 0

	var valueDeserFunc func(value []byte) string
	switch param.ValueSerializer {

	case types.String:
		valueDeserFunc = func(value []byte) string {
			return string(value)
		}
	default:
		valueDeserFunc = func(value []byte) string {
			return base64.StdEncoding.EncodeToString(value)
		}
	}

	for {
		if len(messages) >= totalSize || times > 5 {
			break
		}
		ev := c.Poll(30000)
		if ev == nil {
			times += 1
			continue
		}
		switch e := ev.(type) {
		case *kafka.Message:
			messages = append(messages, KafkaMessage{
				Partition: e.TopicPartition.Partition,
				Offset:    e.TopicPartition.Offset,
				Key:       valueDeserFunc(e.Key),
				Value:     valueDeserFunc(e.Value),
				Timestamp: e.Timestamp.Format("2006-01-02 15:04:05"),
			})

		case kafka.Error:
			times += 1
		default:
			runtime.LogInfo(b.ctx, fmt.Sprintf("Ignored %v\n", e))
			times += 1
		}
	}
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp > messages[j].Timestamp
	})
	return messages
}

func (b *browserService) fetchNewestMessages(item *connectionItem, topic string, param types.FetchRequest) []KafkaMessage {
	broker := item.brokers[0]
	bootstrap := fmt.Sprintf("%s:%d", broker.Host, broker.Port)
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrap})
	if err != nil {
		runtime.LogErrorf(b.ctx, fmt.Sprintf("Failed to create Admin client: %s\n", err))
		return nil
	}
	defer adminClient.Close()

	metadata, err := adminClient.GetMetadata(&topic, false, KafkaTimeOut)
	if err != nil {
		runtime.LogErrorf(b.ctx, fmt.Sprintf("get topic metadata error %v", err))
		return nil
	}
	if t, ok := metadata.Topics[topic]; ok {
		topicPartitionOffsets := make(map[kafka.TopicPartition]kafka.OffsetSpec)
		for partition := range t.Partitions {
			tp := kafka.TopicPartition{Topic: &topic, Partition: int32(partition)}
			topicPartitionOffsets[tp] = kafka.LatestOffsetSpec
		}
		results, err := adminClient.ListOffsets(b.ctx, topicPartitionOffsets,
			kafka.SetAdminIsolationLevel(kafka.IsolationLevelReadCommitted))
		if err != nil {
			runtime.LogErrorf(b.ctx, fmt.Sprintf("get topic metadata error %v", err))
			return nil
		}

		var topics []kafka.TopicPartition
		for tp, v := range results.ResultInfos {
			tp.Offset = kafka.Offset(math.Max(float64(int64(v.Offset)-int64(param.Size)), float64(0)))
			topics = append(topics, tp)
		}

		totalSize := param.Size * len(topics)
		return b.fetchMessages(item, topic, topics, true, totalSize, param)
	}
	return nil
}

func (b *browserService) FetchMessages(server string, topic string, param types.FetchRequest) (resp types.JSResp) {
	if len(server) == 0 || len(topic) == 0 {
		resp.Msg = "server or topic can not be empty"
		return
	}

	item, ok := b.connMap[server]
	if ok {
		var messages []KafkaMessage
		if param.MessageOrder == types.Oldest {
			messages = b.fetchOldestMessages(item, topic, param)
		} else {
			messages = b.fetchNewestMessages(item, topic, param)
		}

		resp.Success = true
		resp.Data = map[string]any{
			"messages": messages,
		}
	} else {
		resp.Msg = fmt.Sprintf("can not find connection for server %s topic %s", server, topic)
	}
	return
}
