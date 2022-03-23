package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"io/ioutil"
	"net/http"
	"time"
)

type JWT struct {
	jwt      *jwtauth.JWTAuth
	audience string
	issuer   string
}

func New(privateKeyLocation, publicKeyLocation, audience, issuer string) (*JWT, error) {
	if privateKeyLocation == "" {
		return nil, errors.New("private key location is not set in config")
	}
	if publicKeyLocation == "" {
		return nil, errors.New("public key location is not set in config")
	}

	priv, err := ioutil.ReadFile(privateKeyLocation)
	if err != nil {
		return nil, err
	}

	privPem, _ := pem.Decode(priv)
	if privPem.Type != "RSA PRIVATE KEY" && privPem.Type != "PRIVATE KEY" {
		return nil, errors.New(fmt.Sprintf("private key is of wrong type: %s", privPem.Type))
	}

	var parsedKey interface{}
	//PKCS1
	parsedKey, err = x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		//If what you are sitting on is a PKCS#8 encoded key
		parsedKey, err = x509.ParsePKCS8PrivateKey(privPem.Bytes)
		if err != nil { // note this returns type `interface{}`
			return nil, errors.New("unable to parse RSA private key")
		}
	}

	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey
	var ok bool

	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("unable to parse RSA private key")
	}

	pub, err := ioutil.ReadFile(publicKeyLocation)
	if err != nil {
		return nil, err
	}

	pubPem, _ := pem.Decode(pub)
	if pubPem.Type != "RSA PUBLIC KEY" && pubPem.Type != "PUBLIC KEY" {
		return nil, errors.New(fmt.Sprintf("public key is of wrong type: %s", pubPem.Type))
	}

	parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes)
	if err != nil {
		return nil, errors.New("unable to parse RSA public key")
	}

	publicKey, ok = parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("unable to parse RSA public key")
	}

	privateKey.PublicKey = *publicKey

	tokenAuth := jwtauth.New("RS256", privateKey, publicKey)
	j := JWT{
		jwt:      tokenAuth,
		audience: audience,
		issuer:   issuer,
	}
	return &j, nil
}

func (s *JWT) Create(claims map[string]interface{}) string {
	now := time.Now().UTC()

	expiry := now.Add(10 * time.Hour)

	allClaims := map[string]interface{}{
		jwt.AudienceKey:   s.audience,
		jwt.IssuerKey:     s.issuer,
		jwt.IssuedAtKey:   now.Unix(),
		jwt.NotBeforeKey:  now.Unix(),
		jwt.ExpirationKey: expiry.Unix(),
	}

	for k, v := range claims {
		allClaims[k] = v
	}

	_, tokenString, _ := s.jwt.Encode(allClaims)
	return tokenString
}

func (s *JWT) Get(tokenString string) (jwt.Token, error) {
	return s.jwt.Decode(tokenString)
}

func (s *JWT) Set(w http.ResponseWriter, content map[string]interface{}) {
	value := s.Create(content)
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    value,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func (s *JWT) GetVerifierMiddleware() func(http.Handler) http.Handler {
	return jwtauth.Verifier(s.jwt)
}
