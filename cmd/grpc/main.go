package main

import (
	"github.com/swarnakumar/go-identity/config"
	"github.com/swarnakumar/go-identity/grpc"
)

func main() {
	grpc.StartGRPCServer(config.GRPCServerPort, config.APIServerPort)
}
