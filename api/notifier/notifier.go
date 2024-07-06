package main

import (
    "log"
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"time"
	"github.com/slack-go/slack"
)

func main() {
    topic := "quickstart-events"
    
    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
        "group.id":          "go",
        "auto.offset.reset": "smallest"})
    if err != nil {
        log.Println(err)
    }
    
    err = consumer.Subscribe(topic, nil)
    if err != nil {
        log.Println(err)
    }
    for {
        ev := consumer.Poll(1000)
        switch e := ev.(type) {
        case *kafka.Message:
            log.Println(string(e.Value))
			attachment := slack.Attachment{
				Color:         "good",
				Fallback:      "alert from accio",
				AuthorName:    "updater",
				AuthorSubname: "some event occured in cluster",
				AuthorLink:    "https://github.com/slack-go/slack",
				AuthorIcon:    "https://avatars2.githubusercontent.com/u/652790",
				Text:          string(e.Value),
				Footer:        "accio",
				FooterIcon:    "https://platform.slack-edge.com/img/default_application_icon.png",
				Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
			}
			msg := slack.WebhookMessage{
				Attachments: []slack.Attachment{attachment},
			}
		
			err := slack.PostWebhook("https://hooks.slack.com/services/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", &msg)
			if err != nil {
				fmt.Println(err)
			}
       
        case kafka.Error:
            fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
        }
    }

	
}
