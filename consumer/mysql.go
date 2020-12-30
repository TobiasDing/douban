package consumer

import (
	"douban/model"
	"douban/wait"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func SaveToMysql() {

	dsn := "root:19950213@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local"
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		panic(err)
	}
	partitionList, err := consumer.Partitions("douban_book")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.BookDetail{})
	defer db.Close()
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
				db.Create(&book)
				fmt.Println("save2mysql")
			}
		}(consumePartition)
	}
	select {}
	wait.Wg.Done()
}
