package database

import (
	"database/sql"
	"testing"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type ClientDbTestSuite struct {
	suite.Suite
	db       *sql.DB
	ClientDb *ClientDb
}

func (s *ClientDbTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec(
		"CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date, updated_at date)",
	)
	s.ClientDb = NewClientDb(db)
}

func (s *ClientDbTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
}

func TestClientDbTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDbTestSuite))
}

func (s *ClientDbTestSuite) TestSave() {
	client := &entity.Client{
		Id:    "1",
		Name:  "John Doe",
		Email: "j@j.com",
	}
	err := s.ClientDb.Save(client)
	s.Nil(err)
}

func (s *ClientDbTestSuite) TestGet() {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	s.ClientDb.Save(client)

	clientDb, err := s.ClientDb.Get(client.Id)
	s.Nil(err)
	s.Equal(client.Id, clientDb.Id)
	s.Equal(client.Name, clientDb.Name)
	s.Equal(client.Email, clientDb.Email)
	s.Equal(client.CreatedAt.Unix(), clientDb.CreatedAt.Unix())
	s.Equal(client.UpdatedAt.Unix(), clientDb.UpdatedAt.Unix())
}
