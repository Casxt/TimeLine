package mail

import "testing"

func TestSendMail(t *testing.T) {
	//Start()
	if err := SendMail("774714620@qq.com", "Go Mail Test", "<h1>Go Mail Test Content</h1>", nil, Attachment{"a.txt", []byte("asdasd")}, Attachment{"b.txt", []byte("sadves")}); err != nil {
		t.Error("TestSendMail Error", err)
	}
}

func TestSendMailWithAttch(t *testing.T) {
	//Start()
	if err := SendMailWithAttch("774714620@qq.com", "Go Mail Test", "<h1>Go Mail Test Content</h1>", []Attachment{Attachment{"a.txt", []byte("asdasd")}, Attachment{"b.txt", []byte("sadves")}}, nil); err != nil {
		t.Error("TestSendMail Error", err)
	}
}
