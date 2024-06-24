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
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	mocks "github.com/ainsleyclark/go-mail/internal/mocks/client"
	"github.com/flightaware/go-mail/internal/httputil"
	"github.com/flightaware/go-mail/mail"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// DriversTestSuite defines the helper used for mail
// testing.
type DriversTestSuite struct {
	suite.Suite
	base string
}

// Assert testing has begun.
func TestMail(t *testing.T) {
	suite.Run(t, new(DriversTestSuite))
}

// Assigns test base.
func (t *DriversTestSuite) SetupSuite() {
	wd, err := os.Getwd()
	t.NoError(err)
	t.base = wd
}

const (
	// DataPath defines where the test data resides.
	DataPath = "testdata"
)

var (
	// Trans is the transmission used for testing.
	Trans = &mail.Transmission{
		Recipients: []string{"recipient@test.com"},
		CC:         []string{"cc@test.com"},
		BCC:        []string{"bcc@test.com"},
		Subject:    "Subject",
		HTML:       "<h1>HTML</h1>",
		PlainText:  "PlainText",
		Headers: map[string]string{
			"X-Go-Mail": "Test",
		},
	}
	// Trans is the transmission with an
	// attachment used for testing.
	TransWithAttachment = &mail.Transmission{
		Recipients:  []string{"recipient@test.com"},
		Subject:     "Subject",
		HTML:        "<h1>HTML</h1>",
		PlainText:   "PlainText",
		Attachments: []mail.Attachment{{Filename: "test.jpg"}},
	}
	// Config is the default configuration used
	// for testing.
	Comfig = mail.Config{
		URL:         "my-url",
		APIKey:      "my-key",
		FromAddress: "hello@gophers.com",
		FromName:    "Gopher",
		Domain:      "my-domain",
	}
)

// Returns a PNG attachment for testing.
func (t *DriversTestSuite) Attachment(name string) mail.Attachment {
	path := filepath.Join(t.base, DataPath, name)
	file, err := os.ReadFile(path)
	if err != nil {
		t.Fail("error getting attachment with the path: "+path, err)
	}
	return mail.Attachment{
		Filename: name,
		Bytes:    file,
	}
}

func (t *DriversTestSuite) UtilTestUnmarshal(r httputil.Responder, buf []byte) {
	errBuf := []byte("wrong")
	err := r.Unmarshal(errBuf)
	t.Error(err)
	err = r.Unmarshal(buf)
	t.NoError(err)
}

func (t *DriversTestSuite) UtilTestMeta(r httputil.Responder, message, id string) {
	got := r.Meta()
	t.Equal(message, got.Message)
	t.Equal(id, got.ID)
}

func (t *DriversTestSuite) UtilTestSend(fn func(m *mocks.Requester) mail.Mailer, json bool) {
	res := mail.Response{
		StatusCode: http.StatusOK,
		Body:       []byte("body"),
		Headers:    nil,
		ID:         "1",
		Message:    "success",
	}

	tt := map[string]struct {
		input  *mail.Transmission
		mock   func(m *mocks.Requester)
		jsonFn func(obj interface{}) (*httputil.JSONData, error)
		want   interface{}
	}{
		"Success": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(res, nil)
			},
			httputil.NewJSONData,
			res,
		},
		"With Attachment": {
			TransWithAttachment,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(res, nil)
			},
			httputil.NewJSONData,
			res,
		},
		"Validation Failed": {
			nil,
			nil,
			httputil.NewJSONData,
			"can't validate a nil transmission",
		},
		"JSON Error": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(mail.Response{}, errors.New("send error"))
			},
			func(obj interface{}) (*httputil.JSONData, error) {
				return nil, errors.New("json error")
			},
			"json error",
		},
		"Send Error": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(mail.Response{}, errors.New("send error"))
			},
			httputil.NewJSONData,
			"send error",
		},
	}

	for name, test := range tt {
		if name == "JSON Error" && !json {
			continue
		}

		t.Run(name, func() {
			orig := newJSONData
			defer func() { newJSONData = orig }()
			newJSONData = test.jsonFn

			requester := &mocks.Requester{}
			if test.mock != nil {
				test.mock(requester)
			}

			m := fn(requester)

			got, err := m.Send(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.want, got)
		})
	}
}
