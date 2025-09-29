package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	Id      int     `json:"id"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	Id            int     `json:"id"`
	FromAccountId int     `json:"from_account_id"`
	ToAccountId   int     `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

//题目2：事务语句
//假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）
// transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
//要求 ：
//编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。
//在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，
//并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

// transfer 实现转账事务
// fromID: 转出账户ID, toID: 转入账户ID, amount: 转账金额
func Transfer(ctx context.Context, fromID int, toID int, amount float64) error {
	// 1. 开启事务（带上下文，支持超时/取消）
	tx, err := db.BeginTx(ctx, nil) // nil 表示使用默认隔离级别
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	// 延迟处理：若未提交则回滚（确保异常时不遗留未处理事务）
	defer func() {
		if r := recover(); r != nil {
			// 发生panic时回滚
			tx.Rollback()
			log.Printf("事务panic，已回滚: %v", r)
		} else if err != nil {
			// 已有错误时回滚
			tx.Rollback()
		}
	}()

	// 2. 查询转出账户的余额（加行锁，防止并发修改，MySQL 中 for update 会锁定行）
	var fromAcount Account
	queryBalance := "SELECT id,balance FROM accounts WHERE id = ? FOR UPDATE"
	rowContext := tx.QueryRowContext(ctx, queryBalance, fromID)
	err = rowContext.Scan(&fromAcount.Id, &fromAcount.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("转出账户 %d 不存在", fromID)
		}
		return fmt.Errorf("查询余额失败: %v", err)
	}
	fmt.Println(&fromAcount)
	if fromAcount.Balance < amount {
		return fmt.Errorf("余额不足，当前余额: %.2f, 需要: %.2f", fromAcount.Balance, amount)
	}

	// 4. 扣减转出账户余额
	updateFrom := "UPDATE accounts SET balance = balance - ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, updateFrom, amount, fromID)
	if err != nil {
		return fmt.Errorf("扣减余额失败: %v", err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("转出账户 %d 不存在（扣减失败）", fromID)
	}

	// 5. 增加转入账户余额
	updateTo := "UPDATE accounts SET balance = balance + ? WHERE id = ?"
	result, err = tx.ExecContext(ctx, updateTo, amount, toID)
	if err != nil {
		return fmt.Errorf("增加余额失败: %v", err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("转入账户 %d 不存在（增加失败）", toID)
	}

	// 6. 记录转账记录
	insertTx := `
		INSERT INTO transactions (from_account_id, to_account_id, amount)
		VALUES (?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, insertTx, fromID, toID, amount)
	if err != nil {
		return fmt.Errorf("记录转账失败: %v", err)
	}

	// 7. 所有操作成功，提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}
