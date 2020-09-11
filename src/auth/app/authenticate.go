package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/helixauth/helix/src/entity"
	"github.com/helixauth/helix/src/lib/database"
	"github.com/helixauth/helix/src/lib/mapper"
	"github.com/helixauth/helix/src/lib/utils"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
)

type authenticateRequest struct {
	Email    string  `form:"email" binding:"required"`
	Password *string `form:"password"`
}

func (a *app) Authenticate(c *gin.Context) {
	ctx := a.context(c)
	oauthReq := oauthRequest{}
	if err := c.BindQuery(&oauthReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authnReq := authenticateRequest{}
	if err := c.Bind(&authnReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO validate the clientID
	// TODO validate the response type
	// TODO validate the scopes
	// TODO validate the redirect URI is authorized
	// TODO validate the prompt

	isSignUp := oauthReq.Prompt != nil && *oauthReq.Prompt == "create"
	userNotFound := false
	user := &entity.User{}
	err := a.Database.QueryItem(ctx, user, "SELECT * FROM users WHERE email = $1", authnReq.Email)
	if err == sql.ErrNoRows {
		userNotFound = true
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	txn, err := a.Database.BeginTxn(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if userNotFound {
		if isSignUp {
			user, err = a.newUser(ctx, authnReq.Email, authnReq.Password, txn)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password do not match"})
			return
		}
	} else {
		// TODO validate password
		// TODO send email verification
		// return
	}

	sess, err := a.newUserSession(ctx, *user, oauthReq, txn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = txn.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dest := "https://" + oauthReq.RedirectURI + fmt.Sprintf("?code=%v&state=%v", sess.Code, sess.State)
	c.Redirect(http.StatusFound, dest)
}

func (a *app) newUser(ctx context.Context, email string, password *string, txn database.Txn) (*entity.User, error) {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &entity.User{
		ID:            uniuri.NewLen(40),
		TenantID:      a.TenantID,
		Email:         &email,
		EmailVerified: mapper.BoolPtr(false),
		PasswordHash:  passwordHash,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	err = txn.Insert(ctx, user)

	// TODO email verification

	return user, err
}

func (a *app) newUserSession(ctx context.Context, user entity.User, req oauthRequest, txn database.Txn) (*entity.Session, error) {
	code := uniuri.NewLen(uniuri.UUIDLen * 2)
	session := &entity.Session{
		ID:           utils.Hash(code),
		TenantID:     a.TenantID,
		ClientID:     req.ClientID,
		UserID:       &user.ID,
		ResponseType: req.ResponseType,
		Scope:        req.Scope,
		State:        req.State,
		Nonce:        req.Nonce,
		Code:         code,
		RedirectURI:  req.RedirectURI,
		CreatedAt:    time.Now().UTC(),
	}
	err := txn.Insert(ctx, session)
	return session, err
}
