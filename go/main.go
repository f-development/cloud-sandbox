package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type error1 struct {
}

func (*error1) Error() string {
	return "error 1"
}

func Error1() error {
	return &error1{}
}

func main() {
	err1 := errors.WithStack(fmt.Errorf("err 1"))
	err2 := errors.Wrapf(err1, "")
	err3 := errors.Wrapf(err2, "err 3")
	fmt.Printf("%+v\n", err1)
	fmt.Printf("%+v\n", err2)
	fmt.Printf("%+v\n", err3)
}
