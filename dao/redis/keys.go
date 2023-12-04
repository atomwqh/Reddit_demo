package redis

// redis key

// redis key自定义命名空间，方便后期查询和拆分

const (
	Prefix                 = "reddit:"
	KeyPostTimeZset        = "post:time"   // zset;帖子及发帖实践
	KeyPostScoreZset       = "post:score"  // zset;帖子及投票的分数
	KeyPostVotedZsetPrefix = "post:voted:" // zset;记录用户及投票类型:参数是post_id

	KeyCommunitySetPrefix = "community:" // set;保存每个分区下的帖子id
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
