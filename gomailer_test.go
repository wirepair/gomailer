package gomailer

import (
	"testing"
)

func TestNew(t *testing.T) {
	_ = New("localhost:2500", "testdata\\")
}

func TestAdd(t *testing.T) {
	m := New("localhost:2500", "testdata\\")
	m.Add([]string{"success.txt", "attempt.txt"})
	if len(m.Templates) != 2 {
		t.Fatalf("Error we don't have two templates we have %d! data: %#v", len(m.Templates), m.Templates)
	}
}

func TestSend(t *testing.T) {
	m := New("localhost:2500", "testdata\\")
	m.Add([]string{"success.txt", "attempt.txt"})
	d := attemptData()
	if err := m.Send(d, d.From, d.To, "attempt.txt"); err != nil {
		t.Fatalf("Error sending mail: %v\n", err)
	}
}

func TestSendFail(t *testing.T) {
	m := New("localhost:2501", "testdata\\")
	m.Add([]string{"success.txt", "attempt.txt"})
	d := attemptData()
	if err := m.Send(d, d.From, d.To, "attempt.txt"); err == nil {
		t.Fatal("we had no error even though we couldn't connect!\n")
	}
}

func TestMD5Auth(t *testing.T) {
	m := New("localhost:2500", "testdata\\")
	m.Add([]string{"success.txt", "attempt.txt"})
	d := attemptData()
	m.MD5Auth("test", "password")
	d.From = "adminMD5@localhost"
	if err := m.Send(d, d.From, d.To, "attempt.txt"); err != nil {
		t.Fatalf("Error sending mail: %v\n", err)
	}
}

func TestPlainAuth(t *testing.T) {
	m := New("localhost:2500", "testdata\\")
	m.Add([]string{"success.txt", "attempt.txt"})
	d := attemptData()
	m.PlainAuth("", "test", "password", "localhost")
	d.From = "adminPLAIN@localhost"
	if err := m.Send(d, d.From, d.To, "attempt.txt"); err != nil {
		t.Fatalf("Error sending mail: %v\n", err)
	}
}

func TestSendAddSend(t *testing.T) {
	m := New("localhost:2500", "testdata\\")
	m.Add([]string{"attempt.txt"})
	d := attemptData()
	if err := m.Send(d, d.From, d.To, "attempt.txt"); err != nil {
		t.Fatalf("Error sending mail: %v\n", err)
	}

	m.Add([]string{"success.txt"})
	if len(m.Templates) != 2 {
		t.Fatalf("Error we don't have two templates we have %d! data: %#v", len(m.Templates), m.Templates)
	}
	s := successData()
	s.From = "successAfterAdd@localhost"
	if err := m.Send(s, s.From, s.To, "success.txt"); err != nil {
		t.Fatalf("Error sending mail: %v\n", err)
	}
}

type AttemptData struct {
	From    string
	To      string
	Subject string
	Ip      string
}

func attemptData() *AttemptData {
	d := new(AttemptData)
	d.From = "admin@localhost"
	d.To = "x@localhost"
	d.Subject = "test"
	d.Ip = "127.0.0.1"
	return d
}

type SuccessData struct {
	From    string
	To      string
	Subject string
	Ip      string
	Url     string
	Token   string
}

func successData() *SuccessData {
	s := new(SuccessData)
	s.From = "admin@localhost"
	s.To = "x@localhost"
	s.Subject = "test"
	s.Ip = "127.0.0.1"
	s.Url = "http://localhost/success?token"
	s.Token = "abc1234"
	return s
}
