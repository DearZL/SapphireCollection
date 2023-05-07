package utils

import (
	"gopkg.in/gomail.v2"
	"log"
)

// SendMail 发送邮件
func SendMail(mailTo, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress("2628008929@qq.com", "数藏"))
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer("smtp.qq.com", 587, "2628008929@qq.com", "tlirfuqldozgdife")
	err := d.DialAndSend(m)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
