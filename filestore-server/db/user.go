package db

import (
	mydb "github.com/bigpengry/netdisc/filestore-server/db/mysql"
	"fmt"
)

// User : 用户信息
type User struct{
	Username string
	Email string
	Phone string
	SignUpAt string
	Status int
}

// UserSignUp : 通过用户名及密码完成user表的注册操作
func UserSignUp(username, encPwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user(user_name,user_pwd)values(?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("Failed to insert,err:", err.Error())
		return false
	}
	ret, err := stmt.Exec(username, encPwd)
	if err != nil {
		fmt.Println("Failed to insert,err:", err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected == 0 {
		return false
	}
	return true
}

// UserSignIn : 用户登录操作
func UserSignIn(username, encPwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encPwd {
		return true
	}
	return false
}

// GetUserInfo : 获取用户信息
func GetUserInfo(username string)(*User,error){
	user:=new(User)
	stmt,err:=mydb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	defer stmt.Close()
	if err!=nil{
		fmt.Println(err.Error())
		return user,err
	}
	
	stmt.QueryRow(username).Scan(user.Username,user.SignUpAt)
	if err!=nil{
		return user,err
	}
	
	return user,nil
}

// UpdateToken : 刷新用户登陆的token
func UpdateToken(username, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`,`user_token`)value(?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
