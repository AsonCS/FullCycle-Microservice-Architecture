package main

import (
	"database/sql"
	"fmt"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/database"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/event"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/usecase/create_account"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/usecase/create_client"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/usecase/create_transaction"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/web"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/web/webserver"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/pkg/events"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "localhost", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS clients (id varchar(255), name varchar(255), email varchar(255), created_at date, updated_at date);",
	)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS accounts (id varchar(255), client_id varchar(255), balance float, created_at date, updated_at date);",
	)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date);",
	)
	if err != nil {
		panic(err)
	}

	// configMap := ckafka.ConfigMap{
	// 	"bootstrap.servers": "kafka:29092",
	// 	"group.id":          "wallet",
	// }
	// kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	// eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	// eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	// balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDb(db)
	accountDb := database.NewAccountDb(db)

	// ctx := context.Background()
	// uow := uow.NewUow(ctx, db)

	// uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
	// 	return database.NewAccountDB(db)
	// })

	// uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
	// 	return database.NewTransactionDB(db)
	// })
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(
		accountDb,
		eventDispatcher,
		transactionCreatedEvent,
		database.NewTransactionDb(db),
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
