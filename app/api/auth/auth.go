package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	Issuer    string
	KeyLookup KeyLookup
}

type Auth struct {
	signingMethod jwt.SigningMethod
	issuer        string
	jwtParser     *jwt.Parser
	keyLookup     KeyLookup
}

type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

type KeyLookup interface {
	PrivateKeyPEM(kid string) (key string, err error)
	PublicKeyPEM(kid string) (key string, err error)
}

func New(cfg Config) (*Auth, error) {
	return &Auth{
		keyLookup:     cfg.KeyLookup,
		issuer:        cfg.Issuer,
		signingMethod: jwt.GetSigningMethod(jwt.SigningMethodRS256.Name),
		jwtParser:     jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
	}, nil
}

func (a *Auth) GenerateToken(kid string, claims Claims) (string, error) {

	privateKeyPemStr, err := a.keyLookup.PrivateKeyPEM(kid)
	if err != nil {
		return "", fmt.Errorf("failed to get private key for this kid[%s] : %w", kid, err)
	}

	privateKeyPem, _ := pem.Decode([]byte(privateKeyPemStr))
	if privateKeyPem == nil {
		return "", fmt.Errorf("failed to decode private key pem string. invalid pem string")
	}

	var pvtKeyParsed any
	pvtKeyParsed, err = x509.ParsePKCS1PrivateKey(privateKeyPem.Bytes)
	if err != nil {
		pvtKeyParsed, err = x509.ParsePKCS8PrivateKey(privateKeyPem.Bytes)
		if err != nil {
			return "", fmt.Errorf("failed to parse private key pem: %w", err)
		}
	}

	privateKey, ok := pvtKeyParsed.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("invalid private key. should be of type rsa.PrivateKey")
	}

	token := jwt.NewWithClaims(a.signingMethod, claims)

	token.Header["kid"] = kid

	tokenStr, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenStr, nil
}

// conciously staying at go-layer or service layer to do authentication, instead of delegating it to OPA
// the authorization is delegated to OPA though, instead of keeping in the go-layer
//
// the authorization logic is bit more complex and is more or less going to be identical across services
// thus, I decided to make authorization language agnostic by leveraging OPA,
// and doing simple token validation (authentication) in the go-layer instead of relying on OPA
func (a *Auth) Authenticate(token string) (Claims, error) {
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return Claims{}, fmt.Errorf("invalid token. must be of format: Bearer <token>")
	}

	var claims Claims
	t, err := a.jwtParser.ParseWithClaims(parts[1], &claims, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token. kid either not found or is of invalid format. should be a string")
		}

		publicKeyPemStr, err := a.keyLookup.PublicKeyPEM(kid)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key for the given kid[%s]: %w", kid, err)
		}

		publicKeyPem, _ := pem.Decode([]byte(publicKeyPemStr))
		if publicKeyPem == nil {
			return nil, fmt.Errorf("failed to decode public key pem str to a pem block. invalid public key pem string.")
		}

		publicKeyParsed, err := x509.ParsePKIXPublicKey(publicKeyPem.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key pem bytes to rsa.PublicKey")
		}

		publicKey, ok := publicKeyParsed.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("parsed public key is not of type *rsa.PublicKey")
		}

		return publicKey, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if !t.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	}

	if err := claims.Valid(); err != nil {
		return Claims{}, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func (a *Auth) Authorize(claims Claims) bool {
	return true
}
