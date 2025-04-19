package easemob_server_go

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// PushMsgMap 推送消息内容
type PushMsgMap map[string]any

func (bm PushMsgMap) Set(key string, value any) PushMsgMap {
	bm[key] = value
	return bm
}
func (bm PushMsgMap) SetBodyMap(key string, value func(b PushMsgMap)) PushMsgMap {
	_bm := make(PushMsgMap)
	value(_bm)
	bm[key] = _bm
	return bm
}

// PushStrategy 推送策略
type PushStrategy int

const (
	PushStrategyThirdPartyFirst PushStrategy = 0 // 第三方厂商通道优先，失败时走环信通道
	PushStrategyEaseMobOnly     PushStrategy = 1 // 只走环信通道
	PushStrategyThirdPartyOnly  PushStrategy = 2 // 只走第三方厂商通道
	PushStrategyOnlineEaseMob   PushStrategy = 3 // 在线走环信通道，离线走第三方厂商通道
	PushStrategyOnlineOnly      PushStrategy = 4 // 只走环信通道且只推在线用户
)

// SyncPushResultData 同步推送结果数据
type SyncPushResultData struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 消息
	Data    struct {
		ExpireTokens []string `json:"expireTokens"` // 过期的token
		SendResult   bool     `json:"sendResult"`   // 发送结果
		RequestID    string   `json:"requestId"`    // 请求ID
		FailTokens   []string `json:"failTokens"`   // 失败的token
		MsgCode      int      `json:"msgCode"`      // 消息码
	} `json:"data"` // 数据
}

// SyncPushResultItem 同步推送结果项
type SyncPushResultItem struct {
	PushStatus string              `json:"pushStatus"`     // 推送状态
	Desc       string              `json:"desc,omitempty"` // 描述（失败时）
	Data       *SyncPushResultData `json:"data,omitempty"` // 数据（成功时）
}

// SyncPushNotification 同步方式发送推送通知
func (c *Client) SyncPushNotification(ctx context.Context, target string, pushMessage PushMsgMap, strategy PushStrategy) (res *BaseRes[[]SyncPushResultItem], err error) {
	if target == "" {
		return nil, errors.New("target is empty")
	}

	data := map[string]any{"pushMessage": pushMessage, "strategy": strategy}

	pathSuffix := fmt.Sprintf("push/sync/%s", target)
	res = new(BaseRes[[]SyncPushResultItem])
	if err = c.doReq(ctx, http.MethodPost, pathSuffix, nil, data, res); err != nil {
		return nil, err
	}
	return
}

// AsyncPushResultItem 异步推送结果项
type AsyncPushResultItem struct {
	Id         string `json:"id,omitempty"`   // 推送的目标用户 ID
	PushStatus string `json:"pushStatus"`     // 推送状态
	Desc       string `json:"desc,omitempty"` // 描述（失败时）
	Data       string `json:"data,omitempty"` // 数据（成功时）
}

// AsyncPushNotification 异步方式向单个用户发送推送通知
func (c *Client) AsyncPushNotification(ctx context.Context, target string, pushMessage PushMsgMap, strategy PushStrategy) (res *BaseRes[[]AsyncPushResultItem], err error) {
	if target == "" {
		return nil, errors.New("target is empty")
	}

	data := map[string]any{"pushMessage": pushMessage, "strategy": strategy}

	pathSuffix := fmt.Sprintf("push/async/%s", target)
	res = new(BaseRes[[]AsyncPushResultItem])
	if err = c.doReq(ctx, http.MethodPost, pathSuffix, nil, data, res); err != nil {
		return nil, err
	}
	return
}

// BatchAsyncPushNotification 异步方式批量发送推送通知
func (c *Client) BatchAsyncPushNotification(ctx context.Context, targets []string, pushMessage PushMsgMap, strategy PushStrategy) (res *BaseRes[[]AsyncPushResultItem], err error) {
	if len(targets) == 0 {
		return nil, errors.New("targets is empty")
	} else if len(targets) > 100 {
		return nil, errors.New("targets is too much")
	}

	data := map[string]any{"targets": targets, "pushMessage": pushMessage, "strategy": strategy}

	res = new(BaseRes[[]AsyncPushResultItem])
	if err = c.doReq(ctx, http.MethodPost, "push/single", nil, data, res); err != nil {
		return nil, err
	}
	return
}

type LabelPushResData struct {
	TaskId int64 `json:"taskId,omitempty"` // 推送任务 ID
}

// LabelPushNotification 使用标签推送接口发送推送通知
func (c *Client) LabelPushNotification(ctx context.Context, labels []string, pushMessage PushMsgMap, strategy PushStrategy, startAt *time.Time) (res *BaseRes[LabelPushResData], err error) {
	if len(labels) == 0 {
		return nil, errors.New("labels is empty")
	} else if len(labels) > 5 {
		return nil, errors.New("labels is too much")
	}

	data := map[string]any{"targets": labels, "pushMessage": pushMessage, "strategy": strategy}
	if startAt != nil {
		data["startDate"] = startAt.Format("2006-01-02 15:04:05")
	}

	res = new(BaseRes[LabelPushResData])
	if err = c.doReq(ctx, http.MethodPost, "push/list/label", nil, data, res); err != nil {
		return nil, err
	}
	return
}

// CreateFullPushTask 创建全量推送任务
func (c *Client) CreateFullPushTask(ctx context.Context, pushMessage PushMsgMap, strategy PushStrategy, startAt *time.Time) (res *BaseRes[int64], err error) {
	data := map[string]any{"pushMessage": pushMessage, "strategy": strategy}

	if startAt != nil {
		data["startDate"] = startAt.Format("2006-01-02 15:04:05")
	}

	res = new(BaseRes[int64])
	if err = c.doReq(ctx, http.MethodPost, "push/task", nil, data, res); err != nil {
		return nil, err
	}
	return
}
