package token

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT generates a new JSON Web Token
func JWT(ctx context.Context, claims map[string]interface{}, expiresAt time.Time, sig jwt.SigningMethod) (string, error) {
	tkn := jwt.New(sig)
	for k, v := range claims {
		tkn.Claims.(jwt.MapClaims)[k] = v
	}
	tkn.Claims.(jwt.MapClaims)["exp"] = expiresAt.UTC().Unix()

	// TODO issuer

	switch sig {
	case jwt.SigningMethodHS256:
		return tkn.SignedString([]byte("supersecret"))
	default:
		return "", fmt.Errorf("Signing method %v not supported", sig)
	}
}
