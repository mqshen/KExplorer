package types

type ConnectionConfig struct {
	Name      string `json:"name" yaml:"name"`
	Group     string `json:"group,omitempty" yaml:"-"`
	Addr      string `json:"addr,omitempty" yaml:"addr,omitempty"`
	Port      int    `json:"port,omitempty" yaml:"port,omitempty"`
	Bootstrap string `json:"bootstrap,omitempty" yaml:"bootstrap,omitempty"`
	Root      string `json:"root,omitempty" yaml:"root,omitempty"`
}

type Connection struct {
	ConnectionConfig `json:",inline" yaml:",inline"`
	Type             string       `json:"type,omitempty" yaml:"type,omitempty"`
	Connections      []Connection `json:"connections,omitempty" yaml:"connections,omitempty"`
}

type Connections []Connection

type Broker struct {
	Id      int64  `json:"id" yaml:"id"`
	Host    string `json:"host" yaml:"host"`
	Port    int64  `json:"port" yaml:"port"`
	Version int64  `json:"version" yaml:"version"`
}
