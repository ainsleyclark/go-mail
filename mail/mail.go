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

import "errors"

var (
	// Debug - Set true to write the HTTP requests in curl for to stdout.
	// Additional information will also be displayed in the errors such as
	// method operations.
	Debug = false
	// ErrEmptyBody is returned by Send when there is nobody attached to the
	// request.
	ErrEmptyBody = errors.New("error, empty body")
)

// Mailer defines the sender for go-mail returning a
// Response or error when an email is sent.
//
// Below is an example of creating and sending a transmission:
// 	cfg := mail.Config{
//    		URL:         "https://api.eu.sparkpost.com",
//    		APIKey:      "my-key",
//    		FromAddress: "hello@gophers.com",
//    		FromName:    "Gopher",
//	}
//
//	mailer, err := drivers.NewSparkPost(cfg)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	tx := &mail.Transmission{
//  		Recipients:  []string{"hello@gophers.com"},
//    		Subject:     "My email",
//    		HTML:        "<h1>Hello from Go Mail!</h1>",
//	}
//
//	result, err := mailer.Send(tx)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	fmt.Printf("%+v\n", result)
type Mailer interface {
	// Send accepts a mail.Transmission to send an email through a particular
	// driver/provider. Transmissions will be validated before sending.
	//
	// A mail.Response or an error will be returned. In some circumstances
	// the body and status code will be attached to the response for debugging.
	//
	Send(t *Transmission) (Response, error)
}
