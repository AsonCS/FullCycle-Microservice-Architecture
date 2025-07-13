package database

import (
	"database/sql"
	"testing"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDbTestSuite struct {
	suite.Suite
	db            *sql.DB
	client        *entity.Client
	client2       *entity.Client
	accountFrom   *entity.Account
	accountTo     *entity.Account
	transactionDb *TransactionDb
}

func (s *TransactionDbTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec(
		"CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date, updated_at date)",
	)
	db.Exec(
		"CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at date, updated_at date)",
	)
	db.Exec(
		"CREATE TABLE transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date)",
	)
	s.transactionDb = NewTransactionDb(db)
	s.client, err = entity.NewClient("John Doe", "j@j.com")
	s.Nil(err)
	s.client2, err = entity.NewClient("John Doe 2", "j2@j2.com")
	s.Nil(err)
	s.accountFrom = entity.NewAccount(s.client)
	s.accountFrom.Credit(1000)
	s.accountTo = entity.NewAccount(s.client2)
	s.accountTo.Credit(1000)
}

func (s *TransactionDbTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE transactions")
}

func TestTransactionDbTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDbTestSuite))
}

func (s *TransactionDbTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 100)
	s.Nil(err)
	err = s.transactionDb.Create(transaction)
	s.Nil(err)
}
