package main

import (
	"database/sql"
)

func getComments() (comments []frontComments, err error){
	comments = make([]frontComments,0)
	rows, err := db.Query("SELECT id,content,mail,time FROM comments")
	if err != nil{
		return
	}
	defer rows.Close()
	for rows.Next(){
		var comment frontComments
		rows.Scan(&comment.CommentId, &comment.Content, &comment.Email, &comment.Date)
		comments = append(comments, comment)
	}
	if err = rows.Err();err!=nil{
		return
	}
	return
}

func truncateTable() (_ sql.Result, err error) {
	return db.Exec("TRUNCATE TABLE comments")
}

func insertComments(content string, email string, nickname string, ip string) (id int, err error){
	rs, err := db.Exec("INSERT INTO comments(content,mail,ipaddr,time,nickname) VALUES (?,?,?,now(),?)", content, email, ip, nickname)
	if err != nil {
		return
	}
	iid, err := rs.LastInsertId()
	if err != nil {
		return
	}
	id = int(iid)
	return
}

func deleteComments(id int)(ra int64,err error){
	rs, err := db.Exec("DELETE FROM comments WHERE id=?", id)
	ra, err = rs.RowsAffected()
	if err != nil {
		return
	}
	return
}
