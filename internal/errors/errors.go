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

package errors

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/flightaware/go-mail/mail"
)

// Application error codes.
const (
	// CONFLICT - An action cannot be performed.
	CONFLICT = "conflict"
	// INTERNAL - Error within Go Mail
	INTERNAL = "internal" // Internal error
	// INVALID - Validation failed
	INVALID = "invalid" // Validation failed
	// API - Error in the http request.
	API = "api"
	// Prefix is the string prefixed to an error message.
	Prefix = "go-mail"
	// GlobalError is a general message when no error message
	// has been found.
	GlobalError = "An error has occurred."
)

// Error defines a standard application error.
type Error struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
	Err       error  `json:"error"`
}

// Error returns the string representation of the error
// message.
func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the prefix
	fmt.Fprintf(&buf, "%s: ", Prefix)

	// Print the current operation in our stack, if any.
	if e.Operation != "" && mail.Debug {
		fmt.Fprintf(&buf, "%s: ", e.Operation)
	}

	// If wrapping an error, print its Error() message.
	// Otherwise, print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			_, _ = fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(e.Message)
	}

	return buf.String()
}

// Code returns the code of the root error, if available.
// Otherwise, returns INTERNAL.
func Code(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Code != "" {
		return e.Code
	} else if ok && e.Err != nil {
		return Code(e.Err)
	}
	return INTERNAL
}

// Message returns the human-readable message of the error,
// if available. Otherwise, returns a generic error
// message.
func Message(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return Message(e.Err)
	}
	return GlobalError
}

// ToError Returns a Go Mail error from input. If The type
// is not of type Error, nil will be returned.
func ToError(err interface{}) *Error {
	switch v := err.(type) {
	case *Error:
		return v
	case Error:
		return &v
	case error:
		return &Error{Err: fmt.Errorf(v.Error())}
	case string:
		return &Error{Err: fmt.Errorf(v)}
	default:
		return nil
	}
}

// New is a wrapper for the stdlib new function.
func New(text string) error {
	return errors.New(text)
}
