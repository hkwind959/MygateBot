package constant

var Headers = map[string]string{
	"Accept":             "application/json",
	"Accept-Encoding":    "gzip, deflate, br, zstd",
	"Accept-Language":    "en-US,en;q=0.9,id;q=0.8",
	"Origin":             "chrome-extension://hajiimgolngmlbglaoheacnejbnnmoco",
	"Priority":           "u=1, i",
	"Referer":            "https://app.mygate.network/",
	"Sec-CH-UA":          "Not A(Brand\";v=\"8\", \"Chromium\";v=\"132\", \"Google Chrome\";v=\"132",
	"Sec-CH-UA-Mobile":   "?0",
	"Sec-CH-UA-Platform": "Windows",
	"Sec-Fetch-Dest":     "empty",
	"Sec-Fetch-Mode":     "cors",
	"Sec-Fetch-Site":     "none",
	"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
}

// SecretKey 生成签名 key
var SecretKey = "|`8S%QN9v&/J^Za"

var baseUrl = "https://api.mygate.network/api"

// GetNodeUrl 获取节点信息
var GetNodeUrl = baseUrl + "/front/nodes?limit=100&page=1"

// RegisterNodeUrl 注册节点
var RegisterNodeUrl = baseUrl + "/front/nodes"

// GetUserInfoUrl 获取用户信息
var GetUserInfoUrl = baseUrl + "/front/users/me"

// GetWssUrl 获取wss地址
// wss://api.mygate.network/socket.io/?nodeId=6993ce82-1bd6-42ac-a254-6291a534cf1f&signature=0f40f206e70f608aec5813d1aae8d8046986fb8f56ff8fc493ce96e994895b44&timestamp=1741189478754&version=2&EIO=4&transport=websocket
var GetWssUrl = "wss://api.mygate.network/socket.io/?nodeId="

// CheckProxyURL 校验代理 URL
var CheckProxyURL = "https://api.ipify.org/?format=json"
