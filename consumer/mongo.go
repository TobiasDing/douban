package consumer

import (
	"context"
	"douban/model"
	"douban/wait"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveToMongo() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	defer client.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}
	collection := client.Database("douban").Collection("books")

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
		defer consumePartition.AsyncClose()
		go func(sarama.PartitionConsumer) {
			for msg := range consumePartition.Messages() {
				fmt.Println(msg)
				book := new(model.BookDetail)
				err = json.Unmarshal([]byte(msg.Value), &book)
				if err != nil {
					panic(err)
				}
				fmt.Println(123)
				insertOneResult, err := collection.InsertOne(context.Background(), book)
				if err != nil {
					panic(err)
				}
				fmt.Println("save2mongo", insertOneResult)
			}
		}(consumePartition)
	}
	select {}
	wait.Wg.Done()
}
