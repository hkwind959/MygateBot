package utils

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type WebSocketProxyClient struct {
	proxyURL *url.URL
	dialer   *websocket.Dialer
	conn     *websocket.Conn
}

// NewWebSocketProxyClient 创建一个新的WebSocketProxyClient实例
func NewWebSocketProxyClient(proxyIP string) (*WebSocketProxyClient, error) {
	parseProxyURL, err := url.Parse(proxyIP)
	if err != nil {
		return nil, fmt.Errorf("无法解析代理URL：%v", err)
	}
	dialer := websocket.DefaultDialer
	if strings.HasPrefix(proxyIP, "socks5") {
		password, _ := parseProxyURL.User.Password()
		auth := &proxy.Auth{
			User:     parseProxyURL.User.Username(),
			Password: password,
		}
		socksDialer, err := proxy.SOCKS5("tcp", strings.TrimPrefix(parseProxyURL.Host, "socks5"), auth, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("创建 SOCKS5 连接错误 ：%v", err)
		}
		dialer = &websocket.Dialer{
			NetDial: func(network, addr string) (net.Conn, error) {
				return socksDialer.Dial(network, addr)
			},
		}
	} else {
		dialer.Proxy = http.ProxyURL(parseProxyURL)
	}
	return &WebSocketProxyClient{
		proxyURL: parseProxyURL,
		dialer:   dialer,
	}, nil
}

// GetConn 返回WebSocket连接
func (wspc *WebSocketProxyClient) GetConn() *websocket.Conn {
	return wspc.conn
}

// Connect 连接到指定的WebSocket服务器
func (wspc *WebSocketProxyClient) Connect(wsURL string, requestHeader http.Header) error {
	var err error
	wspc.conn, _, err = wspc.dialer.Dial(wsURL, requestHeader)
	if err != nil {
		return fmt.Errorf("无法连接到WebSocket服务器：%v", err)
	}
	return nil
}

func (wspc *WebSocketProxyClient) Close() {
	wspc.conn.Close()
}
