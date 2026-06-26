package email

import (
	"crypto/tls"
	"fmt"

	"github.com/helays/gomail/v2"
	cfg_proxy "helay.net/go/utils/v3/config/cfg-proxy"
	cfg_tls "helay.net/go/utils/v3/config/cfg-tls"
)

type EmailSender struct {
	SmtpSenderMail   string   `json:"smtp_sender_mail" yaml:"smtp_sender_mail"`
	SmtpSenderAlias  string   `json:"smtp_sender_alias" yaml:"smtp_sender_alias"`
	EnableExternalTo bool     `json:"enable_external_to" yaml:"enable_external_to"` // 是否启用外部邮件接收人，可用在调用层，是否决定重写 SmtpSenderTo
	SmtpSenderTo     []string `json:"smtp_sender_to" yaml:"smtp_sender_to"`
	SmtpSenderAddr   string   `json:"smtp_sender_addr" yaml:"smtp_sender_addr"`
	SmtpSenderPort   int      `json:"smtp_sender_port" yaml:"smtp_sender_port"`
	Account          string   `json:"account" yaml:"account"`
	Password         string   `json:"password" yaml:"password"`

	Proxy cfg_proxy.Proxy `json:"proxy" yaml:"proxy"`
	TLS   *cfg_tls.TLS    `json:"tls" yaml:"tls"`
}

func (e *EmailSender) Clone(dst []string) *EmailSender {
	out := &EmailSender{
		SmtpSenderMail:  e.SmtpSenderMail,
		SmtpSenderAlias: e.SmtpSenderAlias,
		SmtpSenderTo:    e.SmtpSenderTo,
		SmtpSenderPort:  e.SmtpSenderPort,
		SmtpSenderAddr:  e.SmtpSenderAddr,
		Account:         e.Account,
		Password:        e.Password,
		Proxy:           e.Proxy,
		TLS:             e.TLS,
	}
	if len(dst) > 0 {
		out.SmtpSenderTo = dst
	}
	return out
}

func (e *EmailSender) NewDialer() (*gomail.Dialer, error) {
	d := gomail.NewDialer(e.SmtpSenderAddr, e.SmtpSenderPort, e.Account, e.Password)
	err := e.Proxy.Valid()
	if err != nil {
		return nil, fmt.Errorf("代理配置有误 %v", err)
	}
	if d.Dialer, err = e.Proxy.AutoDialer(); err != nil {
		return nil, fmt.Errorf("自动代理配置有误 %v", err)
	}
	if e.TLS != nil {
		d.TLSConfig = &tls.Config{
			ServerName:         e.SmtpSenderAddr,
			MaxVersion:         e.TLS.MaxVersion.Version(),
			MinVersion:         e.TLS.MinVersion.Version(),
			InsecureSkipVerify: true,
		}
		if len(e.TLS.CipherSuites) > 0 {
			d.TLSConfig.CipherSuites = make([]uint16, len(e.TLS.CipherSuites))
			for idx, _ := range e.TLS.CipherSuites {
				d.TLSConfig.CipherSuites[idx] = e.TLS.CipherSuites[idx].CipherSuite()
			}
		}
	}

	return d, nil

}

func (e *EmailSender) Send(subject, body string) error {
	em := gomail.NewMessage()
	em.SetAddressHeader("From", e.SmtpSenderMail, e.SmtpSenderAlias) // 设置发件箱 以及别称
	em.SetHeader("To", e.SmtpSenderTo...)                            // 发给谁
	em.SetHeader("Subject", subject)                                 // 主题

	em.SetBody("text/html", body) // 设置邮件正文
	d, err := e.NewDialer()
	if err != nil {
		return err
	}
	return d.DialAndSend(em)
}
