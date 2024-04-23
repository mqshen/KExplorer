package explorer

import (
	"context"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type KafkaAdminClient struct {
	adminClient   *kadm.Client
	version       string
	MetricsConfig *MetricsConfig
	shutdownHook  func()
}

func (k *KafkaAdminClient) DescribeTopics(ctx context.Context, topics []string) (kadm.TopicDetails, error) {
	result, err := k.adminClient.Metadata(ctx, topics...)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to describe topics: %v", err)
		return nil, err
	}
	return result.Topics, nil
}

func (k *KafkaAdminClient) GetTopicsConfig(ctx context.Context, topics []string) (kadm.ResourceConfigs, error) {

	result, err := k.adminClient.DescribeTopicConfigs(ctx, topics...)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to describe topics: %v", err)
		return nil, err
	}
	return result, nil
}

func (k *KafkaAdminClient) ListOffsets(ctx context.Context, topics []string) (kadm.ListedOffsets, error) {

	result, err := k.adminClient.ListStartOffsets(ctx, topics...)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to list offsets: %v", err)
		return nil, err
	}
	return result, nil
}

func (k *KafkaAdminClient) GetVersion() string {
	return k.version
}

func (k *KafkaAdminClient) Metadata(ctx context.Context) (kadm.Metadata, error) {
	return k.adminClient.Metadata(ctx)
}

func (k *KafkaAdminClient) DescribeCluster(ctx context.Context) (kadm.BrokerDetails, error) {
	result, err := k.adminClient.Metadata(ctx)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to describe cluster: %v", err)
		return nil, err
	}

	return result.Brokers, nil
}

func (k *KafkaAdminClient) DescribeLogDirs(ctx context.Context, topics []string) (kadm.DescribedAllLogDirs, error) {
	reqTopicsSet := make(map[string]map[int32]struct{})
	for _, topic := range topics {
		partitions := make(map[int32]struct{})
		reqTopicsSet[topic] = partitions
	}
	result, err := k.adminClient.DescribeAllLogDirs(ctx, reqTopicsSet)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to describe log dirs: %v", err)
		return nil, err
	}

	return result, nil
}
