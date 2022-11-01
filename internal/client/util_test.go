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

package client

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIs2XX(t *testing.T) {
	tt := map[string]struct {
		input int
		want  bool
	}{
		"< 200": {
			http.StatusContinue,
			false,
		},
		"200": {
			http.StatusOK,
			true,
		},
		"300 >": {
			http.StatusMultipleChoices,
			false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := Is2XX(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
