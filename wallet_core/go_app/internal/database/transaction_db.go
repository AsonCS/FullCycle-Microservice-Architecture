package database

import (
	"database/sql"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
)

type TransactionDb struct {
	Db *sql.DB
}

func NewTransactionDb(
	db *sql.DB,
) *TransactionDb {
	return &TransactionDb{
		Db: db,
	}
}

func (t *TransactionDb) Create(transaction *entity.Transaction) error {
	stmt, err := t.Db.Prepare(
		"INSERT INTO transactions (id, account_id_from, account_id_to, amount, created_at) VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		transaction.Id,
		transaction.AccountFrom.Id,
		transaction.AccountTo.Id,
		transaction.Amount,
		transaction.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
