package mysql

import (
	"database/sql"
	"main/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id, title, content, author_id, community_id) values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostDetailByID 根据帖子ID查询帖子详情
func GetPostDetailByID(id int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select 
    post_id, title, content, author_id, community_id, create_time
	from post 
	where post_id = ?
	`
	if err := db.Get(post, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvaildID
		}
	}
	return post, err
}

// GetPostList 查询帖子列表
func GetPostList(page, size int64) (postlist []*models.Post, err error) {
	sqlStr := `select 
    post_id, title, content, author_id, community_id, create_time
	from post 
	Order By create_time
	DESC 
	limit ?,?
	`
	postlist = make([]*models.Post, 0, 2)
	err = db.Select(&postlist, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据给定的ID列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?) 
	order by FIND_IN_SET(post_id, ?)
	`

	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	err = db.Select(&postList, query, args...) // 这里注意会出现不匹配的空接口问题，注意在最后加上。。。
	return
}
