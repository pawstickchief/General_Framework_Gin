package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResCode 定义错误码类型
type ResCode int64

// 定义错误码常量
const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidAuth
	CodeServerApiType
	CodeHostlist
	CodeAlarminfo
	CodeSelectSwitch
	CodeNoClientsConnected
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:            "success",
	CodeInvalidParam:       "请求参数错误",
	CodeUserExist:          "用户已存在",
	CodeUserNotExist:       "用户不存在",
	CodeInvalidPassword:    "用户名或密码输入错误",
	CodeServerBusy:         "服务器忙",
	CodeNeedLogin:          "需要登录",
	CodeInvalidAuth:        "无效的token",
	CodeServerApiType:      "接口参数错误",
	CodeHostlist:           "主机已存在",
	CodeAlarminfo:          "报警接口参数错误",
	CodeSelectSwitch:       "交换机上联信息错误",
	CodeNoClientsConnected: "获取WebSocket失败",
}

// Msg 获取错误码对应的消息
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}

// ResponseData 自定义响应结构体
type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// ResponseError 错误响应
func ResponseError(c *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// ResponseSuccess 成功响应
func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}

// ResponseSystemDataSuccess 系统数据成功响应
func ResponseSystemDataSuccess(c *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}

// ResponseErrorWithMsg 自定义错误消息响应
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}
