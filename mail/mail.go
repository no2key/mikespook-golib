package mail

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/mail"
	"path/filepath"
)

const (
	boundary = "f46d043c813270fc6b04c2d223da"
)

type Part interface {
	Bytes() ([]byte, error)
}

type Mail struct {
	Header
	bodies []Part
}

func NewMail() (mail *Mail) {
	mail = &Mail{
		bodies: make([]Part, 0, 1),
	}
	return
}

func addrsToString(addrs []mail.Address) string {
	buf := bytes.NewBuffer(nil)
	for _, a := range addrs {
		buf.WriteString(a.String())
		buf.WriteRune(',')
	}
	return buf.String()
}

func (m *Mail) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.Write(m.Header.Bytes())
	for _, v := range m.bodies {
		b, err := v.Bytes()
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	buf.WriteString("--" + boundary + "--\n")
	return buf.Bytes(), nil
}

func (m *Mail) Append(p Part) {
	m.bodies = append(m.bodies, p)
}

type Subject string

func (s Subject) String() string {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(s)) + "?="
}

type Header struct {
	From    mail.Address
	To      []mail.Address
	Cc      []mail.Address
	Bcc     []mail.Address
	Subject Subject
}

func (h *Header) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("From: " + h.From.String() + "\n")
	buf.WriteString("To: " + addrsToString(h.To) + "\n")
	if len(h.Cc) > 0 {
		buf.WriteString("Cc: " + addrsToString(h.Cc) + "\n")
	}
	if len(h.Bcc) > 0 {
		buf.WriteString("Bcc: " + addrsToString(h.Cc) + "\n")
	}
	buf.WriteString("Subject: " + h.Subject.String() + "\n")
	buf.WriteString("MIME-Version: 1.0\n")
	return buf.Bytes()
}

type Body struct {
	Content     []byte
	ContentType string
	Charset     string
}

func (b *Body) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("\n--" + boundary + "\n")
	buf.WriteString("Content-Type: " + b.ContentType + "; charset=" + b.Charset + "\n")
	buf.WriteString("Content-Transfer-Encoding: base64\n\n")
	buf.WriteString(base64.StdEncoding.EncodeToString(b.Content))
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

type Attachment string

func (a Attachment) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	_, fileName := filepath.Split(string(a))
	buf.WriteString("\n--" + boundary + "\n")
	buf.WriteString("Content-Type: application/octet-stream\n")
	buf.WriteString("Content-Transfer-Encoding: base64\n")
	buf.WriteString("Content-Disposition: attachment; filename=\"" + fileName + "\"\n\n")
	b, err := ioutil.ReadFile(string(a))
	if err != nil {
		return nil, err
	}
	buf.WriteString(base64.StdEncoding.EncodeToString(b))
	buf.WriteString("\n")
	return buf.Bytes(), nil
}
