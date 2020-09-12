package token

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func JWT(ctx context.Context, claims map[string]interface{}, expiresAt time.Time, sig jwt.SigningMethod) (string, error) {
	tkn := jwt.New(sig)
	for k, v := range claims {
		tkn.Claims.(jwt.MapClaims)[k] = v
	}

	switch sig {
	case jwt.SigningMethodHS256:
		return tkn.SignedString([]byte("supersecret"))
	default:
		return "", fmt.Errorf("Signing method %v not supported", sig)
	}

}
