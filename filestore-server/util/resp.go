package util

// RespMsg : http响应的通用数据结构
type RespMsg struct{
	Code int  `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
