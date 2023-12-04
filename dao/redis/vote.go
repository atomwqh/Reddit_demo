package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票多少分
)

/*
投票的几种情况：
dir = 1时，有两种情况：
	1.之前没投票，现在赞同 差值的绝对值：1 +432
	2.之前反对，现在赞同 差值的绝对值：2 +432 * 2
dir = 0，两种情况：
	1.之前反对，现在取消 差值的绝对值：1 +432
	2.之前赞同，现在取消 差值的绝对值：1 -432
dir = -1，两种情况：
	1.之前没投，现在反对 差值的绝对值：1 -432
	2.之前赞同，现在反对 差值的绝对值：2 -432 * 2
 投票的限制：帖子一周之后就不允许投票了（持久化操作），便于后端数据的处理
	1. 到期之后将redis中的数据保存到mysql里面
	2. 到期之后删除KeyPostVotedZsetPrefix
*/

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

// 创建帖子的时候在redis里添加时间戳

func CreatePost(postID, communityID int64) error {

	pipeline := client.TxPipeline()
	// 帖子时间，默认值
	pipeline.ZAdd(getRedisKey(KeyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 帖子分数，随着投票改变
	pipeline.ZAdd(getRedisKey(KeyPostScoreZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 把帖子ID加到社区的set
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}
func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票限制
	// 取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZset), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 其中2和3也需要放在一个事务中去
	// 2. 更新分数(较为复杂)
	// 先查询当前用户之前的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZsetPrefix+postID), userID).Val()
	// 不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value)
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZset), op*diff*scorePerVote, postID)
	//if ErrVoteTimeExpire != nil {
	//	return err
	//}
	// 3. 记录用户为该帖子投了什么票
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZsetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZsetPrefix+postID), redis.Z{
			Score:  value, // 赞成或者反对
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
