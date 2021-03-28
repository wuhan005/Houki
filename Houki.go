package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/ca"
	"github.com/wuhan005/Houki/internal/conf"
	"github.com/wuhan005/Houki/internal/module"
	"github.com/wuhan005/Houki/internal/proxy"
	"github.com/wuhan005/Houki/internal/route"
	"github.com/wuhan005/Houki/internal/sse"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	port := flag.String("port", ":8000", "Web panel port")
	flag.Parse()

	err = conf.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize config: %v", err)
	}

	err = ca.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize CA: %v", err)
	}

	err = module.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize modules: %v", err)
	}
	module.Load()

	sse.Initialize()

	_, err = proxy.Initialize()
	if err != nil {
		log.Fatal("Failed to create proxy: %v", err)
	}
	proxy.Start()

	w := route.New()
	if err != nil {
		log.Fatal("Failed to create web server: %v", err)
	}
	w.Run(*port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Info("Bye!")
}
