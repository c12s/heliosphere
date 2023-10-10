package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"rate-limiter-service/config"
	"rate-limiter-service/handlers"
	pb "rate-limiter-service/proto/ratelimiter"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.GetConfig()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)

	// Bootstrap gRPC server.
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Bootstrap gRPC service server and respond to request.
	rateLimiterHandler, err := handlers.NewRateLimiterHandler()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	pb.RegisterRateLimitServiceServer(grpcServer, rateLimiterHandler)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	grpcServer.Stop()
}
