package database

import (
	"database/sql"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
)

type AccountDb struct {
	Db *sql.DB
}

func NewAccountDb(db *sql.DB) *AccountDb {
	return &AccountDb{
		Db: db,
	}
}

func (a *AccountDb) FindById(id string) (*entity.Account, error) {
	var account entity.Account
	var client entity.Client
	account.Client = &client

	stmt, err := a.Db.Prepare(
		"SELECT a.id, a.client_id, a.balance, a.created_at, a.updated_at, c.id, c.name, c.email, c.created_at, c.updated_at FROM accounts a INNER JOIN clients c ON a.client_id = c.id WHERE a.id = ?",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(
		&account.Id,
		&account.Client.Id,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
		&client.Id,
		&client.Name,
		&client.Email,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (a *AccountDb) Save(account *entity.Account) error {
	stmt, err := a.Db.Prepare(
		"INSERT INTO accounts (id, client_id, balance, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		account.Id,
		account.Client.Id,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountDb) Get(id string) (*entity.Account, error) {
	account := &entity.Account{}
	stmt, err := a.Db.Prepare(
		"SELECT id, client_id, balance, created_at, updated_at FROM accounts WHERE id = ?",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	if err := row.Scan(
		&account.Id,
		&account.Client.Id,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return account, nil

}
