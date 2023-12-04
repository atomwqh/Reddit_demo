package redis

import (
	"main/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	// 确定查询的索引起始点
	start := (page - 1) * size
	end := start + size - 1
	// 3. ZRevRange 按分数从大到小指定元素数量查询
	return client.ZRevRange(key, start, end).Result()
}

// GetPostIDsInOrder 按顺序（时间，分数）获取帖子ID
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis中获取id（按照p中的数据来判断）
	// 根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZset)
	}
	// 确定查询的索引起始点
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZsetPrefix + id)
	//	// 查找key中分数是 1 的元素的数量相当于就是统计每篇帖子的赞同票数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}
	// 使用pipeline减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZsetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmders := range cmders {
		v := cmders.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区根据ids查找每篇帖子赞成票的数据
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZset)
	}
	// 使用 zinterstore 把分区的帖子zset与帖子分数的zet生成一个新的zset
	// 再针对新的zset按之前逻辑取值

	// 社区的key
	communitykey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Weights:   nil,
			Aggregate: "MAX",
		}, communitykey, orderKey) // zinterstore 计算
		pipeline.Expire(key, 60*time.Second) // 计算超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在的话直接根据key查询
	return getIDsFromKey(key, p.Page, p.Size)

}
