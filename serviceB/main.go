package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "mins_EduCenter/proto/reference"

	"google.golang.org/grpc"
)

func main() {
	// Инициализируем сервер с тестовыми данными
	srv := NewReferenceServer()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(TraceInterceptor),
	)
	pb.RegisterReferenceServiceServer(grpcServer, srv)

	log.Println("Reference Service listening on :50051")

	// Graceful shutdown
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Reference Service...")
	grpcServer.GracefulStop()
}
