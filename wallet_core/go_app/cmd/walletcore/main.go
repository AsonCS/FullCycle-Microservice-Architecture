package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/database"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/event"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/event/handler"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/usecase/create_account"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/usecase/create_client"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/usecase/create_transaction"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/web"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/web/webserver"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/pkg/events"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/pkg/kafka"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// mysql -u root -proot wallet
	// select * from clients;select * from accounts;select * from transactions;
	// delete from clients;delete from accounts;delete from transactions;
	// update accounts set balance = 1000 where id = '4251701a-d554-4e38-8ed2-66034d152e8a';
	// go run cmd/walletcore/main.go
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	initDb(db)

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDb(db)
	accountDb := database.NewAccountDb(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDb", func(tx *sql.Tx) any {
		return database.NewAccountDb(db)
	})

	uow.Register("TransactionDb", func(tx *sql.Tx) any {
		return database.NewTransactionDb(db)
	})
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(
		balanceUpdatedEvent,
		eventDispatcher,
		transactionCreatedEvent,
		uow,
	)
	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")
	webserver.Start()
}

func initDb(
	db *sql.DB,
) {
	execute(
		db,
		"CREATE TABLE IF NOT EXISTS clients (id varchar(255), name varchar(255), email varchar(255), created_at date, updated_at date);",
	)
	execute(
		db,
		"CREATE TABLE IF NOT EXISTS accounts (id varchar(255), client_id varchar(255), balance float, created_at date, updated_at date);",
	)
	execute(
		db,
		"CREATE TABLE IF NOT EXISTS transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date);",
	)

	fulano, err := entity.NewClient("Fulano", "fulano@email.com")
	if err != nil {
		panic(err)
	}
	ciclano, err := entity.NewClient("Ciclano", "ciclano@email.com")
	if err != nil {
		panic(err)
	}
	execute(
		db,
		"INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?);",
		fulano.Id,
		fulano.Name,
		fulano.Email,
		fulano.CreatedAt,
		fulano.UpdatedAt,
	)
	execute(
		db,
		"INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?);",
		ciclano.Id,
		ciclano.Name,
		ciclano.Email,
		ciclano.CreatedAt,
		ciclano.UpdatedAt,
	)

	fulanoAccount := entity.NewAccount(fulano)
	fulanoAccount.Credit(1000)
	ciclanoAccount := entity.NewAccount(ciclano)
	execute(
		db,
		"INSERT INTO accounts (id, client_id, balance, created_at, updated_at) VALUES (?, ?, ?, ?, ?);",
		fulanoAccount.Id,
		fulano.Id,
		fulanoAccount.Balance,
		fulanoAccount.CreatedAt,
		fulanoAccount.UpdatedAt,
	)
	execute(
		db,
		"INSERT INTO accounts (id, client_id, balance, created_at, updated_at) VALUES (?, ?, ?, ?, ?);",
		ciclanoAccount.Id,
		ciclano.Id,
		ciclanoAccount.Balance,
		ciclanoAccount.CreatedAt,
		ciclanoAccount.UpdatedAt,
	)
}

func execute(
	db *sql.DB,
	query string,
	args ...any,
) {
	_, err := db.Exec(query, args...)
	if err != nil {
		panic(err)
	}
}
