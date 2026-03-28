package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/vishal2098govind/service/foundations/keystore"
)

func main() {

	ks := keystore.New()
	ks.LoadKeys(os.DirFS("zarf/keys/"))

	// currently active kid is usually fetched from a central storage
	const kid = "8234c8c5-0508-4301-bfdf-da1d515399a1"

	token, err := GenToken(ks, kid)
	if err != nil {
		log.Fatalf("failed to gen key pair: %v", err)
	}
	err = ValidateTokenOPA(ks, token)
	if err != nil {
		log.Fatalf("failed to validate token: %v", err)
	}
	err = ValidateTokenParse(ks, token)
	if err != nil {
		log.Fatalf("failed to validate token parse: %v", err)
	}
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Roles []string
}

func GenToken(ks keystore.KeyStore, kid string) (string, error) {

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "vishal-user-id",
			Issuer:    "services-go service",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{
			"ADMIN",
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(jwt.SigningMethodRS256.Name), claims)
	token.Header["kid"] = kid

	privatePem, err := ks.PrivateKeyPEM(kid)
	if err != nil {
		return "", fmt.Errorf("private key not found: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privatePem))
	if err != nil {
		return "", fmt.Errorf("parsing private pem bytes: %w", err)
	}

	tokenStr, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	fmt.Println("-------------TOKEN----------")
	fmt.Println(tokenStr)
	fmt.Println("----------------------------")

	return tokenStr, nil
}

//go:embed rego/authentication.rego
var opaAuthenticationPolicy string

func OpaEvaluation(input any) error {
	query := fmt.Sprintf("x = data.%s.%s", "service.rego", "auth")
	// fmt.Println(opaAuthenticationPolicy)
	ctx := context.Background()
	q, err := rego.New(
		rego.Query(query),
		rego.Module("policy.rego", opaAuthenticationPolicy),
	).PrepareForEval(ctx)
	if err != nil {
		return fmt.Errorf("prepare rego query: %w", err)
	}

	results, err := q.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("rego eval: %w", err)
	}

	if len(results) == 0 {
		return errors.New("no results")
	}

	result, ok := results[0].Bindings["x"].(bool)
	if !ok {
		return fmt.Errorf("invalid rego result: %w", err)
	}

	if !result {
		fmt.Println(results[0].Bindings)
		return fmt.Errorf("invalid token")
	}

	return nil
}

func ValidateTokenOPA(ks keystore.KeyStore, tokenStr string) error {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))
	var claims CustomClaims
	token, _, err := parser.ParseUnverified(tokenStr, &claims)
	if err != nil {
		return fmt.Errorf("parse token string: %w", err)
	}
	kid, ok := token.Header["kid"]
	if !ok {
		return fmt.Errorf("invalid token (invalid kid)")
	}

	kidStr, ok := kid.(string)
	if !ok {
		return fmt.Errorf("invalid token (kid not string)")
	}

	publicKey, err := ks.PublicKeyPEM(kidStr)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	input := map[string]any{
		"Token": tokenStr,
		"Key":   publicKey,
		"ISS":   "services-go service",
	}

	err = OpaEvaluation(input)
	if err != nil {
		return fmt.Errorf("eval opa policy: %w", err)
	}

	fmt.Println("token is valid")

	return nil
}

func ValidateTokenParse(ks keystore.KeyStore, token string) error {

	var claims CustomClaims
	t, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid kid")
		}

		publicKeyPem, err := ks.PublicKeyPEM(kid)
		if err != nil {
			return nil, fmt.Errorf("reading public key pem file: %w", err)
		}

		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPem))
		if err != nil {
			return nil, fmt.Errorf("parse publick key bytes: %w", err)
		}

		return pubKey, nil
	}, jwt.WithValidMethods([]string{
		jwt.SigningMethodRS256.Name,
	}))
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return err
	}
	if !t.Valid {
		return fmt.Errorf("invalid token")
	}
	if err := claims.Valid(); err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	fmt.Println(claims.Roles)
	fmt.Println(claims.Subject)

	return nil
}

func GenKey() error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to rsa.GenerateKey: %w", err)
	}

	privatePem, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("failed to os.Create private.pem file: %w", err)
	}
	defer privatePem.Close()

	if err = pem.Encode(privatePem, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return fmt.Errorf("failed to write to private.pem file: %w", err)
	}

	publicKey := privateKey.PublicKey
	publicPem, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("failed to create os.Create public.pem file: %w", err)
	}
	defer publicPem.Close()

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return fmt.Errorf("marshal public key: %w", err)
	}
	if err = pem.Encode(publicPem, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}); err != nil {
		return fmt.Errorf("failed to write to public.pem file: %w", err)
	}

	return nil
}
