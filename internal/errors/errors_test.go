package errors

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError_Error(t *testing.T) {
	tt := map[string]struct {
		input *Error
		want  string
	}{
		"Normal": {
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			"op: err",
		},
		"Nil Operation": {
			&Error{Code: INTERNAL, Message: "test", Operation: "", Err: fmt.Errorf("err")},
			"err",
		},
		"Nil Err": {
			&Error{Code: INTERNAL, Message: "test", Operation: "", Err: nil},
			"<internal> test",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.input.Error(), test.want)
		})
	}
}

func TestError_Code(t *testing.T) {
	tt := map[string]struct {
		input error
		want  string
	}{
		"Normal": {
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			"internal",
		},
		"Nil Input": {
			nil,
			"",
		},
		"Nil Code": {
			&Error{Code: "", Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			"internal",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, Code(test.input))
		})
	}
}

func Test_Message(t *testing.T) {
	tt := map[string]struct {
		input error
		want  string
	}{
		"Normal": {
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			"test",
		},
		"Nil Input": {
			nil,
			"",
		},
		"Nil Message": {
			&Error{Code: "", Message: "", Operation: "op", Err: fmt.Errorf("err")},
			GlobalError,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, Message(test.input))
		})
	}
}

func TestError_ToError(t *testing.T) {
	tt := map[string]struct {
		input interface{}
		want  *Error
	}{
		"Pointer": {
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
		},
		"Non Pointer": {
			Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
		},
		"Error": {
			fmt.Errorf("err"),
			&Error{Err: fmt.Errorf("err")},
		},
		"String": {
			"err",
			&Error{Err: fmt.Errorf("err")},
		},
		"Default": {
			nil,
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, ToError(test.input))
		})
	}
}

func TestNew(t *testing.T) {
	want := fmt.Errorf("error")
	got := New("error")
	assert.Errorf(t, want, got.Error())
}
