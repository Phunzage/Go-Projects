package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 导入包但不直接使用
)

var db *sql.DB // 是一个连接池

func initDB() (err error) {
	// 数据库信息
	// Data Source Name
	dsn := "root:3143@tcp(127.0.0.1:3306)/go0408"
	// 检测数据库格式
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		// dsn格式不正确会报错
		return
	}
	// 连接数据库
	err = db.Ping()
	if err != nil {
		return
	}
	// 设置最大连接数
	db.SetMaxOpenConns(10)
	// 设置最大空闲连接数
	db.SetMaxIdleConns(5)
	return
}

type user struct {
	id   int
	name string
	age  int
}

func QuaryOne(id int) {
	var u1 user

	// 1.写查询单条记录的sql语句
	sqlStr := `select id, name, age from user where id=?`
	// 2.执行并拿到结果
	// 从连接池拿一个连接出来，去数据库查询单条信息
	// // Scan具有自动释放连接的功能,rowObj必须要调用Scan方法
	db.QueryRow(sqlStr, id).Scan(&u1.id, &u1.name, &u1.age)
	// 3.打印输出
	fmt.Printf("u1: %v\n", u1)
}

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("初始化数据库失败，err: %v\n", err)
		return
	}
	fmt.Println("连接数据库成功")
	for i := 1; i < 4; i++ {
		QuaryOne(i)
	}
}
