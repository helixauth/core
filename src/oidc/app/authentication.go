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

// TODO authentication request
// - email
// - password

func (a *app) Authentication(c *gin.Context) {
	ctx := a.context(c)
	req := authorizationRequest{}
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

	// TODO authenticate the email/password

	txn, err := a.Database.Txn(ctx)
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

	if err = txn.Insert(ctx, session); err != nil {
		txn.Rollback()
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}

	destination := req.RedirectURI + fmt.Sprintf("?code=%v&state=%v", code, session.State)
	c.Redirect(http.StatusFound, destination)
}
