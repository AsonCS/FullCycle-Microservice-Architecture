package gateway

import "github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
