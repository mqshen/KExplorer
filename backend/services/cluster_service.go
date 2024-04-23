package services

import (
	"kafkaexplorer/backend/explorer"
	"kafkaexplorer/backend/storage"
	"kafkaexplorer/backend/types"
	"sync"
)

type ClusterService struct {
	clusterStorage  *storage.ClusterStorage
	statisticsCache *explorer.StatisticsCacheService
}

var clusterServices *ClusterService
var onceCluster sync.Once

func ClusterServiceInstance() *ClusterService {
	if clusterServices == nil {
		onceCluster.Do(func() {
			clusterServices = &ClusterService{
				clusterStorage: storage.NewClusterStorage(),
			}
		})
	}
	return clusterServices
}

func (c *ClusterService) ClusterStats(name string) (resp types.JSResp) {
	cluster := c.clusterStorage.GetCluster(name)
	if cluster == nil {
		resp.Msg = "cluster not found"
		return
	}
	statistic := c.statisticsCache.Get(cluster)
	resp.Data = statistic
	resp.Success = true
	return
}

func (c *ClusterService) SaveCluster(name string, param types.ClusterConfig) (resp types.JSResp) {
	var err error
	if len(name) > 0 {
		err = c.clusterStorage.UpdateCluster(name, param)
	} else {
		err = c.clusterStorage.CreateCluster(param)
	}
	if err != nil {
		resp.Msg = err.Error()
	} else {
		resp.Success = true
	}
	return
}

func (c *ClusterService) ListCluster() (resp types.JSResp) {
	resp.Data = c.clusterStorage.GetClusters()
	resp.Success = true
	return
}

func (c *ClusterService) GetCluster(name string) *types.Cluster {
	return c.clusterStorage.GetCluster(name)
}
