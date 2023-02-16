package mq

import (
	"fmt"
	"liujun/Time_Micro_GateWay/server/common"
	"time"
)

var (
	MQ *Session
)

func init() {
	config := common.Config
	rabbitmq_addr := config.String("rabbitmq_addr")
	rabbitmq_port, _ := config.Int("rabbitmq_port")
	rabbitmq_user := config.String("rabbitmq_user")
	rabbitmq_pwd := config.String("rabbitmq_pwd")
	rabbitmq_vhost := config.String("rabbitmq_vhost")
	rabbitmq_queue := config.String("rabbitmq_queue")
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", rabbitmq_user,
		rabbitmq_pwd, rabbitmq_addr, rabbitmq_port, rabbitmq_vhost)
	MQ = New(rabbitmq_queue, addr)
}

func Push(message string) error {
	body := []byte(message)
	time.Sleep(time.Second * 3)
	if err := MQ.Push(body); err != nil {
		fmt.Printf("Push failed: %s\n", err)
		return err
	} else {
		fmt.Println("Push succeeded!")
		return nil
	}
}
