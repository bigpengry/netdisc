package db

import (
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
)

// OnFileUploadFinished : 文件上传
func OnFileUploadFinished(filehash, filename, fileaddr string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`,`file_name`,`file_addr`,`file_size`,`status`) values(?,?,?,?,1)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("Failed to prepare satement,err:" + err.Error())
		return false
	}
	ret, err := stmt.Exec(filehash, filename, fileaddr, filesize)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been upload before", filehash)
		}
		return true
	}
	return false
}

// TableFile : a
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMeta :
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1=? and status=1 limit 1")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	tfile := new(TableFile)
	err = stmt.QueryRow(filehash).Scan(
		&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return tfile, nil
}
