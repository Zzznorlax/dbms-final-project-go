package domain

import "time"

type Product struct {
	ID            int
	Name          string
	Description   string
	Picture       string
	Inventory     int
	Price         int
	StartSaleTime *time.Time
	EndSaleTime   *time.Time
	OwnerID       int
}

type ProductDTO struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Inventory     int        `json:"inventory"`
	Price         int        `json:"price"`
	StartSaleTime *time.Time `json:"startSaleTime"`
	EndSaleTime   *time.Time `json:"endSaleTime"`
	Picture       string     `json:"picture"`
}

type ProductRespDTO struct {
	Data ProductDTO `json:"data"`
}

type ProductListRespDTO struct {
	Data []ProductDTO `json:"data"`
}

func (product *Product) ToRespDTO() *ProductRespDTO {
	return &ProductRespDTO{
		Data: ProductDTO{
			ID:            product.ID,
			Name:          product.Name,
			Description:   product.Description,
			Inventory:     product.Inventory,
			Price:         product.Price,
			StartSaleTime: product.StartSaleTime,
			EndSaleTime:   product.EndSaleTime,
			Picture:       product.Picture,
		},
	}
}

func (product *Product) ToDTO() *ProductDTO {
	return &ProductDTO{
		ID:            product.ID,
		Name:          product.Name,
		Description:   product.Description,
		Inventory:     product.Inventory,
		Price:         product.Price,
		StartSaleTime: product.StartSaleTime,
		EndSaleTime:   product.EndSaleTime,
		Picture:       product.Picture,
	}
}
