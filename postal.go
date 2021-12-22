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
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type postal struct {
	cfg    Config
	client *http.Client
}

//// sparkSendFunc defines the function for ending SparkPost
//// transmissions.
//type sparkSendFunc func(t *sp.Transmission) (id string, res *sp.Response, err error)

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
	}, nil
}

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

func (p *postal) Send(t *Transmission) (Response, error) {
	err := t.Validate()
	if err != nil {
		return Response{}, err
	}

	m := postalMessage{
		To:        t.Recipients,
		From:      p.cfg.FromAddress,
		Sender:    p.cfg.FromName,
		HTML:      t.HTML,
		PlainText: t.PlainText,
	}

	data, err := json.Marshal(m)
	if err != nil {
		return Response{}, err
	}

	fmt.Println(string(data))

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/send/message", p.cfg.URL), bytes.NewBuffer(data))
	if err != nil {
		return Response{}, err
	}

	// Ensure the token is set for authorisation
	// on the API along with the app name.
	req.Header.Set("X-Server-API-Key", p.cfg.APIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	fmt.Println(string(body))

	// TODO: Unmarshal postal response.
	return Response{
		StatusCode: resp.StatusCode,
		Body:       "",
		Headers:    nil,
		ID:         "",
		Message:    nil,
	}, nil
}
