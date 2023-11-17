package main

import (
	"goapps/simplecli"
)

func main() {
	opts := simplecli.ParseFlags()
	simplecli.Init(opts)
}
