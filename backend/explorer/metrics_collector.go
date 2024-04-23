package explorer

import (
	"github.com/twmb/franz-go/pkg/kadm"
	"kafkaexplorer/backend/types"
)

type SimpleMetric struct {
}

type MetricsConfig struct {
}

type MetricsCollector struct {
	adminClientService *AdminClientService
}

func (c MetricsCollector) GetBrokerMetrics(cluster *types.Cluster, nodes kadm.BrokerDetails) {
	for _, node := range nodes {
		c.getMetrics(cluster, node)
	}

}

func (c MetricsCollector) getMetrics(cluster *types.Cluster, node kadm.BrokerDetail) []*SimpleMetric {
	metricsConfig := c.adminClientService.Get(cluster).MetricsConfig
	if metricsConfig == nil {
		return make([]*SimpleMetric, 0)
	}
	return make([]*SimpleMetric, 0)

}
