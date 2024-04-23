package explorer

import (
	"context"
	"kafkaexplorer/backend/types"
	"time"
)

type ClustersStatisticsScheduler struct {
	ctx            context.Context
	clusters       []*types.Cluster
	statisticCache *StatisticsCacheService
}

func (c *ClustersStatisticsScheduler) Start(ctx context.Context) {
	c.ctx = ctx
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return

		case _, ok := <-ticker.C:
			if !ok {
				return
			}
			for _, cluster := range c.clusters {
				c.statisticCache.UpdateCache(ctx, cluster)
			}
		}
	}
}
