package gateway

import "github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"

type AccountGateway interface {
	Save(account *entity.Account) error
	FindById(id string) (*entity.Account, error)
	UpdateBalance(account *entity.Account) error
}
