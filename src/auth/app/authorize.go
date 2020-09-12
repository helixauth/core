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
	"golang.org/x/crypto/bcrypt"
)

type oauthParams struct {
	ClientID     string  `form:"client_id" binding:"required"`
	ResponseType string  `form:"response_type" binding:"required"`
	Scope        string  `form:"scope" binding:"required"`
	State        string  `form:"state" binding:"required"`
	Nonce        string  `form:"nonce" binding:"required"`
	RedirectURI  string  `form:"redirect_uri" binding:"required"`
	Prompt       *string `form:"prompt"`
}

type formInput struct {
	Email           string  `form:"email" binding:"required"`
	Password        *string `form:"password"`
	ConfirmPassword *string `form:"confirm_password"`
}

// Authorize is the handler for the /authorize endpoint
func (a *app) Authorize(c *gin.Context) {
	form := formInput{}
	oauth := oauthParams{}
	if err := c.BindQuery(&oauth); err != nil {
		c.HTML(
			http.StatusBadRequest,
			"error.html",
			gin.H{"error": err.Error()},
		)
		return
	}

	// TODO validate the oauth params

	switch c.Request.Method {
	case http.MethodGet:
		render(c, oauth, form, nil)

	case http.MethodPost:
		if err := c.Bind(&form); err != nil {
			render(c, oauth, form, err)
			return
		}
		a.processForm(c, oauth, form)

	default:
		render(c, oauth, form, nil)
	}
}

func (a *app) processForm(c *gin.Context, oauth oauthParams, form formInput) {
	ctx := a.context(c)

	// Query for existing users with the email address in the form
	userNotFound := false
	user := &entity.User{}
	err := a.Database.QueryItem(ctx, user, "SELECT * FROM users WHERE email = $1", form.Email)
	if err == sql.ErrNoRows {
		userNotFound = true
	} else if err != nil {
		render(c, oauth, form, err)
		return
	}

	txn, err := a.Database.BeginTxn(ctx)
	if err != nil {
		render(c, oauth, form, err)
		return
	}

	// Create a new user or authenticate the existing user
	isSignUp := oauth.Prompt != nil && *oauth.Prompt == "create"
	if userNotFound {
		if isSignUp {
			user, err = a.newUser(ctx, form.Email, form.Password, txn)
			if err != nil {
				render(c, oauth, form, err)
				return
			}
		} else {
			render(c, oauth, form, fmt.Errorf("Email or password do not match"))
			return
		}
	} else {
		if err = authenticateUser(user, form); err != nil {
			render(c, oauth, form, err)
		}
	}

	// Start a new user session
	sess, err := a.newUserSession(ctx, *user, oauth, txn)
	if err != nil {
		render(c, oauth, form, err)
		return
	}

	err = txn.Commit()
	if err != nil {
		render(c, oauth, form, err)
		return
	}

	// Redirect to the provided redirect URI with OAuth session code and state
	dest := "https://" + oauth.RedirectURI + fmt.Sprintf("?code=%v&state=%v", sess.Code, sess.State)
	c.Redirect(http.StatusFound, dest)
}

// newUser creates a new User entity
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

// newUserSession starts a new user-based OAuth session
func (a *app) newUserSession(ctx context.Context, user entity.User, req oauthParams, txn database.Txn) (*entity.Session, error) {
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

// authenticateUser validates the form input against an existing user's records
func authenticateUser(user *entity.User, form formInput) error {
	if user.PasswordHash != nil {
		if form.Password != nil {
			if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(*form.Password)); err != nil {
				return fmt.Errorf("Incorrect email or password")
			}
		} else {
			return fmt.Errorf("Password required")
		}
	} else {
		// TODO send email verification
	}
	return nil
}

// render renders the authorization form on screen
func render(c *gin.Context, oauth oauthParams, form formInput, err error) {
	params := gin.H{
		"action":   c.Request.URL.RawPath + "?" + c.Request.URL.RawQuery,
		"email":    form.Email,
		"password": form.Password,
	}
	if err != nil {
		params["error"] = err.Error()
	}

	// Render the 'sign up' page
	if oauth.Prompt != nil && *oauth.Prompt == "create" {
		params["title"] = "Sign up"
		c.HTML(
			http.StatusOK,
			"signUp.html",
			params,
		)
		return
	}

	// Render the 'sign in' page
	params["title"] = "Sign in"
	c.HTML(
		http.StatusOK,
		"signIn.html",
		params,
	)
}
