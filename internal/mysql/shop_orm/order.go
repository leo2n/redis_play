package shop_orm

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	SelfDefine
	OrderNumber string `gorm:"order_number" json:"orderNumber"`
	Username    string `gorm:"username" json:"username"`
	ProductID   string `gorm:"product_id" json:"productID"`
	PurchaseNum int    `gorm:"purchase_num" json:"purchaseNum"`
	Status      string `gorm:"status" json:"status"`
	Price       int    `gorm:"price" json:"price"`
}

// 这里应该选择严格模式吗? 只选择需要创建的东西
// check params in http request maybe is best :)
func (o *Order) CreateOrder(tx *gorm.DB) error {
	if err := tx.Create(o).Error; err != nil {
		return err
	}
	return nil
}

// find all orders
// just for admin
func (o *Order) QueryOrders() ([]*Order, error) {
	orders := make([]*Order, 128)
	if err := conn.Model(&Order{}).Find(&orders).Error; err != nil {
		return orders, err
	}
	return orders, nil
}

// find specific order, by username
func (o *Order) QueryOrderByUsername(username string) ([]*Order, error) {
	orders := make([]*Order, 128)
	if err := conn.Model(&Order{}).Where("username = ?", username).Find(&orders).Error; err != nil {
		return orders, err
	}
	return orders, nil
}

// find specific order, by

func (o *Order) UpdateOrderStatus(tx *gorm.DB) error {
	if err := tx.Model(&Order{}).Where("order_number = ?", o.OrderNumber).Update("status", o.Status).Error; err != nil {
		return err
	}
	return nil
}

func (o *Order) DeleteOrder(tx *gorm.DB) error {
	if err := tx.Model(&Order{}).Where("order_number = ?", o.OrderNumber).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}
