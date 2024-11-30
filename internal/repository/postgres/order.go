package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/minishop/internal/domain"
	"gorm.io/gorm"
	"math"
	"time"
)

type OrderRepository struct {
	db *gorm.DB
}

func (o *OrderRepository) Create(ctx context.Context, params domain.Order) (createResp domain.OrderCreateResponse, err error) {
	repoOrder := order{
		OrderConsignmentId: uuid.New().String(),
		OrderCreatedAt:     time.Now(),
		OrderDescription:   params.OrderDescription,
		MerchantOrderId:    params.MerchantOrderId,
		RecipientName:      params.RecipientName,
		RecipientAddress:   params.RecipientAddress,
		RecipientPhone:     params.RecipientPhone,
		OrderAmount:        params.OrderAmount,
		TotalFee:           params.TotalFee,
		Instruction:        params.Instruction,
		OrderTypeId:        params.OrderTypeId,
		CodFee:             params.CodFee,
		PromoDiscount:      params.PromoDiscount,
		Discount:           params.Discount,
		DeliveryFee:        params.DeliveryFee,
		OrderStatus:        string(params.OrderStatus),
		OrderType:          params.OrderType,
		ItemType:           params.ItemType,
		TransferStatus:     params.TransferStatus,
		Archive:            params.Archive,
		UpdatedAt:          time.Now(),
		CreatedBy:          params.CreatedBy,
		UpdatedBy:          params.UpdatedBy,
	}

	result := o.db.WithContext(ctx).Create(&repoOrder)
	if result.Error != nil {
		return createResp, result.Error
	}

	return domain.OrderCreateResponse{
		ConsignmentId:   repoOrder.OrderConsignmentId,
		MerchantOrderId: repoOrder.MerchantOrderId,
		OrderStatus:     repoOrder.OrderStatus,
		DeliveryFee:     repoOrder.DeliveryFee,
	}, err
}

func (o *OrderRepository) Cancel(ctx context.Context, consignmentId string, userID uint64) error {
	result := o.db.WithContext(ctx).Model(&order{}).Where(
		"order_consignment_id = ?", consignmentId).Updates(map[string]interface{}{
		"order_status": domain.OrderStatusCanceled,
		"updated_by":   userID,
		"updated_at":   time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (o *OrderRepository) List(ctx context.Context, parameters domain.OrderListParameters) (domain.OrderListResponse, error) {
	var (
		limit, page, offset int64
		total, totalPages   int64
	)

	limit = parameters.Limit
	if limit < 1 || limit > 100 {
		limit = 100
	}
	page = parameters.Page
	if page < 1 {
		page = 1
	}

	offset = (page - 1) * limit

	query := o.db.WithContext(ctx).Model(&order{})

	if countRes := query.Where("created_by = ? AND transfer_status = ? AND archive = ?",
		parameters.CreatedBy, parameters.TransferStatus, parameters.Archive).Count(&total); countRes.Error != nil {
		return domain.OrderListResponse{}, countRes.Error
	}

	var orders []order
	result := o.db.WithContext(ctx).Limit(int(limit)).Offset(int(offset)).Find(&orders,
		"created_by = ? AND transfer_status = ? AND archive = ?",
		parameters.CreatedBy, parameters.TransferStatus, parameters.Archive,
	)

	if result.Error != nil {
		return domain.OrderListResponse{}, result.Error
	}

	totalPages = int64(math.Ceil(float64(total) / float64(limit)))

	var domainOrders = make([]domain.Order, len(orders))
	for idx, ord := range orders {
		domainOrders[idx] = ord.convertToDomainOrder()
	}

	return domain.OrderListResponse{
		Data:        domainOrders,
		Total:       uint64(total),
		CurrentPage: uint64(page),
		PerPage:     uint64(limit),
		TotalInPage: uint64(len(domainOrders)),
		LastPage:    uint64(totalPages),
	}, nil
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
	TotalFee           float64   `json:"total_fee"`
	Instruction        string    `json:"instruction"`
	OrderTypeId        int       `json:"order_type_id"`
	CodFee             float64   `json:"cod_fee"`
	PromoDiscount      int       `json:"promo_discount"`
	Discount           int       `json:"discount"`
	DeliveryFee        float64   `json:"delivery_fee"`
	OrderStatus        string    `json:"order_status"`
	OrderType          string    `json:"order_type"`
	ItemType           string    `json:"item_type"`
	TransferStatus     uint8     `json:"transfer_status"`
	Archive            uint8     `json:"archive"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uint64    `json:"created_by"`
	UpdatedBy uint64    `json:"updated_by"`
}

func (o order) convertToDomainOrder() domain.Order {
	return domain.Order{
		OrderConsignmentId: o.OrderConsignmentId,
		OrderCreatedAt:     o.OrderCreatedAt,
		OrderDescription:   o.OrderDescription,
		MerchantOrderId:    o.MerchantOrderId,
		RecipientName:      o.RecipientName,
		RecipientAddress:   o.RecipientAddress,
		RecipientPhone:     o.RecipientPhone,
		OrderAmount:        o.OrderAmount,
		TotalFee:           o.TotalFee,
		Instruction:        o.Instruction,
		OrderTypeId:        o.OrderTypeId,
		CodFee:             o.CodFee,
		PromoDiscount:      o.PromoDiscount,
		Discount:           o.Discount,
		DeliveryFee:        o.DeliveryFee,
		OrderStatus:        domain.OrderStatus(o.OrderStatus),
		OrderType:          o.OrderType,
		ItemType:           o.ItemType,
		TransferStatus:     o.TransferStatus,
		Archive:            o.Archive,
		UpdatedAt:          o.UpdatedAt,
		CreatedBy:          o.CreatedBy,
		UpdatedBy:          o.UpdatedBy,
	}
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) AutoMigrate() error {
	return o.db.AutoMigrate(&order{})
}
