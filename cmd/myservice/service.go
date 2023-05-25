package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	taskmanagementv1 "github.com/marekbuild/annotations/gen/taskmanagement/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcAddr = "127.0.0.1:9090"
	httpAddr = "127.0.0.1:8080"
)

func main() {
	log.Printf("Starting service")

	g := &errgroup.Group{}

	log.Printf("Listening http on: %v", httpAddr)
	httpSrv, err := runHTTP(g)
	dieOnError(err, "serve http failed")

	log.Printf("Listening grpc on: %v", grpcAddr)
	grpcSrv, err := runGRPC(g)
	dieOnError(err, "serve grpc failed")

	err = g.Wait()
	dieOnError(err, "service failed")

	err = httpSrv.Shutdown(context.Background())
	dieOnError(err, "shutdown http failed")

	grpcSrv.Stop()
}

func runHTTP(g *errgroup.Group) (*http.Server, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := taskmanagementv1.RegisterTaskManagementServiceHandlerFromEndpoint(context.Background(), mux, grpcAddr, opts)
	if err != nil {
		return nil, err
	}
	s := &http.Server{Addr: httpAddr, Handler: mux}
	g.Go(s.ListenAndServe)
	return s, nil
}

func runGRPC(g *errgroup.Group) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, err
	}
	s := grpc.NewServer()
	taskmanagementv1.RegisterTaskManagementServiceServer(s, &service{})
	g.Go(func() error { return s.Serve(lis) })
	return s, nil
}

func dieOnError(err error, context string, vs ...any) {
	if err != nil {
		log.Printf("error: %v: %v", fmt.Sprintf(context, vs...), err)
		os.Exit(1)
	}
}

type service struct {
	taskmanagementv1.UnimplementedTaskManagementServiceServer
}

func (s *service) Health(ctx context.Context, req *taskmanagementv1.HealthRequest) (*taskmanagementv1.HealthResponse, error) {
	return &taskmanagementv1.HealthResponse{}, nil
}
