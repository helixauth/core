package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type authorizeRequest struct {
	ClientID     string  `form:"client_id" binding:"required"`
	ResponseType string  `form:"response_type" binding:"required"`
	Scope        string  `form:"scope" binding:"required"`
	State        string  `form:"state" binding:"required"`
	Nonce        string  `form:"nonce" binding:"required"`
	RedirectURI  string  `form:"redirect_uri" binding:"required"`
	Prompt       *string `form:"prompt"`
}

func (a *app) Authorize(c *gin.Context) {
	req := authorizeRequest{}
	if err := c.BindQuery(&req); err != nil {
		c.HTML(
			http.StatusBadRequest,
			"error.html",
			gin.H{"error": err.Error()},
		)
		return
	}

	// TODO validate the clientID
	// TODO validate the response type
	// TODO validate the scopes
	// TODO validate the redirect URI is authorized
	// TODO validate the prompt

	if req.Prompt != nil && *req.Prompt == "create" {
		c.HTML(
			http.StatusOK,
			"signUp.html",
			gin.H{"title": "Sign up"},
		)
		return
	}

	c.HTML(
		http.StatusOK,
		"signIn.html",
		gin.H{"title": "Sign in"},
	)
}
