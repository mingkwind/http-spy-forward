package main

import (
	"http-spy-forward/cmd"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	cmd.Execute()
}
