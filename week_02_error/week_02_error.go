package week_02_error

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func conn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:qt123@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Println("connect fail")
		return nil, err
	}

	return db, nil
}

// dao层只需要加上
func daoGetUserNameById(id int) (string, error) {
	var name string

	db, err := conn()
	if err != nil {
		return name, err
	}

	defer db.Close()

	rowErr := db.QueryRow("select name from users where id = ?", id).Scan(&name)

	if rowErr != nil {
		return name, errors.Wrap(rowErr, fmt.Sprintf("select user fail,id : %d", id))
	}

	return name, nil
}

func GetUserNameById(id int) {
	userName, err := daoGetUserNameById(id)

	if err != nil {
		// 如果为查询不到数据，则直接打印错误信息，其他错误则写入日志
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("original :%T %v\n", errors.Cause(err), errors.Cause(err))
			fmt.Printf("trace :\n%+v\n", err)

			return
		} else {
			log.Printf("original :%T %v\n", errors.Cause(err), errors.Cause(err))
			log.Printf("trace :\n%+v\n", err)

			return
		}
	}

	fmt.Println("get success", userName)
}
