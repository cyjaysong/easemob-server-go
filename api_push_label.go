package easemob_server_go

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type PushLabelData struct {
	Name        string `json:"name"`        // 推送标签的名称
	Description string `json:"description"` // 推送标签的描述
	CreatedAt   int64  `json:"createdAt"`   // 推送标签的创建时间。该时间为 Unix 时间戳，单位为毫秒
	Count       int64  `json:"count"`       // 该推送标签下的用户数量
}

// CreatePushLabel  创建推送标签
func (c *Client) CreatePushLabel(ctx context.Context, labelName, description string) (res *BaseRes[PushLabelData], err error) {
	data := map[string]any{"name": labelName}
	if len(description) > 0 {
		data["description"] = description
	}
	res = new(BaseRes[PushLabelData])
	if err = c.doReq(ctx, http.MethodPost, "push/label", nil, data, res); err != nil {
		return nil, err
	}
	return
}

// DeletePushLabel  删除推送标签
func (c *Client) DeletePushLabel(ctx context.Context, labelName string) (res *BaseRes[string], err error) {
	pathSuffix := fmt.Sprintf("push/label/%s", labelName)
	res = new(BaseRes[string])
	if err = c.doReq(ctx, http.MethodDelete, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

// GetPushLabel  查询推送标签
func (c *Client) GetPushLabel(ctx context.Context, labelName string) (res *BaseRes[PushLabelData], err error) {
	pathSuffix := fmt.Sprintf("push/label/%s", labelName)
	res = new(BaseRes[PushLabelData])
	if err = c.doReq(ctx, http.MethodGet, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

// GetPushLabelList  分页查询推送标签
func (c *Client) GetPushLabelList(ctx context.Context, limit int, cursor string) (res *PageRes[[]PushLabelData], err error) {
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	params := map[string]any{"limit": limit, "cursor": cursor}
	res = new(PageRes[[]PushLabelData])
	if err = c.doReq(ctx, http.MethodGet, "push/label", params, nil, res); err != nil {
		return nil, err
	}
	return
}

// EditPushLabelUserResult 编辑推送标签用户结果
type EditPushLabelUserResult struct {
	Success []string          `json:"success"` // 添加成功的用户列表
	Fail    map[string]string `json:"fail"`    // 添加失败的结果,key为添加失败的用户名,value为失败原因
}

// AddPushLabelUser 在推送标签下添加用户
func (c *Client) AddPushLabelUser(ctx context.Context, labelName string, usernames []string) (res *BaseRes[EditPushLabelUserResult], err error) {
	if len(usernames) == 0 {
		return nil, errors.New("usernames is empty")
	} else if len(usernames) > 100 {
		return nil, errors.New("too many username, maximum count is 100")
	}
	data := map[string]any{"usernames": usernames}
	pathSuffix := fmt.Sprintf("push/label/%s/user", labelName)
	res = new(BaseRes[EditPushLabelUserResult])
	if err = c.doReq(ctx, http.MethodPost, pathSuffix, nil, data, res); err != nil {
		return nil, err
	}
	return
}

// DelPushLabelUser 批量移出指定推送标签下的用户
func (c *Client) DelPushLabelUser(ctx context.Context, labelName string, usernames []string) (res *BaseRes[EditPushLabelUserResult], err error) {
	if len(usernames) == 0 {
		return nil, errors.New("usernames is empty")
	} else if len(usernames) > 100 {
		return nil, errors.New("too many username, maximum count is 100")
	}
	data := map[string]any{"usernames": usernames}
	pathSuffix := fmt.Sprintf("push/label/%s/user", labelName)
	res = new(BaseRes[EditPushLabelUserResult])
	if err = c.doReq(ctx, http.MethodDelete, pathSuffix, nil, data, res); err != nil {
		return nil, err
	}
	return
}

// PushLabelUserData 推送标签用户数据
type PushLabelUserData struct {
	Username string `json:"username"` // 查询的用户 ID
	Created  int64  `json:"created"`  // 添加用户的 Unix 时间戳，单位为毫秒。
}

// GetPushLabelUser 查询指定标签下的指定用户
func (c *Client) GetPushLabelUser(ctx context.Context, labelName string, username string) (res *BaseRes[PushLabelUserData], err error) {
	pathSuffix := fmt.Sprintf("push/label/%s/user/%s", labelName, username)
	res = new(BaseRes[PushLabelUserData])
	if err = c.doReq(ctx, http.MethodGet, pathSuffix, nil, nil, res); err != nil {
		return nil, err
	}
	return
}

// GetPushLabelUserList 分页查询指定标签下的用户
func (c *Client) GetPushLabelUserList(ctx context.Context, labelName string, limit int, cursor string) (res *PageRes[[]PushLabelUserData], err error) {
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	params := map[string]any{"limit": limit, "cursor": cursor}
	pathSuffix := fmt.Sprintf("push/label/%s/user", labelName)
	res = new(PageRes[[]PushLabelUserData])
	if err = c.doReq(ctx, http.MethodGet, pathSuffix, params, nil, res); err != nil {
		return nil, err
	}
	return
}
