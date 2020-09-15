package token

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/helixauth/helix/src/lib/secrets"
)

// Validate validates a JWT and returns its claims
func Validate(
	ctx context.Context,
	jwtStr string,
	sig jwt.SigningMethod,
	secrets secrets.Manager,
) (map[string]interface{}, error) {

	// Parse the JWT string
	tkn, err := jwt.Parse(jwtStr, getValidationKey(sig, secrets))
	if err != nil {
		return nil, err
	}

	// TODO validate issuer

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Token is invalid")
}

func getValidationKey(sig jwt.SigningMethod, secrets secrets.Manager) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if t.Method != sig {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		switch sig {
		case jwt.SigningMethodHS256:
			sec, err := secrets.GetString("jws.hs256.secret")
			if err != nil {
				return nil, err
			}
			return []byte(sec), nil

		case jwt.SigningMethodRS256:
			secs, err := secrets.GetMap("jws.rs256")
			if err != nil {
				return nil, err
			}

			// Naive strategy: choose first RSA256 key
			// TODO pick an available key at random
			for sec := range secs {
				if kid, ok := sec.(string); ok {
					strBase64, err := secrets.GetString(fmt.Sprintf("jws.rs256.%v.public", kid))
					if err != nil {
						return nil, err
					}
					key, err := parseRSA256PublicKey(strBase64)
					if err != nil {
						return nil, err
					}
					return key, err
				}
			}
			return nil, fmt.Errorf("No RSA256 public key found")

		default:
			return nil, fmt.Errorf("Validation method '%v' not supported", sig)
		}
	}
}

func parseRSA256PublicKey(str string) (*rsa.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	keyPem, _ := pem.Decode(keyBytes)
	if keyPem == nil {
		return nil, fmt.Errorf("Failed to get secret")
	}
	if keyPem.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("Decoded key is of the wrong type (%v)", keyPem.Type)
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(keyPem.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS1PublicKey(keyPem.Bytes); err != nil {
			return nil, err
		}
	}
	pubKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Failed to parse public key")
	}

	return pubKey, nil
}
