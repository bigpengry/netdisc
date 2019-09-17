package util

import (
	"log"
	"encoding/json"

)
// RespMsg : http响应的通用数据结构
type RespMsg struct{
	Code int  `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

// NewRespMsg : http消息构造函数
func NewRespMsg(code int,msg string,data interface{})*RespMsg{
	return &RespMsg{
		Code:code,
		Msg:msg,
		Data:data,
	}
}

// JSONBytes : 对象转json格式的byte
func (resp *RespMsg) JSONBytes() []byte {
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return r
}

// JSONString : 对象转json格式的string
func (resp *RespMsg)JSONString()string{
	r,err:=json.Marshal(resp)
	if err!=nil{
		log.Println(err)
	}
	return string(r)
}