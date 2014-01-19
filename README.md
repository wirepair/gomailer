# gomailer 
This library is for sending simple templated text based emails to users. It is NOT an SMTP server.

### installation
go get github.com/wirepair/gomailer

### example  
Given the template of:
```
From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

Hello,

Someone at {{.Ip}} has attempted to register your email address in our system. If this was you, please note you already have an account and disregard this email. If this wasn't you, please make sure your password is safe and has not been used in any other systems as an attacker may be targetting you. 
```

You would create the necessary structure required by the template, add optional authentication information and then call Send.
```Go
package main

import (
	"github.com/wirepair/gomailer"
	"log"
)

type AttemptData struct {
	From    string
	To      string
	Subject string
	Ip      string
}

func main() {
	m := gomailer.New("localhost:2500", "testdata\\")
	m.Add([]string{"attempt.txt"})

	// Create the struct necessary for the email template.
	d := new(AttemptData)
	d.From = "example@localhost"
	d.To = "x@localhost"
	d.Subject = "test"
	d.Ip = "127.0.0.1"
	// add auth, if necessary
	m.MD5Auth("test", "password")
	if err := m.Send(d, d.From, d.To, "attempt.txt"); err != nil {
		log.Fatalf("Error sending mail: %v\n", err)
	}
}
```