package types

type ConnectionConfig struct {
	Name  string `json:"name" yaml:"name"`
	Group string `json:"group,omitempty" yaml:"-"`
	Addr  string `json:"addr,omitempty" yaml:"addr,omitempty"`
	Port  int    `json:"port,omitempty" yaml:"port,omitempty"`
}

type Connection struct {
	ConnectionConfig `json:",inline" yaml:",inline"`
	Type             string       `json:"type,omitempty" yaml:"type,omitempty"`
	Connections      []Connection `json:"connections,omitempty" yaml:"connections,omitempty"`
}

type Connections []Connection
