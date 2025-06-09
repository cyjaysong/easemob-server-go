package easemob_server_go

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GetAppTokenRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Application string `json:"application"`
}

// GetAppToken 获取AppToken
// 若ttl < 0 以环信即时通讯云控制台的用户认证页面的 token 有效期的设置为准
// 若ttl = 0 则 token 永久有效
// 注意: VIP5集群 ttl 单位为毫秒, 其他集群 ttl 单位为秒
func (c *Client) GetAppToken(ctx context.Context, ttl int64) (appToken *GetAppTokenRes, err error) {
	data := map[string]any{"grant_type": "client_credentials", "client_id": c.clientId, "client_secret": c.clientSecret}
	if ttl >= 0 {
		data["ttl"] = ttl
	}
	appToken = new(GetAppTokenRes)
	if err = c.doReq(ctx, http.MethodPost, "token", nil, data, appToken); err != nil {
		return nil, err
	}
	c.SetAppToken(appToken.AccessToken)
	return
}

// SetAppToken 设置AppToken, 一般用于分布式部署的时候，业务层统一维护AppToken，为节点设置AppToken
func (c *Client) SetAppToken(appToken string) {
	c.appToken = appToken
}

type GetUserTokenRes struct {
	AccessToken string               `json:"access_token"`
	ExpiresIn   int                  `json:"expires_in"`
	User        GetUserTokenUserInfo `json:"user"`
}

type GetUserTokenUserInfo struct {
	UUID      string `json:"uuid"`
	Type      string `json:"type"`
	Created   int64  `json:"created"`
	Modified  int64  `json:"modified"`
	Username  string `json:"username"`
	Activated bool   `json:"activated"`
}

// GetUserToken 获取UserToken
// 若ttl < 0 以环信即时通讯云控制台的用户认证页面的 token 有效期的设置为准
// 若ttl = 0 则 token 永久有效
// 注意: VIP5集群 ttl 单位为毫秒, 其他集群 ttl 单位为秒
func (c *Client) GetUserToken(ctx context.Context, username, password string, autoCreateUser bool, ttl int64) (userToken *GetUserTokenRes, err error) {
	data := map[string]any{"grant_type": "inherit", "username": username}
	if len(password) > 0 {
		data["grant_type"], data["password"] = "password", password
	} else if autoCreateUser {
		data["autoCreateUser"] = true
	}
	if ttl >= 0 {
		data["ttl"] = ttl
	}
	userToken = new(GetUserTokenRes)
	if err = c.doReq(ctx, http.MethodPost, "token", nil, data, userToken); err != nil {
		return nil, err
	}
	return
}

// CreateUserToken 生成动态的UserToken
// ttl 单位为秒
func (c *Client) CreateUserToken(username string, ttl int64) (userToken string) {
	if ttl <= 0 {
		panic("`ttl` must greater than 0")
	}
	timeNow := time.Now().Unix()
	signature := fmt.Sprintf("%s%s%s%d%d%s", c.clientId, c.appKey, username, timeNow, ttl, c.clientSecret)
	jsonMap := map[string]any{"appkey": c.appKey, "userId": username, "curTime": timeNow, "ttl": ttl,
		"signature": fmt.Sprintf("%x", sha256.Sum256([]byte(signature)))}
	jsonBytes, _ := json.Marshal(jsonMap)
	userToken = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("dt-%s", jsonBytes)))
	return
}
