package config

import (
	"fmt"
	"os"
	"strings"
)

var (
	pgUser     = os.Getenv("PG_USER")
	pgPassword = os.Getenv("PG_PASSWORD")
	pgHost     = os.Getenv("PG_HOST")
	pgPort     = os.Getenv("PG_PORT")
	pgDatabase = os.Getenv("PG_DATABASE")
)

var PgConnStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
	pgUser, pgPassword, pgHost, pgPort, pgDatabase)

var (
	CookieHashKey  = os.Getenv("COOKIE_HASH_KEY")
	CookieBlockKey = os.Getenv("COOKIE_BLOCK_KEY")
)

var CSRFAuthKey = os.Getenv("CSRF_AUTH_KEY")

var ENV = os.Getenv("ENV")

var UseHttps = strings.HasPrefix(strings.ToLower(ENV), "p")

var (
	JWTPrivateKeyFile = os.Getenv("JWT_PRIVATE_KEY")
	JWTPublicKeyFile  = os.Getenv("JWT_PUBLIC_KEY")
	JWTAudience       = os.Getenv("JWT_AUDIENCE")
	JWTIssuer         = os.Getenv("JWT_ISSUER")
)

var (
	GRPCServerPort = os.Getenv("GRPC_SERVER_PORT")
	APIServerPort  = os.Getenv("API_SERVER_PORT")
)
