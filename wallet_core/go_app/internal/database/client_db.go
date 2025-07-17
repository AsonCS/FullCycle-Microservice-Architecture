package database

import (
	"database/sql"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
)

type ClientDb struct {
	Db *sql.DB
}

func NewClientDb(
	db *sql.DB,
) *ClientDb {
	return &ClientDb{
		Db: db,
	}
}

func (c *ClientDb) Get(id string) (*entity.Client, error) {
	client := &entity.Client{}
	stmt, err := c.Db.Prepare(
		"SELECT id, name, email, created_at, updated_at FROM clients WHERE id = ?",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	if err := row.Scan(
		&client.Id,
		&client.Name,
		&client.Email,
		&client.CreatedAt,
		&client.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *ClientDb) Save(client *entity.Client) error {
	stmt, err := c.Db.Prepare(
		"INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		client.Id,
		client.Name,
		client.Email,
		client.CreatedAt,
		client.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
