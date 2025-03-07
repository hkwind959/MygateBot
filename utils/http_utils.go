package utils

import (
	"MygateBot/logs"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout      = 100 * time.Second
	tlsHandshakeTimeout = 50 * time.Second
)

type HttpClientUtil struct {
	client   *resty.Client
	proxyUrl string
}

// NewHttpClient 创建一个新的 http client ，如果不传 proxy ,则创建默认的client
func NewHttpClient(proxy string) *HttpClientUtil {
	if proxy == "" {
		// 创建默认 client
		client := newClient()
		return &HttpClientUtil{
			client:   client,
			proxyUrl: proxy,
		}
	} else {
		// 创建 proxy client
		client := newProxyClient(proxy)
		return &HttpClientUtil{
			client:   client,
			proxyUrl: client.HostURL,
		}
	}
}

// newProxyClient 创建代理客户端
func newProxyClient(proxyIp string) *resty.Client {
	// SOCKS5 代理地址
	proxyURL, err := url.Parse(proxyIp)
	if err != nil {
		return nil
	}
	var transport *http.Transport
	if strings.HasPrefix(proxyIp, "socks") {
		password, _ := proxyURL.User.Password()
		auth := &proxy.Auth{
			User:     proxyURL.User.Username(),
			Password: password,
		}
		dialer, err := proxy.SOCKS5("tcp", strings.TrimPrefix(proxyURL.Host, "socks5://"), auth, proxy.Direct)
		if err != nil {
			return nil
		}
		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
	} else {
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			// 设置超时时间
			TLSHandshakeTimeout: tlsHandshakeTimeout, // TLS 握手超时
		}
	}
	client := resty.New().SetTimeout(defaultTimeout)
	client.SetTransport(transport)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return client
}

// newClient 创建客户端
func newClient() *resty.Client {
	// 创建自定义的 Transport
	transport := &http.Transport{
		// 设置超时时间
		TLSHandshakeTimeout: tlsHandshakeTimeout, // TLS 握手超时
	}
	client := resty.New().SetTimeout(defaultTimeout)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTransport(transport)
	logs.I().Info("默认客户端创建成功")
	return client
}

// Get 发送Get 请求
func (hcu *HttpClientUtil) Get(url string, body interface{}, headers map[string]string, result interface{}) (*resty.Response, error) {
	resp, err := hcu.client.R().SetHeaders(headers).SetBody(body).SetResult(result).Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取信息出错：%v", err)
	}
	return resp, nil
}

// Post 发送 Post 请求
func (hcu *HttpClientUtil) Post(url string, headers map[string]string, body interface{}, result interface{}) (*resty.Response, error) {
	resp, err := hcu.client.R().SetBody(body).SetResult(result).SetHeaders(headers).Post(url)
	if err != nil {
		logs.I().Error("", zap.Error(err))
		return nil, fmt.Errorf("post 获取信息出错：%v", err)
	}
	return resp, nil
}
