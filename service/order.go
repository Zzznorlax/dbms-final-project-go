package service

import (
	"dbmsbackend/domain"
	"dbmsbackend/repository"
	"dbmsbackend/util"
	"fmt"
	"time"
)

type OrderRepository interface {
	Create(*domain.Order) (*domain.Order, error)
	Update(int, []*domain.OrderItem) (*domain.Order, error)
	FetchByID(int) (*domain.Order, error)
	FetchByBuyerID(int) ([]domain.Order, error)
	Query(map[string]interface{}) ([]domain.Order, error)
	Delete(int) error
	Initialize(*util.Config) error
}

type Order struct {
	repo OrderRepository
}

func (orderService *Order) Initialize(config *util.Config) (err error) {

	if config.DBDriver == "sqlite" {

		orderService.repo = new(repository.OrderSqldbRepository)
		err = orderService.repo.Initialize(config)

		if err != nil {
			err = fmt.Errorf("initializing Order service: %w", err)
		}
	}

	return

}

func (orderService *Order) New(buyerID int, products []domain.OrderItemDTO) (*domain.Order, error) {

	buyer := domain.User{
		ID: buyerID,
	}

	now := time.Now()

	entity := domain.Order{
		Buyer:     &buyer,
		Products:  domain.OrderItemDTOToEntity(products),
		Timestamp: &now,
	}

	return orderService.repo.Create(&entity)
}

func (orderService *Order) Update(id int, products []domain.OrderItemDTO) (*domain.Order, error) {

	return orderService.repo.Update(id, domain.OrderItemDTOToEntity(products))
}

func (orderService *Order) Query(condition map[string]interface{}) ([]domain.Order, error) {

	Orders, err := orderService.repo.Query(condition)

	return Orders, err
}

func (orderService *Order) Delete(id int) error {

	err := orderService.repo.Delete(id)

	return err
}

func (orderService *Order) GetByID(id int) (*domain.Order, error) {
	entity, err := orderService.repo.FetchByID(id)

	return entity, err
}

func (orderService *Order) GetByBuyerID(buyerID int) ([]domain.Order, error) {

	Orders, err := orderService.repo.FetchByBuyerID(buyerID)

	return Orders, err
}
