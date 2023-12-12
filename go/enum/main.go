package main

type SomeEnum string

const (
	SomeVal1 SomeEnum = "1"
	SomeVal2 SomeEnum = "2"
)

type SomeType interface {
	SomeToString()
}

func (e SomeEnum) SomeToString() {
}

func F(s SomeType) {
}

func main() {
	F(SomeVal1)

}
