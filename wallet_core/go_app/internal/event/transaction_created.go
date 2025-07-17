package event

import "time"

type TransactionCreated struct {
	Name    string `json:"name"`
	Payload any    `json:"payload"`
}

func NewTransactionCreated() *TransactionCreated {
	return &TransactionCreated{
		Name: "TransactionCreated",
	}
}

func (e *TransactionCreated) GetName() string {
	return e.Name
}

func (e *TransactionCreated) GetPayload() interface{} {
	return e.Payload
}

func (e *TransactionCreated) SetPayload(payload any) {
	e.Payload = payload
}

func (e *TransactionCreated) GetDateTime() time.Time {
	return time.Now()
}
