package model

// Tokens 结构体
type Token struct {
	Token   string   `json:"token"`
	Proxies []string `json:"proxies"`
	Remark  string   `json:"remark"`
}

type Tokens []Token
