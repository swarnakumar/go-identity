package grpc

import (
	"context"

	"github.com/swarnakumar/go-identity/config"
	"github.com/swarnakumar/go-identity/db"
	"github.com/swarnakumar/go-identity/web/jwt"

	userspb "github.com/swarnakumar/go-identity/proto/users/v1"
)

type server struct {
	db  *db.Client
	jwt *jwt.JWT

	userspb.UnimplementedUsersServiceServer
}

func NewServer(ctx context.Context) *server {
	j, _ := jwt.New(config.JWTPrivateKeyFile, config.JWTPublicKeyFile, config.JWTAudience, config.JWTIssuer)
	client := db.New(ctx)

	n := &server{jwt: j, db: client}
	return n
}
