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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TODO
// https://docs.postalserver.io/developer/api
// https://apiv1.postalserver.io/controllers/send/message.html
type postal struct {
	cfg        Config
	client     *http.Client
	marshaller func(v interface{}) ([]byte, error)
	bodyReader func(r io.Reader) ([]byte, error)
}

// Creates a new Postal client. Configuration is
// validated before initialisation.
func newPostal(cfg Config) (*postal, error) {
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

const (
	postalErrorMessage = "error sending message to Postal api"
)

// {"status":"success","time":0.07,"flags":{},"data":{"message_id":"2cdf0f8f-18e5-4286-bb66-a22fb0c3c30a@rp.postal.example.com","messages":{"ainsley@reddico.co.uk":{"id":3,"token":"y5ChzHNHWVnR"}}}}
// {"status":"error","time":0.0,"flags":{},"data":{"code":"FromAddressMissing","message":"The From address is missing and is required"}}

type postalMessage struct {
	To          []string           `json:"to"`
	CC          []string           `json:"cc"`
	BCC         []string           `json:"bcc"`
	From        string             `json:"from"`
	Sender      string             `json:"sender"`
	Subject     string             `json:"subject"`
	HTML        string             `json:"html_body"`
	PlainText   string             `json:"plain_body"`
	Attachments []postalAttachment `json:"attachments"`
}

type postalAttachment struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"`
}

type postalResponse struct {
	Status string                 `json:"status"`
	Time   float32                `json:"time"`
	Flags  map[string]interface{} `json:"flags"`
	Data   map[string]interface{} `json:"data"`
}

func (p *postal) Send(t *Transmission) (Response, error) {
	err := t.Validate()
	if err != nil {
		return Response{}, err
	}

	m := postalMessage{
		To:          t.Recipients,
		CC:          t.CC,
		BCC:         t.BCC,
		From:        p.cfg.FromAddress,
		Sender:      p.cfg.FromName,
		Subject:     t.Subject,
		HTML:        t.HTML,
		PlainText:   t.PlainText,
		Attachments: nil,
	}

	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			m.Attachments = append(m.Attachments, postalAttachment{
				Name:        v.Filename,
				ContentType: v.Mime(),
				Data:        v.B64(),
			})
		}
	}

	data, err := p.marshaller(m)
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/send/message", p.cfg.URL), bytes.NewBuffer(data))
	if err != nil {
		return Response{}, err
	}

	// Ensure the API Key is set for authorisation
	// and add the JSON content type.
	req.Header.Set("X-Server-API-Key", p.cfg.APIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Response{}, errors.New(postalErrorMessage)
	}

	buf, err := p.bodyReader(resp.Body)
	if err != nil {
		return Response{}, err
	}

	pResponse := postalResponse{}
	err = json.Unmarshal(buf, &pResponse)
	if err != nil {
		return Response{}, err
	}

	if pResponse.Status != "success" {
		msg := postalErrorMessage
		if code, ok := pResponse.Data["code"]; ok {
			msg = fmt.Sprintf("%s - code: %s", msg, code)
		}
		if message, ok := pResponse.Data["message"]; ok {
			msg = fmt.Sprintf("%s, message: %s", msg, message)
		}
		return Response{}, errors.New(msg)
	}

	response := Response{
		StatusCode: http.StatusOK,
		Body:       string(buf),
		Message:    "Successfully sent Postal email",
	}

	if val, ok := pResponse.Data["message_id"]; ok {
		response.ID = fmt.Sprintf("%v", val)
	}

	return response, nil
}
