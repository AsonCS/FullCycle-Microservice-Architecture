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
	AccountIdFrom string  `json:"account_id_from"`
	AccountIdTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
	TransactionId string  `json:"transaction_id"`
}

type CreateTransactionUseCase struct {
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	Uow                uow.UowInterface
}

func NewCreateTransactionUseCase(
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	uow uow.UowInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
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
	//balanceUpdatedOutput := &BalanceUpdatedOutputDto{}
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
		output.AccountIdFrom = input.AccountIdFrom
		output.AccountIdTo = input.AccountIdTo
		output.Amount = input.Amount
		output.TransactionId = transaction.Id

		// balanceUpdatedOutput.AccountIDFrom = input.AccountIdFrom
		// balanceUpdatedOutput.AccountIDTo = input.AccountIdTo
		// balanceUpdatedOutput.BalanceAccountIDFrom = accountFrom.Balance
		// balanceUpdatedOutput.BalanceAccountIDTo = accountTo.Balance
		return nil
	})
	if err != nil {
		return nil, err
	}

	uc.TransactionCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)

	//uc.BalanceUpdated.SetPayload(balanceUpdatedOutput)
	//uc.EventDispatcher.Dispatch(uc.BalanceUpdated)
	return output, nil
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := uc.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGateway)
}
