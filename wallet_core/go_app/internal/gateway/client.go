package gateway

import "github.com/AsonCS/FullCycle-Microservice-Architecture/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}
