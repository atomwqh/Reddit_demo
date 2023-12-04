package controllers

import (
	"main/logic"
	"main/models"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PostVoteController 给帖子投票
// @Summary 给帖子投票
// @Description 描述
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param obj query models.ParamVoteData true "帖子投票参数"
// @Security ApiKeyAuth
// @Success 200
// @Router /api/v1/vote [post]
func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) // 翻译并去除错误提示中的结构体
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	// 获取当前请求的用户ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 具体投票业务的实现
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.voteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
