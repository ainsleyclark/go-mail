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

import "github.com/ainsleyclark/go-mail/internal/httputil"

var (
	// newJSONData is an alias for httputil.NewJSONData
	// for creating JSON payloads.
	newJSONData = httputil.NewJSONData
	// formDataFn is an alias for httputil.NewFormData
	// for creating form data payloads.
	newFormData = httputil.NewFormData
)
