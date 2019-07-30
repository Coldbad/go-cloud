package db

import (
	"fmt"
	"go-cloud/db/mysql"
)

func UserSignup(username string, passwd string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`) values(?,?)")
	if err != nil {
		fmt.Println("sql prepare err:", err)
		return false
	}
	defer stmt.Close()
	result, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("exec err:", err)
		return false
	}
	if rows, err := result.RowsAffected(); nil == err && rows > 0 {
		return true
	}
	fmt.Println("rows affected err:", err)
	return false
}

func UserSignin(username string, encpwd string) bool {
	stmt, err := mysql.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println("sql prepare err:", err)
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("null" + username)
		return false
	}
	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

func UpdateToken(username string, token string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`,`user_token`) values(?,?)")
	if err != nil {
		fmt.Println("prepare err:", err)
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println("exec err:", err)
		return false
	}
	return true
}
