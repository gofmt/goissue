package oauth2

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Github struct {
	id          string
	secret      string
	scope       string
	redirectURI string
}

func (oauth *Github) Scope() string {
	return oauth.scope
}

func (oauth *Github) SetScope(scope string) {
	oauth.scope = scope
}

func (oauth *Github) RedirectURI() string {
	return oauth.redirectURI
}

func (oauth *Github) SetRedirectURI(redirectURI string) {
	oauth.redirectURI = redirectURI
}

func NewGithubClient(id, secret string) *Github {
	return &Github{
		id:     id,
		secret: secret,
		scope:  "read:user,user:email",
	}
}

// Authorize 请求用户授权
func (oauth *Github) Authorize(state string) string {
	params := url.Values{}
	params.Set("client_id", oauth.id)
	params.Set("scope", oauth.scope)
	params.Set("state", state)
	params.Set("redirect_uri", oauth.redirectURI)
	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", params.Encode())
}

type GithubAccessToken struct {
	AccessToken string `json:"access_token" xml:"access_token"`
	TokenType   string `json:"token_type" xml:"token_type"`
	Scope       string `json:"scope" xml:"scope"`
}

// AccessToken 请求授权服务器授权
func (oauth *Github) AccessToken(code string) (token *GithubAccessToken, err error) {
	params := url.Values{}
	params.Set("client_id", oauth.id)
	params.Set("client_secret", oauth.secret)
	params.Set("code", code)
	urlAddr := fmt.Sprintf("https://github.com/login/oauth/access_token?%s", params.Encode())

	resp, err := http.Post(urlAddr, "", nil)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	token = new(GithubAccessToken)
	accept := resp.Header.Get("Accept")
	switch accept {
	case "application/json":
		err = json.Unmarshal(body, token)
	case "application/xml":
		err = xml.Unmarshal(body, token)
	default:
		v, err2 := url.ParseQuery(string(body))
		if err2 != nil {
			err = err2
			return
		}
		token.AccessToken = v.Get("access_token")
		token.Scope = v.Get("scope")
		token.TokenType = v.Get("token_type")
	}
	if err != nil {
		return
	}

	return
}

type GithubUser struct {
	Login                   string    `json:"login"`
	ID                      int       `json:"id"`
	NodeID                  string    `json:"node_id"`
	AvatarURL               string    `json:"avatar_url"`
	GravatarID              string    `json:"gravatar_id"`
	URL                     string    `json:"url"`
	HTMLURL                 string    `json:"html_url"`
	FollowersURL            string    `json:"followers_url"`
	FollowingURL            string    `json:"following_url"`
	GistsURL                string    `json:"gists_url"`
	StarredURL              string    `json:"starred_url"`
	SubscriptionsURL        string    `json:"subscriptions_url"`
	OrganizationsURL        string    `json:"organizations_url"`
	ReposURL                string    `json:"repos_url"`
	EventsURL               string    `json:"events_url"`
	ReceivedEventsURL       string    `json:"received_events_url"`
	Type                    string    `json:"type"`
	SiteAdmin               bool      `json:"site_admin"`
	Name                    string    `json:"name"`
	Company                 string    `json:"company"`
	Blog                    string    `json:"blog"`
	Location                string    `json:"location"`
	Email                   string    `json:"email"`
	Hireable                bool      `json:"hireable"`
	Bio                     string    `json:"bio"`
	PublicRepos             int       `json:"public_repos"`
	PublicGists             int       `json:"public_gists"`
	Followers               int       `json:"followers"`
	Following               int       `json:"following"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	PrivateGists            int       `json:"private_gists"`
	TotalPrivateRepos       int       `json:"total_private_repos"`
	OwnedPrivateRepos       int       `json:"owned_private_repos"`
	DiskUsage               int       `json:"disk_usage"`
	Collaborators           int       `json:"collaborators"`
	TwoFactorAuthentication bool      `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		PrivateRepos  int    `json:"private_repos"`
		Collaborators int    `json:"collaborators"`
	} `json:"plan"`

	GithubRespErr
}

type GithubRespErr struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

// User 获取用户信息
func (oauth *Github) User(token *GithubAccessToken) (user *GithubUser, err error) {
	urlAddr := fmt.Sprintf("https://api.github.com/user?access_token=%s", token.AccessToken)

	resp, err := http.Get(urlAddr)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	user = new(GithubUser)
	err = json.Unmarshal(body, user)
	if err != nil {
		return
	}

	if user.Message != "" {
		err = fmt.Errorf(user.Message)
		return
	}
	return
}
