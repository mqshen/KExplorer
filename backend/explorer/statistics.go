package explorer

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"kafkaexplorer/backend/types"
)

type Metrics struct {
}
type KafkaLogDirStats struct {
}
type Statistics struct {
	Version           string
	Metrics           Metrics
	topicDescriptions map[string]*kafka.TopicDescription
	topicConfigs      map[string]*kafka.ConfigEntryResult
	logDirInfo        KafkaLogDirStats
}

type StatisticsCacheService struct {
	cache              map[string]*Statistics
	adminClientService *AdminClientService
	metricsCollector   *MetricsCollector
}

func (c *StatisticsCacheService) Get(cluster *types.Cluster) *Statistics {
	return c.cache[cluster.Name]
}

func (c *StatisticsCacheService) Replace(name string, stats *Statistics) {
	c.cache[name] = stats
}

func (c *StatisticsCacheService) Update(cluster *types.Cluster, descriptions []kafka.TopicDescription, configs []kafka.ConfigResourceResult) {
	stats := c.Get(cluster)
	for _, description := range descriptions {
		stats.topicDescriptions[description.Name] = &description
	}

	for _, config := range configs {
		for k, v := range config.Config {
			stats.topicConfigs[k] = &v
		}
	}

}

func (c *StatisticsCacheService) UpdateCache(ctx context.Context, cluster *types.Cluster) {
	adminClient := c.adminClientService.Get(cluster)
	if adminClient == nil {
		return
	}
	//version := adminClient.GetVersion()
	description, err := adminClient.DescribeCluster(ctx)
	if err != nil {
		runtime.LogErrorf(ctx, "Error describing cluster %s: %s", cluster, err)
		return
	}
	c.metricsCollector.GetBrokerMetrics(cluster, description)
	result, err := c.getLogDirInfo(ctx, description, adminClient)
}

func (c *StatisticsCacheService) getLogDirInfo(ctx context.Context, des kadm.BrokerDetails, client *KafkaAdminClient) (kadm.DescribedAllLogDirs, error) {
	var brokerIds []int32
	for _, node := range des {
		brokerIds = append(brokerIds, node.NodeID)
	}
	result, err := client.DescribeLogDirs(ctx, []string{""})
	if err != nil {
		runtime.LogErrorf(ctx, "Error describing log dirs: %s", err)
		return nil, err
	}
	return result, nil
}
