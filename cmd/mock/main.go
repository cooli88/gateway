package main

import (
	"log"
	"net/http"

	"github.com/demo/contracts/gen/go/order/v1/orderv1connect"
	"github.com/demo/gateway/test/mock"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	orderService := mock.NewOrderService()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	path, handler := orderv1connect.NewOrderServiceHandler(orderService)
	mux.Handle(path, loggingMiddleware(handler))

	addr := ":8081"
	log.Printf("Mock Order Service listening on %s", addr)

	err := http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{}))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
