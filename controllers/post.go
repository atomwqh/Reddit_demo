package controllers

import (
	"main/logic"
	"main/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建新帖子，存入数据库并在redis中记录该帖子的分数和所处社区
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT_AToken"
// @Param obj body models.Post false "参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseCreatePost
// @Router /api/v1/post [post]
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数以及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		// 这里出错的可能性很大，这里错误处理会详细一些
		zap.L().Debug("c.ShouldBindJSON() err", zap.Any("err", err))
		zap.L().Error("create post with invaild param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从 c（content）中读取到当前请求的用户ID,哪位用户要创建帖子
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
	}
	p.AuthorID = userID
	// 2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// PostDetailHandler 获取帖子详情
// @Summary 通过post id获取post详情
// @Description 通过post id获取post内容以及所所在社区和作者名
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int64 true "帖子id"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostDetail
// @Router /api/v1/post/{id} [get]
func PostDetailHandler(c *gin.Context) {
	// 1.获取帖子ID参数以及参数校验
	idStr := c.Param("id")                     // 获取URL参数
	id, err := strconv.ParseInt(idStr, 10, 64) // 参数校验
	if err != nil {
		zap.L().Error("invalid param when get post detail", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2.根据ID获取帖子详情
	data, err := logic.GetPostDetail(id)
	if err != nil {
		zap.L().Error("logic.GetPostDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 这里不轻易把服务端的错误暴露给外部
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 或缺帖子列表接口
// @Summary 概况
// @Description 描述
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param page path string false "页码"
// @Param size path string false "页面大小"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /api/v1/posts [post]
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数

	page, size := getPageInfo(c)
	// 查询数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// GetPostListHandler2 根据时间或者分数获取帖子列表， 升级版
// @Summary 获取帖子分页数据
// @Description 根据社区id（可以为空）、页码、数量返回分页数据
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /api/v2/posts [get]
func GetPostListHandler2(c *gin.Context) {
	/*
		升级版查询帖子，帖子列表接口，根据前端传过来的选择动态获取帖子列表
		这里选择包括 按照创建时间排序 或者 分数排序
			1.获取参数
			2.从redis中查询id列表
			3.根据id去数据库查询帖子详细信息
	*/
	// GET请求参数：/api/v1/post2?page=1&size=10&order=time query string参数
	// 初始化结构体指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	// 1. 获取参数
	//c.ShouldBindJSON() 如果请求中携带的是json格式数据，才用这个方法获取参数
	//c.ShouldBind() 自动获取参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 init with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	} // query string数据
	data, err := logic.GetPostListNew(p) // 更新后的接口（将之前的合二为一）
	// 2. 查询数据
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// GetCommunityPostListHandler 根据社区去查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	p := &models.ParamCommunityPostList{
//		ParamPostList: &models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime,
//		},
//	}
//	// 1. 获取参数
//	//c.ShouldBindJSON() 如果请求中携带的是json格式数据，才用这个方法获取参数
//	//c.ShouldBind() 自动获取参数
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler init with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	} // query string数据
//
//	// 2. 查询数据
//	data, err := logic.GetCommunityListPostList(p)
//	if err != nil {
//		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	ResponseSuccess(c, data)
//}
