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

package mail

import (
	"errors"
	"net/http"
)

// Config represents the configuration passed when a new
// client is constructed. Dependant on what driver is used,
// different options are required to be present.
type Config struct {
	URL         string
	APIKey      string
	Domain      string
	FromAddress string
	FromName    string
	Password    string
	Port        int
	Client      *http.Client
}

// Validate runs sanity checks of a Config struct.
// This is run before a new client is created
// to ensure there are no invalid API
// calls.
func (c *Config) Validate() error {
	if c.FromAddress == "" {
		return errors.New("driver requires from address")
	}
	if c.FromName == "" {
		return errors.New("driver requires from name")
	}
	if c.APIKey == "" {
		return errors.New("driver requires api key")
	}
	return nil
}
