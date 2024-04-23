package storage

import (
	"errors"
	"gopkg.in/yaml.v3"
	"kafkaexplorer/backend/types"
	"sync"
)

type ClusterStorage struct {
	storage *localStorage
	mutex   sync.Mutex
}

func NewClusterStorage() *ClusterStorage {
	return &ClusterStorage{
		storage: NewLocalStore("cluster.yaml"),
	}
}

func (c *ClusterStorage) GetCluster(name string) *types.Cluster {

	clusters := c.getClusters()

	for i, cluster := range clusters {
		if cluster.Name == name {
			return clusters[i]
		}
	}
	return nil

}

func (c *ClusterStorage) defaultClusters() []*types.Cluster {
	return []*types.Cluster{}
}

func (c *ClusterStorage) getClusters() (ret []*types.Cluster) {
	b, err := c.storage.Load()
	ret = c.defaultClusters()
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(b, &ret); err != nil {
		ret = c.defaultClusters()
		return
	}
	if len(ret) <= 0 {
		ret = c.defaultClusters()
	}
	return
}

func (c *ClusterStorage) CreateCluster(param types.ClusterConfig) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cluster := c.GetCluster(param.Name)
	if cluster != nil {
		return errors.New("duplicated cluster name")
	}
	clusters := c.getClusters()
	clusters = append(clusters, &types.Cluster{
		ClusterConfig: param,
	})

	return c.saveClusters(clusters)
}

func (c *ClusterStorage) saveClusters(cluster []*types.Cluster) error {
	b, err := yaml.Marshal(cluster)
	if err != nil {
		return err
	}
	if err = c.storage.Store(b); err != nil {
		return err
	}
	return nil
}

func (c *ClusterStorage) UpdateCluster(name string, param types.ClusterConfig) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	clusters := c.getClusters()
	var updated bool

	for _, cluster := range clusters {
		if cluster.Name == name {
			updated = true
			cluster.Bootstrap = param.Bootstrap
			cluster.ClientTimeout = param.ClientTimeout
		}
	}

	if !updated {
		return errors.New("connection not found")
	}

	return c.saveClusters(clusters)
}

func (c *ClusterStorage) GetClusters() []*types.Cluster {
	return c.getClusters()
}
