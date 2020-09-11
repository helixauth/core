package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/helixauth/helix/src/entity"
	"github.com/helixauth/helix/src/lib/utils"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
)

// TODO authentication request
// - email
// - password

type authenticateRequest struct {
	Email    string  `form:"email" binding:"required"`
	Password *string `form:"password"`
}

func (a *app) Authenticate(c *gin.Context) {
	ctx := a.context(c)
	authzReq := authorizeRequest{}
	if err := c.BindQuery(&authzReq); err != nil {
		c.HTML(
			http.StatusBadRequest,
			"error.html",
			gin.H{"error": err.Error()},
		)
		return
	}
	authnReq := authenticateRequest{}
	if err := c.BindJSON(&authnReq); err != nil {
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

	user := &entity.User{}
	err := a.Database.QueryItem(ctx, user, "SELECT * FROM users WHERE email = $1", authnReq.Email)
	if err == sql.ErrNoRows {
		// TODO create new user
	} else if err != nil {
		panic(err)
	}

	txn, err := a.Database.BeginTxn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	userID := "foo"

	code := uniuri.NewLen(uniuri.UUIDLen * 2)
	session := &entity.Session{
		ID:           utils.Hash(code),
		TenantID:     a.TenantID,
		ClientID:     authzReq.ClientID,
		UserID:       &userID,
		ResponseType: authzReq.ResponseType,
		Scope:        authzReq.Scope,
		State:        authzReq.State,
		Nonce:        authzReq.Nonce,
		Code:         code,
		RedirectURI:  authzReq.RedirectURI,
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

	destination := authzReq.RedirectURI + fmt.Sprintf("?code=%v&state=%v", code, session.State)
	c.Redirect(http.StatusFound, destination)
}
