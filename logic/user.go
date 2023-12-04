package logic

import (
	"main/dao/mysql"
	"main/models"
	"main/pkg/jwt"
	"main/pkg/snowflake"
)

// 业务逻辑的代码

func Signup(p *models.ParamSignUp) (err error) {
	// 1. 判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		// 数据库查询出错
		return err
	}
	// 2. 生成uid
	userID := snowflake.GenID()
	// 需要一个User实例（结构体）才能传入数据库
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3. 保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 由于传递的是指针，user能拿到UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成jwt
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
