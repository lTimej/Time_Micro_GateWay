package rabbitmq

import (
	"encoding/json"
	"fmt"
	"liujun/Time_Micro_GateWay/mail/common"
	"liujun/Time_Micro_GateWay/mail/utils"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var (
	MQ             *Session
	rabbitmq_addr  string
	rabbitmq_port  int
	rabbitmq_user  string
	rabbitmq_pwd   string
	rabbitmq_vhost string
	rabbitmq_queue string
)

func init() {
	config := common.Config
	rabbitmq_addr = config.String("rabbitmq_addr")
	rabbitmq_port, _ = config.Int("rabbitmq_port")
	rabbitmq_user = config.String("rabbitmq_user")
	rabbitmq_pwd = config.String("rabbitmq_pwd")
	rabbitmq_vhost = config.String("rabbitmq_vhost")
	rabbitmq_queue = config.String("rabbitmq_queue")
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", rabbitmq_user,
		rabbitmq_pwd, rabbitmq_addr, rabbitmq_port, rabbitmq_vhost)
	MQ = New(rabbitmq_queue, addr)
}

func SendMail() error {
	time.Sleep(time.Second * 1)
	msgs, err := MQ.Stream()
	if err != nil {
		fmt.Printf("获取消息失败: %s\n", err)
		return err
	}
	var forever chan struct{}
	go func() {
		for d := range msgs {
			f := successHandler(d)
			if f != nil {
				f(d.Body)
			}
		}
	}()
	<-forever
	return nil
}
func successHandler(consumer amqp.Delivery) func(data []byte) bool {
	switch consumer.RoutingKey {
	case rabbitmq_queue:
		return user_reg_handler()
	}
	return nil
}

func user_reg_handler() func(data []byte) bool {
	return func(data []byte) bool { //如果返回true 则无需重试
		log.Printf("user_reg_handler data:%s\n", string(data))

		//解析json数据
		var user struct {
			Id       int
			Username string
			Phone    string
			Email    string
		}
		err := json.Unmarshal(data, &user)
		if err != nil {
			log.Println("user_reg_handler Unmarshal Err:", err)
			return false
		}

		//推送邮件
		Body := fmt.Sprintf("Hello. The name is %v and the phone is %v register success!!", user.Username, user.Phone)
		err = utils.SendMail(user.Email, Body)
		if err != nil {
			log.Println("user_reg_handler Err:", err)
			return false
		}

		// retryClient.Ack()
		return true
	}
}
