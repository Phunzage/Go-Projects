package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// 全局数据库连接实例
var db *sql.DB

// 初始化数据库连接
func initDB() (err error) {
	// 数据库信息
	dsn := "root:Phz81114002@@tcp(127.0.0.1:3306)/rose_shop"
	// 检测数据库格式
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		// dsn格式不正确会报错
		return fmt.Errorf("无法连接到数据库 %v", err)
	}
	// 连接数据库
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("数据库连接失败 %v", err)
	}

	log.Println("成功连接到数据库")
	return nil
}

// 查询玫瑰价格的处理函数
func getRosePrice(w http.ResponseWriter, r *http.Request) {
	// 添加 CORS 头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	// 从查询参数中获取日期
	date := r.URL.Query().Get("date")

	// 如果日期不存在或无效
	if date == "" {
		http.Error(w, "缺少日期参数", http.StatusBadRequest)
		return
	}

	// 查询数据库中的价格
	var price float64
	err := db.QueryRow("select price from rose_prices where date = ?", date).Scan(&price)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "未找到该日期的玫瑰价格", http.StatusInternalServerError)
		} else {
			http.Error(w, "数据库查询出错", http.StatusInternalServerError)
		}
		return
	}

	// 设置响应头
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 创建响应结构体
	response := map[string]interface{}{
		"success": true,
		"price":   price,
	}

	// 将响应写入
	json.NewEncoder(w).Encode(response)
}

// 启动服务器并设置路由
func main() {
	// 连接数据库
	err := initDB()
	if err != nil {
		log.Fatal("无法连接到数据库：", err)
	}
	defer db.Close()

	// 使用 gorilla/mux 路由器
	r := mux.NewRouter()

	// 定义路由
	r.HandleFunc("/api/rose-price", getRosePrice).Methods("GET")

	// 添加静态文件服务
	fs := http.FileServer(http.Dir("../"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	// 启动HTTP服务器
	fmt.Println("服务器正在监听端口 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
