package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/swarnakumar/go-identity/web/server"
)

func StartServer(port int) {
	ctx := context.Background()

	s := server.New(ctx)
	defer s.Close()

	s.InitRouter()
	setupRoutes(s)

	// Bind to a port and pass our router in
	l := s.GetLogger()
	l.Infow("Starting Server", "Port", port)
	l.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), s.GetRouter()))
}
