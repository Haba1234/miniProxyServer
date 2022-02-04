package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	proxyserver "github.com/Haba1234/miniProxyServer/internal/proxy-server"
)

func main() {
	h := proxyserver.NewService()

	mux := http.NewServeMux()
	mux.HandleFunc("/proxy", h.RedirectReq)

	logger := proxyserver.NewLogger(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: logger,
	}

	log.Println("server start on port 8080")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	log.Fatal(err)
}
