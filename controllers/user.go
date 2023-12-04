package controllers

import (
	"errors"
	"fmt"
	"main/dao/mysql"
	"main/logic"
	"main/models"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求的函数
// @Summary 注册
// @Description 注册
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param obj body models.ParamSignUp true "用户注册参数"
// @Success 200
// @Router /api/v1/signup [post]
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamSignUp)
	// 这里只能校验参数的格式和类型
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Sign up with invalid params", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型，不是的话就不用翻译了
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 手动对请求参数进行详细的参数校验
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//	zap.L().Error("Sign up with invaild params")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "Sign up with invaild params",
	//	})
	//	return
	//}
	//fmt.Println(p)

	// 2. 业务处理
	if err := logic.Signup(p); err != nil {
		zap.L().Error("login.Signup failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 1.获取参数以及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invaild params", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))

		return
	}
	// 2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), // 这里用json有可能出现精度丢失问题
		"user_name": user.Username,
		"token":     user.Token,
	})
}
