package contract

import "os"

const (
	Success = iota
	NoApiKey
	NoCred
	NoBot
	)

func Finish(code int) {
	os.Exit(code)
}