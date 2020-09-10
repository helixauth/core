package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/helixauth/helix/src/shared/entity"
	"github.com/helixauth/helix/src/shared/utils"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
)

func (a *app) Authenticate(c *gin.Context) {
	ctx := c.Request.Context()
	req := authorizationRequest{}
	if err := c.BindQuery(&req); err != nil {
		c.HTML(
			http.StatusBadRequest,
			"error.html",
			gin.H{"error": err.Error()},
		)
		return
	}

	// email := c.PostForm("email")
	// password := c.PostForm("password")
	// TODO authenticate the email/password

	tx, err := a.Gateways.Database.BeginTx(ctx)
	if err != nil {
		log.Fatal(err)
	}

	code := uniuri.NewLen(uniuri.UUIDLen * 2)
	session := &entity.Session{
		ID:           utils.Hash(code),
		UserID:       "foo",
		ClientID:     req.ClientID,
		ResponseType: req.ResponseType,
		Scope:        req.Scope,
		State:        req.State,
		Nonce:        req.Nonce,
		Code:         code,
		RedirectURI:  req.RedirectURI,
		CreatedAt:    time.Now().UTC(),
	}

	if err = utils.SQLInsert(ctx, session, "sessions", tx); err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// TODO validate the redirect URI is an authorized domain

	destination := req.RedirectURI + fmt.Sprintf("?code=%v&state=%v", code, session.State)
	c.Redirect(http.StatusFound, destination)
}
