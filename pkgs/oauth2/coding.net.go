package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CodingNet struct {
	id     string
	secret string
	scope  string
}

func (oauth *CodingNet) Scope() string {
	return oauth.scope
}

func (oauth *CodingNet) SetScope(scope string) {
	oauth.scope = scope
}

func NewCodingNetClient(id, secret string) *CodingNet {
	return &CodingNet{
		id:     id,
		secret: secret,
		scope:  "user,user:email",
	}
}

// Authorize 请求用户授权
func (oauth *CodingNet) Authorize(state string) string {
	parse := url.Values{}
	parse.Set("client_id", oauth.id)
	parse.Set("scope", oauth.scope)
	parse.Set("state", state)
	parse.Set("response_type", "code")
	return fmt.Sprintf("https://coding.net/oauth_authorize.html?%s", parse.Encode())
}

type CodingNetAccessToken struct {
	ExpiresIn    int64  `json:"expires_in,string"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	CodingNetRespErr
}

type CodingNetRespErr struct {
	Code int               `json:"code"`
	Msg  map[string]string `json:"msg"`
}

func (oauth *CodingNet) AccessToken(code string) (token *CodingNetAccessToken, err error) {
	parse := url.Values{}
	parse.Set("client_id", oauth.id)
	parse.Set("client_secret", oauth.secret)
	parse.Set("code", code)
	parse.Set("grant_type", "authorization_code")
	urlAddr := fmt.Sprintf("https://coding.net/api/oauth/access_token?%s", parse.Encode())

	resp, err := http.Post(urlAddr, "", nil)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	token = new(CodingNetAccessToken)
	err = json.Unmarshal(body, token)
	if err != nil {
		return
	}

	if token.Code != 0 {
		msg := ""
		for _, e := range token.Msg {
			msg = e
			break
		}
		err = fmt.Errorf(msg)
		return
	}

	return
}

type CodingNetUserResp struct {
	CodingNetRespErr
	Data *CodingNetUser `json:"data"`
}

type CodingNetUser struct {
	TagsStr        string  `json:"tags_str"`
	Tags           string  `json:"tags"`
	JobStr         string  `json:"job_str"`
	Job            int     `json:"job"`
	Sex            int     `json:"sex"`
	Phone          string  `json:"phone"`
	Birthday       string  `json:"birthday"`
	Location       string  `json:"location"`
	Company        string  `json:"company"`
	Slogan         string  `json:"slogan"`
	Website        string  `json:"website"`
	Introduction   string  `json:"introduction"`
	Avatar         string  `json:"avatar"`
	Gravatar       string  `json:"gravatar"`
	Lavatar        string  `json:"lavatar"`
	CreatedAt      int64   `json:"created_at"`
	LastLoginedAt  int64   `json:"last_logined_at"`
	LastActivityAt int64   `json:"last_activity_at"`
	GlobalKey      string  `json:"global_key"`
	Name           string  `json:"name"`
	NamePinyin     string  `json:"name_pinyin"`
	UpdatedAt      int64   `json:"updated_at"`
	Path           string  `json:"path"`
	Status         int     `json:"status"`
	Email          string  `json:"email"`
	IsMember       int     `json:"is_member"`
	ID             int     `json:"id"`
	PointsLeft     float64 `json:"points_left"`
	Vip            int     `json:"vip"`
	VipExpiredAt   int64   `json:"vip_expired_at"`
	Skills         []struct {
		SkillName string `json:"skillName"`
		SkillID   int    `json:"skillId"`
		Level     int    `json:"level"`
	} `json:"skills"`
	Degree           int    `json:"degree"`
	School           string `json:"school"`
	FollowsCount     int    `json:"follows_count"`
	FansCount        int    `json:"fans_count"`
	TweetsCount      int    `json:"tweets_count"`
	PhoneCountryCode string `json:"phone_country_code"`
	Country          string `json:"country"`
	Followed         bool   `json:"followed"`
	Follow           bool   `json:"follow"`
	IsPhoneValidated bool   `json:"is_phone_validated"`
	EmailValidation  int    `json:"email_validation"`
	PhoneValidation  int    `json:"phone_validation"`
	TwofaEnabled     int    `json:"twofa_enabled"`
	WechatName       string `json:"wechat_name"`
	QcloudName       string `json:"qcloud_name"`
	IsTencentUser    bool   `json:"is_tencent_user"`
	DevUserID        int    `json:"dev_user_id"`
	IsWelcomed       bool   `json:"is_welcomed"`
	HaveDemo         bool   `json:"have_demo"`
}

// User 获取用户信息
func (oauth *CodingNet) User(token *CodingNetAccessToken) (user *CodingNetUser, err error) {
	urlAddr := fmt.Sprintf("https://coding.net/api/current_user?access_token=%s", token.AccessToken)

	resp, err := http.Get(urlAddr)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	data := new(CodingNetUserResp)
	err = json.Unmarshal(body, data)
	if err != nil {
		return
	}

	if data.Code != 0 {
		msg := ""
		for _, e := range data.Msg {
			msg = e
			break
		}
		err = fmt.Errorf(msg)
		return
	}

	return data.Data, err
}
