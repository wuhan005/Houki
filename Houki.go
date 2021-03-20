package main

import (
	"net/http"

	"github.com/wuhan005/Houki/internal/ca"
	"github.com/wuhan005/Houki/internal/conf"
	"github.com/wuhan005/Houki/internal/module"
	"github.com/wuhan005/Houki/internal/proxy"
	log "unknwon.dev/clog/v2"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

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

	err = http.ListenAndServe(":8080", p)
	if err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
