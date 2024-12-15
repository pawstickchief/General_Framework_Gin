package controllers

import (
	"General_Framework_Gin/database/casbin"
	"General_Framework_Gin/schemas/request"
	"github.com/gin-gonic/gin"
	"strings"
)

func HandlePolicy(c *gin.Context) {
	var policyOption request.PolicyOption
	if err := c.ShouldBindJSON(&policyOption); err != nil {
		ResponseError(c, CodeInvalidParam) // 处理无效参数错误
		return
	}

	// 根据操作类型执行相应的操作
	switch strings.ToLower(policyOption.Option) {
	case "add":
		if err := casbin.AddPolicy(casbin.Enforcer, policyOption); err != nil {
			ResponseSuccess(c, err)
		} else {
			ResponseSuccess(c, "策略添加成功")
		}
	case "remove":
		if err := casbin.RemovePolicy(casbin.Enforcer, policyOption); err != nil {
			ResponseSuccess(c, err)
		} else {
			ResponseSuccess(c, "策略删除成功")
		}
	case "select":
		if err := casbin.GetRolePolicy(casbin.Enforcer, policyOption); err != nil {
			ResponseSuccess(c, err)
		} else {
			ResponseSuccess(c, "策略获取成功")
		}
	case "modify":
		if err := casbin.EditPolicy(casbin.Enforcer, policyOption); err != nil {
			ResponseSuccess(c, err)
		} else {
			ResponseSuccess(c, "策略编辑成功成功")
		}

	default:
		ResponseError(c, CodeInvalidParam) // 操作无效
	}
}
