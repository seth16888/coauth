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
	ClientID     string `json:"client_id"`     // 	OAuth2客户ID
	ClientSecret string `json:"client_secret"` // OAuth2密钥
	GrantType    string `json:"grant_type"`    // 授权方式：authorization_code或者refresh_token
	Code         string `json:"code"`          // 调用 /action/oauth2/authorize 接口返回的授权码(grant_type为authorization_code时必选)
	RedirectURI  string `json:"redirect_uri"`  // 回调地址
	DataType     string `json:"data_type"`     // 返回数据类型['json'|'jsonp'|'xml']
	RefreshToken string `json:"refresh_token"` // 上次调用 /action/oauth2/token 接口返回的refresh_token(grant_type为refresh_token时必选)
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
