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

package httputil

// A Request represents an HTTP request received by a server
// or to be sent by a client. It is an extension of the
// std http.Request for Go Mail.
type Request struct {
	Method            string
	URL               string
	Headers           map[string]string
	BasicAuthUser     string
	BasicAuthPassword string
}

// NewHTTPRequest returns a new Request given a method and URL.
func NewHTTPRequest(method, url string) *Request {
	return &Request{
		Method: method,
		URL:    url,
	}
}

// AddHeader adds the key, value pair to the request headers.
func (r *Request) AddHeader(name, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[name] = value
}

// SetBasicAuth sets the request's Authorization header to use HTTP
// Basic Authentication with the provided username and password.
//
// With HTTP Basic Authentication the provided username and password
// are not encrypted.
func (r *Request) SetBasicAuth(user, password string) {
	r.BasicAuthUser = user
	r.BasicAuthPassword = password
}
