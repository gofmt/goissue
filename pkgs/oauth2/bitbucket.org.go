package oauth2

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type BitBucket struct {
	id     string
	secret string
	scope  string
}

func NewBitBucketClient(id, secret string) *BitBucket {
	return &BitBucket{
		id:     id,
		secret: secret,
		scope:  "account",
	}
}

// Authorize 请求用户授权
func (oauth *BitBucket) Authorize(state string) string {
	return fmt.Sprintf("https://bitbucket.org/site/oauth2/authorize?client_id=%s&response_type=code&state=%s", oauth.id, state)
}

type BitBucketAccessToken struct {
	AccessToken string `json:"access_token" xml:"access_token"`
	TokenType   string `json:"token_type" xml:"token_type"`
	Scope       string `json:"scope" xml:"scope"`
}

func (oauth *BitBucket) AccessToken(code string) (token *BitBucketAccessToken, err error) {
	urlAddr := fmt.Sprintf("https://bitbucket.org/site/oauth2/access_token")
	parse := fmt.Sprintf("grant_type=authorization_code&code=%s", code)

	req, err := http.NewRequest(http.MethodGet, urlAddr, strings.NewReader(parse))
	if err != nil {
		return
	}
	req.SetBasicAuth(oauth.id, oauth.secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	_ = body

	return
}
