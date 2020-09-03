package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) Index(c *gin.Context) {
	c.Redirect(http.StatusFound, fmt.Sprintf("/authorization%v", c.Request.URL.String()))
}
