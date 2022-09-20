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

import "net/http"

// Response represents the data passed back from a successful transmission.
type Response struct {
	StatusCode int         // e.g. 200
	Body       []byte      // e.g. {"result: success"}
	Headers    http.Header // e.g. map[X-Ratelimit-Limit:[600]]
	ID         string      // e.g "100"
	Message    string      // e.g "Email sent successfully"
}
