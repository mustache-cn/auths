package auths

import (
	"errors"
	"fmt"
	"github.com/mustache-cn/https"
	"strconv"
)

type github struct {
	auths Auths
}

const (
	githubAccessTokenUrl = "https://github.com/login/oauth/access_token"
	githubUserInfoUrl    = "https://api.github.com/user"
)

func NewGithub(a Auths) *github {
	return &github{auths: a}
}

func (g *github) AccessToken() (OAuthAccessToken, error) {
	result := OAuthAccessToken{}
	if err := validator(string(g.auths.platform), tagToken, g.auths); err != nil {
		return result, err
	}
	// request data
	res, err := https.NewClient(githubAccessTokenUrl).
		SetContentType(https.JsonType).
		AddHeader("Accept", "application/json").
		AddParam("client_id", g.auths.clientId).
		AddParam("client_secret", g.auths.clientSecret).
		AddParam("code", g.auths.code).
		AddParam("redirect_uri", g.auths.callBack).
		Post()
	// request failed
	if err != nil {
		return result, err
	}

	resData := struct {
		// error result
		Error            string `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		ErrorURI         string `json:"error_uri,omitempty"`
		// success result
		AccessToken string `json:"access_token,omitempty"`
		Scope       string `json:"scope,omitempty"`
		TokenType   string `json:"token_type,omitempty"`
	}{}

	// response to json
	if err := res.JSON(&resData); err != nil {
		return result, err
	}
	if resData.Error != "" {
		return result, errors.New(resData.ErrorDescription)
	}

	// auths value
	g.auths.accessToken = resData.AccessToken

	// result value
	result.AccessToken = resData.AccessToken
	return result, nil
}

func (g *github) UserInfo() (OAuthUserInfo, error) {
	result := OAuthUserInfo{Platform: string(g.auths.platform)}
	if err := validator(string(g.auths.platform), tagInfo, g.auths); err != nil {
		return result, err
	}
	// request data
	res, err := https.NewClient(githubUserInfoUrl).
		AddHeader("Authorization", fmt.Sprintf("bearer %s", g.auths.accessToken)).
		Get()
	// request failed
	if err != nil {
		return result, err
	}

	// detail value
	result.Detail = res.String()
	// response struct
	resData := struct {
		ID        int64  `json:"id"`
		AvatarURL string `json:"avatar_url"`
		Name      string `json:"login"`
		Message   string `json:"message,omitempty"`
	}{}
	// response to json
	if err := res.JSON(&resData); err != nil {
		return result, err
	}
	if resData.Message != "" {
		return result, errors.New(resData.Message)
	}

	// result value
	result.OpenId = strconv.FormatInt(resData.ID, 10)
	result.Name = resData.Name
	result.Avatar = resData.AvatarURL
	return result, nil
}
