package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// 存储玫瑰价格的模拟数据库
var rosePrice = map[string]float64{
	"2025-05-01": 50.00,
	"2025-05-02": 45.00,
	"2025-05-03": 48.00,
}

// 查询玫瑰价格的处理函数
func getRosePrice(w http.ResponseWriter, r *http.Request) {
	// 从查询参数中获取日期
	date := r.URL.Query().Get("date")

	// 如果日期不存在或无效
	if date == "" {
		http.Error(w, "缺少日期参数", http.StatusBadRequest)
		return
	}

	// 查询价格
	price, found := rosePrice[date]
	if !found {
		http.Error(w, "未找到该日期的玫瑰价格", http.StatusBadRequest)
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
	// 使用 gorilla/mux 路由器
	r := mux.NewRouter()

	// 定义路由
	r.HandleFunc("/api/rose-price", getRosePrice).Methods()

	// 启动HTTP服务器
	fmt.Println("服务器正在监听端口 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
