package jwt

import (
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateToken(t *testing.T) {
	signer := New("01234567890123456789012345678901", 120)
	token, expiresIn, err := signer.GenerateToken(42, "admin", "web-admin", "web")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if expiresIn != 120 {
		t.Fatalf("expiresIn = %d, want 120", expiresIn)
	}
	claims, err := signer.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if claims.UserId != 42 || claims.ClientId != "web-admin" || claims.ID == "" {
		t.Fatalf("claims = %#v", claims)
	}
}

func TestValidateTokenRejectsWrongAlgorithmAndIssuer(t *testing.T) {
	signer := New("01234567890123456789012345678901", 120)

	wrongAlgorithm := jwtlib.NewWithClaims(jwtlib.SigningMethodHS384, Claims{
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    "microservice-kit",
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	encoded, err := wrongAlgorithm.SignedString(signer.secret)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := signer.ValidateToken(encoded); err == nil {
		t.Fatal("ValidateToken() accepted HS384 token")
	}

	wrongIssuer := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, Claims{
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    "another-service",
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	encoded, err = wrongIssuer.SignedString(signer.secret)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := signer.ValidateToken(encoded); err == nil {
		t.Fatal("ValidateToken() accepted token from another issuer")
	}
}
