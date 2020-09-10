package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *app) Index(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{},
	)
}
