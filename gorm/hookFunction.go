package gorm

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//题目3：钩子函数
//继续使用博客系统的模型。
//要求 ：
//为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
//为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

func (p *Post) AfterCreate(tx *gorm.DB) error {
	userID := p.UserID
	if 0 == userID {
		return fmt.Errorf("用户 ID 不能为空")
	}

	result := tx.Model(&User{}).
		Where("id = ?", userID).
		UpdateColumn("posts_quantity", gorm.Expr("posts_quantity + ?", 1))

	if result.Error != nil {
		return fmt.Errorf("更新用户文章数量失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到 ID 为 %d 的用户", p.UserID)
	}

	return nil
}

// AfterDelete 评论删除后，检查对应文章的剩余评论数，若为0则更新状态为"无评论"
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 1. 校验文章 ID 合法性
	if c.PostID == 0 {
		return fmt.Errorf("文章 ID 不能为空")
	}

	// 2. 统计该文章当前的评论总数
	var commentCount int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&commentCount).Error; err != nil {
		return fmt.Errorf("统计评论数量失败: %v", err)
	}

	// 3. 若评论数为 0，更新文章的 CommentStatus 为"无评论"；否则保持"有评论"
	var status string
	if commentCount == 0 {
		status = "无评论"
	} else {
		status = "有评论" // 若原状态为其他值（如"热门"），可根据业务调整
	}

	// 4. 更新文章状态
	result := tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_status", status)

	if result.Error != nil {
		return fmt.Errorf("更新文章评论状态失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到 ID 为 %d 的文章", c.PostID)
	}

	return nil
}

func main() {
	// 连接数据库
	dsn := "root:你的密码@tcp(127.0.0.1:3306)/blogdb?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 自动迁移表结构（确保新增字段生效）
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	// 测试 1：创建文章，验证用户文章数量自动更新
	user := User{Name: "测试用户", Email: "test@example.com"}
	db.Create(&user) // 初始 posts_quantity = 0

	post := Post{Title: "测试文章", Content: "测试内容", UserID: user.ID}
	if err := db.Create(&post).Error; err != nil {
		fmt.Printf("创建文章失败: %v\n", err)
	} else {
		var updatedUser User
		db.First(&updatedUser, user.ID)
		fmt.Printf("用户文章数量更新后: %d（预期 1）\n", updatedUser.PostsQuantity)
	}

	// 测试 2：删除评论，验证文章评论状态更新
	// 先创建一篇文章和一条评论
	testPost := Post{Title: "评论测试文章", Content: "评论测试内容", UserID: user.ID}
	db.Create(&testPost)
	comment := Comment{CommentText: "测试评论", UserID: user.ID, PostID: testPost.ID}
	db.Create(&comment)

	// 删除评论
	if err := db.Delete(&comment).Error; err != nil {
		fmt.Printf("删除评论失败: %v\n", err)
	} else {
		var updatedPost Post
		db.First(&updatedPost, testPost.ID)
		fmt.Printf("文章评论状态更新后: %s（预期 无评论）\n", updatedPost.CommentStatus)
	}
}

func CreatePostsHookUser(db *gorm.DB) {
	user := User{Name: "测试用户", Email: "test@example.com"}
	db.Create(&user) // 初始 posts_quantity = 0

	post := Post{Title: "测试文章", Content: "测试内容", UserID: user.ID}
	if err := db.Create(&post).Error; err != nil {
		fmt.Printf("创建文章失败: %v\n", err)
	} else {
		var updatedUser User
		db.First(&updatedUser, user.ID)
		fmt.Printf("用户文章数量更新后: %d（预期 1）\n", updatedUser.PostsQuantity)
	}
}

func DeleteCommentHookPost(db *gorm.DB) {
	testPost := Post{Title: "评论测试文章14", Content: "评论测试内容14", UserID: 1}
	db.Create(&testPost)
	comment := Comment{CommentText: "测试评论14", UserID: 1, PostID: testPost.ID}
	db.Create(&comment)
	tx := db.Model(&Post{}).Where("id = ?", testPost.ID).Update("CommentStatus", "有评论")
	if tx.Error != nil {
		fmt.Printf("更新文章评论状态失败: %v\n", tx.Error)
	} else {
		fmt.Printf("文章评论状态更新后: %s（预期 有评论）\n", testPost.CommentStatus)
	}

	time.Sleep(2 * time.Second)

	// 删除评论
	if err := db.Delete(&comment).Error; err != nil {
		fmt.Printf("删除评论失败: %v\n", err)
	} else {
		var updatedPost Post
		db.First(&updatedPost, testPost.ID)
		fmt.Printf("文章评论状态更新后: %s（预期 无评论）\n", updatedPost.CommentStatus)
	}
}
