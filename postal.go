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
	"fmt"
	sp "github.com/SparkPost/gosparkpost"
	"io/ioutil"
	"net/http"
	"strings"
)


type postal struct {
	cfg    Config
	client sp.Client
	send   sparkSendFunc
}

//// sparkSendFunc defines the function for ending SparkPost
//// transmissions.
//type sparkSendFunc func(t *sp.Transmission) (id string, res *sp.Response, err error)


// Creates a new Postal client. Configuration is
// validated before initialisation.
func newPostal(cfg Config) (*sparkPost, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}



	return &sparkPost{
		cfg:    cfg,

	}, nil
}


func (s *postal) Send(t *Transmission) (Response, error) {
	err := t.Validate()
	if err != nil {
		return Response{}, err
	}

	return Response{}, nil
}


func main() {

	url := "https://postal.reddico.io/api/v1/send/message"
	method := "POST"

	payload := strings.NewReader(`{
    "to": [
        "ainsley@reddico.co.uk"
    ],
    "from": "ainsley@reddico.io",
    "html_body": "This is a test"
}`)

	client := &http.Client{}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("X-Server-API-Key", "Iw0mO9rOsRjKU1thvEmZbGXm")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
