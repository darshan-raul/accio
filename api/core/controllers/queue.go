package controllers

import (
    "fmt"
    "log"
    "os"
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/darshan-raul/accio/api/utils"
)


func Sendmsg() {
   // creating producer
   kafka_url := utils.GetEnv("KAFKA_URL")
   fmt.Println(kafka_url)
    p, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": "broker:9092",
        "client.id":         "cli",
        "acks":              "all"})

    if err != nil {
        fmt.Printf("Failed to create producer: %s\n", err)
        os.Exit(1)
    }
    // defining the name of the topic
    topic := "quickstart-events"

    delivery_chan := make(chan kafka.Event, 10000)

    
    value := "msg from producer"

    
    fmt.Println("sending message")
    err = p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          []byte(value)},
        delivery_chan,
    )
    if err != nil {
        log.Println(err)
    }

    <-delivery_chan
        

}