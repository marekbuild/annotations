package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	taskmanagementv1 "github.com/marekbuild/annotations/gen/taskmanagement/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcAddr = "127.0.0.1:9090"
	httpAddr = "127.0.0.1:8080"
)

func main() {
	log.Printf("Starting myservice")

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := taskmanagementv1.RegisterTaskManagementServiceHandlerFromEndpoint(context.Background(), mux, grpcAddr, opts)
	dieOnError(err, "failed to register task management handler")

	log.Printf("Listening grpc: %v", grpcAddr)
	log.Printf("Listening http: %v", httpAddr)

	err = http.ListenAndServe(httpAddr, mux)
	if err != http.ErrServerClosed {
		dieOnError(err, "failed to listen or shutdown")
	}
}

func dieOnError(err error, context string, vs ...any) {
	if err != nil {
		log.Printf("error: %v: %v", fmt.Sprintf(context, vs...), err)
		os.Exit(1)
	}
}
