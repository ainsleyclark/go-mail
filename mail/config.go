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

// Config represents the configuration passed when a new
// client is constructed. FromAddress, FromName and
// an APIKey are all required to create a new
// client.
type Config struct {
	URL         string
	APIKey      string
	Domain      string
	FromAddress string
	FromName    string
	Password    string
	Port        int
}

// Validate runs sanity checks of a Config struct.
// This is run before a new client is created
// to ensure there are no invalid API
// calls.
func (c *Config) Validate() error {
	if c.FromAddress == "" {
		return errors.New("mailer requires from address")
	}

	if c.FromName == "" {
		return errors.New("mailer requires from name")
	}

	if c.APIKey == "" {
		return errors.New("mailer requires api key")
	}

	return nil
}
