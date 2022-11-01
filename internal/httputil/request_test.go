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

package httputil

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewHTTPRequest(t *testing.T) {
	req := NewHTTPRequest(http.MethodGet, "https://gomail.example.com")
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "https://gomail.example.com", req.URL)
}

func TestRequest_AddHeader(t *testing.T) {
	req := NewHTTPRequest(http.MethodGet, "https://gomail.example.com")
	req.AddHeader("header", "value")
	want := map[string]string{"header": "value"}
	assert.Equal(t, want, req.Headers)
}

func TestRequest_SetBasicAuth(t *testing.T) {
	req := NewHTTPRequest(http.MethodGet, "https://gomail.example.com")
	req.SetBasicAuth("user", "pass")
	assert.Equal(t, "user", req.BasicAuthUser)
	assert.Equal(t, "pass", req.BasicAuthPassword)
}
