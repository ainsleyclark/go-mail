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

package drivers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ainsleyclark/go-mail/internal/httputil"
	"github.com/ainsleyclark/go-mail/mail"
	"net/http"
	"strings"
	"time"
)

// postal represents the entity for sending mail via the
// Postmark API.
//
// See: https://postmarkapp.com/developer/api/email-api
type postmark struct {
	cfg    mail.Config
	client httputil.Requester
}

const (
	// postalEndpoint defines the endpoint to POST to.
	postmarkEndpoint = "https://api.postmarkapp.com/email"
	// postmarkErrorMessage defines the message when an error occurred
	// when sending mail via the Postmark API.
	postmarkErrorMessage = "error sending transmission to Postmark API"
)

// NewPostmark creates a new Postmark client. Configuration
// is validated before initialisation.
func NewPostmark(cfg mail.Config) (mail.Mailer, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &postmark{
		cfg:    cfg,
		client: httputil.NewClient(),
	}, nil
}

type (
	// postmarkTransmission defines the data to be sent to the Postmark API.
	postmarkTransmission struct {
		From      string `json:"From"`
		To        string `json:"To"`
		CC        string `json:"Cc"`
		BCC       string `json:"Bcc"`
		Subject   string `json:"Subject"`
		Tag       string `json:"Tag"`
		HTML      string `json:"HtmlBody"`
		PlainText string `json:"TextBody"`
		ReplyTo   string `json:"ReplyTo"`
		//Headers   []struct {
		//	Name  string `json:"Name"`
		//	Value string `json:"Value"`
		//} `json:"headers"`
		TrackOpens  bool                 `json:"TrackOpens"`
		TrackLinks  string               `json:"TrackLinks"`
		Attachments []postmarkAttachment `json:"Attachments"`
		//Metadata    struct {
		//	Color    string `json:"color"`
		//	ClientID string `json:"client-id"`
		//} `json:"Metadata"`
		MessageStream string `json:"MessageStream"`
	}
	// postmarkAttachment defines a singular Postmark mail attachment.
	postmarkAttachment struct {
		Name        string `json:"Name"`
		Content     string `json:"Content"`
		ContentType string `json:"ContentType"`
		ContentID   string `json:"ContentID,omitempty"`
	}
	// postmarkResponse defines the data sent back from the Postmark API.
	// An error code of 0 represents a successful transmission.
	postmarkResponse struct {
		To          string    `json:"To"`
		SubmittedAt time.Time `json:"SubmittedAt"`
		MessageID   string    `json:"MessageID"`
		ErrorCode   int       `json:"ErrorCode"`
		Message     string    `json:"Message"`
	}
)

// HasError determines if the Postal call was successful
// by comparing the status.
func (p *postmarkResponse) HasError() bool {
	return p.ErrorCode != 0
}

// Error returns a formatted response error.
func (p *postmarkResponse) Error() error {
	return fmt.Errorf("%s - code: %d, message: %s", postmarkErrorMessage, p.ErrorCode, p.Message)
}

// Send posts the go mail Transmission to the Postal
// API. Transmissions are validated before sending
// and attachments are added. Returns an error
// upon failure.
func (d *postmark) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	m := postmarkTransmission{
		To:            strings.Join(t.Recipients, ","),
		CC:            strings.Join(t.CC, ","),
		BCC:           strings.Join(t.BCC, ","),
		From:          fmt.Sprintf("%s <%s>", d.cfg.FromName, d.cfg.FromAddress),
		Subject:       t.Subject,
		HTML:          t.HTML,
		PlainText:     t.PlainText,
		MessageStream: "outbound",
	}

	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			m.Attachments = append(m.Attachments, postmarkAttachment{
				Name:        v.Filename,
				ContentType: v.Mime(),
				Content:     v.B64(),
			})
		}
	}

	pl, err := httputil.NewJSONData(m)
	if err != nil {
		return mail.Response{}, err
	}

	req := httputil.NewHTTPRequest(http.MethodPost, postmarkEndpoint)
	req.AddHeader("X-Postmark-Server-Token", d.cfg.APIKey)

	buf, resp, err := d.client.Do(context.Background(), req, pl)
	if err != nil {
		return mail.Response{}, err
	}

	// Unmarshal the buffer into a postmarkResponse.
	response := postmarkResponse{}
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return mail.Response{}, err
	}

	if response.HasError() {
		return mail.Response{}, response.Error()
	}

	return mail.Response{
		StatusCode: resp.StatusCode,
		Body:       string(buf),
		Headers:    resp.Header,
		ID:         response.MessageID,
		Message:    response.Message,
	}, nil
}
