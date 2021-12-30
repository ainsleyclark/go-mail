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

package mime

import (
	"bytes"
	"net/http"
)

const (
	// sniffLength is the amount of bytes to read to
	// detect the MIME type.
	sniffLength uint32 = 512
)

// DetectBuffer returns the MIME type found from the provided byte slice.
//
// The result is always a valid MIME type, with application/octet-stream
// returned when identification failed.
// Uses http.DetectContentType with a layer for SVG detection.
func DetectBuffer(buf []byte) string {
	header := make([]byte, sniffLength)
	copy(header, buf)

	// Detect for SVGs
	// See https://github.com/golang/go/issues/15888
	if bytes.Contains(header, []byte("<svg")) {
		return "image/svg+xml"
	}

	return http.DetectContentType(header)
}
