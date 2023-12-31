package logic

import (
	"main/dao/mysql"
	"main/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查数据库 查找所有的Community 并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
