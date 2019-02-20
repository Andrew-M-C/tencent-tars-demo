/**
 */
package main

import (
	"fmt"
	// "time"
	// "strings"
	"runtime"
	"../GoLogger/log"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type Consumer struct {
	GroupId		int
	GroupIdStr	string
	Object		*cluster.Consumer
}

func initClusterConsumer(consumer *Consumer) (err error) {
	brokers := []string{"127.0.0.1:9092"}
	topics := []string{"hello-tars-kafka"}
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer.Object, err = cluster.NewConsumer(brokers, consumer.GroupIdStr, topics, config)
	return
}

func consumerRoutine(consumer *Consumer) {
	group := fmt.Sprintf("GP-%02d", consumer.GroupId)
	consumer.GroupIdStr = group
	err := initClusterConsumer(consumer)
	if err != nil {
		log.Error("Failed to init kafka consumer: " + err.Error())
		return
	}
	defer consumer.Object.Close()

	go func() {
        for err := range consumer.Object.Errors() {
            log.Infof("%s - Error: %s", group, err.Error())
        }
	}()

	go func() {
        for ntf := range consumer.Object.Notifications() {
            log.Infof("%s - Rebalanced: %+v", group, ntf)
        }
	}()

	// consume messages
	for {
		select {
			case msg, ok := <-consumer.Object.Messages():
				if ok {
					log.Infof("Topic: %s; Partition: %d; Offset: %d; Key: <%s>", msg.Topic, msg.Partition, msg.Offset, msg.Key)
					log.Info("Context: " + string(msg.Value));
					consumer.Object.MarkOffset(msg, "")	// mark message as processed
					log.Debug("== Consume kafka message ends ==")
				} else {
					log.Error("Failed to read message")
				}
		}
	}
	return	// should never goes here
}

func startConsumers(count int) (err error) {
	if count <= 0 {
		count = runtime.NumCPU()
	}

	for i := 0; i < count; i++ {
		consumer := Consumer{GroupId: i}
		go consumerRoutine(&consumer)
	}

	log.Infof("%d consumer(s) running", count)
	err = nil
	return
}
