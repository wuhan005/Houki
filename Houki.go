package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/wuhan005/Houki/internal/ca"
	"github.com/wuhan005/Houki/internal/conf"
	"github.com/wuhan005/Houki/internal/module"
	"github.com/wuhan005/Houki/internal/proxy"
	"github.com/wuhan005/Houki/internal/web"
	log "unknwon.dev/clog/v2"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	proxyPort := flag.String("proxy-port", ":8080", "Proxy port")
	webPort := flag.String("web-port", ":8000", "Web panel port")
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

	p, err := proxy.New()
	if err != nil {
		log.Fatal("Failed to create proxy: %v", err)
	}
	p.Run(*proxyPort)

	w := web.New()
	if err != nil {
		log.Fatal("Failed to create web server: %v", err)
	}
	w.Run(*webPort)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Info("Bye!")
}
