package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	discovery "github.com/vaidik-bajpai/mallwalk/common/discovery"
	"github.com/vaidik-bajpai/mallwalk/common/discovery/consul"
	"github.com/vaidik-bajpai/mallwalk/gateway/internal/gateway"
	"github.com/vaidik-bajpai/mallwalk/gateway/internal/handler"
)

type config struct {
	addr        string
	consulAddr  string
	serviceName string
}

func main() {
	var cfg config
	cfg.serviceName = "gateway"
	flag.StringVar(&cfg.addr, "http-addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.consulAddr, "consul-addr", "localhost:8500", "consul network address")
	flag.Parse()
	registry, err := consul.NewRegistry(cfg.consulAddr)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(cfg.serviceName)
	if err := registry.Register(ctx, instanceID, cfg.serviceName, cfg.addr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, cfg.serviceName); err != nil {
				log.Fatal("Failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.DeRegister(ctx, instanceID, cfg.serviceName)

	mux := http.NewServeMux()

	gateway := gateway.NewGRPCGateway(registry)
	validate := validator.New()

	handler := handler.NewHandler(gateway, validate)
	handler.RegisterRoutes(mux)

	log.Printf("Starting HTTP server at %s", cfg.addr)

	http.ListenAndServe(cfg.addr, mux)
}
