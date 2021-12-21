package service

import (
	"dbmsbackend/domain"
	"dbmsbackend/repository"
	"dbmsbackend/util"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"time"
)

type ProductRepository interface {
	Create(*domain.Product) *domain.Product
	Update(*domain.Product) *domain.Product
	FetchByID(int) (*domain.Product, error)
	Query(map[string]interface{}) ([]domain.Product, error)
	Delete(int) error
	Initialize(*util.Config) error
}

type Product struct {
	repo ProductRepository
}

func (productService *Product) Initialize(config *util.Config) (err error) {

	if config.DBDriver == "sqlite" {

		productService.repo = new(repository.ProductSqldbRepository)
		err = productService.repo.Initialize(config)

		if err != nil {
			err = fmt.Errorf("initializing Product service: %w", err)
		}
	}

	return

}

func (productService *Product) New(name string, description string, picture string, inventory int, price int, startSaleTime *time.Time, endSaleTime *time.Time, ownerID int) *domain.Product {

	entity := domain.Product{
		Name:          name,
		Description:   description,
		Picture:       picture,
		Inventory:     inventory,
		Price:         price,
		StartSaleTime: startSaleTime,
		EndSaleTime:   endSaleTime,
		OwnerID:       ownerID,
	}

	return productService.repo.Create(&entity)
}

func (productService *Product) Update(id int, name string, description string, picture string, inventory int, price int, startSaleTime *time.Time, endSaleTime *time.Time) *domain.Product {

	entity := domain.Product{
		ID:            id,
		Name:          name,
		Description:   description,
		Picture:       picture,
		Inventory:     inventory,
		Price:         price,
		StartSaleTime: startSaleTime,
		EndSaleTime:   endSaleTime,
	}

	return productService.repo.Update(&entity)
}

func (productService *Product) Query(condition map[string]interface{}) ([]domain.Product, error) {

	Products, err := productService.repo.Query(condition)

	return Products, err
}

func (productService *Product) Delete(id int) error {

	err := productService.repo.Delete(id)

	return err
}

func (productService *Product) GetByID(id int) (*domain.Product, error) {
	entity, err := productService.repo.FetchByID(id)

	return entity, err
}

func (productService *Product) UploadToImgur(url string, clientID string, file *multipart.FileHeader) (string, error) {

	fileContent, _ := file.Open()
	byteContainer, _ := ioutil.ReadAll(fileContent)

	headers := make(map[string]string)

	headers["Authorization"] = fmt.Sprintf("Client-ID %v", clientID)
	target := util.ImgurUploadResp{}

	err := util.UploadFile(url, headers, file.Filename, byteContainer, &target)

	if err != nil {
		err = fmt.Errorf("uploading to imgur: %w", err)
	}

	return target.Data.Link, err
}
