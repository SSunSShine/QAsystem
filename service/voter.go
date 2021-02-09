package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/SSunSShine/QAsystem/database"
	"github.com/SSunSShine/QAsystem/model"
	"strconv"
)

var ctx = context.Background()

// AddVoter 点赞或点踩
func AddVoter(AnswerID, UserID uint, UpOrDown bool) (err error) {

	redisKey := getRedisKey(UserID, UpOrDown)
	isMember, err := IsMember(AnswerID, UserID)
	if err != nil {
		return
	}
	if isMember {
		err = errors.New("Has voted.")
		return
	}
	err = database.RDB.SAdd(ctx, redisKey, AnswerID).Err()
	return
}

// RemoveVoter 取消点赞或点踩
func RemoveVoter(AnswerID, UserID uint, UpOrDown bool) (err error) {

	redisKey := getRedisKey(UserID, UpOrDown)
	isMember, err := IsMember(AnswerID, UserID)
	if err != nil {
		return
	}
	if !isMember {
		err = errors.New("Hasn't voted.")
		return
	}
	err = database.RDB.SRem(ctx, redisKey, AnswerID).Err()
	return
}

// WrapVoted 标记回答是否被点赞或点踩
func WrapVoted(answer *model.Answer, UserID uint) (err error) {

	memMap, err := getMemMap(UserID)
	if err != nil {
		return
	}
	upID := strconv.Itoa(int(answer.ID))
	answer.Voted = memMap[upID]

	return
}

// WrapListSupported 标记回答列表是否被点赞或点踩
func WrapListSupported(answers []model.Answer, UserID uint) (err error) {

	memMap, err := getMemMap(UserID)
	if err != nil {
		return
	}

	for key, answer := range answers {
		upID := strconv.Itoa(int(answer.ID))
		answers[key].Voted = memMap[upID]
	}
	return
}

func getRedisKey(UserID uint, UpOrDown bool) (redisKey string) {
	if UpOrDown {
		redisKey = fmt.Sprintf("upvoted:%v", UserID)
	} else {
		redisKey = fmt.Sprintf("downvoted:%v", UserID)
	}
	return
}

// IsMember 是否点赞或者点踩
func IsMember(AnswerID, UserID uint) (isMember bool, err error) {

	up, err := database.RDB.SIsMember(ctx, getRedisKey(UserID, true), AnswerID).Result()
	if err != nil {
		return
	}
	down, err := database.RDB.SIsMember(ctx, getRedisKey(UserID, false), AnswerID).Result()
	if err != nil {
		return
	}
	isMember = up || down
	return
}

// getMemMap 获取用户点赞或点踩的redis数据，用map封装, 1 赞， -1 踩
func getMemMap(UserID uint) (memMap map[string]int, err error) {

	upList, err := database.RDB.SMembers(ctx, getRedisKey(UserID, true)).Result()
	if err != nil {
		return
	}
	memMap = make(map[string]int)
	for _, AnswerID := range upList {
		memMap[AnswerID] = 1
	}

	downList, err := database.RDB.SMembers(ctx, getRedisKey(UserID, false)).Result()
	if err != nil {
		return
	}
	for _, AnswerID := range downList {
		memMap[AnswerID] = -1
	}
	return
}
