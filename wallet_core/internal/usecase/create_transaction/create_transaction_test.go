package createtransaction

import (
	"testing"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) FindById(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)
	client2, _ := entity.NewClient("John Doe 2", "j2@j2.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	accountGatewayMock := &AccountGatewayMock{}
	accountGatewayMock.On("FindById", account1.Id).Return(account1, nil)
	accountGatewayMock.On("FindById", account2.Id).Return(account2, nil)

	transactionGatewayMock := &TransactionGatewayMock{}
	transactionGatewayMock.On("Create", mock.Anything).Return(nil)

	inputDto := CreateTransactionInputDto{
		AccountIdFrom: account1.Id,
		AccountIdTo:   account2.Id,
		Amount:        100,
	}

	uc := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)
	output, err := uc.Execute(inputDto)

	assert.Nil(t, err)
	assert.NotNil(t, output)

	accountGatewayMock.AssertExpectations(t)
	transactionGatewayMock.AssertExpectations(t)

	accountGatewayMock.AssertNumberOfCalls(t, "FindById", 2)
	transactionGatewayMock.AssertNumberOfCalls(t, "Create", 1)
}
