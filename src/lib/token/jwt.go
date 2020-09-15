package token

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/helixauth/helix/src/lib/secrets"
)

// JWT generates a new JSON Web Token
func JWT(ctx context.Context,
	claims map[string]interface{},
	expiresAt time.Time,
	sig jwt.SigningMethod,
	secrets secrets.Manager) (string, error) {

	// Create new token
	tkn := jwt.New(sig)
	for k, v := range claims {
		tkn.Claims.(jwt.MapClaims)[k] = v
	}
	tkn.Claims.(jwt.MapClaims)["exp"] = expiresAt.UTC().Unix()

	// TODO add issuer

	switch sig {
	case jwt.SigningMethodHS256:
		sec, err := secrets.GetString("jws.hs256.secret")
		if err != nil {
			return "", err
		}
		return tkn.SignedString([]byte(sec))

	default:
		return "", fmt.Errorf("Signing method %v not supported", sig)
	}
}
