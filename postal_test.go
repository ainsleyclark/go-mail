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
	"net/http"
	"net/http/httptest"
)

func (t *MailTestSuite) TestPostal_Send() {
	tt := map[string]struct {
		input   *Transmission
		handler http.HandlerFunc
		url string
		marshaller func(v interface{}) ([]byte, error)
		want    interface{}
	}{
		//"Success": {
		//	Trans,
		//	func(w http.ResponseWriter, r *http.Request) {
		//		w.Write([]byte("ok"))
		//	},
		//	"",
		//	json.Marshal,
		//	Response{
		//		StatusCode: 200,
		//		Body:       "",
		//		Headers:    nil,
		//		//ID:         "",
		//		//Message:    "success",
		//	},
		//},
		"Validation Failed": {
			nil,
			nil,
			"",
			json.Marshal,
			"can't validate a nil transmission",
		},
		"Marshal Error": {
			Trans,
			nil,
			"",
			func(v interface{}) ([]byte, error) {
				return nil, fmt.Errorf("marshal error")
			},
			"marshal error",
		},
		"Bad Request": {
			Trans,
			nil,
			"@#@#$$%$",
			json.Marshal,
			"invalid URL",
		},
		"Do Error": {
			Trans,
			nil,
			"@#@#$$%$",
			json.Marshal,
			"invalid URL",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			server := httptest.NewServer(test.handler)
			defer server.Close()

			url := server.URL
			if test.url != "" {
				url = test.url
			}

			ptl := postal{
				cfg: Config{
					URL:         url,
					FromAddress: "from",
				},
				client: server.Client(),
				marshaller: test.marshaller,
			}

			resp, err := ptl.Send(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.want, resp)
		})
	}
}
