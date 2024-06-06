package redis

import (
	"begin/goweb/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	return client.ZRevRange(key, start, end).Result()
}

func GetPostDsInOrder(p *models.ParamPostList) ([]string, error) {
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇贴子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostDsInOrder 按社区查询ids
func GetCommunityPostDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	//社区的key
	ckey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zintersstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(orderKey).Val() < 1 {
		pipline := client.Pipeline()
		pipline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, ckey, orderKey)
		pipline.Expire(key, 60*time.Second)
		_, err := pipline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFromKey(key, p.Page, p.Size)
}
