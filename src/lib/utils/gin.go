package utils

import (
	"github.com/gin-gonic/gin"
)

func GetScheme(c *gin.Context) string {
	var scheme = "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme
}
