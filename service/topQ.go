package service

import (
	"github.com/SSunSShine/QAsystem/database"
	"github.com/SSunSShine/QAsystem/model"
	"log"
	"strconv"
	"sync"
	"time"
)

// 热门问题
var topQ *sync.Map
var once sync.Once

// GetTopQ 单例模式，系统只维护这一个榜单
func GetTopQ() *sync.Map {
	once.Do(func() {
		topQ = &sync.Map{}
	})
	return topQ
}

// UpdateTopQ 更新热榜
func UpdateTopQ(d time.Duration, n int64)  {
	for  {
		result, err := database.RDB.ZRevRangeWithScores(ctx, ZSetKey, 0, n).Result()
		if err != nil {
			log.Print("[ERROR]更新热门问题出现错误: "+err.Error()+": zset")
		}
		var q model.Question
		for i, z := range result {
			id, err := strconv.Atoi(z.Member.(string))
			if err != nil {
				log.Print("[ERROR]更新热门问题出现错误: "+err.Error()+": id")
			}
			q.ID = uint(id)
			question, err := q.Get()
			if err != nil {
				log.Print("[ERROR]更新热门问题出现错误: "+err.Error()+": get question")
			}
			question.Hot = z.Score
			GetTopQ().Store(strconv.Itoa(i+1), question)
		}
		log.Print("热门问题更新成功! ")
		log.Print(time.Now())
		time.Sleep(d)
	}
}
