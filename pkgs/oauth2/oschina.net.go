package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type OSChinaNet struct {
	id     string
	secret string
}

func NewOSChinaNetClient(id, secret string) *OSChinaNet {
	return &OSChinaNet{
		id:     id,
		secret: secret,
	}
}

// Authorize 请求用户授权
// redirect_uri 必填参数 回调地址
// state 可选参数，回调时，将会带回
func (oauth *OSChinaNet) Authorize(parse url.Values) string {
	parse.Set("client_id", oauth.id)
	parse.Set("response_type", "code")
	return fmt.Sprintf("https://www.oschina.net/action/oauth2/authorize?%s", parse.Encode())
}

type OSChinaNetAccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	UID          int    `json:"uid"`
	OSChinaNetRespErr
}

type OSChinaNetRespErr struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// AccessToken 请求获取服务器授权
// code 必填 Authorize 回调后的参数
func (oauth *OSChinaNet) AccessToken(code string) (token *OSChinaNetAccessToken, err error) {
	parse := url.Values{}
	parse.Set("client_id", oauth.id)
	parse.Set("client_secret", oauth.secret)
	parse.Set("code", code)
	parse.Set("redirect_uri", "")
	parse.Set("dataType", "json")

	urlAddr := fmt.Sprintf("https://www.oschina.net/action/openapi/token?%s", parse.Encode())

	resp, err := http.Post(urlAddr, "", nil)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	token = new(OSChinaNetAccessToken)
	err = json.Unmarshal(body, token)
	if err != nil {
		return
	}

	if token.Error != "" {
		err = fmt.Errorf(token.ErrorDescription)
	}
	return
}

type OSChinaNetUser struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Avatar   string `json:"avatar"`
	Location string `json:"location"`
	URL      string `json:"url"`

	OSChinaNetRespErr
}

// User 获取用户信息
func (oauth *OSChinaNet) User(token *OSChinaNetAccessToken) (user *OSChinaNetUser, err error) {
	urlAddr := fmt.Sprintf("https://www.oschina.net/action/openapi/user?dataType=json&access_token=%s", token.AccessToken)

	resp, err := http.Get(urlAddr)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	user = new(OSChinaNetUser)
	err = json.Unmarshal(body, user)
	if err != nil {
		return
	}

	if user.Error != "" {
		err = fmt.Errorf(user.ErrorDescription)
	}
	return user, err
}
