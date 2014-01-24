/*
The MIT License (MIT)

Copyright (c) 2014 isaac dawson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// a very simple mail client for sending notifications/text based email.
package gomailer

import (
	"bytes"
	"net/smtp"
	"sync"
	"text/template"
)

// Gomail interface
type Gomailer interface {
	Add(templates []string) error
	MD5Auth(username, secret string)
	PlainAuth(identity, username, password, host string)
	Send(mailData interface{}, sender, recipient, templateName string) error
}

// Gomail structure to hold template and server information for sending simple
// templated emails.
type Gomail struct {
	sync.RWMutex
	Server    string             // the smtp server, should be in format: server:port (ie, localhost:25)
	Path      string             // path to the directory containing templates
	Templates []string           // a slice of template files
	auth      smtp.Auth          // authentication information, if necessary.
	tmpl      *template.Template // the template reference.
}

// New takes a server address and path and returns a new Gomail object.
// Note this path should point to where the text based templates are stored.
func New(server, path string) *Gomail {
	gomail := new(Gomail)
	gomail.Server = server
	gomail.Path = path
	gomail.tmpl = new(template.Template)
	return gomail
}

// Safely adds new templates to the mailer object. Returns an error if the template
// fails to parse properly.
func (g *Gomail) Add(templates []string) error {
	defer g.Unlock()

	var err error

	g.Lock()
	if g.Templates == nil {
		g.Templates = make([]string, 0)
	}
	g.Templates = append(g.Templates, templates...)

	for _, v := range templates {
		if g.tmpl, err = g.tmpl.ParseFiles(g.Path + v); err != nil {
			return err
		}
	}
	return err
}

// Safely adds md5 based authentication when sending emails, overwrites any
// previous calls to MD5Auth or PlainAuth.
func (g *Gomail) MD5Auth(username, secret string) {
	g.addAuth(smtp.CRAMMD5Auth(username, secret))
}

// Safely adds plaintext based authentication when sending emails, overwrites any
// previous calls to MD5Auth or PlainAuth.
func (g *Gomail) PlainAuth(identity, username, password, host string) {
	g.addAuth(smtp.PlainAuth(identity, username, password, host))
}

// adds the passed in auth to our Gomail object with synchronization.
func (g *Gomail) addAuth(auth smtp.Auth) {
	defer g.Unlock()

	g.Lock()
	g.auth = auth
}

// Send sends a mail with an arbitrary object with properties the template expects.
// returns an error if the template fails to execute or the mail itself fails to be
// sent.
func (g *Gomail) Send(mailData interface{}, sender, recipient, templateName string) error {
	defer g.RUnlock()

	mail := new(bytes.Buffer)
	g.RLock()
	if err := g.tmpl.ExecuteTemplate(mail, templateName, mailData); err != nil {
		return err
	}
	return smtp.SendMail(g.Server, g.auth, sender, []string{recipient}, mail.Bytes())
}
