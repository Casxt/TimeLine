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

//SendMail sending an email with attchment, it just a little faster and shorter but its calling is not graceful enough
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

/*

var mailCh chan *gomail.Message

//Start initial the sender
func Start() {
	mailCh = make(chan *gomail.Message)
	go mailSender()
}

func mailSender() {
	d := gomail.NewDialer("smtp.qq.com", 587, "774714620@qq.com", "blwcpdcjpxftbdji")
	var sender gomail.SendCloser
	var err error

	defer close(mailCh)
	defer sender.Close() //这里是否要关闭？

	open := false
	for {
		select {
		case mail, ok := <-mailCh:
			if !ok {
				log.Println("Mail Server Closed")
				return
			}
			if !open {
				if sender, err = d.Dial(); err != nil {
					log.Fatalln(err)
				}
				open = true
			}
			if err := gomail.Send(sender, mail); err != nil {
				log.Fatalln(err)
			}
		// Close the connection to the SMTP server if no email was sent in
		// the last 30 seconds.
		case <-time.After(30 * time.Second):
			if open {
				if err := sender.Close(); err != nil {
					log.Fatalln(err)
				}
				open = false
			}
		}
	}
}
*/
