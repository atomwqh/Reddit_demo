package logic

import (
	"main/dao/redis"
	"main/models"
	"strconv"

	"go.uber.org/zap"
)

/*
 投票的限制：帖子一周之后就不允许投票了（持久化操作），便于后端数据的处理
	1. 到期之后将redis中的数据保存到mysql里面
	2. 到期之后删除KeyPostVotedZsetPrefix
*/

// voteForPost 为帖子投票

func VoteForPost(userID int64, p *models.ParamVoteData) error {
	// 在这里记录日志还不错
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
