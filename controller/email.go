package controllers

import (
	"bytes"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
	"html/template"
	"io/ioutil"
)

func SendEmail(recipient string, subject string, tmplName string, tmplValue map[string]string) error {
	f, err := ioutil.ReadFile(fmt.Sprintf("templates/%s.tmpl", tmplName))
	if err != nil {
		main.LogService.Printf("email: failed to read mail template %s: %s", tmplName, err)
		return fmt.Errorf("failed to read mail template %s", tmplName)
	}

	tmpl := template.Must(template.New(tmplName).Parse(string(f)))

	buf := bytes.NewBufferString("")

	err = tmpl.Execute(buf, tmplValue)
	if err != nil {
		main.LogService.Printf("email: failed to execute mail template %s: %s", tmplName, err)
		return fmt.Errorf("failed to execute mail template %s", tmplName)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", main.Conf.Email.SMTP.Username, "Teamer")
	m.SetAddressHeader("To", recipient, "")
	m.SetHeader("Subject", subject)
	m.SetBody(echo.MIMETextHTML, buf.String())
	m.SetHeader("Message-ID", fmt.Sprintf("<%s-%s@teamer.cc>", tmplName, uniuri.NewLen(32)))

	err = main.MailDialer.DialAndSend(m)
	if err != nil {
		main.LogService.Printf("email: failed to send email %s", err)
	}
	return err
}