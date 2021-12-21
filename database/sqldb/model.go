package sqldb

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string
	Email          string `gorm:"index:UX_user_email,unique"`
	Phone          string
	Kind           string
	HashedPassword string
	Salt           string
}

type UserToken struct {
	gorm.Model
	User     User `gorm:"foreignKey:UserID"`
	UserID   int
	Token    string
	ExpireAt time.Time
}

type Product struct {
	gorm.Model
	Name          string
	Description   string
	Picture       string
	Inventory     int
	Price         int
	StartSaleTime *time.Time
	EndSaleTime   *time.Time
	Owner         User `gorm:"foreignKey:OwnerID"`
	OwnerID       int
}

type OrderItem struct {
	Product   Product `gorm:"foreignKey:ProductID"`
	ProductID int     `gorm:"primaryKey;autoIncrement:false"`
	Amount    int
	OrderID   int `gorm:"primaryKey;autoIncrement:false"`
}

type Order struct {
	gorm.Model
	Buyer     User `gorm:"foreignKey:BuyerID"`
	BuyerID   int
	TimeStamp *time.Time
	Products  []OrderItem `gorm:"foreignKey:OrderID"`
}
