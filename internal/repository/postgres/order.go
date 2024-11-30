package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/minishop/internal/domain"
	"gorm.io/gorm"
	"time"
)

type OrderRepository struct {
	db *gorm.DB
}

func (o *OrderRepository) Create(ctx context.Context, params domain.OrderCreateParameters) (domain.Order, error) {
	repoOrder := order{
		OrderConsignmentId: uuid.New().String(),
		OrderCreatedAt:     time.Now(),
		OrderDescription:   params.ItemDescription,
		MerchantOrderId:    params.MerchantOrderId,
		RecipientName:      params.RecipientName,
		RecipientAddress:   params.RecipientAddress,
		RecipientPhone:     params.RecipientPhone,
		OrderAmount:        1200,
		TotalFee:           72,
		Instruction:        "",
		OrderTypeId:        0,
		CodFee:             0,
		PromoDiscount:      0,
		Discount:           0,
		DeliveryFee:        0,
		OrderStatus:        "",
		OrderType:          "",
		ItemType:           "",
		TransferStatus:     0,
		Archive:            0,
		UpdatedAt:          time.Time{},
		CreatedBy:          0,
		UpdatedBy:          0,
	}

	result := o.db.WithContext(ctx).Create(&repoOrder)
	if result.Error != nil {
		return domain.Order{}, result.Error
	}

	return domain.Order{
		OrderConsignmentId: repoOrder.OrderConsignmentId,
		OrderCreatedAt:     repoOrder.OrderCreatedAt,
		OrderDescription:   repoOrder.OrderDescription,
		MerchantOrderId:    repoOrder.MerchantOrderId,
		RecipientName:      repoOrder.RecipientName,
		RecipientAddress:   repoOrder.RecipientAddress,
		RecipientPhone:     repoOrder.RecipientPhone,
		OrderAmount:        0,
		TotalFee:           0,
		Instruction:        "",
		OrderTypeId:        0,
		CodFee:             0,
		PromoDiscount:      0,
		Discount:           0,
		DeliveryFee:        0,
		OrderStatus:        "",
		OrderType:          "",
		ItemType:           "",
		TransferStatus:     0,
		Archive:            0,
		UpdatedAt:          time.Time{},
		CreatedBy:          0,
		UpdatedBy:          0,
	}, nil
}

func (o *OrderRepository) Cancel(ctx context.Context, consignmentId string, userID uint64) error {
	result := o.db.WithContext(ctx).Model(&user{}).Where("consignment_id = ?", consignmentId).Update("order_status", domain.OrderStatusCanceled)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (o *OrderRepository) List(ctx context.Context, parameters domain.OrderListParameters) ([]domain.Order, error) {
	var orders []order
	result := o.db.WithContext(ctx).Limit(parameters.Limit).Offset(parameters.Page).Find(&orders,
		"created_by = ? AND transfer_status = ? AND archive = ?",
		parameters.CreatedBy, parameters.TransferStatus, parameters.Archive,
	)

	if result.Error != nil {
		return nil, result.Error
	}

	return []domain.Order{}, nil
}

type order struct {
	ID                 uint64    `json:"id,omitempty" gorm:"primarykey"`
	OrderConsignmentId string    `json:"order_consignment_id" gorm:"uniqueIndex"`
	OrderCreatedAt     time.Time `json:"order_created_at"`
	OrderDescription   string    `json:"order_description"`
	MerchantOrderId    string    `json:"merchant_order_id"`
	RecipientName      string    `json:"recipient_name"`
	RecipientAddress   string    `json:"recipient_address"`
	RecipientPhone     string    `json:"recipient_phone"`
	OrderAmount        int       `json:"order_amount"`
	TotalFee           int       `json:"total_fee"`
	Instruction        string    `json:"instruction"`
	OrderTypeId        int       `json:"order_type_id"`
	CodFee             int       `json:"cod_fee"`
	PromoDiscount      int       `json:"promo_discount"`
	Discount           int       `json:"discount"`
	DeliveryFee        int       `json:"delivery_fee"`
	OrderStatus        string    `json:"order_status"`
	OrderType          string    `json:"order_type"`
	ItemType           string    `json:"item_type"`
	TransferStatus     uint8     `json:"transfer_status"`
	Archive            uint8     `json:"archive"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uint64    `json:"created_by"`
	UpdatedBy uint64    `json:"updated_by"`
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) AutoMigrate() error {
	return o.db.AutoMigrate(&order{})
}
