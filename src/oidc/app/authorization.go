package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type authorizationRequest struct {
	ClientID     string  `form:"client_id" binding:"required"`
	ResponseType string  `form:"response_type" binding:"required"`
	Scope        string  `form:"scope" binding:"required"`
	State        string  `form:"state" binding:"required"`
	Nonce        string  `form:"nonce" binding:"required"`
	RedirectURI  string  `form:"redirect_uri" binding:"required"`
	Prompt       *string `form:"prompt"`
}

func (a *app) Authorization(c *gin.Context) {
	req := authorizationRequest{}
	if err := c.BindQuery(&req); err != nil {
		c.HTML(
			http.StatusBadRequest,
			"error.html",
			gin.H{"error": err.Error()},
		)
		return
	}

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
