package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Gitlab struct {
	id          string
	secret      string
	redirectURI string
	baseURL     string
}

func (oauth *Gitlab) BaseURL() string {
	return oauth.baseURL
}

func (oauth *Gitlab) SetBaseURL(baseURL string) {
	oauth.baseURL = baseURL
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

type GitlabUser struct {
	ID                int       `json:"id"`
	Login             string    `json:"login"`
	Name              string    `json:"name"`
	AvatarURL         string    `json:"avatar_url"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	Stared            int       `json:"stared"`
	Watched           int       `json:"watched"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Email             string    `json:"email"`
	PrivateToken      string    `json:"private_token"`
	TotalRepos        int       `json:"total_repos"`
	OwnedRepos        int       `json:"owned_repos"`
	TotalPrivateRepos int       `json:"total_private_repos"`
	OwnedPrivateRepos int       `json:"owned_private_repos"`
	PrivateGists      int       `json:"private_gists"`

	Message string `json:"message"`
}

func (oauth *Gitlab) User(token *GitlabAccessToken) (user *GitlabUser, err error) {
	resp, err := http.Get(oauth.baseURL + "/api/v4/user?access_token=" + token.AccessToken)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	user = new(GitlabUser)
	err = json.Unmarshal(body, user)
	if err != nil {
		return
	}

	if user.Message != "" {
		err = fmt.Errorf(user.Message)
		return
	}

	return user, err
}
