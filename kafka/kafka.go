package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

type logData struct {
	topic string
	data string
}
var (
	client sarama.SyncProducer
	logDataChan chan *logData
)

func Init(addrs []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true //Success messages return to the success channel

	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("Connecting to Kafka failed, err:", err)
		return
	}
	logDataChan = make(chan *logData, maxSize)
	go sendToKafka()
	return
}

func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data: data,
	}
	logDataChan <- msg
}
func sendToKafka()  {
	for {
		select {
		case data:= <- logDataChan:
			msg := &sarama.ProducerMessage{}
			msg.Topic = data.topic
			msg.Value = sarama.StringEncoder(data.data)

			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println(msg.Topic, ": Sending message to Kafka failed, err:", err)
				return
		}
			fmt.Printf("%s: Sending message to Kafka succeed, pid:%v, offset:%v\n", msg.Topic, pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
	}

	}

}
