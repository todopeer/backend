package internallog

import (
	"log"
)

type Logger interface {
	Println(...any)
	Printf(string, ...any)

	Fatalln(...any)
	Fatalf(string, ...any)
}

type NoPrintLogger struct{}

func (nl *NoPrintLogger) Println(...any)        {}
func (nl *NoPrintLogger) Printf(string, ...any) {}

func (nl *NoPrintLogger) Fatalln(args ...any) {
	log.Fatalln(args...)
}
func (nl *NoPrintLogger) Fatalf(msg string, args ...any) {
	log.Fatalf(msg, args...)
}
