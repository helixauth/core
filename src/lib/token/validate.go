package token

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func Validate(ctx context.Context, tknStr string, sig jwt.SigningMethod) (map[string]interface{}, error) {
	tkn, err := jwt.Parse(tknStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != sig {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("supersecret"), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Token is invalid")
}
