package types

type SerializerType int

const (
	ByteArray SerializerType = iota
	String
	Avro
)

type FetchOrder int

const (
	Oldest FetchOrder = iota
	Newest
)

type TopicConfig struct {
	Server          string         `json:"server" yaml:"server"`
	Topic           string         `json:"topic" yaml:"topic"`
	KeySerializer   SerializerType `json:"keySerializer" yaml:"keySerializer"`
	ValueSerializer SerializerType `json:"valueSerializer" yaml:"valueSerializer"`
}

type Topics []TopicConfig

type FetchRequest struct {
	Server          string         `json:"server" yaml:"server"`
	Topic           string         `json:"topic" yaml:"topic"`
	KeySerializer   SerializerType `json:"keySerializer" yaml:"keySerializer"`
	ValueSerializer SerializerType `json:"valueSerializer" yaml:"valueSerializer"`
	MessageOrder    FetchOrder     `json:"messageOrder" yaml:"messageOrder"`
	Size            int
}
