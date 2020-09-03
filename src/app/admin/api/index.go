package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) Index(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{},
	)
}
