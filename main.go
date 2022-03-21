package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/swarnakumar/go-identity/cmd/cli"
	//"github.com/swarnakumar/go-identity/web"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No Command Passed!")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "run-server":
		runServer := flag.NewFlagSet("run-server", flag.ExitOnError)
		//serverPortPtr := runServer.Int("port", 8000, "port to run the server on")

		runServer.Parse(os.Args[2:])
		//web.StartServer(*serverPortPtr)
	case "create-user":
		cli.CreateFromCli()
	default:
		fmt.Println("Unsupported Command! Supported Commands: run-server, create-user")
		os.Exit(1)

	}
}
