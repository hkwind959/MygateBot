package config

import (
	"MygateBot/logs"
	"MygateBot/model"
	"encoding/json"
	"io/ioutil"
	"os"
)

func GetTokens() model.Tokens {
	fileContent, err := os.Open("./tokens.json")
	if err != nil {
		logs.I().Fatal("读取配置文件失败")
		return nil
	}
	defer fileContent.Close()
	var tokens model.Tokens
	byteResult, _ := ioutil.ReadAll(fileContent)
	err = json.Unmarshal([]byte(byteResult), &tokens)
	return tokens
}
