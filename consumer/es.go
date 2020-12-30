package consumer

import (
	"context"
	"douban/model"
	"douban/wait"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/olivere/elastic/v7"
)

func SaveToES() {

	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		panic(err)
	}
	partitionList, err := consumer.Partitions("douban_book")
	if err != nil {
		panic(err)
	}

	//fmt.Println(partitionList)
	for partition := range partitionList {
		consumePartition, err := consumer.ConsumePartition("douban_book", int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}

		client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://localhost:9200"))
		if err != nil {
			panic(err)
		}

		defer consumePartition.AsyncClose()
		go func(sarama.PartitionConsumer) {
			for msg := range consumePartition.Messages() {
				book := new(model.BookDetail)
				err = json.Unmarshal([]byte(msg.Value), &book)
				if err != nil {
					panic(err)
				}
				_, err := client.Index().
					Index("douban").
					BodyJson(book).
					Do(context.Background())
				if err != nil {
					// Handle error
					panic(err)
				}
				fmt.Println("save2es")
			}

		}(consumePartition)
	}
	select {

	}
	wait.Wg.Done()
}
