package sql

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
)

// 题目1：基本CRUD操作
// 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
// 要求 ：
// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
type Student struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade string `json:"grade"`
}

// 1. 创建用户（Create）
func CreateStudent(ctx context.Context, name string, age int, grade string) (int64, error) {
	sql := "INSERT INTO students (name, age, grade) VALUES (?, ?, ?)"
	result, err := db.ExecContext(ctx, sql, name, age, grade)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// 2. 查询用户（Read）
func GetStudent(ctx context.Context) ([]*Student, error) {
	sql := "SELECT id,  name, age, grade FROM students WHERE age > 18"
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*Student
	for rows.Next() {
		var u Student
		err := rows.Scan(&u.Id, &u.Name, &u.Age, &u.Grade)
		if err != nil {
			return nil, err
		}
		students = append(students, &u)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return students, nil
}

// 3. 更新用户（Update）
func UpdateStudent(ctx context.Context, name string, grade string) (int64, error) {
	sql := "UPDATE students SET grade = ? WHERE name like ?"
	result, err := db.ExecContext(ctx, sql, grade, name)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// 4. 删除用户（Delete）
func DeleteStudent(ctx context.Context, age int) (int64, error) {
	sql := "DELETE FROM students WHERE age < ?"
	result, err := db.ExecContext(ctx, sql, age)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
