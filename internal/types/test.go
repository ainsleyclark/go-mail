package types

import (
	"fmt"
	"github.com/ainsleyclark/go-mail"
)

type Test struct {
	Hello string
}

func (t *Test) Hey(tx mail.Transmission) {
	fmt.Println(t.Hello)
}
