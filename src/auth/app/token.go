package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/helixauth/helix/src/auth/app/oauth"
	"github.com/helixauth/helix/src/lib/token"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Token is the handler for the /token endpoint
func (a *app) Token(c *gin.Context) {

	// Parse request
	ctx := a.context(c)
	req := oauth.TokenRequest{}
	if err := c.BindJSON(&req); err != nil {
		panic(err)
	}

	// Validate the client
	client, err := a.getClient(ctx, req.ClientID)
	if err != nil {
		panic(err)
	}
	if client.Secret != nil && req.ClientSecret != client.Secret {
		panic(fmt.Errorf("client_secret is invalid"))
	}

	// Handle grant type
	switch req.GrantType {
	case oauth.GrantTypeAuthorizationCode:
		if err = a.validateAuthorizationCode(ctx, req); err != nil {
			panic(err)
		}

		// Generate ID token
		claims := map[string]interface{}{}
		exp := time.Now().UTC().Add(5 * 60 * time.Second)
		idToken, err := token.JWT(ctx, claims, exp, jwt.SigningMethodRS256, a.Secrets)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// TODO prevent replay attacks by marking the authorization code as used

		// Respond
		resp := oauth.TokenResponse{
			IDToken: idToken,
		}
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.JSON(http.StatusOK, resp)

	default:
		panic(fmt.Errorf("Grant type '%v' is unsupported", req.GrantType))
	}
}

// validateAuthorizationCode validates an authorization code for issuing tokens
func (a *app) validateAuthorizationCode(ctx context.Context, req oauth.TokenRequest) error {
	claims, err := token.Validate(ctx, req.Code, jwt.SigningMethodHS256, a.Secrets)
	if err != nil {
		panic(err)
	}

	isAuthorizationCodeValid := claims["client_id"] == req.ClientID &&
		claims["redirect_uri"] == req.RedirectURI
	if !isAuthorizationCodeValid {
		panic(fmt.Errorf("code is invalid"))
	}

	return nil
}
