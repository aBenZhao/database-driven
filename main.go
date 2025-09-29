package main

import (
	"context"
	"database-driven/sql"
	"log"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// 初始化数据库连接
	if err := sql.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer sql.CloseDb()

	ctx := context.Background()

	// 题目一：
	//1. 创建用户
	//studentID, err := sql.CreateStudent(ctx, "张三", 20, "三年级")
	//if err != nil {
	//	log.Printf("创建失败: %v", err)
	//} else {
	//	log.Printf("创建用户成功，ID: %d", studentID)
	//}

	// 2. 查询年龄大于18岁的学生信息
	//studentSlice, err := sql.GetStudent(ctx)
	//if err != nil {
	//	log.Printf("查询失败: %v", err)
	//} else {
	//	for _, student := range studentSlice {
	//		log.Printf("查询到用户: %+v", student)
	//	}
	//}

	// 3. 更新用户
	//rowsAffected, err := sql.UpdateStudent(ctx, "张三", "四年级")
	//if err != nil {
	//	log.Printf("更新失败: %v", err)
	//} else {
	//	log.Printf("更新成功，影响行数: %d", rowsAffected)
	//}

	// 4. 删除用户
	//deleteStudent, err := sql.DeleteStudent(ctx, 15)
	//if err != nil {
	//	log.Printf("删除失败: %v", err)
	//} else {
	//	log.Printf("删除成功，影响行数: %d", deleteStudent)
	//}

	// 题目二：
	//转账参数：从账户1向账户2转账100元
	fromID := 1
	toID := 2
	amount := 100.0

	// 创建上下文（可设置超时，例如5秒）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 实现从账户 A 向账户 B 转账 100 元的操作。
	if err := sql.Transfer(ctx, fromID, toID, amount); err != nil {
		log.Printf("转账失败: %v", err)
	} else {
		log.Printf("转账成功: 从账户%d向账户%d转账%.2f元", fromID, toID, amount)
	}
}
