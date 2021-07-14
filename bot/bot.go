package main

import (
	"encoding/json"
	"fmt"
	"log"
	"nsq-chat/config"
	"nsq-chat/models"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

func subscribeToNsq(channelName string) error {

	cfg := nsq.NewConfig()
	cfg.Set("LookupdPollInterval", config.LookupdPollInterval*time.Second)
	cfg.Set("MaxInFlight", config.MaxInFlight)
	cfg.UserAgent = fmt.Sprintf("Bot client go-nsq/%s", nsq.VERSION)

	nsqConsumer, err := nsq.NewConsumer(config.TopicName, channelName, cfg)

	if err != nil {
		log.Println("subscribeToNsq error: ", err)
		return err
	}

	nsqConsumer.AddHandler(nsq.HandlerFunc(handleMessage))
	nsqErr := nsqConsumer.ConnectToNSQLookupd(config.AddrNsqlookupd)

	if nsqErr != nil {
		log.Println("NSQ connection error: ", nsqErr)
		return err
	}

	log.Printf("Subscribe to NSQ success to channel %s", channelName)
	return nil
}

// handleMessage pushes messages from NSQ to Mongodb
func handleMessage(msg *nsq.Message) error {
	message := models.Message{}
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Println("NSQ HandleMessage ERROR: invalid JSON subscribe data")
		return err
	}

	// Simple reply for now
	// TODO: Some NLP here?
	if strings.Index(message.Body, "@bot") != -1 {
		replyMessage := &models.Message{
			Name:      config.BotChannelName,
			Channel:   message.Channel,
			User:      config.BotChannelName,
			TimeStamp: time.Now(),
			Body:      "Hi human!!!",
		}
		msgJSON, _ := json.Marshal(replyMessage)
		err = models.SendMessageToTopic(config.TopicName, []byte(string(msgJSON)))
	}
	return err
}

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go subscribeToNsq(config.BotChannelName)

	// Wait here for SigInt or SigTerm
	<-sigs
}
