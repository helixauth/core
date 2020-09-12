package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/helixauth/helix/src/auth/app/oauth"
	"github.com/helixauth/helix/src/entity"
	"github.com/helixauth/helix/src/lib/database"
	"github.com/helixauth/helix/src/lib/mapper"
	"github.com/helixauth/helix/src/lib/token"
	"github.com/helixauth/helix/src/lib/utils"

	"github.com/dchest/uniuri"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type formInput struct {
	Email           string  `form:"email" binding:"required"`
	Password        *string `form:"password"`
	ConfirmPassword *string `form:"confirm_password"`
}

// Authorize is the handler for the /authorize endpoint
func (a *app) Authorize(c *gin.Context) {
	form := formInput{}
	params := oauth.Params{}
	if err := c.BindQuery(&params); err != nil {
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
		render(c, params, nil, nil)

	case http.MethodPost:
		if err := c.Bind(&form); err != nil {
			render(c, params, nil, err)
			return
		}
		a.processForm(c, params, form)

	default:
		render(c, params, nil, nil)
	}
}

func (a *app) processForm(c *gin.Context, params oauth.Params, form formInput) {
	ctx := a.context(c)

	// Query for existing users with the email address in the form
	userNotFound := false
	user := &entity.User{}
	err := a.Database.QueryItem(ctx, user, "SELECT * FROM users WHERE email = $1", form.Email)
	if err == sql.ErrNoRows {
		userNotFound = true
	} else if err != nil {
		render(c, params, &form, err)
		return
	}

	txn, err := a.Database.BeginTxn(ctx)
	if err != nil {
		render(c, params, &form, err)
		return
	}

	// Register a new user or authenticate the existing user
	if userNotFound {
		if params.IsSignUp() {
			user, err = a.registerUser(ctx, params, form, txn)
		} else {
			render(c, params, &form, fmt.Errorf("Incorrect email or password"))
			return
		}
	} else {
		err = a.authenticateUser(user, form)
	}
	if err != nil {
		render(c, params, &form, err)
		return
	}

	err = txn.Commit()
	if err != nil {
		render(c, params, &form, err)
		return
	}

	// Start a new user session
	code, err := a.generateAuthorizationCode(ctx, params, user)
	if err != nil {
		render(c, params, &form, err)
		return
	}

	// Redirect to the provided redirect URI with session code and state
	dest := "https://" + mapper.String(params.RedirectURI) + fmt.Sprintf("?code=%v&state=%v", code, params.State)
	c.Redirect(http.StatusFound, dest)
}

// registerUser creates a new user
func (a *app) registerUser(ctx context.Context, params oauth.Params, form formInput, txn database.Txn) (*entity.User, error) {
	passwordHash, err := utils.HashPassword(form.Password)
	if err != nil {
		return nil, err
	}
	user := &entity.User{
		ID:            uniuri.NewLen(40),
		TenantID:      a.TenantID,
		Email:         &form.Email,
		EmailVerified: mapper.BoolPtr(false),
		PasswordHash:  passwordHash,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	err = txn.Insert(ctx, user)
	return user, err
}

// authenticateUser validates the form input against an existing user's records
func (a *app) authenticateUser(user *entity.User, form formInput) error {
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

func (a *app) validateOAuthParams(ctx context.Context, params oauth.Params) error {
	// TODO validate the client and request
	return nil
}

// render renders the authorization form on screen
func render(c *gin.Context, params oauth.Params, form *formInput, err error) {
	tmplParams := gin.H{
		"action":   c.Request.URL.RawPath + "?" + c.Request.URL.RawQuery,
		"email":    nil,
		"password": nil,
		"error":    nil,
	}
	if form != nil {
		tmplParams["email"] = form.Email
		tmplParams["password"] = form.Password
	}
	if err != nil {
		tmplParams["error"] = err.Error()
	}

	// Render the 'sign up' page
	if params.IsSignUp() {
		tmplParams["title"] = "Sign up"
		c.HTML(http.StatusOK, "signUp.html", tmplParams)
		return
	}

	// Render the 'sign in' page
	tmplParams["title"] = "Sign in"
	c.HTML(http.StatusOK, "signIn.html", tmplParams)
}

func (a *app) generateAuthorizationCode(ctx context.Context, params oauth.Params, user *entity.User) (string, error) {
	claims := map[string]interface{}{
		"jti":          uniuri.NewLen(uniuri.UUIDLen),
		"client_id":    params.ClientID,
		"redirect_uri": params.RedirectURI,
		"user_id":      user.ID,
	}
	exp := time.Now().UTC().Add(30 * time.Second)
	return token.JWT(ctx, claims, exp, jwt.SigningMethodHS256)
}
