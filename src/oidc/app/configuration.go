package app

import (
	"fmt"
	"net/http"

	"github.com/helixauth/helix/src/shared/utils"

	"github.com/gin-gonic/gin"
)

func (a *app) Configuration(c *gin.Context) {
	var msg struct {
		Issuer                            string   `json:"issuer"`
		AuthorizationEndpoint             string   `json:"authorization_endpoint"`
		TokenEndpoint                     string   `json:"token_endpoint"`
		UserInfoEndpoint                  string   `json:"userinfo_endpoint"`
		RegistrationEndpoint              string   `json:"registration_endpoint"`
		JWKsURI                           string   `json:"jwks_uri"`
		ResponseTypesSupported            []string `json:"response_types_supported"`
		SubjectTypesSupported             []string `json:"subject_types_supported"`
		IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
		ScopesSupported                   []string `json:"scopes_supported"`
		TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
		ClaimsSupported                   []string `json:"claims_supported"`
		CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
		GrantTypesSupported               []string `json:"grant_types_supported"`
		// TODO
	}
	scheme := utils.GetScheme(c)
	baseURL := fmt.Sprintf("%v://%v", scheme, c.Request.Host)
	msg.AuthorizationEndpoint = baseURL + "/authorize"
	msg.TokenEndpoint = baseURL + "/token"
	msg.UserInfoEndpoint = baseURL + "/userinfo"
	msg.JWKsURI = baseURL + "/jwks"
	msg.ResponseTypesSupported = []string{"code"}
	msg.SubjectTypesSupported = []string{"public"}
	msg.IDTokenSigningAlgValuesSupported = []string{"RS256"}
	msg.ScopesSupported = []string{"openid", "profile", "email", "phone"}
	// msg.TokenEndpointAuthMethodsSupported
	// msg.ClaimsSupported
	msg.CodeChallengeMethodsSupported = []string{"S256"}
	msg.GrantTypesSupported = []string{"authorization_code"}
	c.JSON(http.StatusOK, msg)
}
