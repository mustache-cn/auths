package auths

import (
	"errors"
	"reflect"
	"sort"
	"strings"
)

type Platform string

// Support platform
const (
	Wechat     Platform = "wechat"
	WechatMini Platform = "wechatMini"
	Weibo      Platform = "weibo"
	QQ         Platform = "qq"
	Github     Platform = "github"

	tagToken = "token"
	tagInfo  = "info"
)

/*
OAuthAccessToken
General AccountToken Return structure
*/
type OAuthAccessToken struct {
	AccessToken string
	ExpiresIn   string
	OpenId      string
	UnionId     string
}

/*
OAuthUserInfo
Common user information structure
*/
type OAuthUserInfo struct {
	Platform string
	Name     string
	OpenId   string
	UnionId  string
	Avatar   string
	Detail   string
}

/*
Auths
Define the Auths structure
*/
type Auths struct {
	clientId     string `token:"github,qq,weibo,wechat,wechatMini" info:"qq,wechatMini"`
	clientSecret string `token:"github,qq,weibo,wechat,wechatMini" info:"wechatMini"`
	platform     Platform

	code        string `token:"github,qq,weibo,wechat"`
	callBack    string `token:"qq,weibo" info:""`
	accessToken string `token:"" info:"github,qq,weibo,wechat,wechatMini"`
	openId      string `token:"" info:"qq,weibo,wechat"`
}

/*
AuthBuilder
Unified builder interface
*/
type AuthBuilder interface {
	AccessToken() (OAuthAccessToken, error)
	UserInfo() (OAuthUserInfo, error)
}

// NewAuths Instantiate the Auths
func NewAuths(clientId, clientSecret string, platform Platform) Auths {
	return Auths{
		clientId:     clientId,
		clientSecret: clientSecret,
		platform:     platform,
	}
}

// SetCode set code
func (a Auths) SetCode(code string) Auths {
	a.code = code
	return a
}

// SetCallBack set callback
func (a Auths) SetCallBack(callBack string) Auths {
	a.callBack = callBack
	return a
}

// SetAccessToken set account token
func (a Auths) SetAccessToken(accessToken string) Auths {
	a.accessToken = accessToken
	return a
}

// SetOpenId set openid
func (a Auths) SetOpenId(openId string) Auths {
	a.openId = openId
	return a
}

/*
Build
Instantiate the AuthBuilder according to the Auths
*/
func (a Auths) Build() AuthBuilder {
	var builder AuthBuilder
	switch a.platform {
	case Weibo:
		builder = NewWeibo(a)
		break
	case Github:
		builder = NewGithub(a)
		break
	case QQ:
		builder = NewQQ(a)
		break
	case Wechat:
		builder = NewWechat(a)
		break
	case WechatMini:
		builder = NewWechatMini(a)
		break
	}
	return builder
}

func validator(platform string, method string, s interface{}) error {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		arr := strings.Split(field.Tag.Get(method), ",")
		sort.Strings(arr)
		index := sort.SearchStrings(arr, platform)
		if index < len(arr) && arr[index] == platform && val.Field(i).IsZero() {
			return errors.New(field.Name + " is not empty")
		}
	}
	return nil
}
