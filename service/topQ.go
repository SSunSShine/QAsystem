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
func UpdateTopQ(d time.Duration)  {
	for  {
		result, err := database.RDB.ZRevRangeWithScores(ctx, ZSetKey, 0, 10).Result()
		if err != nil {
			log.Print("[ERROR]更新热门问题出现错误: "+err.Error())
		}
		var q model.Question
		for i, z := range result {
			id, err := strconv.Atoi(z.Member.(string))
			if err != nil {
				log.Print("[ERROR]更新热门问题出现错误: "+err.Error())
			}
			q.ID = uint(id)
			question, err := q.Get()
			if err != nil {
				log.Print("[ERROR]更新热门问题出现错误: "+err.Error())
			}
			GetTopQ().Store(strconv.Itoa(i+1), question)
		}
		log.Print("热门问题更新成功! ")
		log.Print(time.Now())
		time.Sleep(d)
	}
}
