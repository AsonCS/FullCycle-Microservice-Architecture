package web

import (
	"encoding/json"
	"net/http"

	"github.com/AsonCS/FullCycle-Microservice-Architecture/internal/usecase/create_transaction"
)

type WebTransactionHandler struct {
	CreateTransactionUseCase create_transaction.CreateTransactionUseCase
}

func NewWebTransactionHandler(createTransactionUseCase create_transaction.CreateTransactionUseCase) *WebTransactionHandler {
	return &WebTransactionHandler{
		CreateTransactionUseCase: createTransactionUseCase,
	}
}

func (h *WebTransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var dto create_transaction.CreateTransactionInputDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	output, err := h.CreateTransactionUseCase.Execute(
		r.Context(),
		dto,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
