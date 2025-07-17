package create_account

import (
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/gateway"
)

type CreateAccountInputDto struct {
	ClientId string `json:"client_id"`
}

type CreateAccountOutputDto struct {
	AccountId string `json:"account_id"`
}

type CreateAccountUseCase struct {
	AccountGateway gateway.AccountGateway
	ClientGateway  gateway.ClientGateway
}

func NewCreateAccountUseCase(a gateway.AccountGateway, c gateway.ClientGateway) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountGateway: a,
		ClientGateway:  c,
	}
}

func (uc *CreateAccountUseCase) Execute(input CreateAccountInputDto) (*CreateAccountOutputDto, error) {
	client, err := uc.ClientGateway.Get(input.ClientId)
	if err != nil {
		return nil, err
	}
	account := entity.NewAccount(client)
	err = uc.AccountGateway.Save(account)
	if err != nil {
		return nil, err
	}
	return &CreateAccountOutputDto{
		AccountId: account.Id,
	}, nil
}
