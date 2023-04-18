package auths

import (
	"errors"
	"github.com/mustache-cn/https"
)

type qq struct {
	auths Auths
}

const (
	qqAccessTokenUrl = "https://graph.qq.com/oauth2.0/token"
	qqOpenIdUrl      = "https://graph.qq.com/oauth2.0/me"
	qqUserInfoUrl    = "https://graph.qq.com/user/get_user_info"
)

func NewQQ(a Auths) *qq {
	return &qq{auths: a}
}

func (q *qq) AccessToken() (OAuthAccessToken, error) {
	result := OAuthAccessToken{}
	if err := validator(string(q.auths.platform), tagToken, q.auths); err != nil {
		return result, err
	}
	// request data
	res, err := https.NewClient(qqAccessTokenUrl).
		//AddHeader("Accept", "application/json").
		AddParam("grant_type", "authorization_code").
		AddParam("client_id", q.auths.clientId).
		AddParam("client_secret", q.auths.clientSecret).
		AddParam("code", q.auths.code).
		AddParam("redirect_uri", q.auths.callBack).
		AddParam("fmt", "json").
		Get()
	// request failed
	if err != nil {
		return result, err
	}
	// response struct
	resData := struct {
		Error            int32  `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		// success
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    string `json:"expires_in"`
	}{}
	// response to json
	if err := res.JSON(&resData); err != nil {
		return result, errors.New(resData.ErrorDescription)
	}

	// get openid and set auths
	if q.auths.openId, err = q.openId(resData.AccessToken); err != nil {
		return result, err
	}

	// result value
	result.AccessToken = resData.AccessToken
	result.OpenId = q.auths.openId
	return result, nil
}

func (q *qq) UserInfo() (OAuthUserInfo, error) {
	result := OAuthUserInfo{Platform: string(q.auths.platform)}
	if err := validator(string(q.auths.platform), tagInfo, q.auths); err != nil {
		return result, err
	}
	//request data
	res, err := https.NewClient(qqUserInfoUrl).
		SetContentType(https.FormType).
		AddParam("access_token", q.auths.accessToken).
		AddParam("oauth_consumer_key", q.auths.clientId).
		AddParam("openid", q.auths.openId).
		Get()
	// request failed
	if err != nil {
		return result, err
	}

	// detail
	result.Detail = res.String()

	//response struct
	resData := struct {
		// error result
		Ret int    `json:"ret,omitempty"`
		Msg string `json:"Met,omitempty"`
		// success result
		ID       string `json:"id"`
		Avatar   string `json:"figureurl_qq"`
		AvatarHD string `json:"figureurl_qq_2"`
		Nickname string `json:"nickname"`
		Gender   string `json:"gender"`
	}{}

	// request ro json
	if err := res.JSON(&resData); err != nil {
		return result, err
	}

	// failed
	if resData.Ret != 0 {
		return result, errors.New(resData.Msg)
	}

	// return data
	result.OpenId = q.auths.openId
	result.Name = resData.Nickname
	result.Avatar = resData.Avatar

	return result, nil
}

// openId return openid by access token
func (q *qq) openId(accessToken string) (string, error) {
	// request data
	res, err := https.NewClient(qqOpenIdUrl).
		AddParam("access_token", accessToken).
		AddParam("fmt", "json").
		Get()
	// request failed
	if err != nil {
		return "", err
	}
	// response struct
	resData := struct {
		Error            int    `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`

		ClientID string `json:"client_id"`
		OpenId   string `json:"openid"`
	}{}

	// response to json
	if err := res.JSON(&resData); err != nil {
		return "", err
	}

	// error
	if resData.Error > 0 {
		return "", errors.New(resData.ErrorDescription)
	}

	return resData.OpenId, nil
}
