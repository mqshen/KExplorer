package explorer

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"kafkaexplorer/backend/types"
)

type KafkaTopic struct {
	TopicDescription  *kafka.TopicDescription
	Config            *kafka.ConfigEntryResult
	partitionsOffsets *kafka.ListOffsetsResult
	metrics           Metrics
	logDirInfo        KafkaLogDirStats
}

type TopicService struct {
	adminClientService *AdminClientService
	statisticsCache    *StatisticsCacheService
}

func (t *TopicService) LoadTopics(ctx context.Context, c *types.Cluster, topics []string) []*KafkaTopic {
	adminClint := t.adminClientService.Get(c)
	descriptions, err := adminClint.DescribeTopics(ctx, topics)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to describe topics: %s", err.Error())
		return []*KafkaTopic{}
	}
	configs, err := t.adminClientService.Get(c).GetTopicsConfig(ctx, topics)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to describe topic configs: %s", err.Error())
		return []*KafkaTopic{}
	}
	t.statisticsCache.Update(c, descriptions, configs)
	offsets, err := getPartitionOffsets(ctx, descriptions, adminClint)
	metrics := t.statisticsCache.Get(c)
	return createList(
		topics,
		metrics.topicDescriptions,
		metrics.topicConfigs,
		offsets,
		metrics.Metrics,
		metrics.logDirInfo,
	)
}

func createList(topics []string, descriptions map[string]*kafka.TopicDescription, configs map[string]*kafka.ConfigEntryResult,
	offsets *kafka.ListOffsetsResult, metrics Metrics, info KafkaLogDirStats) []*KafkaTopic {
	var list []*KafkaTopic

	for _, topic := range topics {
		if description, ok := descriptions[topic]; ok {
			config, ok := configs[topic]
			if !ok {
				config = &kafka.ConfigEntryResult{}
			}

			list = append(list, &KafkaTopic{
				TopicDescription:  description,
				Config:            config,
				partitionsOffsets: offsets,
				metrics:           metrics,
				logDirInfo:        info,
			})

		}
	}
	return list
}

func getPartitionOffsets(ctx context.Context, descriptions []kafka.TopicDescription, ac *KafkaAdminClient) (*kafka.ListOffsetsResult, error) {
	result, err := ac.ListOffsets(ctx, descriptions, kafka.EarliestOffsetSpec)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to list offsets: %s", err.Error())
		return nil, err
	}
	return result, nil
}
