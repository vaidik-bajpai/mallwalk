package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/vaidik-bajpai/mallwalk/common/discovery"
	"github.com/vaidik-bajpai/mallwalk/common/discovery/consul"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
)

type config struct {
	serviceName string
	serviceAddr string
	consulAddr  string
	dbDSN       string
}

func main() {
	var cfg config
	cfg.serviceName = "stock"
	flag.StringVar(&cfg.serviceAddr, "service-addr", "localhost:4040", "stocks service address")
	flag.StringVar(&cfg.consulAddr, "consul-addr", "localhost:8500", "consul address")
	flag.StringVar(&cfg.dbDSN, "DB-DSN", "mongodb://localhost:27017", "database dsn")
	flag.Parse()

	registry, err := consul.NewRegistry(cfg.consulAddr)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(cfg.serviceName)
	err = registry.Register(ctx, instanceID, cfg.serviceName, cfg.serviceAddr)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err = registry.HealthCheck(instanceID, cfg.serviceName)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.DeRegister(ctx, instanceID, cfg.serviceName)

	mongoClient, err := connectToMongoDB(cfg.dbDSN)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", cfg.serviceAddr)
	if err != nil {
		panic(err)
	}

	store := NewStore(mongoClient)
	svc := NewService(store)

	NewGRPCHandler(grpcServer, svc)

	fmt.Printf("Starting the %s service at the port %s\n", cfg.serviceName, cfg.serviceAddr)
	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	return client, err
}
