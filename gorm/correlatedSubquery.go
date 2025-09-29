package gorm

import (
	"fmt"

	"gorm.io/gorm"
)

// 题目2：关联查询,
//基于上述博客系统的模型定义。
//要求 ：
//编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。,
//编写Go代码，使用Gorm查询评论数量最多的文章信息。

func QueryUserPosts(db *gorm.DB, userId uint) ([]User, error) {
	var user []User
	//Preload("Posts")：查询用户时，同时加载其关联的 Posts 切片（所有文章）。
	//Preload("Posts.Comments")：嵌套预加载，在加载 Posts 的同时，为每篇文章加载其关联的 Comments 切片（所有评论）。
	//最终返回的 user 结构体中，user.Posts 包含该用户的所有文章，每个 Post 的 Comments 字段包含对应评论。
	tx := db.Preload("Posts.Comments").First(&user, userId)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to query user: %w", tx.Error)
	}
	return user, nil
}

// 查询评论数量最多的文章
func GetPostWithMostComments(db *gorm.DB) (*Post, error) {
	var post Post
	// 1. 子查询：统计每篇文章的评论数，别名 comment_count
	// 2. 关联查询：按 comment_count 降序，取第一条
	result := db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON posts.id = comments.post_id").
		Group("posts.id").
		Order("comment_count DESC").
		First(&post)

	if result.Error != nil {
		return nil, fmt.Errorf("查询失败: %v", result.Error)
	}
	return &post, nil
}
