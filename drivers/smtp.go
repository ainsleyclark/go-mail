// Copyright 2022 Ainsley Clark. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package drivers

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/ainsleyclark/go-mail/mail"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strings"
)

// smtpClient represents the data for sending mail via
// plain ol SMTP. Configuration, the client and the
// main send function are parsed for sending
// data.
type smtpClient struct {
	cfg  mail.Config
	send smtpSendFunc
}

// smtpSendFunc defines the function for ending
// SMTP mail.Transmissions.
type smtpSendFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

// NewSMTP creates a new smtp client. Configuration
// is validated before initialisation.
func NewSMTP(cfg mail.Config) (mail.Mailer, error) {
	if cfg.URL == "" {
		return nil, errors.New("driver requires a url")
	}
	if cfg.FromAddress == "" {
		return nil, errors.New("driver requires from address")
	}
	if cfg.FromName == "" {
		return nil, errors.New("driver requires from name")
	}
	if cfg.Password == "" {
		return nil, errors.New("driver requires a password")
	}
	return &smtpClient{
		cfg:  cfg,
		send: smtp.SendMail,
	}, nil
}

// Send mail via plain SMTP. mail.Transmissions are validated
// before sending and attachments are added. Returns
// an error upon failure.
func (m *smtpClient) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	servername := fmt.Sprintf("%s:%d", m.cfg.URL, m.cfg.Port)

	// TLS config for sending SMTP mail.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.cfg.URL,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require a ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", m.cfg.URL, tlsConfig)
	if err != nil {
		return mail.Response{}, err
	}

	mailer, err := smtp.NewClient(conn, m.cfg.URL)
	if err != nil {
		return mail.Response{}, err
	}

	// Authenticate the SMTP client.
	auth := smtp.PlainAuth("", m.cfg.FromAddress, m.cfg.Password, m.cfg.URL)

	// Auth
	err = mailer.Auth(auth)
	if err != nil {
		return mail.Response{}, err
	}

	// To && From
	err = mailer.Mail(fr)
	if err != nil {
		return mail.Response{}, err
	}

	if err = mailer.Rcpt(to.Address) err != nil {
		return mail.Response{}, err
	}

	err = m.send(servername, auth, m.cfg.FromAddress, m.getTo(t), m.bytes(t))
	if err != nil {
		return mail.Response{}, err
	}

	return mail.Response{
		StatusCode: http.StatusOK,
		Message:    "Email sent successfully",
	}, nil
}

// getTo returns the merged mail.Transmission recipients, CC and
// BCC email addresses.
func (m *smtpClient) getTo(t *mail.Transmission) []string {
	var to []string
	to = append(t.Recipients, t.CC...)
	to = append(to, t.BCC...)
	return to
}

// Processes the mail.Transmission and returns the bytes for
// sending. Mime types are set dependent on the
// content passed.
// See: https://gist.github.com/tylermakin/d820f65eb3c9dd98d58721c7fb1939a8?permalink_comment_id=2703291
func (m *smtpClient) bytes(t *mail.Transmission) []byte {
	buf := bytes.NewBuffer(nil)

	for k, v := range t.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	if t.HasAttachments() {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", boundary))
	} else {
		buf.WriteString(fmt.Sprintf("Content-Type: text/html; charset=UTF-8; boundary=%s\r\n", boundary))
	}

	buf.WriteString(fmt.Sprintf("Subject: %s\n", t.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(t.Recipients, ",")))

	if t.HasCC() {
		buf.WriteString(fmt.Sprintf("CC: %s\n", strings.Join(t.CC, ",")))
	}

	buf.WriteString("\n")

	if t.PlainText != "" {
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
		buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		buf.WriteString(fmt.Sprintf("\r\n%s\r\n\n", strings.TrimSpace(t.PlainText)))
	}

	if t.HTML != "" {
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
		buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		buf.WriteString(fmt.Sprintf("\r\n%s\r\n\n", t.HTML))
	}

	if t.HasAttachments() {
		for _, v := range t.Attachments {
			buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", v.Mime()))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", v.Filename))
			buf.WriteString(fmt.Sprintf("\r\n--%s", v.B64()))
		}
		buf.WriteString("--")
	}

	buf.WriteString(fmt.Sprintf("--%s--\n", boundary))

	return buf.Bytes()
}
