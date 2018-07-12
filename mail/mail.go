package mail

import (
	"io"

	gomail "gopkg.in/gomail.v2"
)

//Attachment Define the Attachment
type Attachment struct {
	Name    string
	Content []byte
}

//Dialer return a new dialer
func Dialer() (dialer *gomail.Dialer) {
	return gomail.NewDialer("smtp.qq.com", 465, "774714620@qq.com", "wwnqjvhiuqdbbebe")
}

//Sender return a new sender
func Sender() (sender gomail.SendCloser, err error) {
	return Dialer().Dial()
}

//SendMail sending an email with attchment, it just a little shorter but its calling is not graceful enough
//注意区别，sender 在 attachment之前
func SendMail(To string, Subject string, Content string, sender gomail.SendCloser, Attachments ...Attachment) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "maple@forer.cn")
	m.SetHeader("To", To)
	m.SetHeader("Subject", Subject)
	m.SetBody("text/html", Content)
	for _, attachment := range Attachments {
		m.Attach(attachment.Name, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachment.Content)
			return err
		}))
	}

	if sender == nil {
		dialer := Dialer()
		return dialer.DialAndSend(m)
	}
	return gomail.Send(sender, m)
}

//SendMailWithAttch sending an email with attchment which tell you can add attch with it
func SendMailWithAttch(To string, Subject string, Content string, Attachments []Attachment, sender gomail.SendCloser) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "maple@forer.cn")
	m.SetHeader("To", To)
	m.SetHeader("Subject", Subject)
	m.SetBody("text/html", Content)
	for _, attachment := range Attachments {
		m.Attach(attachment.Name, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachment.Content)
			return err
		}))
	}

	if sender == nil {
		dialer := Dialer()
		return dialer.DialAndSend(m)
	}
	return gomail.Send(sender, m)
}
