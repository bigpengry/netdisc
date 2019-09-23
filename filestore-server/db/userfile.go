package db

import (
	"fmt"
	"time"
	mydb "netdisc/filestore-server/db/mysql"
)

// UserFile : 用户文件表结构体
type UserFile struct{
	UserName string
	FileHash string
	FileName string
	FileSize int
	UpdateAt string
	LastUpdated string
}

// UploadUserFileFinished : 
func UploadUserFileFinished(userName ,fileHash ,fileName string,fileSize int64)bool{
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
			"`file_size`,`upload_at`,`status`) values (?,?,?,?,?,1)")
	defer stmt.Close()
	if err != nil{
		return false
	}
	_,err = stmt.Exec(userName ,fileHash ,fileName ,fileSize,time.Now())
	if err!=nil{
		return false
	}
	return true
}

// QueryUserFileMetas : 批量获取用户文件信息
func QueryUserFileMetas(username string,limit int)([]UserFile,error){
	stmt,err:=mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name =? limit ?")
	defer stmt.Close()
	if err!=nil{
		return nil,err
	}
	rows,err:=stmt.Query(username,limit)
	if err!=nil{
		return nil,err
	}

	ufs:=make([]UserFile,0)
	for rows.Next(){
		uf:=UserFile{}
		err=rows.Scan(&uf.FileHash,&uf.FileName,&uf.FileSize,&uf.UpdateAt,&uf.LastUpdated)
		if err!=nil{
			fmt.Println(err)
			break
		}
		ufs=append(ufs,uf)
	}
	return ufs,nil
}