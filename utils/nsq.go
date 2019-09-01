package utils

import (
	"fmt"
	"time"

	nsq_client "github.com/nsqio/go-nsq"
)

var nsqProducer *nsq_client.Producer

type ConsumerHandle struct {
	q       *nsq_client.Consumer
	handler func(data []byte) error
}

//NsqPublish 发布消息
func NsqPublish(topicName string, data []byte) error {
	return nsqProducer.Publish(topicName, data)
}

//InitNsq 初始化NSQ
func InitNsq() {
	var err error
	nsqCfg := GetNsq()
	config := nsq_client.NewConfig()
	addr := nsqCfg.NsqdServer + ":" + nsqCfg.NsqdPort
	SysLog.Debugf("Nsq 连接地址 %s", addr)
	// 随便给哪个ip发都可以
	nsqProducer, err = nsq_client.NewProducer(addr, config)
	err = nsqProducer.Ping()
	if err != nil {
		SysLog.Panic("连接NSQ失败")
		return
	}

}

func (h *ConsumerHandle) HandleMessage(message *nsq_client.Message) error {
	return h.handler(message.Body)
}

//ReadMessage 获取NSQ消息
func ReadMessage(topicName string, handler func(data []byte) error) {

	// defer func() {
	// 	if err := recover(); err != nil {

	// 		SysLog.Errorf("接收NSQ出错 %v", err)
	// 	}
	// }()

	config := nsq_client.NewConfig()
	nsqCfg := GetNsq()
	config.MaxInFlight = 1000
	config.MaxBackoffDuration = 500 * time.Second
	chanID := nsqCfg.ChainID
	SysLog.Debugf("开始监听NSQ %s %s", topicName, chanID)
	q, err := nsq_client.NewConsumer(topicName, chanID, config)
	if err != nil {
		fmt.Println(err)
	}
	h := &ConsumerHandle{q: q, handler: handler}
	q.AddHandler(h)
	addr := nsqCfg.LookupdServer + ":" + nsqCfg.LookupdPort

	err = q.ConnectToNSQLookupd(addr)
	if err != nil {
		SysLog.Errorf("连接NSQLookupd错误 %v", err)
	}
	<-q.StopChan
	SysLog.Warningln("NSQ 服务端已关闭")
}
