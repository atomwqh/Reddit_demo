package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 简单写一个创建帖子的单元测试
func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
		"community_id":1,
		"title":"test",
		"content":"this is a test"
	}`
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// 方法一：判断响应的内容是不是按预期返回了需要登录的错误
	//assert.Contains(t, w.Body.String(), "需要登录")
	// 方法二：将相应的内容反序列化到ResponseData, 然后判断字段与预期是否一致
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.body failed, err:%v\n", err)
	}
	assert.Equal(t, res.Code, CodeNeedLogin)
}
