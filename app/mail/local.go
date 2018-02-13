package mail

import (
	"fmt"
	"github.com/revel/revel"
)

type MailSender interface {
	SendMail(to string, cc []string, body []byte) error
}

type ConsoleMailSender struct{}

func (sender *ConsoleMailSender) SendMail(to string, cc []string, body []byte) (err error) {
	err = nil
	message := fmt.Sprintf("\n******\nto: %s\ncc: %i\nbody:\n%s\n******\n", to, cc, string(body))
	revel.INFO.Println(message)
	return
}
