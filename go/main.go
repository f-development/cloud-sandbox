package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	err1 := fmt.Errorf("err 1")
	err2 := errors.Wrapf(err1, "err 2")
	err3 := errors.Wrapf(err2, "err 3")
	fmt.Printf("%+v\n", err1)
	fmt.Printf("%+v\n", err2)
	fmt.Printf("%+v\n", err3)
}
