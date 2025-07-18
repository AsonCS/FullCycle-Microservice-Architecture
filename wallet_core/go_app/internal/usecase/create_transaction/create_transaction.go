package create_transaction

import (
	"context"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/gateway"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/pkg/events"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/pkg/uow"
)

type CreateTransactionInputDto struct {
	AccountIdFrom string  `json:"account_id_from"`
	AccountIdTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDto struct {
	OriginAccountId      string  `json:"origin_account_id"`
	DestinationAccountId string  `json:"destination_account_id"`
	Amount               float64 `json:"amount"`
	TransactionId        string  `json:"transaction_id"`
}

type BalanceUpdatedOutputDto struct {
	OriginAccountBalance      float64 `json:"origin_account_balance"`
	OriginAccountId           string  `json:"origin_account_id"`
	DestinationAccountBalance float64 `json:"destination_account_balance"`
	DestinationAccountId      string  `json:"destination_account_id"`
}

type CreateTransactionUseCase struct {
	BalanceUpdated     events.EventInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	Uow                uow.UowInterface
}

func NewCreateTransactionUseCase(
	balanceUpdated events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	uow uow.UowInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		BalanceUpdated:     balanceUpdated,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
		Uow:                uow,
	}
}

func (uc *CreateTransactionUseCase) Execute(
	ctx context.Context,
	input CreateTransactionInputDto,
) (*CreateTransactionOutputDto, error) {
	output := &CreateTransactionOutputDto{}
	balanceUpdatedOutput := &BalanceUpdatedOutputDto{}
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := uc.getAccountRepository(ctx)
		transactionRepository := uc.getTransactionRepository(ctx)

		accountFrom, err := accountRepository.FindById(input.AccountIdFrom)
		if err != nil {
			return err
		}
		accountTo, err := accountRepository.FindById(input.AccountIdTo)
		if err != nil {
			return err
		}
		transaction, err := entity.NewTransaction(
			accountFrom,
			accountTo,
			input.Amount,
		)
		if err != nil {
			return err
		}

		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}

		err = accountRepository.UpdateBalance(accountTo)
		if err != nil {
			return err
		}

		err = transactionRepository.Create(transaction)
		if err != nil {
			return err
		}
		output.OriginAccountId = input.AccountIdFrom
		output.DestinationAccountId = input.AccountIdTo
		output.Amount = input.Amount
		output.TransactionId = transaction.Id

		balanceUpdatedOutput.OriginAccountId = input.AccountIdFrom
		balanceUpdatedOutput.DestinationAccountId = input.AccountIdTo
		balanceUpdatedOutput.OriginAccountBalance = accountFrom.Balance
		balanceUpdatedOutput.DestinationAccountBalance = accountTo.Balance
		return nil
	})
	if err != nil {
		return nil, err
	}

	uc.TransactionCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)

	uc.BalanceUpdated.SetPayload(balanceUpdatedOutput)
	uc.EventDispatcher.Dispatch(uc.BalanceUpdated)
	return output, nil
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDb")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := uc.Uow.GetRepository(ctx, "TransactionDb")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGateway)
}
