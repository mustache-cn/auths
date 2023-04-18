package auths

import (
	"errors"
	"github.com/mustache-cn/https"
)

const (
	wechatAccessTokenUrl = "https://api.weixin.qq.com/sns/oauth2/access_token"
	wechatUserInfoUrl    = "https://api.weixin.qq.com/sns/userinfo"
)

type wechat struct {
	auths Auths
}

func NewWechat(a Auths) *wechat {
	return &wechat{auths: a}
}

func (w *wechat) AccessToken() (OAuthAccessToken, error) {
	result := OAuthAccessToken{}
	if err := validator(string(w.auths.platform), tagToken, w.auths); err != nil {
		return result, err
	}
	// request data
	res, err := https.NewClient(wechatAccessTokenUrl).
		SetContentType(https.JsonType).
		AddParam("appid", w.auths.clientId).
		AddParam("secret", w.auths.clientSecret).
		AddParam("code", w.auths.code).
		AddParam("grant_type", "authorization_code").
		Post()

	// request failed
	if err != nil {
		return result, err
	}

	//response struct
	resData := struct {
		// error result
		ErrCode int32  `json:"errcode,omitempty"`
		ErrMsg  string `json:"errmsg,omitempty"`
		// success result
		AccessToken  string `json:"access_token,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
		ExpiresIn    int64  `json:"expires_in"`
		OpenId       string `json:"openid"`
		UnionId      string `json:"unionid"`
	}{}

	// response to json
	if err := res.JSON(&resData); err != nil {
		return result, err
	}

	// failed
	if resData.ErrCode > 0 {
		return result, errors.New(resData.ErrMsg)
	}

	// auths value
	w.auths.accessToken = resData.AccessToken
	w.auths.openId = resData.OpenId
	// result value
	result.AccessToken = resData.AccessToken
	result.OpenId = resData.OpenId
	result.UnionId = resData.UnionId
	return result, nil
}

func (w *wechat) UserInfo() (OAuthUserInfo, error) {
	result := OAuthUserInfo{Platform: string(w.auths.platform)}
	if err := validator(string(w.auths.platform), tagInfo, w.auths); err != nil {
		return result, err
	}
	// request data
	res, err := https.NewClient(wechatUserInfoUrl).
		AddParam("access_token", w.auths.accessToken).
		AddParam("openid", w.auths.openId).
		AddParam("lang", "zh_CN").
		Get()
	// request failed
	if err != nil {
		return result, err
	}

	// detail
	result.Detail = res.String()

	// response struct
	respData := struct {
		// error result
		ErrCode int64  `json:"errcode,omitempty"`
		ErrMsg  string `json:"errmsg,omitempty"`
		// success reslut
		OpenId     string   `json:"openid,omitempty"`
		Nickname   string   `json:"nickname,omitempty"`
		Sex        int      `json:"sex,omitempty"`
		Province   string   `json:"province,omitempty"`
		City       string   `json:"city,omitempty"`
		Country    string   `json:"country,omitempty"`
		HeadImgURL string   `json:"headimgurl,omitempty"`
		Privilege  []string `json:"privilege,omitempty"`
		UnionId    string   `json:"unionid,omitempty"`
	}{}

	// response to json
	if err := res.JSON(&respData); err != nil {
		return result, err
	}

	// failed
	if respData.ErrCode != 0 {
		return result, errors.New(respData.ErrMsg)
	}

	// return data
	result.OpenId = respData.OpenId
	result.UnionId = respData.UnionId
	result.Name = respData.Nickname
	result.Avatar = respData.HeadImgURL
	return result, err
}
