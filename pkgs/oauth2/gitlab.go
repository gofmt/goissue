package oauth2

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Gitlab struct {
	id          string
	secret      string
	redirectURI string
	baseURL     string
}

func (oauth *Gitlab) RedirectURI() string {
	return oauth.redirectURI
}

func (oauth *Gitlab) SetRedirectURI(redirectURI string) {
	oauth.redirectURI = redirectURI
}

func NewGitlabClient(id, secret string) *Gitlab {
	return &Gitlab{
		id:      id,
		secret:  secret,
		baseURL: "https://gitlab.com",
	}
}

// Authorize 请求用户授权
func (oauth *Gitlab) Authorize(state string) string {
	parse := url.Values{}
	parse.Set("client_id", oauth.id)
	parse.Set("redirect_uri", oauth.redirectURI)
	parse.Set("response_type", "code")
	parse.Set("state", state)
	return oauth.baseURL + "/oauth/authorize?" + parse.Encode()
}

type GitlabAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (oauth *Gitlab) AccessToken(code string) (token *GitlabAccessToken, err error) {
	parse := url.Values{}
	parse.Set("client_id", oauth.id)
	parse.Set("client_secret", oauth.secret)
	parse.Set("code", code)
	parse.Set("grant_type", "authorization_code")
	parse.Set("redirect_uri", oauth.redirectURI)

	resp, err := http.PostForm(oauth.baseURL+"/oauth/token", parse)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	token = new(GitlabAccessToken)
	err = json.Unmarshal(body, token)
	if err != nil {
		return
	}

	return
}

func (oauth *Gitlab) User(token *GitlabAccessToken) (user, err error) {
	resp, err := http.Get(oauth.baseURL + "/api/v4/user?access_token=" + token.AccessToken)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// TODO 返回值

	err = json.Unmarshal(body, user)
	if err != nil {
		return
	}

	return user, err
}
