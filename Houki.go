package main

import (
	"net/http"
	"os"

	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/proxy"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	_ = os.Mkdir(".certificate", 0644)

	p, err := proxy.New()
	if err != nil {
		log.Fatal("Failed to create proxy: %v", err)
	}

	err = http.ListenAndServe(":8080", p)
	if err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
