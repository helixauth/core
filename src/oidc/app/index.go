package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *app) Index(c *gin.Context) {
	c.Redirect(http.StatusFound, fmt.Sprintf("/authorize%v", c.Request.URL.String()))
}
