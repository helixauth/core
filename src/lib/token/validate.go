package token

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/helixauth/helix/src/lib/secrets"
)

// Validate validates a JWT and returns its claims
func Validate(ctx context.Context,
	jwtStr string,
	sig jwt.SigningMethod,
	secrets secrets.Manager) (map[string]interface{}, error) {

	// Parse the JWT string
	tkn, err := jwt.Parse(jwtStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != sig {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		switch sig {
		case jwt.SigningMethodHS256:
			sec, err := secrets.GetString("jws.hs256.secret")
			if err != nil {
				return "", err
			}
			return []byte(sec), nil

		default:
			return "", fmt.Errorf("Signing method %v not supported", sig)
		}
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
