package logic

import (
	"main/dao/mysql"
	"main/dao/redis"
	"main/models"
	"main/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = int64(snowflake.GenID())
	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	// 创建帖子的时候需要添加一个时间——方便之后投票功能实现
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	// 3. 返回
	return
}

func GetPostDetail(id int64) (data *models.ApiPostDetail, err error) {

	// 查询并组合我们接口想用的数据
	post, err := mysql.GetPostDetailByID(id)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailByID(id) failed", zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	// 根据社区id查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
			zap.Int64("author_id", post.CommunityID),
			zap.Error(err))
		return
	}
	// 拼接数据 + 这里注意要指针初始化
	data = &models.ApiPostDetail{
		Authorname:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}

	return
}

// GetPostList 在mysql数据库里面查询获取post列表
func GetPostList(offset, limit int64) (data []*models.ApiPostDetail, err error) {

	posts, err := mysql.GetPostList(offset, limit)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	// 根据帖子来查询作者，社区信息
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("author_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		// 之后需要像上面一样进行拼装
		postDetail := &models.ApiPostDetail{
			Authorname:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2 按一定顺序查询帖子详细信息 GetPostList的升级版
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//2.从redis中查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	//3.根据id去mysql数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	// 这里其实可以提前查好每篇帖子投票数据，不用for遍历
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 根据帖子来查询作者，社区信息，参考之前数据
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("author_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		// 之后需要像上面一样进行拼装
		postDetail := &models.ApiPostDetail{
			Authorname:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetCommunityListPostList 按照社区查询帖子
func GetCommunityListPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	//2.从redis中查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	//3.根据id去mysql数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	// 这里其实可以提前查好每篇帖子投票数据，不用for遍历
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 根据帖子来查询作者，社区信息，参考之前数据
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("author_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		// 之后需要像上面一样进行拼装
		postDetail := &models.ApiPostDetail{
			Authorname:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListNew 将上面两个查询逻辑函数合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 请求参数的不同，执行不同的逻辑
	if p.CommunityID == 0 {
		// 查询所有
		data, err = GetPostList2(p)
	} else {
		// 查询社区id查询
		data, err = GetCommunityListPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}

	return
}
