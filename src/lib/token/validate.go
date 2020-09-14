package token

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// Validate validates a JWT and returns its claims
// TODO accept secrets manager and use it to validate signatures
func Validate(ctx context.Context, jwtStr string, sig jwt.SigningMethod) (map[string]interface{}, error) {
	tkn, err := jwt.Parse(jwtStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != sig {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("supersecret"), nil
	})
	if err != nil {
		return nil, err
	}

	// TODO validate issuer

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Token is invalid")
}
