package main

import (
	"log"
	"net/http"

	"github.com/demo/contracts/gen/go/order/v1/orderv1connect"
	"github.com/demo/gateway/internal/domain/orders"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	orderClient := orderv1connect.NewOrderServiceClient(
		http.DefaultClient,
		"http://localhost:8081",
	)

	orderServer := orders.NewServer(orderClient)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	path, handler := orderv1connect.NewOrderServiceHandler(orderServer)
	mux.Handle(path, handler)

	addr := ":8080"
	log.Printf("Gateway listening on %s", addr)

	err := http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{}))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
