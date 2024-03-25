package services

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-zookeeper/zk"
	"log"
	"time"
)

type OffsetPoller struct {
	client   *zk.Conn
	stepSize int64
	root     string // current database index
}

func (p *OffsetPoller) Run(ctx context.Context) {
	messageChan := p.Start(ctx)

	go func() {

		for {
			select {
			case message, ok := <-messageChan:
				if !ok {
					return
				}
				log.Printf("test %v", message)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (p *OffsetPoller) Start(ctx context.Context) <-chan *kafka.Message {

	var outputChan = make(chan *kafka.Message)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"myTopic", "^aRegex.*[Tt]opic"}, nil)

	go func() {
		defer close(outputChan)
		for {
			select {
			default:
				msg, err := c.ReadMessage(time.Second)
				if err != nil {
					panic(err)
				}
				select {
				case outputChan <- msg:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return outputChan

}
