package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"time"
)

var v dbConfig        // 数据库配置文件
var dbCfg = "db.json" // 数据库配置文件路径
var ReleaseMode = true

func LoadDbConfig() {
	JsonParse := NewJsonStruct()
	JsonParse.Load(dbCfg, &v)
}
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[%v] Addr %v called path %v use %v\n", time.Now(), c.ClientIP(), c.FullPath(), c.Request.Method)
		c.Next()
	}
}
func main() {
	// MySQL 部分
	LoadDbConfig() // 读取数据库配置文件
	// 调试信息：输出数据库连接信息
	// fmt.Printf("Read database config from json file. User:%v,Password:%v,URL:%v,DB:%v",v.SqlUser,v.SqlPasswd,v.SqlURL,v.SqlDatabase)
	db, err := sql.Open("mysql", v.SqlUser+":"+v.SqlPasswd+"@tcp("+v.SqlURL+")/"+v.SqlDatabase)
	log.Printf("Attempting connect to MySQL Server: %v\n", v.SqlUser+":"+v.SqlPasswd+"@tcp("+v.SqlURL+")/"+v.SqlDatabase)
	if err := db.Ping(); err != nil {
		log.Fatalf("Error raised when using MySQL: %v\n", err.Error())
	}
	// Gin 部分
	if ReleaseMode{
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New() //设定主路由
	r.Use(Cors(), Handler(), gin.Recovery())
	r.LoadHTMLGlob("templates/**")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "text.html", gin.H{})
	})
	r.GET("/clear",func(c *gin.Context){
		_,err := db.Exec("TRUNCATE TABLE comments")
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"result":1})
			log.Fatalln(err)
		}
		if err != nil{
			log.Fatalln(err)
			c.JSON(http.StatusInternalServerError,gin.H{"result":1})
			return
		}
		c.JSON(http.StatusOK,gin.H{"result":0})
	})
	// RESTful API POST
	r.POST("/comments", func(c *gin.Context) {
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
			rs, err := db.Exec("INSERT INTO comments(content,mail,ipaddr,time,nickname) VALUES (?,?,?,now(),?)", content, mail, c.ClientIP(), nickname)
			// TODO 多一个字段，也就是昵称
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"result": 1, "msg": err.Error()})
				return
			}
			id, err := rs.LastInsertId()
			if err != nil {
				log.Fatalln(err)
				c.JSON(http.StatusInternalServerError, gin.H{"result": 1, "msg": err.Error()})
				return
			}
			log.Println("insert person Id {}", id)
			msg := fmt.Sprintf("insert successful %d", id) // 插入了第几个数据？
			// c.SetCookie("commented", "yes",10,"/","127.0.0.1",false,true)
			c.JSON(http.StatusOK, gin.H{
				"result": 0,
				"msg":    msg,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"result": 1, "msg": "Invalid values."})
		}
	})
	// RESTful API GET
	r.GET("/comments", func(c *gin.Context) {
		rs, err := db.Query("SELECT id,content,mail,time FROM comments")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"result": 1, "msg": err})
		}
		comments := make([]frontComments, 0)
		for rs.Next() {
			var comment frontComments
			err = rs.Scan(&comment.CommentId, &comment.Content, &comment.Email, &comment.Date)
			comments = append(comments, comment)
		}
		if err != nil {
			log.Fatalln(err)
		}
		// c.SetCookie("comment", time.Now().String(),86400,"/","127.0.0.1",false,true)
		c.JSON(http.StatusOK, gin.H{"result": 0, "comments": comments})
	})
	// RESTful API Delete
	r.DELETE("/comments/:id",func(c *gin.Context){
		cid := c.Param("id")
		id, err := strconv.Atoi(cid)
		rs,err := db.Exec("DELETE FROM comments WHERE id=?",id)
		ra, err := rs.RowsAffected()
		if err != nil {
			log.Fatalln(err)
		}
		msg := fmt.Sprintf("Delete person %d successful %d", id, ra)
		log.Println(msg)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	})
	err = r.Run(":8088")
	if err != nil {
		log.Fatalf("Error raised when router is running: %v\n", err.Error())
	} //运行服务端
	err = db.Close()
	if err != nil {
		log.Fatalf("Error raised when using MySQL: %v\n", err.Error())
	}
}
