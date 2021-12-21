package domain

import "time"

type OrderItem struct {
	Product *Product
	Amount  int
	OrderID int
}

type Order struct {
	ID        uint
	Buyer     *User
	Products  []*OrderItem
	Timestamp *time.Time
}

type OrderDTO struct {
	ID         int            `json:"id"`
	BuyerName  string         `json:"buyerName"`
	BuyerEmail string         `json:"buyerEmail"`
	BuyerPhone string         `json:"buyerPhone"`
	Timestamp  string         `json:"timestamp"`
	Products   []OrderItemDTO `json:"products"`
}

type OrderRespDTO struct {
	Data OrderDTO `json:"data"`
}

type OrderListRespDTO struct {
	Data []OrderDTO `json:"data"`
}

type OrderItemDTO struct {
	ProductID int `json:"productId"`
	Amount    int `json:"amount"`
}

func OrderItemDTOToEntity(orderItems []OrderItemDTO) []*OrderItem {
	var orderItemEntities []*OrderItem

	for _, item := range orderItems {

		productEntity := Product{
			ID: item.ProductID,
		}

		orderItemEntities = append(orderItemEntities, &OrderItem{
			Amount:  item.Amount,
			Product: &productEntity,
		})
	}

	return orderItemEntities
}

func (order *Order) ToDTO() *OrderDTO {
	var orderItemDTOs []OrderItemDTO

	for _, item := range order.Products {
		orderItemDTOs = append(orderItemDTOs, OrderItemDTO{
			ProductID: item.Product.ID,
			Amount:    item.Amount,
		})
	}
	return &OrderDTO{
		ID:         int(order.ID),
		BuyerName:  order.Buyer.Name,
		BuyerEmail: order.Buyer.Email,
		BuyerPhone: order.Buyer.Phone,
		Timestamp:  order.Timestamp.Format(time.RFC3339),
		Products:   orderItemDTOs,
	}
}

func (order *Order) ToRespDTO() *OrderRespDTO {
	var orderItemDTOs []OrderItemDTO

	for _, item := range order.Products {
		orderItemDTOs = append(orderItemDTOs, OrderItemDTO{
			ProductID: item.Product.ID,
			Amount:    item.Amount,
		})
	}
	return &OrderRespDTO{
		Data: OrderDTO{
			ID:         int(order.ID),
			BuyerName:  order.Buyer.Name,
			BuyerEmail: order.Buyer.Email,
			BuyerPhone: order.Buyer.Phone,
			Timestamp:  order.Timestamp.Format(time.RFC3339),
			Products:   orderItemDTOs,
		},
	}
}
