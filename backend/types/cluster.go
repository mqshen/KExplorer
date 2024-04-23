package types

import "time"

type ClusterConfig struct {
	Name          string `json:"name" yaml:"name"`
	Bootstrap     string `json:"bootstrap,omitempty" yaml:"bootstrap,omitempty"`
	ClientTimeout time.Duration
}

type Cluster struct {
	ClusterConfig `json:",inline" yaml:",inline"`
}
