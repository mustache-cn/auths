package auths

import (
	"errors"
	"github.com/mustache-cn/https"
	"strconv"
)

const (
	weiboAccessTokenUrl = "https://api.weibo.com/oauth2/access_token"
	weiboUserInfoUrl    = "https://api.weibo.com/2/users/show.json"
)

type weibo struct {
	auths Auths
}

func NewWeibo(a Auths) *weibo {
	return &weibo{auths: a}
}

func (w *weibo) AccessToken() (OAuthAccessToken, error) {
	result := OAuthAccessToken{}
	if err := validator(string(w.auths.platform), tagToken, w.auths); err != nil {
		return result, err
	}
	// request data
	res, err := https.NewClient(weiboAccessTokenUrl).
		SetContentType(https.FormType).
		AddHeader("Accept", "application/x-www-form-urlencoded").
		AddParam("client_id", w.auths.clientId).
		AddParam("client_secret", w.auths.clientSecret).
		AddParam("code", w.auths.code).
		AddParam("redirect_uri", w.auths.callBack).
		AddParam("grant_type", "authorization_code").
		Post()
	// request failed
	if err != nil {
		return result, err
	}

	// response struct
	resData := struct {
		// error result
		Error            string `json:"error,omitempty"`
		ErrorCode        int64  `json:"error_code,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		// success result
		AccessToken string `json:"access_token"`
		RemindIn    string `json:"remind_in"`
		ExpiresIn   int64  `json:"expires_in"`
		UID         string `json:"uid"`
		IsRealName  string `json:"isRealName"`
	}{}
	// response to json
	if err := res.JSON(&resData); err != nil {
		return result, err
	}
	// error
	if resData.ErrorCode > 0 {
		return result, errors.New(resData.ErrorDescription)
	}
	// auths value
	w.auths.openId = resData.UID
	// result value
	result.AccessToken = resData.AccessToken
	result.OpenId = resData.UID
	return result, nil
}

func (w *weibo) UserInfo() (OAuthUserInfo, error) {
	result := OAuthUserInfo{Platform: string(w.auths.platform)}
	if err := validator(string(w.auths.platform), tagInfo, w.auths); err != nil {
		return result, err
	}
	// request data
	res, err := https.NewClient(weiboUserInfoUrl).
		AddParam("access_token", w.auths.accessToken).
		AddParam("uid", w.auths.openId).
		Get()
	// request failed
	if err != nil {
		return result, err
	}

	// detail
	result.Detail = res.String()

	resData := struct {
		// error result
		Error     string `json:"error"`
		ErrorCode int64  `json:"error_code"`
		Request   string `json:"request"`
		// success result
		ID     int64  `json:"id"`
		Name   string `json:"screen_name"`
		Avatar string `json:"avatar_hd"`
		Gender string `json:"gender"`
	}{}

	// response to json
	if err := res.JSON(&resData); err != nil {
		return result, err
	}
	// error
	if resData.ErrorCode > 0 {
		return result, errors.New(resData.Error)
	}

	// return data
	result.OpenId = strconv.FormatInt(resData.ID, 10)
	result.Name = resData.Name
	result.Avatar = resData.Avatar
	return result, err
}
