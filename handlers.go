package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[%v] Addr %v called path %v use %v\n", time.Now(), c.ClientIP(), c.FullPath(), c.Request.Method)
		c.Next()
	}
}

func clearHandler(c *gin.Context){
	_,err := truncateTable()
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"result":1})
		log.Fatalln(err)
	}
	c.JSON(http.StatusOK,gin.H{"result":0})
}

func getCommentsHandler(c *gin.Context) {
	comments := make([]frontComments, 0)
	comments,err := getComments()
	if err != nil{
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": 0, "comments": comments})
	return
}

func postCommentsHandler(c *gin.Context){
		var content, mail, nickname string
		var comment CommentsPost
		if err := c.ShouldBind(&comment); err == nil {
			content = comment.Text
			mail = comment.Mail
			nickname = comment.Nickname
			fmt.Printf("Mail: %v Content: %v\n", mail, content)
		} else {
			panic(err.Error())
		}
		if mail != "" && VerifyEmailFormat(mail) {
			// 前端一个判空，后端需要判断空值还要判断邮箱格式
			id,err := insertComments(content, mail, nickname, c.ClientIP())
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"result": 1})
			}
			log.Println("insert person Id {}", id)
			msg := fmt.Sprintf("insert successful %d", id) // 插入了第几个数据？
			c.JSON(http.StatusOK, gin.H{
				"result": 0,
				"msg":    msg,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"result": 1, "msg": "Invalid values."})
		}
}

func deleteCommentsHandler(c *gin.Context){
	cid := c.Param("id")
	id, _ := strconv.Atoi(cid)
	ra, err := deleteComments(id)
	if err != nil {
		log.Fatalln(err)
		return
	}
	msg := fmt.Sprintf("Delete person %d successful %d", id, ra)
	log.Println(msg)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}