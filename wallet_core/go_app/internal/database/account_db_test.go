package database

import (
	"database/sql"
	"testing"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type AccountDbTestSuite struct {
	suite.Suite
	db        *sql.DB
	AccountDb *AccountDb
	client    *entity.Client
}

func (s *AccountDbTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec(
		"CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date, updated_at date)",
	)
	db.Exec(
		"CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at date, updated_at date)",
	)
	s.AccountDb = NewAccountDb(db)
	s.client, err = entity.NewClient("John Doe", "j@j.com")
	s.Nil(err)
}

func (s *AccountDbTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
}

func TestAccountDbTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDbTestSuite))
}

func (s *AccountDbTestSuite) TestSave() {
	account := entity.NewAccount(s.client)
	err := s.AccountDb.Save(account)
	s.Nil(err)
}

func (s *AccountDbTestSuite) TestFindById() {
	s.db.Exec(
		"INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		s.client.Id,
		s.client.Name,
		s.client.Email,
		s.client.CreatedAt,
		s.client.UpdatedAt,
	)
	account := entity.NewAccount(s.client)
	err := s.AccountDb.Save(account)
	s.Nil(err)

	accountDb, err := s.AccountDb.FindById(account.Id)
	s.Nil(err)
	s.Equal(account.Id, accountDb.Id)
	s.Equal(account.Client.Id, accountDb.Client.Id)
	s.Equal(account.Client.Name, accountDb.Client.Name)
	s.Equal(account.Client.Email, accountDb.Client.Email)
	s.Equal(account.Balance, accountDb.Balance)
	s.Equal(account.CreatedAt.Unix(), accountDb.CreatedAt.Unix())
	s.Equal(account.UpdatedAt.Unix(), accountDb.UpdatedAt.Unix())
}
