package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// Sender sends email via SMTP STARTTLS (typical port 587). No card data is ever handled here.
type Sender struct {
	Host, Port, User, Password, From string
}

func (s *Sender) Configured() bool {
	return strings.TrimSpace(s.Host) != "" && strings.TrimSpace(s.From) != ""
}

func (s *Sender) SendPlain(to, subject, body string) error {
	if !s.Configured() || strings.TrimSpace(to) == "" {
		return nil
	}
	port := s.Port
	if port == "" {
		port = "587"
	}
	addr := net.JoinHostPort(s.Host, port)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s\r\n",
		s.From, to, subject, body))

	host, _, _ := net.SplitHostPort(addr)
	tlsCfg := &tls.Config{ServerName: host}

	if s.User != "" && s.Password != "" {
		auth := smtp.PlainAuth("", s.User, s.Password, host)
		return smtp.SendMail(addr, auth, s.From, []string{to}, msg)
	}

	// Opportunistic STARTTLS without auth (rare); most providers require auth.
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err := c.StartTLS(tlsCfg); err != nil {
		return err
	}
	if err := c.Mail(s.From); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
}
