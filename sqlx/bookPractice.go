package sqlx

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//题目2：实现类型安全映射
//假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
//要求 ：
//定义一个 Book 结构体，包含与 books 表对应的字段。
//编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

type Book struct {
	Id     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

func QueryPricesGreaterThan(ctx context.Context, i int, db *sqlx.DB) ([]Book, error) {
	query := "SELECT * FROM books WHERE price > ?"
	var books []Book
	err := db.SelectContext(ctx, &books, query, i)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	return books, nil
}
