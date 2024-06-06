package redis

// redis key

const (
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "post:time"   //zset:贴子级发帖时间
	KeyPostScoreZSet   = "post:score"  //zset:贴子级投票的分数
	KeyPostVotedZSetPF = "post:voted:" //zset:记录用户及投票类型；参数是post id
	KeyCommunitySetPF  = "community:"  //set:保存每个分区下帖子的id
)

func getRedisKey(key string) string {
	return Prefix + key
}
