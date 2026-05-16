package usecase

import (
	"context"
	"mins_EduCenter/pkg/errors"
	"strings"
)

type EstimatorUsecase struct {
}

func NewEstimatorUsecase() *EstimatorUsecase {
	return &EstimatorUsecase{}
}

func (e *EstimatorUsecase) EstimateCoursePrice(ctx context.Context, courseName string, hours int) (float64, error) {
	const op = "EstimatorUsecase.EstimateCoursePrice"

	if hours <= 0 {
		return 0, errors.NewValidationError(op, "hours", "must be positive")
	}

	baseRatePerHour := 500.0
	taxPercent := 20.0
	bulkDiscountThreshold := 40
	bulkDiscountPercent := 10.0
	advancedBonusPercent := 30.0

	basePrice := float64(hours) * baseRatePerHour

	if hours > bulkDiscountThreshold {
		basePrice = basePrice * (1 - bulkDiscountPercent/100)
	}

	if strings.Contains(strings.ToLower(courseName), "advanced") {
		basePrice = basePrice * (1 + advancedBonusPercent/100)
	}

	totalPrice := basePrice * (1 + taxPercent/100)

	return totalPrice, nil
}
