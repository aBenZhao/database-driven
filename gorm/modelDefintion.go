package gorm

import (
	"fmt"

	"gorm.io/gorm"
)

//题目1：模型定义
//假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
//要求 ：
//使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
//编写Go代码，使用Gorm创建这些模型对应的数据库表。

type User struct {
	ID            uint   `gorm:"type:bigint;primaryKey;autoIncrement" json:"id"`
	Name          string `gorm:"type:varchar(50);not null;unique" json:"name"`
	Email         string `gorm:"type:varchar(100);not null;uniqueIndex" json:"email"`
	Posts         []Post `gorm:"foreignKey:UserID" json:"posts"`
	PostsQuantity int    `gorm:"type:int;" json:"posts_quantity"`
}

type Post struct {
	ID            uint      `gorm:"type:bigint;primaryKey;autoIncrement" json:"id"`
	Title         string    `gorm:"type:varchar(100);not null" json:"title"`
	Content       string    `gorm:"type:text;not null" json:"content"`
	UserID        uint      `gorm:"type:bigint;not null" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user"`
	Comments      []Comment `gorm:"foreignKey:PostID" json:"comments"`
	CommentStatus string    `gorm:"type:varchar(20);default:'无评论'" json:"comment_status"` // 评论状态
}
type Comment struct {
	ID          uint   `gorm:"type:bigint;primaryKey;autoIncrement" json:"id"`
	CommentText string `gorm:"type:text;not null" json:"comment_text"`
	UserID      uint   `gorm:"type:bigint;not null" json:"user_id"`
	User        User   `gorm:"foreignKey:UserID" json:"user"`
	PostID      uint   `gorm:"type:bigint;not null" json:"post_id"`
	Post        Post   `gorm:"foreignKey:PostID" json:"post"`
}

func AutoMerge(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		return fmt.Errorf("自动迁移失败: %v", err)
	}
	return nil
}
