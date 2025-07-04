package createtransaction

import (
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/gateway"
)

type CreateTransactionInputDto struct {
	AccountIdFrom string
	AccountIdTo   string
	Amount        float64
}

type CreateTransactionOutputDto struct {
	TransactionId string
}

type CreateTransactionUseCase struct {
	TransactionGateway gateway.TransactionGateway
	AccountGateway     gateway.AccountGateway
}

func NewCreateTransactionUseCase(
	transactionGateway gateway.TransactionGateway,
	accountGateway gateway.AccountGateway,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		TransactionGateway: transactionGateway,
		AccountGateway:     accountGateway,
	}
}

func (uc *CreateTransactionUseCase) Execute(input CreateTransactionInputDto) (*CreateTransactionOutputDto, error) {
	accountFrom, err := uc.AccountGateway.FindById(input.AccountIdFrom)
	if err != nil {
		return nil, err
	}

	accountTo, err := uc.AccountGateway.FindById(input.AccountIdTo)
	if err != nil {
		return nil, err
	}

	transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
	if err != nil {
		return nil, err
	}

	err = uc.TransactionGateway.Create(transaction)
	if err != nil {
		return nil, err
	}

	return &CreateTransactionOutputDto{
		TransactionId: transaction.ID,
	}, nil
}
