package auth_test

import (
	"bytes"
	"context"
	"fmt"
	"runtime/debug"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/foundations/logger"
)

func Test_Auth(t *testing.T) {
	_, teardown := newUnit(t)
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}

		teardown()
	}()

	var ks keyStore
	a, err := auth.New(auth.Config{
		Issuer:    "test-issuer",
		KeyLookup: &ks,
	})
	if err != nil {
		t.Fatalf("should be able to create new auth object. got error: %v", err)
	}

	t.Run("case - 1", case1(a))
}

func case1(ath *auth.Auth) func(t *testing.T) {
	return func(t *testing.T) {
		claims := auth.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "test-issuer",
				Subject:   "test-user-id",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
			Roles: []string{"ADMIN"},
		}

		token, err := ath.GenerateToken(kid, claims)
		if err != nil {
			t.Fatalf("should have generated token successfully. got error: %s", err)
		}
		parsedClaims, err := ath.Authenticate(fmt.Sprintf("Bearer %s", token))
		if err != nil {
			t.Fatalf("should have authenticated successfully. go error: %s", err)
		}
		if claims.Issuer != parsedClaims.Issuer {
			t.Fatal("issuer mismatch")
		}
		if claims.Subject != parsedClaims.Subject {
			t.Fatal("sub mismatch")
		}

	}
}

func newUnit(t *testing.T) (*logger.Logger, func()) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", func(ctx context.Context) string { return "" }, "test-build", "test-build-date")

	teardown := func() {
		t.Helper()

		fmt.Println("***************** TEST LOGS ****************")
		fmt.Print(buf.String())
		fmt.Println("***************** TEST LOGS ****************")
	}

	return log, teardown
}

type keyStore struct {
}

func (ks *keyStore) PrivateKeyPEM(kid string) (string, error) {
	return privateKey, nil
}

func (ks *keyStore) PublicKeyPEM(kid string) (string, error) {
	return publicKey, nil
}

const (
	kid = "aeafe"

	privateKey = `-----BEGIN PRIVATE KEY-----
MIIEowIBAAKCAQEAwT36CvbE3HRuZMJ76Ok0G76BaORHXmPpYvLqPH7s+tEny6nS
KdW0k+ZqwZtViICIHUqW5iZxpzk6bsySE+CdpxWq5X7OwbV0+GfvndH4/H9cRRwJ
kbTgtvD0cfpaeqttlA0c9KbV44BHFOKQE2CrxI8dpV/avh1ttJ8LPtnhA7zV1Vns
PwyXPew5S27RtqNzS3YY73v+uW619enLgJ2sdBounKOsiKDTiPmLXcbrpUOpz1oj
cIkfQkcecFCt/nTCYQXJM12qQUDTalCUASneITyA6sJ0nV6zX1pukXiOV+FhWeBx
jr+RXRaQLKlrvITkekvaW8i6W1vzxMEVnaIQbQIDAQABAoIBAD/g3OcC69B0LIni
QFOdwzvonZ6u12i8Xkn3Qrb8vhmrShxo5rAtlKPPZzrYvk2BL31/SqKJ0sgUgtxd
g2xOs86nuvZiFLaz5Ra4RN1Gq6fL8hRmUEN6J05SGHwU5hPm1zI5o3i8Lbgmj1i4
DChbrGwtYv9n0EdIMxuh93WvUFKLZnlfNjdizjiHxIt0dA1WyC93DIHY8CzxTKTh
LztUyfZJY3t4tJB71RdgqAo3KqYvLHlPWdH+/tnhvB/tkVC7Pg3Wz69am3T7ZsOR
3pmqLu0uuj2LeRyqBa7bUjmGsg4PlTDIc9D14+jf1pJ2Uqf8pY7hLrE/hLJMjqIl
SsMnMYECgYEA54/dGQeOgpo/JreI/PXZZgORwjhaJHcsiMV7xQpVutnzMsJCT8RZ
xH8QzjQaNyPkwEyC5swcsRQjNT2WKATtcFPVXL3Lfjw1XUkMsqJIyTyxULMaXzkQ
vg3LxDEvY1C8QclfReqFvAnvbgQPerkFGn4uYB5eSB7g25Owu4bAowMCgYEA1aLS
GallDqTtxNUFOPBOhjYg43eKYlG+qPTmCSZyO5YkaRPnCLYXubr/z/D2kuHegKqY
PWH0vUug0Jl3xGNCUxQgkjNGaptGH4xdkXJwS3JfnVrBh35b6/MCvcsBL1H/RV/A
pLlJd/BqEg995hFPOifmTERuuX5zhUzPiaKUa88CgYBrF0CUg8cnpXhRSo5DFtwc
7sepP/CIbskc0+Aix13QlV2A+xA53b/6PR2jj7sUBziSqo/bd0hJqNuDkuDvzx+U
VGEXi/shfWzcoJ0LvDvXPYdvn/UxVq3kh3LWYDzfPIZkEmJKbmS9Cwc9JQmXoxu2
HecKsuC9j/JD+nDAuVg7ywKBgQDTaAVLGaMEY9daoYZCeyLpkyTmk8PgWY43uejt
gxsltgUf6m8E4tUFmXRN/OC0J0m8v5RZodbDf9SKuGOZdrQKbG9y3cS/+BnjXE63
gKx5LJxLpaS/hR5HljnPQNVSU20t5HxJRYXbZ5A/gQ8QHW7uWM+AB3QeoXCFp6X0
SGthZwKBgFew2F89chWxOubz04S777kNgrqGw5j2kY93JZVFPVPuihAxnOfePJFj
eCClcvPf+rtXbXaMiyIDnzOcmvWW4dymxLxvtkfpZi92g375osrfvmkuhzcsoAt/
zOUckgMWCaCANa4DedYdlVv1EYIcODrKOo0V+K1MDyq4NMKXbvoM
-----END PRIVATE KEY-----
`

	publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwT36CvbE3HRuZMJ76Ok0
G76BaORHXmPpYvLqPH7s+tEny6nSKdW0k+ZqwZtViICIHUqW5iZxpzk6bsySE+Cd
pxWq5X7OwbV0+GfvndH4/H9cRRwJkbTgtvD0cfpaeqttlA0c9KbV44BHFOKQE2Cr
xI8dpV/avh1ttJ8LPtnhA7zV1VnsPwyXPew5S27RtqNzS3YY73v+uW619enLgJ2s
dBounKOsiKDTiPmLXcbrpUOpz1ojcIkfQkcecFCt/nTCYQXJM12qQUDTalCUASne
ITyA6sJ0nV6zX1pukXiOV+FhWeBxjr+RXRaQLKlrvITkekvaW8i6W1vzxMEVnaIQ
bQIDAQAB
-----END PUBLIC KEY-----`
)
