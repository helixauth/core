package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/helixauth/helix/src/auth/app/oauth"
	"github.com/helixauth/helix/src/lib/token"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

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

		resp := oauth.TokenResponse{}

		// TODO prevent replay attacks by marking the authorization code as used

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
