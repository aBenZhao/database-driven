package sqlx

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDB() (*sqlx.DB, error) {
	dsn := "root:bin5201314@tcp(127.0.0.1:3306)/gotestdb?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return db, nil
}

func CloseDb(db *sqlx.DB) {
	err := db.Close()
	if err != nil {
		return
	} // 程序结束时关闭连接
}
