package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/cmd"
)

var (
	Version = "development"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "Houki"
	app.Usage = "Customizable MitM proxy"
	app.Version = Version
	app.Commands = []*cli.Command{
		cmd.Web,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("Failed to start application: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
