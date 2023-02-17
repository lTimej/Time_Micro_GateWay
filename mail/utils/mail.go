package utils

import (
	"liujun/Time_Micro_GateWay/mail/common"

	"gopkg.in/gomail.v2"
)

func SendMail(toMail string, body string) error {
	mail_addr := common.Config.String("mail_addr")
	mail_port, _ := common.Config.Int("mail_port")
	mail_user := common.Config.String("mail_user")
	mail_pwd := common.Config.String("mail_password")
	m := gomail.NewMessage()
	m.SetHeader("From", mail_user)
	m.SetHeader("To", toMail)
	m.SetBody("text/plain", body)
	d := gomail.NewDialer(mail_addr, mail_port, mail_user, mail_pwd)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
