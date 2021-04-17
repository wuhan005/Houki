package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/ca"
	"github.com/wuhan005/Houki/internal/cmd"
	"github.com/wuhan005/Houki/internal/db"
	"github.com/wuhan005/Houki/internal/module"
	"github.com/wuhan005/Houki/internal/proxy"
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

	if err := db.Initialize(); err != nil {
		log.Fatal("Failed to initialize database: %v", err)
	}

	if err := ca.Initialize(); err != nil {
		log.Fatal("Failed to initialize CA: %v", err)
	}

	if _, err = proxy.Initialize(); err != nil {
		log.Fatal("Failed to create proxy: %v", err)
	}

	if err := module.Initialize(); err != nil {
		log.Fatal("Failed to initialize modules: %v", err)
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
	log.Info("Bye!")
}
