package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/database"
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
	// show databases;
	// select * from clients;select * from accounts;select * from transactions;
	// delete from clients;delete from accounts;delete from transactions;
	// update accounts set balance = 1000 where id = '4251701a-d554-4e38-8ed2-66034d152e8a';
	// go run cmd/walletcore/main.go
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
