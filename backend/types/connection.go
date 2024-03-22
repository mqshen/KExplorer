package types

type ConnectionConfig struct {
	Name string `json:"name" yaml:"name"`
	Addr string `json:"addr,omitempty" yaml:"addr,omitempty"`
	Port int    `json:"port,omitempty" yaml:"port,omitempty"`
}
