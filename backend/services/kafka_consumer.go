package services

//
//import (
//	"fmt"
//	"strings"
//	"time"
//)
//
//func NewConsumer(kafkaCfg *config.Kafka, kafkaConsumerCfg *config.Consumer, logger *zap.Logger) (Consumer, error) {
//	readerConfig := segmentio.ReaderConfig{
//		Brokers:        strings.Split(kafkaCfg.Servers, ","),
//		GroupID:        kafkaConsumerCfg.ConsumerGroup,
//		GroupTopics:    kafkaConsumerCfg.Topics,
//		MaxBytes:       10e6, // 10MB
//		MaxWait:        2 * time.Second,
//		CommitInterval: time.Second, // flushes commits to Kafka every second
//		StartOffset:    kafkaConsumerCfg.StartOffset.Value(),
//		ErrorLogger: segmentio.LoggerFunc(func(s string, i ...interface{}) {
//			logger.Error(fmt.Sprintf(s, i...))
//		}),
//		Logger: segmentio.LoggerFunc(func(s string, i ...interface{}) {
//			logger.Info(fmt.Sprintf(s, i...))
//		}),
//	}
//
//	if kafkaCfg.IsSecureCluster {
//		tls, err := newTLSConfig(kafkaCfg)
//		if err != nil {
//			return nil, err
//		}
//
//		saslMechanism, err := newMechanism(kafkaCfg)
//		if err != nil {
//			return nil, err
//		}
//
//		readerConfig.Dialer = &segmentio.Dialer{
//			TLS:           tls,
//			SASLMechanism: saslMechanism,
//		}
//
//		readerConfig.GroupBalancers = []segmentio.GroupBalancer{
//			segmentio.RackAffinityGroupBalancer{Rack: kafkaCfg.Rack},
//		}
//	}
//
//	return &kafkaConsumer{
//		consumer: segmentio.NewReader(readerConfig),
//		logger:   logger,
//	}, nil
//}
