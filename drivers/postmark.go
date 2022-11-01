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
	"context"
	"encoding/json"
	"fmt"
	"github.com/ainsleyclark/go-mail/internal/client"
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
	client client.Requester
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
		client: client.New(cfg.Client),
	}, nil
}

type (
	// postmarkTransmission defines the data to be sent to the Postmark API.
	postmarkTransmission struct {
		From        string               `json:"From"`
		To          string               `json:"To"`
		CC          string               `json:"Cc"`
		BCC         string               `json:"Bcc"`
		Subject     string               `json:"Subject"`
		Tag         string               `json:"Tag"`
		HTML        string               `json:"HtmlBody"`
		PlainText   string               `json:"TextBody"`
		ReplyTo     string               `json:"ReplyTo"`
		Headers     []postmarkHeader     `json:"headers"`
		TrackOpens  bool                 `json:"TrackOpens"`
		TrackLinks  string               `json:"TrackLinks"`
		Attachments []postmarkAttachment `json:"Attachments"`
		Metadata    struct {
			Color    string `json:"color"`
			ClientID string `json:"client-id"`
		} `json:"Metadata"`
		MessageStream string `json:"MessageStream"`
	}
	// postmarkHeaders defines the key value pair of custom headers
	// to send with the email.
	postmarkHeader struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
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
	//
	// Example JSON Responses:
	// {"To":"info@ainsleyclark.com","SubmittedAt":"2021-12-29T15:58:17.8637679Z","MessageID":"947125ed-9e43-4dce-b66c-def49198b3d3","ErrorCode":0,"Message":"OK"}
	// {"ErrorCode":300,"Message":"Zero recipients specified"}
	postmarkResponse struct {
		To          string    `json:"To"`
		SubmittedAt time.Time `json:"SubmittedAt"`
		ID          string    `json:"MessageID"`
		ErrorCode   int       `json:"ErrorCode"`
		Message     string    `json:"Message"`
	}
)

func (r *postmarkResponse) Unmarshal(buf []byte) error {
	resp := &postmarkResponse{}
	err := json.Unmarshal(buf, resp)
	if err != nil {
		return err
	}
	*r = *resp
	return nil
}

func (r *postmarkResponse) CheckError(response *http.Response, buf []byte) error {
	if r.ErrorCode == 0 {
		return nil
	}
	if len(buf) == 0 {
		return mail.ErrEmptyBody
	}
	return fmt.Errorf("%s - code: %d, message: %s", postmarkErrorMessage, r.ErrorCode, r.Message)
}

func (r *postmarkResponse) Meta() httputil.Meta {
	return httputil.Meta{
		Message: r.Message,
		ID:      r.ID,
	}
}

func (d *postmark) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	tx := postmarkTransmission{
		To:            strings.Join(t.Recipients, ","),
		CC:            strings.Join(t.CC, ","),
		BCC:           strings.Join(t.BCC, ","),
		From:          fmt.Sprintf("%s <%s>", d.cfg.FromName, d.cfg.FromAddress),
		Subject:       t.Subject,
		HTML:          t.HTML,
		PlainText:     t.PlainText,
		MessageStream: "outbound",
	}

	if t.HasAttachments() {
		for _, v := range t.Attachments {
			tx.Attachments = append(tx.Attachments, postmarkAttachment{
				Name:        v.Filename,
				ContentType: v.Mime(),
				Content:     v.B64(),
			})
		}
	}

	for k, v := range t.Headers {
		tx.Headers = append(tx.Headers, postmarkHeader{
			Name:  k,
			Value: v,
		})
	}

	pl, err := newJSONData(tx)
	if err != nil {
		return mail.Response{}, err
	}

	req := httputil.NewHTTPRequest(http.MethodPost, postmarkEndpoint)
	req.AddHeader("X-Postmark-Server-Token", d.cfg.APIKey)

	return d.client.Do(context.Background(), req, pl, &postmarkResponse{})
}
