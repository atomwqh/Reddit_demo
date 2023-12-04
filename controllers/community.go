package controllers

import (
	"main/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//--- 跟社区相关 ---

// CommunityHandler 查询所有社区
// @Summary 查询所有社区
// @Description 查询到所有的社区 (community_id, community_name) 以列表的形式返回
// @Tags 社区
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Security ApiKeyAuth
// @Success 200
// @Router /api/v1/community [get]

func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区（community_id, community_name）以列表（切片）的方式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 这里不轻易把服务端的错误暴露给外部
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
// @Summary 概况
// @Description 描述
// @Tags 社区
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param id path string true "社区id"
// @Security ApiKeyAuth
// @Success 200
// @Router /api/v1/community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	// 1,获取社区id（gin框架获取参数) + 参数校验
	idStr := c.Param("id")                     // 获取URL参数
	id, err := strconv.ParseInt(idStr, 10, 64) // 参数校验
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2,根据ID获取社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 这里不轻易把服务端的错误暴露给外部
		return
	}
	ResponseSuccess(c, data)
}
