package token

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/helixauth/helix/src/lib/secrets"

	"github.com/dgrijalva/jwt-go"
)

// JWT generates a new JSON Web Token
func JWT(
	ctx context.Context,
	claims map[string]interface{},
	expiresAt time.Time,
	signingMethod jwt.SigningMethod,
	secrets secrets.Manager,
) (string, error) {

	// Create new token
	tkn := jwt.New(signingMethod)
	for k, v := range claims {
		tkn.Claims.(jwt.MapClaims)[k] = v
	}
	tkn.Claims.(jwt.MapClaims)["exp"] = expiresAt.UTC().Unix()

	// TODO add issuer

	// Depending on the selected signing method...
	switch signingMethod {

	// For HS256, fetch the shared secret
	case jwt.SigningMethodHS256:
		sec, err := secrets.GetString("jws.hs256.secret")
		if err != nil {
			return "", err
		}
		return tkn.SignedString([]byte(sec))

	// For RS256, fetch one of the RSA256 private signing key
	case jwt.SigningMethodRS256:
		secs, err := secrets.GetMap("jws.rs256")
		if err != nil {
			return "", err
		}

		// Naive strategy: choose first RSA256 key
		// TODO pick an available key at random
		for sec := range secs {
			if keyID, ok := sec.(string); ok {
				tkn.Header["kid"] = keyID
				keyStr, err := secrets.GetString(fmt.Sprintf("jws.rs256.%v.private", keyID))
				if err != nil {
					return "", err
				}
				key, err := parseRSA256PrivateKey(keyStr)
				if err != nil {
					return "", err
				}
				return tkn.SignedString(key)
			}
		}
		return "", fmt.Errorf("No RSA256 signing key found")

	// Return error if selected signing method is unsupported
	default:
		return "", fmt.Errorf("Signing method '%v' not supported", signingMethod)
	}
}

// parseRSA256PublicKey parses a RS256 private key from a base64 encoded string
func parseRSA256PrivateKey(str string) (*rsa.PrivateKey, error) {

	// Base64 decode string
	keyBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	// Decode PEM
	keyPem, _ := pem.Decode(keyBytes)
	if keyPem == nil {
		return nil, fmt.Errorf("Failed to get secret")
	}

	// Validate PEM type
	if keyPem.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("Decoded key is of the wrong type (%v)", keyPem.Type)
	}

	// Parse to key bytes
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(keyPem.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(keyPem.Bytes); err != nil {
			return nil, err
		}
	}

	// Cast to RSA private key type
	privKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("Failed to parse private key")
	}
	return privKey, nil
}
