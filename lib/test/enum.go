// Package test ...
package test

//go:generate go run github.com/dmarkham/enumer@latest -type=Enum -values --transform=snake -trimprefix=Enum -json
type Enum int

const (
	EnumOne Enum = iota
	EnumTwo
	EnumThree
)
