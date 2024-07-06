package main

import (
    "fmt"
    "log"
    "os"
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func sendmsg(p *kafka.Producer, delivery_chan chan kafka.Event, value,topic string){
    fmt.Println("sending message")
    err := p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          []byte(value)},
        delivery_chan,
    )
    if err != nil {
        log.Println(err)
    }

    <-delivery_chan

}

func main() {
   // creating producer
    p, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
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

    
    sendmsg(p,delivery_chan,value,topic)
        

    
}