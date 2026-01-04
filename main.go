package main

import (
	"fmt"
	"net/http"

	"github.com/le0developer/traefik-waf/internal"
)

func main() {
	config := internal.NewConfigFromEnv()
	instance, err := internal.New(config)
	if err != nil {
		panic(err)
	}

	httpMux := instance.Mux()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpMux,
	}

	fmt.Println("Starting server on :8080")
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}

}
