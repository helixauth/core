package oauth

type TokenRequest struct {
	ClientID     string    `json:"client_id"`
	ClientSecret *string   `json:"client_secret"`
	Code         string    `json:"code"`
	GrantType    GrantType `json:"grant_type"`
	RedirectURI  string    `json:"redirect_uri"`
}
