package helpers

import (
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPHelper struct {
	Host        string
	Port        string
	Username    string
	Password    string
	MailDetails MailDetails
}

type MailDetails struct {
	Receivers []string
	Subject   string
	Data      []byte
	Sender    string
	MimeType  string
}

func NewSMTPHelper(host, port, username, password string) *SMTPHelper {
	return &SMTPHelper{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// SetSender sets the From address.
func (sh *SMTPHelper) SetSender(sender string) {
	sh.MailDetails.Sender = sender
}

// SetReceivers sets the To addresses.
func (sh *SMTPHelper) SetReceivers(receivers []string) {
	sh.MailDetails.Receivers = receivers
}

// SetSubject sets the email subject.
func (sh *SMTPHelper) SetSubject(subject string) {
	sh.MailDetails.Subject = subject
}

// SetPlainBody sets the email body as plain text.
func (sh *SMTPHelper) SetPlainBody(body []byte) {
	sh.MailDetails.Data = body
	sh.MailDetails.MimeType = "text/plain; charset=\"UTF-8\""
}

// SetHTMLBody sets the email body as HTML.
func (sh *SMTPHelper) SetHTMLBody(body []byte) {
	sh.MailDetails.Data = body
	sh.MailDetails.MimeType = "text/html; charset=\"UTF-8\""
}

// SendMail constructs and sends the email using net/smtp.
func (sh *SMTPHelper) SendMail() error {
	if len(sh.MailDetails.Receivers) == 0 {
		return fmt.Errorf("no receivers specified")
	}

	var auth smtp.Auth
	if sh.Username != "" && sh.Password != "" {
		auth = smtp.PlainAuth("", sh.Username, sh.Password, sh.Host)
	}

	// Build MIME message
	header := make(map[string]string)
	header["From"] = sh.MailDetails.Sender
	header["To"] = strings.Join(sh.MailDetails.Receivers, ",")
	header["Subject"] = sh.MailDetails.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = sh.MailDetails.MimeType

	var msgBuilder strings.Builder
	for k, v := range header {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msgBuilder.WriteString("\r\n")
	msgBuilder.Write(sh.MailDetails.Data)

	addr := fmt.Sprintf("%s:%s", sh.Host, sh.Port)
	return smtp.SendMail(addr, auth, sh.MailDetails.Sender, sh.MailDetails.Receivers, []byte(msgBuilder.String()))
}
