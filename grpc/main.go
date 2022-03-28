package grpc

import (
	"context"

	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userspb "github.com/swarnakumar/go-identity/proto/users/v1"
)

func StartGRPCServer(grpcServerPort, apiServerPort string) {
	ctx := context.Background()
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":"+grpcServerPort)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	n := NewServer(ctx)
	defer n.db.Close()

	// Create a gRPC server object
	authInterceptor := n.getAuthInterceptor()
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor),
	}

	s := grpc.NewServer(opts...)
	// Attach the Greeter service to the server
	userspb.RegisterUsersServiceServer(s, n)
	// Serve gRPC Server in a go-routine. So that this doesnt block whatever follows.
	log.Println("Serving gRPC on 0.0.0.0:" + grpcServerPort)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		ctx,
		"0.0.0.0:"+grpcServerPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	// This is the API gateway
	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = userspb.RegisterUsersServiceHandler(ctx, gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":" + apiServerPort,
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:" + apiServerPort)
	log.Fatalln(gwServer.ListenAndServe())
}
