package easemob_server_go

import (
	"context"

	"github.com/imroc/req/v3"
)

type ApiError struct {
	ErrorInfo        string `json:"error"`
	Exception        string `json:"exception"`
	Timestamp        int64  `json:"timestamp"`
	Duration         int    `json:"duration"`
	ErrorDescription string `json:"error_description"`
}

func (e ApiError) Error() string {
	return e.ErrorDescription
}

type UserBaseRes[T any] struct {
	Action          string `json:"action"`
	Application     string `json:"application"`
	Path            string `json:"path"`
	Uri             string `json:"uri"`
	Count           int    `json:"count"`
	Entities        T      `json:"entities"`
	Timestamp       int64  `json:"timestamp"`
	Duration        int    `json:"duration"`
	Organization    string `json:"organization"`
	ApplicationName string `json:"applicationName"`
}

type BaseRes[T any] struct {
	Timestamp int64 `json:"timestamp,omitempty"`
	Duration  int   `json:"duration,omitempty"`
	Data      T     `json:"data"`
}

// PageRes 通用分页响应结构体
type PageRes[T any] struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	Cursor    string `json:"cursor,omitempty"` // 下次查询的起始位置
	Data      T      `json:"data"`
}

type Response struct {
	Path            string      `json:"path,omitempty"`
	URI             string      `json:"uri,omitempty"`
	Timestamp       int64       `json:"timestamp,omitempty"`
	Count           int         `json:"count,omitempty"`
	Action          string      `json:"action,omitempty"`
	Duration        int         `json:"duration,omitempty"`
	Data            bool        `json:"data,omitempty"`
	ApplicationName string      `json:"applicationName,omitempty"`
	Organization    string      `json:"organization,omitempty"`
	Application     string      `json:"application,omitempty"`
	Cursor          string      `json:"cursor,omitempty"`
	Entities        interface{} `json:"entities,omitempty"`
	Properties      interface{} `json:"properties,omitempty"`
}

type ResultResponse struct {
	Data interface{} `json:"data"`
	Response
}

func (c *Client) doReq(ctx context.Context, method, pathSuffix string, params map[string]any, data any, res any) (err error) {
	r := c.reqClient.R().SetContext(ctx)
	if params != nil {
		r.SetQueryParamsAnyType(params)
	}
	if data != nil {
		r.SetBodyJsonMarshal(data)
	}
	resp, err := r.Send(method, pathSuffix)
	if err != nil {
		return err
	}
	return c.parseResponse(resp, res)
}

func (c *Client) parseResponse(resp *req.Response, res any) (err error) {
	if resp.StatusCode != 200 {
		var apiErr ApiError
		if err = resp.UnmarshalJson(&apiErr); err == nil {
			return apiErr
		}
		return
	}
	return resp.UnmarshalJson(res)
}
