package main

import (
	// "strings"
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"github.com/Andrew-M-C/tencent-tars-demo/GoLogger/log"
	"github.com/Shopify/sarama"
)

var kafkaTopic = "hello-tars-kafka"
var kafkaProducer sarama.SyncProducer

type kafkaMsg struct {
	IP			string 	`json:"IP"`
	Port		int		`json:"port"`
	Url			string	`json:"URL"`
	Time		string	`json:"datetime"`
}

type sendmsgRet struct {
	Msg			string	`json:"msg"`
	Time		string	`json:"datetime"`
}

func getKafkaSyncProducer() (sarama.SyncProducer, error) {
	addresses := []string{"127.0.0.1:9092"}
	kafka_config := sarama.NewConfig()
	kafka_config.Producer.Return.Successes = true
	// kafka_config.Producer.Partitioner = sarama.NewRandomPartitioner
	// kafka_config.Producer.RequiredAcks = sarama.WaitForAll
	kafka_config.Producer.Timeout = 2 * time.Second
	kafka_config.Version = sarama.V2_1_0_0

	return sarama.NewSyncProducer(addresses, kafka_config)
}

func init() {
	var err error
	kafkaProducer, err = getKafkaSyncProducer()
	if err != nil {
		log.Error("Failed to create producer: " + err.Error())
	}
}

func HttpTarsHandler(w http.ResponseWriter, info *HttpRequestInfo, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte("{\"msg\": \"Hello, TarsGo\"}"))
	return
}

func HttpTarsSendMsgHandler(w http.ResponseWriter, info *HttpRequestInfo, r *http.Request) {
	// analyze request
	producer := kafkaProducer
	var kafka_msg *sarama.ProducerMessage
	var err error
	ok := true
	ret := &sendmsgRet{Msg: "Hello, TarsGo"}
	log.Debugf("Request: %s:%d, %s", info.Ip, info.Port, info.Url)
	msg_struct := &kafkaMsg{
		IP: info.Ip,
		Port: info.Port,
		Url: info.Url,
	}
	utc_time := time.Now()
	local_time := utc_time.Local()
	msg_struct.Time = local_time.Format("2006/01/02 15:04:05")
	ret.Time = msg_struct.Time
	msg, err := json.Marshal(msg_struct)
	if err != nil {
		ret.Msg = err.Error()
		ok = false
	}
	// get kafka producer
	if ok {
		kafka_msg = &sarama.ProducerMessage{
			Topic: kafkaTopic,
			Value: sarama.ByteEncoder(msg),
		}
		if nil == producer {
			ret.Msg = "producer not initialized"
			ok = false
		}
	}
	// send data
	if ok {
		partition, offset, err := producer.SendMessage(kafka_msg)
		if err != nil {
			ret.Msg = err.Error()
		} else {
			ret.Msg = fmt.Sprintf("message sent, partition = %d, offset = %d", partition, offset)
		}
	}
	// return
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	ret_str, _ := json.Marshal(ret)
	w.Write([]byte(ret_str))
	return
}
