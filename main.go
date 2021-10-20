package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var v dbConfig        // 数据库配置文件
var dbCfg = "db.json" // 数据库配置文件路径
var db *sql.DB

func main() {
	// MySQL 部分
	LoadDbConfig() // 读取数据库配置文件
	// 调试信息：输出数据库连接信息
	// fmt.Printf("Read database config from json file. User:%v,Password:%v,URL:%v,DB:%v",v.SqlUser,v.SqlPasswd,v.SqlURL,v.SqlDatabase)
	var err error
	db, err = sql.Open("mysql", v.SqlUser+":"+v.SqlPasswd+"@tcp("+v.SqlURL+")/"+v.SqlDatabase)
	log.Printf("Attempting connect to MySQL Server: %v\n", v.SqlUser+":"+v.SqlPasswd+"@tcp("+v.SqlURL+")/"+v.SqlDatabase)
	if err := db.Ping(); err != nil {
		log.Fatalf("Error raised when using MySQL: %v\n", err.Error())
	}
	// Gin 部分
	r := gin.Default() //设定主路由
	r.Use(Cors(), Handler(), gin.Recovery())
	r.LoadHTMLGlob("templates/**")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301,"https://www.yukishirataco.info")
		// c.HTML(http.StatusMovedPermanently, "text.html", gin.H{})
	})
	r.GET("/clear", clearHandler)
	// RESTful API POST
	r.POST("/comments", postCommentsHandler)
	// RESTful API GET
	r.GET("/comments", getCommentsHandler)
	// RESTful API Delete
	r.DELETE("/comments/:id",deleteCommentsHandler)
	err = r.Run(":8088")
	if err != nil {
		log.Fatalf("Error raised when router is running: %v\n", err.Error())
	} //运行服务端
	err = db.Close()
	if err != nil {
		log.Fatalf("Error raised when using MySQL: %v\n", err.Error())
	}
}
