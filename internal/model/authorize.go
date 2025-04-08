package model

type AuthorizeRequest struct {
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	ResponseType string `json:"response_type"`
	State        string `json:"state"`
}

type AuthorizeReply struct {
	Code        string `json:"code"`
	State       string `json:"state"`
	RedirectURI string `json:"redirect_uri"`
}

type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	DataType     string `json:"data_type"`
	Callback     string `json:"callback"`
	RefreshToken string `json:"refresh_token"`
}

type TokenReply struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type AddClientRequest struct {
	ClientName    string   `json:"client_name"`
	HomePage      string   `json:"home_page"`
	ClientSummary string   `json:"client_summary"`
	CallbackURL   string   `json:"callback_url"`
	Scopes        []string `json:"scopes"`
	UserID        string   `json:"user_id"`
}

type AddClientReply struct {
	ClientID string `json:"client_id"`
}
