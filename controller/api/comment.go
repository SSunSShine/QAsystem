package api

import (
	"github.com/SSunSShine/QAsystem/model"
	"github.com/SSunSShine/QAsystem/service"
	"github.com/SSunSShine/QAsystem/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CommentVO struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Commenter Commenter `json:"answerer"`
	AnswerID  uint      `json:"answerId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Commenter struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"userId"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
}

func GetComment(c *gin.Context) {

	var co model.Comment
	var p model.Profile

	id, _ := strconv.Atoi(c.Param("id"))
	co.ID = uint(id)

	comment, err := co.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error() + ": comment",
		})
		return
	}

	p.UserID = co.UserID
	profile, err := p.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error() + ": profile",
		})
		return
	}

	var commentVO CommentVO
	util.SimpleCopyProperties(&commentVO, &comment)
	util.SimpleCopyProperties(&commentVO.Commenter, &profile)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    commentVO,
	})
}

func UpdateComment(c *gin.Context) {

	var co model.Comment

	id, _ := strconv.Atoi(c.Param("id"))
	co.ID = uint(id)

	comment, err := co.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": comment",
		})
		return
	}
	uid, _ := c.Get("uid")
	UID := uid.(uint)
	if UID != comment.UserID {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "无权修改他人的信息",
		})
		return
	}

	if err := c.ShouldBindJSON(&co); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error() + ": bind comment json",
		})
		return
	}
	// 防止json中的id 与 url的id不同
	if comment.ID != co.ID {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "JSON中的id与url中的id不同",
		})
		return
	}

	if err := co.Update(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error() + ": update comment",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
		})
	}
}

func DeleteComment(c *gin.Context) {

	var co model.Comment

	id, _ := strconv.Atoi(c.Param("id"))
	co.ID = uint(id)

	comment, err := co.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error()+": comment",
		})
		return
	}
	uid, _ := c.Get("uid")
	UID := uid.(uint)
	if UID != comment.UserID {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "无权修改他人的信息",
		})
		return
	}

	if err := co.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error() + ": delete comment",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
		})
	}
}

// GetCommentsCount 按用户id统计回答数量
func GetCommentsCount(c *gin.Context) {

	var co model.Comment

	userID, _ := strconv.Atoi(c.Query("userID"))
	co.UserID = uint(userID)

	if count, err := co.Count(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error() + ": comment count",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
			"data":    count,
		})
	}
}

func CreateComment(c *gin.Context) {

	var cc service.CreateCommentInterface
	var p model.Profile

	AnswerID, _ := strconv.Atoi(c.Query("AnswerID"))
	aid := uint(AnswerID)

	UserID, exist := c.Get("uid")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not exist: userid",
		})
		c.Abort()
		return
	}
	uid, ok := UserID.(uint)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not uint: userid",
		})
		c.Abort()
		return
	}

	if err := c.ShouldBindJSON(&cc); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error() + ": bind comment json",
		})
		return
	}

	comment, err := cc.Create(uid, aid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error() + ": userid or answerId",
		})
		return
	}

	p.UserID = comment.UserID
	profile, err := p.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error() + ": profile",
		})
		return
	}

	var commentVO CommentVO
	util.SimpleCopyProperties(&commentVO, &comment)
	util.SimpleCopyProperties(&commentVO.Commenter, &profile)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    commentVO,
	})
}

// GetCommentsByUser 按用户获取评论列表
func GetCommentsByUser(c *gin.Context) {

	var co model.Comment
	var p model.Profile
	var comments []model.Comment
	var err error

	userID, _ := strconv.Atoi(c.Query("userID"))
	co.UserID = uint(userID)

	order := c.Query("order")
	if order == "create_time" {
		if comments, err = co.GetOrderList("created_at desc"); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error() + ": comments",
			})
			return
		}
	} else {
		if comments, err = co.GetList(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error() + ": comments",
			})
			return
		}
	}

	p.UserID = co.UserID
	profile, err := p.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error() + ": profile",
		})
		return
	}

	count := 0
	var commentsVO []CommentVO
	for _, comment := range comments {
		var commentVO CommentVO
		util.SimpleCopyProperties(&commentVO, &comment)
		util.SimpleCopyProperties(&commentVO.Commenter, &profile)

		commentsVO = append(commentsVO, commentVO)
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    commentsVO,
		"total":   count,
	})
}

// GetCommentsByAnswer 按回答获取评论列表
func GetCommentsByAnswer(c *gin.Context) {

	var co model.Comment
	var p model.Profile
	var comments []model.Comment
	var err error

	answerID, _ := strconv.Atoi(c.Query("answerID"))
	co.AnswerID = uint(answerID)

	order := c.Query("order")
	if order == "create_time" {
		if comments, err = co.GetOrderList("created_at desc"); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error() + ": comments",
			})
			return
		}
	} else {
		if comments, err = co.GetList(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error() + ": comments",
			})
			return
		}
	}

	count := 0
	var commentsVO []CommentVO
	for _, comment := range comments {
		p.UserID = comment.UserID
		profile, err := p.Get()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error() + ": profile",
			})
			return
		}

		var commentVO CommentVO
		util.SimpleCopyProperties(&commentVO, &comment)
		util.SimpleCopyProperties(&commentVO.Commenter, &profile)

		commentsVO = append(commentsVO, commentVO)
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    commentsVO,
		"total":   count,
	})
}
