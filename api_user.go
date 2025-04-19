package easemob_server_go

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type UserEntity struct {
	Uuid      string `json:"uuid"`
	Type      string `json:"type"`
	Created   int64  `json:"created"`
	Modified  int64  `json:"modified"`
	Username  string `json:"username"`
	Activated bool   `json:"activated"`
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Nickname string `json:"nickname,omitempty"`
}

// AddUser 添加用户
func (c *Client) AddUser(ctx context.Context, users ...NewUser) (res *UserBaseRes[[]UserEntity], err error) {
	if len(users) == 0 {
		return nil, errors.New("minimum count of user is 1")
	} else if len(users) > 60 {
		return nil, errors.New("maximum count of user is 60")
	}
	var data any
	if len(users) == 1 {
		data = users[0]
	}
	res = new(UserBaseRes[[]UserEntity])
	if err = c.doReq(ctx, http.MethodPost, "users", nil, data, res); err != nil {
		return nil, err
	}
	return
}

// DelUser 删除用户
func (c *Client) DelUser(ctx context.Context, username string) (res *UserBaseRes[[]UserEntity], err error) {
	pathSuffix := fmt.Sprintf("users/%s", username)
	res = new(UserBaseRes[[]UserEntity])
	if err = c.doReq(ctx, http.MethodDelete, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

// BatchDelUser 批量删除用户
func (c *Client) BatchDelUser(ctx context.Context, limit int, cursor string) (res *UserBaseRes[[]UserEntity], err error) {
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	params := map[string]any{"limit": limit, "cursor": cursor}
	res = new(UserBaseRes[[]UserEntity])
	if err = c.doReq(ctx, http.MethodDelete, "users", params, nil, res); err != nil {
		return nil, err
	}
	return
}

// EditUserPassword 修改用户密码
func (c *Client) EditUserPassword(ctx context.Context, username, newPassword string) (res *UserBaseRes[any], err error) {
	data := map[string]string{"newpassword": newPassword}
	pathSuffix := fmt.Sprintf("/users/%s/password", username)
	res = new(UserBaseRes[any])
	if err = c.doReq(ctx, http.MethodPut, pathSuffix, nil, data, res); err != nil {
		return nil, err
	}
	return
}

// GetUser 获取用户详情
func (c *Client) GetUser(ctx context.Context, username string) (res *UserBaseRes[[]UserEntity], err error) {
	pathSuffix := fmt.Sprintf("users/%s", username)
	res = new(UserBaseRes[[]UserEntity])
	if err = c.doReq(ctx, http.MethodGet, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

// BatchGetUser 批量获取用户详情
func (c *Client) BatchGetUser(ctx context.Context, limit int, cursor string) (res *UserBaseRes[[]UserEntity], err error) {
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	params := map[string]any{"limit": limit, "cursor": cursor}
	res = new(UserBaseRes[[]UserEntity])
	if err = c.doReq(ctx, http.MethodGet, "users", params, nil, res); err != nil {
		return nil, err
	}
	return
}

// UserDeactivate 用户账号封禁
func (c *Client) UserDeactivate(ctx context.Context, username string) (res *UserBaseRes[[]UserEntity], err error) {
	pathSuffix := fmt.Sprintf("/users/%s/deactivate", username)
	res = new(UserBaseRes[[]UserEntity])
	if err = c.doReq(ctx, http.MethodPost, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

// UserActivate 用户账号解禁
func (c *Client) UserActivate(ctx context.Context, username string) (res *UserBaseRes[[]UserEntity], err error) {
	pathSuffix := fmt.Sprintf("/users/%s/activate", username)
	res = new(UserBaseRes[[]UserEntity])
	if err = c.doReq(ctx, http.MethodPost, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

type UserOnlineStatus struct {
	Username string `json:"username"`
	Status   string `json:"status"`
}

// GetUserOnlineStatus 获取用户在线状态
func (c *Client) GetUserOnlineStatus(ctx context.Context, username string) (res *BaseRes[UserOnlineStatus], err error) {
	pathSuffix := fmt.Sprintf("users/%s/status", username)
	resTmp := new(BaseRes[map[string]string])
	if err = c.doReq(ctx, http.MethodGet, pathSuffix, nil, nil, resTmp); err != nil {
		return nil, err
	}
	res = &BaseRes[UserOnlineStatus]{Timestamp: resTmp.Timestamp, Duration: resTmp.Duration}
	for key, val := range resTmp.Data {
		res.Data.Username, res.Data.Status = key, val
	}
	return
}

// BatchGetUserOnlineStatus 批量获取用户在线状态
func (c *Client) BatchGetUserOnlineStatus(ctx context.Context, usernames []string) (res *BaseRes[[]UserOnlineStatus], err error) {
	if len(usernames) == 0 {
		return nil, errors.New("minimum count of username is 1")
	} else if len(usernames) > 100 {
		return nil, errors.New("maximum count of username is 100")
	}
	data := map[string]any{"usernames": usernames}
	resTmp := new(BaseRes[[]map[string]string])
	if err = c.doReq(ctx, http.MethodPost, "users/batch/status", nil, data, resTmp); err != nil {
		return nil, err
	}
	res = &BaseRes[[]UserOnlineStatus]{Timestamp: resTmp.Timestamp, Duration: resTmp.Duration}
	for _, item := range resTmp.Data {
		for key, val := range item {
			res.Data = append(res.Data, UserOnlineStatus{Username: key, Status: val})
		}
	}
	return
}

type UserOnlineDevice struct {
	Res        string `json:"res"`
	DeviceUUID string `json:"device_uuid"`
	DeviceName string `json:"device_name"`
}

// GetUserOnlineDevices 获取用户在线设备
func (c *Client) GetUserOnlineDevices(ctx context.Context, username string) (res *BaseRes[[]UserOnlineDevice], err error) {
	pathSuffix := fmt.Sprintf("users/%s/resources", username)
	res = new(BaseRes[[]UserOnlineDevice])
	if err = c.doReq(ctx, http.MethodGet, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

type UserDisconnectResData struct {
	Result bool `json:"result"`
}

// UserDisconnect 强制用户下线
func (c *Client) UserDisconnect(ctx context.Context, username string) (res *BaseRes[UserDisconnectResData], err error) {
	pathSuffix := fmt.Sprintf("/users/%s/disconnect", username)
	res = new(BaseRes[UserDisconnectResData])
	if err = c.doReq(ctx, http.MethodGet, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

// UserDeviceDisconnect 强制用户从单设备下线
func (c *Client) UserDeviceDisconnect(ctx context.Context, username, resourceId string) (res *BaseRes[UserDisconnectResData], err error) {
	pathSuffix := fmt.Sprintf("/users/%s/disconnect/%s", username, resourceId)
	res = new(BaseRes[UserDisconnectResData])
	if err = c.doReq(ctx, http.MethodDelete, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

// SetUserMetadata 设置用户属性
func (c *Client) SetUserMetadata(ctx context.Context, username string, metadata map[string]string) (res *BaseRes[map[string]string], err error) {
	pathSuffix := fmt.Sprintf("metadata/user/%s", username)
	res = new(BaseRes[map[string]string])

	r := c.reqClient.R().SetContext(ctx).SetFormData(metadata)
	resp, err := r.Put(pathSuffix)
	if err != nil {
		return nil, err
	}
	if err = c.parseResponse(resp, res); err != nil {
		return nil, err
	}
	return
}

// DelUserMetadata 删除用户属性
func (c *Client) DelUserMetadata(ctx context.Context, username string) (res *BaseRes[bool], err error) {
	pathSuffix := fmt.Sprintf("metadata/user/%s", username)
	res = new(BaseRes[bool])
	if err = c.doReq(ctx, http.MethodDelete, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

// GetUserMetadata 获取用户属性
func (c *Client) GetUserMetadata(ctx context.Context, username string) (res *BaseRes[map[string]string], err error) {
	pathSuffix := fmt.Sprintf("metadata/user/%s", username)
	res = new(BaseRes[map[string]string])
	if err = c.doReq(ctx, http.MethodGet, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

type UserMetadata struct {
	Username string            `json:"username"`
	Metadata map[string]string `json:"metadata"`
}

// BatchGetUserMetadata 批量获取用户属性
func (c *Client) BatchGetUserMetadata(ctx context.Context, targets, properties []string) (res *BaseRes[[]UserMetadata], err error) {
	data := map[string]any{"targets": targets, "properties": properties}
	resTmp := new(BaseRes[map[string]map[string]string])
	if err = c.doReq(ctx, http.MethodPost, "/metadata/user/get", nil, data, resTmp); err != nil {
		return nil, err
	}
	res = &BaseRes[[]UserMetadata]{Timestamp: resTmp.Timestamp, Duration: resTmp.Duration}
	for username, metadata := range resTmp.Data {
		res.Data = append(res.Data, UserMetadata{Username: username, Metadata: metadata})
	}
	return
}

// GetAppUserMetadataCapacity 获取 App 下用户属性总大小,单位为字节
func (c *Client) GetAppUserMetadataCapacity(ctx context.Context) (res *BaseRes[int64], err error) {
	res = new(BaseRes[int64])
	if err = c.doReq(ctx, http.MethodGet, "metadata/user/capacity", nil, nil, res); err != nil {
		return nil, err
	}
	return
}
