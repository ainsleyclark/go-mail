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
	"errors"
	"github.com/ainsleyclark/go-mail/internal/httputil"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/ainsleyclark/go-mail/mocks/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
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
	Comfig =    mail.Config{
		URL:         "my-url",
		APIKey:      "my-key",
		FromAddress: "hello@gophers.com",
		FromName:    "Gopher",
		Domain:      "my-domain",
	}
)

// Returns a PNG attachment for testing.
func (t *DriversTestSuite) Attachment(name string) mail.Attachment {
	path := t.base + string(os.PathSeparator) + DataPath + string(os.PathSeparator) + name
	file, err := ioutil.ReadFile(path)

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

func (t *DriversTestSuite) UtilTestCheckError(r httputil.Responder, errMsg string, checkBody bool) {
	tt := map[string]struct {
		response *http.Response
		buf      []byte
		want     error
	}{
		"Error": {
			&http.Response{StatusCode: http.StatusInternalServerError},
			[]byte("test"),
			errors.New(errMsg),
		},
		"200": {
			&http.Response{StatusCode: http.StatusOK},
			nil,
			nil,
		},
		"Empty Body": {
			&http.Response{StatusCode: http.StatusInternalServerError},
			nil,
			mail.ErrEmptyBody,
		},
	}

	for name, test := range tt {
		if !checkBody && name == "Empty Body" {
			continue
		}
		t.Run(name, func() {
			err := r.CheckError(test.response, test.buf)
			if err != nil {
				t.Contains(err.Error(), test.want.Error())
				return
			}
			t.Equal(test.want, err)
		})
	}
}

func (t *DriversTestSuite) UtilTestMeta(r httputil.Responder, message, id string) {
	got := r.Meta()
	t.Equal(message, got.Message)
	t.Equal(id, got.ID)
}

func (t *DriversTestSuite) UtilTestSend(fn func(m *mocks.Requester) mail.Mailer) {
	res := mail.Response{
		StatusCode: http.StatusOK,
		Body:       []byte("body"),
		Headers:    nil,
		ID:         "1",
		Message:    "success",
	}

	tt := map[string]struct {
		input *mail.Transmission
		mock  func(m *mocks.Requester)
		want  interface{}
	}{
		"Success": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(res, nil)
			},
			res,
		},
		"With Attachment": {
			TransWithAttachment,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(res, nil)
			},
			res,
		},
		"Validation Failed": {
			nil,
			nil,
			"can't validate a nil transmission",
		},
		"Send Error": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(mail.Response{}, errors.New("send error"))
			},
			"send error",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
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
