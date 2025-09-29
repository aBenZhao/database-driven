package sqlx

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// 题目1：使用SQL扩展库进行查询
// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
// 要求 ：
// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
type employee struct {
	Id         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func GetEmployeesByDepartment(ctx context.Context, department string, db *sqlx.DB) ([]employee, error) {
	query := "SELECT * FROM employees WHERE department = ?"

	var employees []employee
	err := db.SelectContext(ctx, &employees, query, department)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	return employees, nil
}

func GetHighestPaidEmployee(ctx context.Context, db *sqlx.DB) (employee, error) {
	query := "SELECT * FROM employees ORDER BY salary DESC LIMIT 1"

	var employee employee
	err := db.GetContext(ctx, &employee, query)
	if err != nil {
		return employee, fmt.Errorf("查询失败: %v", err)
	}

	return employee, nil
}
