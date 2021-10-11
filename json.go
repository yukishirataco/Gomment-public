package main

import (
	"encoding/json"
	"io/ioutil"
)

type dbConfig struct {
	SqlUser     string
	SqlPasswd   string
	SqlURL      string
	SqlDatabase string
}

type JsonStruct struct {
}

type frontComments struct {
	Email     string
	Content   string
	Date      string
	CommentId string
}

type CommentsPost struct{
	Mail string `json:"mail"`
	Text string `json:"text"`
	Nickname string `json:"nickname"`
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}


