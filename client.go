package easemob_server_go

import (
	"fmt"
	"github.com/imroc/req/v3"
	"path"
	"time"
)

type Client struct {
	host    string
	orgName string
	appName string
	appKey  string

	clientId     string
	clientSecret string
	appToken     string

	reqClient *req.Client
}

func New(host, orgName, appName, clientId, clientSecret string, devMode bool) (client *Client) {
	baseUrl := fmt.Sprintf("https://%s/%s", host, path.Join(orgName, appName))
	reqClient := req.C().SetBaseURL(baseUrl).SetUserAgent("easemob-server-go")
	reqClient.SetCommonRetryCount(2).SetMaxConnsPerHost(5)
	reqClient.SetTimeout(6 * time.Second)
	reqClient.SetIdleConnTimeout(60 * time.Minute)
	reqClient.SetExpectContinueTimeout(2 * time.Second)
	if devMode {
		reqClient.DevMode()
	}
	client = &Client{reqClient: reqClient, appKey: fmt.Sprintf("%s#%s", orgName, appName),
		host: host, orgName: orgName, appName: appName, clientId: clientId, clientSecret: clientSecret}

	// TTL:设为0表示AppToken永久有效
	//data := &ClientParam{GrantType: "client_credentials", ClientId: clientId, ClientSecret: clientSecret, TTL: ttl}
	//if client.appToken, err = client.GetAppToken(context.Background(), data); err != nil {
	//	return nil, err
	//}
	return
}
