# Auths
Third-party platform OAuth login integrated class library, support wechat, QQ, Weibo, Github, wechat small program...

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/mustache-cn/auths) [![GoDoc](https://pkg.go.dev/badge/github.com/mustache-cn/auths?utm_source=godoc)](https://godoc.org/github.com/mustache-cn/auths)[![License MIT](https://img.shields.io/github/license/mustache-cn/auths)](https://github.com/mustache-cn/auths)


License
======

Auths is licensed under the MIT License, Version 2.0. See [LICENSE](LICENSE) for the full license text

Features
========

- Chain call, multiple platforms access
- Pure and independent
- Simple and efficient access
- Support `Wechat, QQ, Weibo, Github, Wechat small program, More...`

Install
=======
`go get -u github.com/mustache-cn/auths`

Usage
======
`import "github.com/mustache-cn/auths"`

Basic Examples
=========
Basic Github oauth:

```go
	builder := auths.NewAuths(client_id, client_secret, auths.Github).SetCode("code").Build()
	token, err := builder.AccessToken()
	if err != nil {
		return
	}
	fmt.Println(token)
	info, err := builder.UserInfo()
	if err != nil {
		return
	}
	fmt.Println(info)
```

- If an error occurs, an error message is returned
- Changing the platform just requires a change platform parameters：`auths.Github auths.Wechat auths.WechatMini auths.Weibo auths.QQ`

Quirks
=======
## AccessToken Quirks

It is designed to return the access token rather than the user information directly because in some scenarios the access token is cached to avoid being restricted by the platform due to the number of requests.


## UserInfo Quirks

User information is used to associate users with their own system or login directly.

The unique ids returned by each platform are collectively called Openid.

Use the user information scenario.

- Associate your own system users with OpenId
- The OpenId returned is used directly as a user

Description of the Openids of each platform

- Github:`id -> OpenId,Name,Avatar `
- QQ: `id -> OpenId,Name,Avatar `
- WeChat: `OpenId,Name,Avatar,UnionId `  
- WechatMini：`OpenId,UnionId `  Only in the wechat open platform binding small program after UnionId
- Weibo：`id -> OpenId,Name,Avatar `


