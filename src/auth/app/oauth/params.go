package oauth

type Params struct {
	ClientID     string       `form:"client_id" binding:"required"`
	ResponseType ResponseType `form:"response_type" binding:"required"`
	Scope        string       `form:"scope" binding:"required"`
	State        *string      `form:"state"`
	Nonce        *string      `form:"nonce"`
	RedirectURI  *string      `form:"redirect_uri"`
	Prompt       *Prompt      `form:"prompt"`
}

func (p Params) IsSignUp() bool {
	return p.Prompt != nil && *p.Prompt == PromptCreate
}
