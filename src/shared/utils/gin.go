package utils

import (
	"github.com/gin-gonic/gin"
)

func GetQueryParam(c *gin.Context, name string) *string {
	values := c.Request.URL.Query()[name]
	if len(values) < 1 {
		return nil
	}
	return &values[0]
}

func GetScheme(c *gin.Context) string {
	var scheme = "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme
}
