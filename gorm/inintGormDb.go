package gorm

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InintGormDb() (*gorm.DB, error) {
	// 数据库连接信息
	dsn := "root:bin5201314@tcp(127.0.0.1:3306)/gotestdb?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}
	//db.Logger = logger.Default.LogMode(logger.Info)
	return db, nil
}
