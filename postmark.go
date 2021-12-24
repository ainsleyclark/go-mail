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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type postmark struct {
	cfg        Config
	client     *http.Client
	marshaller func(v interface{}) ([]byte, error)
	bodyReader func(r io.Reader) ([]byte, error)
}

func newPostmark(cfg Config) (*postal, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &postal{
		cfg: cfg,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		marshaller: json.Marshal,
		bodyReader: io.ReadAll,
	}, nil
}

// postmarkMessage defines the data to be sent to the Postmark API.
type postmarkMessage struct {
	From     string `json:"From"`
	To       string `json:"To"`
	CC       string `json:"Cc"`
	BCC      string `json:"Bcc"`
	Subject  string `json:"Subject"`
	Tag      string `json:"Tag"`
	HTML string `json:"HtmlBody"`
	PlainText string `json:"TextBody"`
	ReplyTo  string `json:"ReplyTo"`
	Headers  []struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	} `json:"Headers"`
	TrackOpens  bool   `json:"TrackOpens"`
	TrackLinks  string `json:"TrackLinks"`
	Attachments []postmarkAttachment `json:"Attachments"`
	Metadata struct {
		Color    string `json:"color"`
		ClientID string `json:"client-id"`
	} `json:"Metadata"`
	MessageStream string `json:"MessageStream"`
}

type postmarkAttachment struct {
	Name        string `json:"Name"`
	Content     string `json:"Content"`
	ContentType string `json:"ContentType"`
	ContentID   string `json:"ContentID,omitempty"`
}

func (p *postmark) Send(t *Transmission) (Response, error) {
	err := t.Validate()
	if err != nil {
		return Response{}, err
	}

	m := postmarkMessage{
		To:          strings.Join(t.Recipients, ","),
		CC:          strings.Join(t.CC, ","),
		BCC:         strings.Join(t.BCC, ","),
		From:        fmt.Sprintf("%s <%s>", p.cfg.FromName, p.cfg.FromAddress),
		Subject:     t.Subject,
		HTML:        t.HTML,
		PlainText:   t.PlainText,
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

	return Response{}, err
}
