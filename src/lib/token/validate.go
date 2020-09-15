package token

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/helixauth/helix/src/lib/secrets"

	"github.com/dgrijalva/jwt-go"
)

// Validate validates a JWT and returns its claims
func Validate(
	ctx context.Context,
	jwtStr string,
	signingMethod jwt.SigningMethod,
	secrets secrets.Manager,
) (map[string]interface{}, error) {

	// Parse the JWT string
	tkn, err := jwt.Parse(jwtStr, getValidationKey(signingMethod, secrets))
	if err != nil {
		return nil, err
	}

	// TODO validate issuer

	// Return extracted claims
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("Token is invalid")
}

// getValidationKey inspects a token and the expected signing method and returns a key for validating the signature
func getValidationKey(
	signingMethod jwt.SigningMethod,
	secrets secrets.Manager,
) jwt.Keyfunc {

	return func(t *jwt.Token) (interface{}, error) {
		if t.Method != signingMethod {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// Depending on the expected signing method...
		switch signingMethod {

		// For HS256, fetch the shared secret
		case jwt.SigningMethodHS256:
			sec, err := secrets.GetString("jws.hs256.secret")
			if err != nil {
				return nil, err
			}
			return []byte(sec), nil

		// For RSA256, fetch the public key for the token's KID
		case jwt.SigningMethodRS256:
			keyID, ok := t.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("Header 'kid' not found")
			}
			keyStr, err := secrets.GetString(fmt.Sprintf("jws.rs256.%v.public", keyID))
			if err != nil {
				return nil, err
			}
			key, err := parseRSA256PublicKey(keyStr)
			if err != nil {
				return nil, err
			}
			return key, err

		// Return error if expected signing method is unsupported
		default:
			return nil, fmt.Errorf("Validation method '%v' not supported", signingMethod)
		}
	}
}

// parseRSA256PublicKey parses a RS256 public key from a base64 encoded string
func parseRSA256PublicKey(str string) (*rsa.PublicKey, error) {

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
	if keyPem.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("Decoded key is of the wrong type (%v)", keyPem.Type)
	}

	// Parse to key bytes
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(keyPem.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS1PublicKey(keyPem.Bytes); err != nil {
			return nil, err
		}
	}

	// Cast to RSA public key type
	pubKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Failed to parse public key")
	}
	return pubKey, nil
}
