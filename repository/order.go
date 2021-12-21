package repository

import (
	"dbmsbackend/database/sqldb"
	"dbmsbackend/domain"
	"dbmsbackend/util"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderSqldbRepository struct {
	db *gorm.DB
}

func (repo *OrderSqldbRepository) Initialize(config *util.Config) (err error) {
	dbGetter := sqldb.NewSQLiteGetter()
	repo.db = dbGetter(config.DBSource)

	if err != nil {
		err = fmt.Errorf("connecting to db: %w", err)
	}

	repo.db.AutoMigrate(&sqldb.Order{})

	return err
}

func toOrderItemDao(entity *domain.OrderItem) *sqldb.OrderItem {
	dao := sqldb.OrderItem{
		ProductID: entity.Product.ID,
	}

	return &dao
}

func toOrderDao(entity *domain.Order) *sqldb.Order {

	var products []sqldb.OrderItem

	for _, item := range entity.Products {
		products = append(products, sqldb.OrderItem{
			ProductID: item.Product.ID,
			Amount:    item.Amount,
		})
	}

	dao := sqldb.Order{
		BuyerID:   entity.Buyer.ID,
		Products:  products,
		TimeStamp: entity.Timestamp,
	}

	return &dao
}

func toOrderEntity(dao *sqldb.Order) *domain.Order {

	var products []*domain.OrderItem

	for _, item := range dao.Products {
		products = append(products, &domain.OrderItem{
			Product: toProductEntity(&item.Product),
			Amount:  item.Amount,
		})
	}

	entity := domain.Order{
		ID:        dao.ID,
		Buyer:     toUserEntity(&dao.Buyer),
		Products:  products,
		Timestamp: dao.TimeStamp,
	}

	return &entity
}

func (repo *OrderSqldbRepository) Create(order *domain.Order) (*domain.Order, error) {

	dao := toOrderDao(order)

	if err := repo.db.Create(dao).Error; err != nil {
		return nil, fmt.Errorf("creating order: %w", err)
	}

	if err := repo.db.Preload(clause.Associations).First(&dao, dao.ID).Error; err != nil {
		return nil, fmt.Errorf("creating order %v: %w", dao.ID, err)
	}

	return toOrderEntity(dao), nil
}

func (repo *OrderSqldbRepository) Update(id int, orderItems []*domain.OrderItem) (*domain.Order, error) {

	dao := new(sqldb.Order)

	if err := repo.db.First(&dao, id).Error; err != nil {
		return nil, fmt.Errorf("updating order %v: %w", id, err)
	}

	var products []sqldb.OrderItem

	for _, item := range orderItems {
		products = append(products, sqldb.OrderItem{
			OrderID:   id,
			ProductID: item.Product.ID,
			Amount:    item.Amount,
		})
	}

	dao.Products = products

	if err := repo.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&dao).Error; err != nil {
		return nil, fmt.Errorf("updating order %v: %w", id, err)
	}

	if err := repo.db.Preload(clause.Associations).First(&dao, id).Error; err != nil {
		return nil, fmt.Errorf("updating order %v: %w", id, err)
	}

	return toOrderEntity(dao), nil
}

func (repo *OrderSqldbRepository) Query(condition map[string]interface{}) ([]domain.Order, error) {

	var daos []sqldb.Order
	var result []domain.Order

	if err := repo.db.Preload(clause.Associations).Where(condition).Find(&daos).Error; err != nil {
		return nil, fmt.Errorf("fetching orders: %w", err)
	}

	for _, item := range daos {
		result = append(result, *toOrderEntity(&item))
	}

	return result, nil
}

func (repo *OrderSqldbRepository) FetchByID(id int) (*domain.Order, error) {

	dao := new(sqldb.Order)

	if err := repo.db.Preload(clause.Associations).First(&dao, id).Error; err != nil {
		return nil, fmt.Errorf("fetching order %v: %w", id, err)
	}

	return toOrderEntity(dao), nil
}

func (repo *OrderSqldbRepository) FetchByBuyerID(buyerID int) ([]domain.Order, error) {

	var daos []sqldb.Order
	var result []domain.Order

	if err := repo.db.Preload(clause.Associations).Where("buyer_id = ?", buyerID).Find(&daos).Error; err != nil {
		return nil, fmt.Errorf("fetching orders by buyer ID %v: %w", buyerID, err)
	}

	for _, item := range daos {
		result = append(result, *toOrderEntity(&item))
	}

	return result, nil
}

func (repo *OrderSqldbRepository) Delete(id int) error {

	err := repo.db.Delete(&sqldb.Order{}, id).Error

	return err
}
