package token

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/helixauth/helix/src/lib/secrets"
)

// JWT generates a new JSON Web Token
func JWT(
	ctx context.Context,
	claims map[string]interface{},
	expiresAt time.Time,
	sig jwt.SigningMethod,
	secrets secrets.Manager,
) (string, error) {

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

	case jwt.SigningMethodRS256:
		secs, err := secrets.GetMap("jws.rs256")
		if err != nil {
			return "", err
		}

		// Naive strategy: choose first RSA256 key
		// TODO pick an available key at random
		for sec := range secs {
			if kid, ok := sec.(string); ok {
				tkn.Claims.(jwt.MapClaims)["kid"] = kid
				strBase64, err := secrets.GetString(fmt.Sprintf("jws.rs256.%v.private", kid))
				if err != nil {
					return "", err
				}
				key, err := parseRSA256PrivateKey(strBase64)
				if err != nil {
					return "", err
				}
				return tkn.SignedString(key)
			}
		}
		return "", fmt.Errorf("No RSA256 signing key found")

	default:
		return "", fmt.Errorf("Signing method %v not supported", sig)
	}
}

func parseRSA256PrivateKey(str string) (*rsa.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	keyPem, _ := pem.Decode(keyBytes)
	if keyPem == nil {
		return nil, fmt.Errorf("Failed to get secret")
	}
	if keyPem.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("Decoded key is of the wrong type (%v)", keyPem.Type)
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(keyPem.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(keyPem.Bytes); err != nil {
			return nil, err
		}
	}
	privKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("Failed to parse private key")
	}

	return privKey, nil
}
