package repository

import (
	"dbmsbackend/database/sqldb"
	"dbmsbackend/domain"
	"dbmsbackend/util"
	"fmt"

	"gorm.io/gorm"
)

type ProductSqldbRepository struct {
	db *gorm.DB
}

func (repo *ProductSqldbRepository) Initialize(config *util.Config) (err error) {
	dbGetter := sqldb.NewSQLiteGetter()
	repo.db = dbGetter(config.DBSource)

	if err != nil {
		err = fmt.Errorf("connecting to db: %w", err)
	}

	repo.db.AutoMigrate(&sqldb.Product{})

	return err
}

func toProductDao(entity *domain.Product) *sqldb.Product {
	dao := sqldb.Product{
		Name:          entity.Name,
		Description:   entity.Description,
		Picture:       entity.Picture,
		Inventory:     entity.Inventory,
		Price:         entity.Price,
		StartSaleTime: entity.StartSaleTime,
		EndSaleTime:   entity.EndSaleTime,
		OwnerID:       entity.OwnerID,
	}

	return &dao
}

func toProductEntity(dao *sqldb.Product) *domain.Product {
	entity := domain.Product{
		ID:            int(dao.ID),
		Name:          dao.Name,
		Description:   dao.Description,
		Picture:       dao.Picture,
		Inventory:     dao.Inventory,
		Price:         dao.Price,
		StartSaleTime: dao.StartSaleTime,
		EndSaleTime:   dao.EndSaleTime,
		OwnerID:       dao.OwnerID,
	}

	return &entity
}

func (repo *ProductSqldbRepository) Create(Product *domain.Product) *domain.Product {

	dao := toProductDao(Product)

	repo.db.Create(dao)

	return toProductEntity(dao)
}

func (repo *ProductSqldbRepository) Update(Product *domain.Product) *domain.Product {

	dao := new(sqldb.Product)

	repo.db.First(&dao, Product.ID)

	dao.Name = Product.Name
	dao.Description = Product.Description
	dao.Picture = Product.Picture
	dao.Inventory = Product.Inventory
	dao.StartSaleTime = Product.StartSaleTime
	dao.EndSaleTime = Product.EndSaleTime

	repo.db.Save(&dao)

	return toProductEntity(dao)
}

func (repo *ProductSqldbRepository) Query(condition map[string]interface{}) ([]domain.Product, error) {

	var daos []sqldb.Product
	var result []domain.Product

	if err := repo.db.Where(condition).Find(&daos).Error; err != nil {
		return nil, fmt.Errorf("fetching products: %w", err)
	}

	for _, item := range daos {
		result = append(result, *toProductEntity(&item))
	}

	return result, nil
}

func (repo *ProductSqldbRepository) FetchByID(id int) (*domain.Product, error) {

	dao := new(sqldb.Product)

	if err := repo.db.First(&dao, id).Error; err != nil {
		return nil, fmt.Errorf("fetching product %v: %w", id, err)
	}

	return toProductEntity(dao), nil
}

func (repo *ProductSqldbRepository) Delete(id int) error {

	err := repo.db.Delete(&sqldb.Product{}, id).Error

	return err
}
