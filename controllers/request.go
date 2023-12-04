package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const ContextUserIDkey = "userID"

var ErrorUserNotLogin = errors.New("用户未登入")

// GetCurrentUser通过这个函数我们就能快速的拿到用户的id

func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContextUserIDkey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

// 获取分页参数
func getPageInfo(c *gin.Context) (int64, int64) {
	// 获取分页参数

	pageNumber := c.Query("page")
	pageSizeStr := c.Query("size")
	var (
		page int64
		size int64
		err  error
	)
	page, err = strconv.ParseInt(pageNumber, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
