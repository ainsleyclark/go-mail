// Copyright 2020 The Go Mail Authors. All rights reserved.
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

package mail

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
)

// smtpClient represents the data for sending mail via
// plain ol SMTP. Configuration, the client and the
// main send function are parsed for sending
// data.
type smtpClient struct {
	cfg  Config
	send smtpSendFunc
}

// smtpSendFunc defines the function for ending
// SMTP transmissions.
type smtpSendFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

// Creates a new smtp client. Configuration is
// validated before initialisation.
func newSMTP(cfg Config) *smtpClient {
	fmt.Println("Warning, using SMTP is insecure, only use for development.")
	return &smtpClient{
		cfg:  cfg,
		send: smtp.SendMail,
	}
}

// Send mail via plain SMTP. Transmissions are validated
// before sending and attachments are added. Returns
// an error upon failure.
func (m *smtpClient) Send(t *Transmission) (Response, error) {
	err := t.Validate()
	if err != nil {
		return Response{}, err
	}

	auth := smtp.PlainAuth("", m.cfg.FromAddress, m.cfg.Password, m.cfg.URL)

	err = m.send(m.cfg.URL+":"+strconv.Itoa(m.cfg.Port), auth, m.cfg.FromAddress, t.Recipients, m.bytes(t))
	if err != nil {
		fmt.Println(err)
		return Response{}, err
	}

	return Response{
		StatusCode: http.StatusOK,
		Message:    "Email sent successfully",
	}, nil
}

// Processes the transmission and returns the bytes for
// sending. Mime types are set dependant on the
// content passed.
func (m *smtpClient) bytes(t *Transmission) []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(fmt.Sprintf("Subject: %s\n", t.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(t.Recipients, ",")))

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	if t.Attachments.Exists() {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/html; charset=\"ascii\"\n")
	}

	if t.PlainText != "" {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
		buf.WriteString(t.PlainText)
	}

	if t.HTML != "" {
		buf.WriteString("Content-Type: text/html; charset=\"ascii\"\n")
		buf.WriteString(t.HTML)
	}

	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", v.Mime()))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", v.Filename))
			buf.WriteString(fmt.Sprintf("\n--%s", v.B64()))
		}
		buf.WriteString("--")
	}

	return buf.Bytes()
}
