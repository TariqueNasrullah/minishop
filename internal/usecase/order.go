package usecase

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/minishop/internal/domain"
	"reflect"
	"regexp"
	"strings"
)

var (
	validate *validator.Validate
)

type orderUsecase struct {
	orderRepo domain.OrderRepository
}

func (o *orderUsecase) Create(ctx context.Context, params domain.OrderCreateParameters) (domain.Order, error) {
	if err := validate.Struct(params); err != nil {
		validationErr := convertValidationError(err, params)
		if validationErr != nil {
			return domain.Order{}, validationErr
		}

		return domain.Order{}, err
	}
	return o.orderRepo.Create(ctx, params)
}

func (o *orderUsecase) Cancel(ctx context.Context, consignmentId string, userID uint64) error {
	return o.orderRepo.Cancel(ctx, consignmentId, userID)
}

func (o *orderUsecase) List(ctx context.Context, parameters domain.OrderListParameters) ([]domain.Order, error) {
	return o.orderRepo.List(ctx, parameters)
}

func NewOrderUsecase(orderRepo domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo}
}

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
