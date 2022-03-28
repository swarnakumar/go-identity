package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/swarnakumar/go-identity/db/sql/sqlc"
)

type authUser string

const authUserKey authUser = "user"

func (s *server) GetUser(ctx context.Context) *sqlc.GetUserByEmailRow {

	user := ctx.Value(authUserKey)
	if user == nil {
		return nil
	}

	u := user.(sqlc.GetUserByEmailRow)
	return &u
}

type validToken string

const validTokenKey validToken = "valid-token"

type tokenValidParams struct {
	present bool
	valid   bool
}

func (s *server) IsValidToken(ctx context.Context) tokenValidParams {
	return ctx.Value(validTokenKey).(tokenValidParams)
}

func (s *server) getAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		tokenP := tokenValidParams{
			present: false,
			valid:   false,
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			newCtx := context.WithValue(ctx, validTokenKey, tokenP)
			return handler(newCtx, req)
		}

		auth := md["authorization"]
		if len(auth) < 1 {
			newCtx := context.WithValue(ctx, validTokenKey, tokenP)
			return handler(newCtx, req)
		}
		token := strings.TrimPrefix(auth[0], "Bearer ")
		t, err := s.jwt.Get(token)
		tokenP.present = true
		if err != nil || t == nil {

			newCtx := context.WithValue(ctx, validTokenKey, tokenP)
			return handler(newCtx, req)
		}

		user, _ := t.PrivateClaims()["user"]
		if user == "" || user == nil {
			newCtx := context.WithValue(ctx, validTokenKey, tokenP)
			return handler(newCtx, req)
		}

		tokenP.valid = true

		u, err := s.db.Users.GetByEmail(ctx, user.(string))
		newCtx := context.WithValue(ctx, authUserKey, *u)
		newCtx = context.WithValue(newCtx, validTokenKey, tokenP)

		return handler(newCtx, req)

	}
}
