package auths

import (
	"errors"
	"github.com/mustache-cn/https"
)

const (
	wechatMiniAccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token"
	wechatMiniUserInfoUrl    = "https://api.weixin.qq.com/sns/jscode2session"
	wechatMiniQRUrl          = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
)

type wechatMini struct {
	auths Auths
}

func NewWechatMini(a Auths) *wechatMini {
	return &wechatMini{auths: a}
}

func (w *wechatMini) AccessToken() (OAuthAccessToken, error) {
	result := OAuthAccessToken{}
	if err := validator(string(w.auths.platform), tagToken, w.auths); err != nil {
		return result, err
	}
	// request data
	res, err := https.NewClient(wechatMiniAccessTokenUrl).
		AddParam("appid", w.auths.clientId).
		AddParam("secret", w.auths.clientSecret).
		AddParam("grant_type", "client_credential").
		Get()
	// request failed
	if err != nil {
		return result, err
	}

	// response struct
	resData := struct {
		// error result
		ErrCode int32  `json:"errcode,omitempty"`
		ErrMsg  string `json:"errmsg,omitempty"`
		// success result
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}{}
	// response to json
	if err := res.JSON(&resData); err != nil {
		return result, err
	}
	// error
	if resData.ErrCode > 0 {
		return result, errors.New(resData.ErrMsg)
	}

	// auths value
	w.auths.accessToken = resData.AccessToken
	// result value
	result.AccessToken = resData.AccessToken
	return result, nil
}

func (w *wechatMini) UserInfo() (OAuthUserInfo, error) {
	result := OAuthUserInfo{Platform: string(w.auths.platform)}
	if err := validator(string(w.auths.platform), tagInfo, w.auths); err != nil {
		return result, err
	}
	// request data
	res, err := https.NewClient(wechatMiniUserInfoUrl).
		AddParam("appid", w.auths.clientId).
		AddParam("secret", w.auths.clientSecret).
		AddParam("grant_type", "authorization_code").
		AddParam("js_code", w.auths.accessToken).
		Get()
	// request failed
	if err != nil {
		return result, err
	}

	// detail
	result.Detail = res.String()
	// response struct
	resData := struct {
		// error result
		ErrCode int32  `json:"errcode,omitempty"`
		ErrMsg  string `json:"errmsg,omitempty"`
		// success result
		SessionKey string `json:"session_key"`
		UnionId    string `json:"unionid"`
		OpenId     string `json:"openid"`
	}{}

	// response to json
	if err := res.JSON(&resData); err != nil {
		return result, err
	}
	// error
	if resData.ErrCode > 0 {
		return result, errors.New(resData.ErrMsg)
	}

	// return data
	result.OpenId = resData.OpenId
	result.UnionId = resData.UnionId
	return result, err
}
