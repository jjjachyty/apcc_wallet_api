package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

func test(topc string, b []byte) {
	fmt.Println("test", topc)
	fmt.Println(string(b))
}
func main() {

	topicName := "MHC2USDT"
	// msgCount := 2
	// for i := 0; i < msgCount; i++ {
	//time.Sleep(time.Millisecond * 20)
	go readMessage(topicName, 1, test)
	go readMessage("MHC2USDT", 1, test)

	//cleanup := make(chan os.Signal, 1)
	cleanup := make(chan os.Signal)
	signal.Notify(cleanup, os.Interrupt)
	fmt.Println("server is running....")

	quit := make(chan bool)
	go func() {

		select {
		case <-cleanup:
			fmt.Println("Received an interrupt , stoping service ...")
			for _, ele := range consumers {
				ele.StopChan <- 1
				ele.Stop()
			}
			quit <- true
		}
	}()
	<-quit
	fmt.Println("Shutdown server....")
}

type ConsumerHandle struct {
	q       *nsq.Consumer
	msgGood int
	topc    string
	fn      func(string, []byte)
}

var consumers []*nsq.Consumer = make([]*nsq.Consumer, 0)
var mux *sync.Mutex = &sync.Mutex{}

func (h *ConsumerHandle) HandleMessage(message *nsq.Message) error {
	msg := string(message.Body) + "  " + strconv.Itoa(h.msgGood)
	h.fn(h.topc, message.Body)
	fmt.Println(msg)
	return nil
}

func readMessage(topicName string, msgCount int, fn func(string, []byte)) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error: ", err)
		}
	}()

	config := nsq.NewConfig()
	config.MaxInFlight = 1000
	config.MaxBackoffDuration = 500 * time.Second

	//q, _ := nsq.NewConsumer(topicName, "ch" + strconv.Itoa(msgCount), config)
	//q, _ := nsq.NewConsumer(topicName, "ch" + strconv.Itoa(msgCount) + "#ephemeral", config)
	q, _ := nsq.NewConsumer(topicName, "ch"+strconv.Itoa(msgCount), config)

	h := &ConsumerHandle{q: q, msgGood: msgCount, topc: topicName, fn: fn}
	q.AddHandler(h)

	// err := q.ConnectToNSQLookupd("119.3.108.19:4161")
	//err := q.ConnectToNSQDs([]string{"192.168.0.105:4161"})
	err := q.ConnectToNSQD("119.3.108.19:4150")
	//err := q.ConnectToNSQD("192.168.0.105:4415")
	if err != nil {
		fmt.Println("conect nsqd error")
		log.Println(err)
	}
	mux.Lock()
	consumers = append(consumers, q)
	mux.Unlock()
	<-q.StopChan
	fmt.Println("end....")
}
