package sql

import (
	"database/sql"
	"fmt"
	"time"
)

var db *sql.DB // 全局数据库连接

// 初始化数据库连接
func InitDB() error {
	// 数据库连接信息
	dsn := "root:bin5201314@tcp(127.0.0.1:3306)/gotestdb?charset=utf8mb4&parseTime=true&loc=Local"

	// 打开连接
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// 测试连接
	return db.Ping()
}

func CloseDb() {
	err := db.Close()
	if err != nil {
		return
	} // 程序结束时关闭连接
}
