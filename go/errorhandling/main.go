package main

import (
	"log/slog"
	"os"

	"github.com/pkg/errors"
)

var (
	slogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
)

type ErrType struct {
	Msg string
}

func (e *ErrType) Error() string {
	return e.Msg
}

func CreateError() error {
	return &ErrType{
		Msg: "msg",
	}
}

func temp() interface{} {
	err := CreateError()
	return &err
}

func main() {
	err1 := errors.Errorf("%s", "error f")
	err2 := errors.Wrap(err1, "err 2")
	err3 := errors.Wrapf(err2, "err 3")
	// fmt.Printf("%+v\n", err1)
	// fmt.Printf("%+v\n", err2)
	// fmt.Printf("%+v\n", err3)
	// fmt.Printf("%+v\n", errors.Unwrap(err3))
	// fmt.Printf("%+v\n", (err3))

	slogger.Debug("hey", "err", err3, "?", temp())
}
