package main

import (
	"flag"
	"github.com/swarnakumar/go-identity/web"
)

func main() {
	portPtr := flag.Int("port", 8000, "Server Port?")
	flag.Parse()

	web.StartServer(*portPtr)
}
