package models

import (
	"encoding/json"
	"fmt"
	"log"
	"nsq-chat/config"
	"os"
	"time"

	"github.com/nsqio/go-nsq"
)

type NsqReader struct {
	channelName string
	consumer    *nsq.Consumer
	rooms       map[*Room]bool
}

func newNsqReader(r *Room, channelName string) error {
	cfg := nsq.NewConfig()
	cfg.Set("LookupdPollInterval", config.LookupdPollInterval*time.Second)
	cfg.Set("MaxFlight", config.MaxInFlight)
	cfg.UserAgent = fmt.Sprintf("Chat client go-nsq/%s", nsq.VERSION)

	nsqConsumer, err := nsq.NewConsumer(config.TopicName, channelName, cfg)
	if err != nil {
		log.Println("could not create newNsqReader:", err)
		return err
	}

	nsqReader := NsqReader{
		channelName: channelName,
		rooms:       map[*Room]bool{r: true},
	}
	r.nsqReaders[channelName] = &nsqReader

	nsqConsumer.AddHandler(&nsqReader)
	err = nsqConsumer.ConnectToNSQLookupd(config.AddrNsqlookupd)
	if err != nil {
		log.Println("could not connect to NSQ:", err)
		return err
	}

	nsqReader.consumer = nsqConsumer
	log.Printf("subscribe to NSQ success to channel %q", channelName)
	return nil
}

func (nr *NsqReader) HandleMessage(msg *nsq.Message) error {
	var msgObj Message
	if err := json.Unmarshal(msg.Body, &msgObj); err != nil {
		log.Println("Handle Message Error:", err)
		return err
	}

	for r := range nr.rooms {
		r.forward <- &msgObj
	}
	return nil
}

func getChannelName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "undefined"
	}
	hostname = "websocket-server-" + hostname
	maxLen := len(hostname)
	if maxLen > 63 {
		maxLen = 63
	}
	return hostname[:maxLen]
}

func subscribeToNsq(r *Room) error {
	nsqChannelName := getChannelName()
	_, ok := r.nsqReaders[nsqChannelName]
	if !ok {
		if err := newNsqReader(r, nsqChannelName); err != nil {
			log.Printf("Failed to subsribe to channel %q\n", nsqChannelName)
			return err
		}
	}
	return nil
}

func SendMessageToTopic(TopicName string, message []byte) error {
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(config.AddrNsqd, cfg)
	if err != nil {
		log.Println("could not create producer:", err)
		return err
	}
	doChan := make(chan *nsq.ProducerTransaction)
	go func() {
		for res := range doChan {
			if res.Error != nil {
				log.Println("send error:", res.Error.Error())
			}
		}
	}()

	return producer.PublishAsync(config.TopicName, message, doChan)
}
