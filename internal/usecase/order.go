package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/minishop/internal/domain"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	validate *validator.Validate
)

type orderUsecase struct {
	orderRepo domain.OrderRepository
}

func (o *orderUsecase) Create(ctx context.Context, params domain.OrderCreateParameters) (createResp domain.OrderCreateResponse, err error) {
	// Basic validation. UserId is required to create Order
	if params.CreatedBy == 0 {
		return createResp, errors.New("user id not passed")
	}

	// validate domain.OrderCreateParameters (params) using validation package
	if err = validate.Struct(params); err != nil {
		// Try to convert the validation error into domain.ValidationError
		validationErr := convertValidationError(err, params)
		if validationErr != nil {
			return createResp, validationErr
		}

		return createResp, err
	}

	// Fee Calculation
	var (
		baseFee     = 60.0 // City 1 and Weight <= .5 kg
		deliveryFee = 0.0
		codFee      = 0.0
	)

	// BaseFee is 100 of city other than 1
	if params.RecipientCity != 1 {
		baseFee = 100.0
	}

	if params.ItemWeight <= 0.5 {
		deliveryFee = baseFee
	} else if params.ItemWeight <= 1.0 {
		deliveryFee = baseFee + 10
	} else {
		deliveryFee = baseFee + 15*math.Ceil(params.ItemWeight-1.0)
	}

	// Cod fee is 1% of the Amount to Collect
	codFee = float64(params.AmountToCollect) * 0.01

	order := domain.Order{
		OrderConsignmentId: uuid.New().String(),
		OrderCreatedAt:     time.Now(),
		OrderDescription:   params.ItemDescription,
		MerchantOrderId:    params.MerchantOrderId,
		RecipientName:      params.RecipientName,
		RecipientAddress:   params.RecipientAddress,
		RecipientPhone:     params.RecipientPhone,
		OrderAmount:        params.AmountToCollect,
		TotalFee:           deliveryFee + codFee,
		Instruction:        params.SpecialInstruction,
		OrderTypeId:        1,
		CodFee:             codFee,
		PromoDiscount:      0,
		Discount:           0,
		DeliveryFee:        deliveryFee,
		OrderStatus:        domain.OrderStatusPending,
		OrderType:          "Delivery",
		ItemType:           "Parcel",
		TransferStatus:     1,
		Archive:            0,
		CreatedBy:          params.CreatedBy,
		UpdatedBy:          params.CreatedBy,
	}

	return o.orderRepo.Create(ctx, order)
}

func (o *orderUsecase) Cancel(ctx context.Context, consignmentId string, userID uint64) error {
	return o.orderRepo.Cancel(ctx, consignmentId, userID)
}

func (o *orderUsecase) List(ctx context.Context, parameters domain.OrderListParameters) (domain.OrderListResponse, error) {
	return o.orderRepo.List(ctx, parameters)
}

func NewOrderUsecase(orderRepo domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo}
}

// convertValidationError converts the validator's validation error into domain.ValidationError
func convertValidationError(err error, obj interface{}) *domain.ValidationError {
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := domain.ValidationError{ErrorMap: make(map[string][]string)}

		for _, fieldError := range validationErrors {
			// Get the field's JSON tag
			jsonTag := getJSONTag(obj, fieldError.StructField())

			var readable = func(j string) string {
				return strings.Join(strings.Split(j, "_"), " ")
			}

			switch fieldError.ActualTag() {
			case "required":
				errors.ErrorMap[jsonTag] = []string{fmt.Sprintf("The %s field is required.", readable(jsonTag))}
			case "eq":
				errors.ErrorMap[jsonTag] = []string{fmt.Sprintf("Wrong %s selected.", readable(jsonTag))}
			default:
				errors.ErrorMap[jsonTag] = []string{fmt.Sprintf("Field %s failed validation.", readable(jsonTag))}
			}
		}

		return &errors
	}

	return nil
}

func init() {
	validate = validator.New()

	// Register the bd phone number validation
	_ = validate.RegisterValidation("regex_bd_phone", func(fl validator.FieldLevel) bool {
		regex := `^(01)[3-9]{1}[0-9]{8}$` // Example: alphanumeric usernames with underscores, 3-16 chars
		matched, _ := regexp.MatchString(regex, fl.Field().String())
		return matched
	})
}

// getJSONTag retrieves the JSON tag for a struct field
func getJSONTag(obj interface{}, fieldName string) string {
	field, found := reflect.TypeOf(obj).FieldByName(fieldName)
	if !found {
		return fieldName // Return fieldName if JSON tag not found
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName // Return fieldName if JSON tag is empty
	}
	return jsonTag
}
