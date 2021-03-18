package main

import (
	"net/http"

	"github.com/wuhan005/Houki/ca"
	"github.com/wuhan005/Houki/module"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/proxy"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	module.Init()
	module.Load()

	ca.Init()

	p, err := proxy.New()
	if err != nil {
		log.Fatal("Failed to create proxy: %v", err)
	}

	err = http.ListenAndServe(":8080", p)
	if err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
