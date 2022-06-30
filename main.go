package main

import (
	"http-spy-forward/cmd"
	"os"
	"runtime"

	"github.com/urfave/cli"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// ./http-spy-forward start -i en0 -f "dst host 192.168.1.100" -u "http://127.0.0.1:8080"
	app := cli.NewApp()
	app.Name = "http-spy-forward"
	app.Author = "M1ngkvv1nd"
	app.Version = "2022/5/13"
	app.Usage = "Packet forward tool, capture HTTP packet and forward request to target URL"
	app.Commands = []cli.Command{cmd.Start}
	app.Flags = append(app.Flags, cmd.Start.Flags...)
	_ = app.Run(os.Args)
}
