package payments

import (
	"context"
	"fmt"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

// CashStrategy handles Cash on Delivery or In-Showroom / In-Studio Payment.
type CashStrategy struct{}

func NewCashStrategy() *CashStrategy {
	return &CashStrategy{}
}

func (s *CashStrategy) ID() string {
	return "cash"
}

func (s *CashStrategy) Name() string {
	return "Cash on Delivery / In-Store Pickup"
}

func (s *CashStrategy) Type() string {
	return "cash"
}

func (s *CashStrategy) Description() string {
	return "Pay in cash upon white-glove delivery or in person at our showroom/sanctuary upon handover."
}

func (s *CashStrategy) IsEnabled(cfg models.PaymentConfig) bool {
	return cfg.CashEnabled
}

func (s *CashStrategy) GetPublicInfo(cfg models.PaymentConfig) map[string]any {
	return map[string]any{
		"instructions": cfg.CashInstructions,
	}
}

func (s *CashStrategy) Process(ctx context.Context, req models.PaymentInitiateRequest, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error) {
	msg := fmt.Sprintf("Order placed for Cash on Delivery / In-Showroom Handover. Total due: KES %d.", req.AmountKES)
	if cfg.CashInstructions != "" {
		msg = cfg.CashInstructions
	}

	return &models.PaymentInitiateResponse{
		Success:        true,
		Status:         "pending",
		TransactionID:  fmt.Sprintf("cash_%s", req.OrderID),
		PromptRequired: false,
		Message:        msg,
		Details: map[string]string{
			"paymentTerms": "Cash on Handover",
			"orderId":      req.OrderID,
			"amountKes":    fmt.Sprintf("%d", req.AmountKES),
		},
	}, nil
}

func (s *CashStrategy) Verify(ctx context.Context, reference string, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error) {
	return &models.PaymentInitiateResponse{
		Success:       true,
		Status:        "completed",
		TransactionID: reference,
		Message:       "Cash payment received and logged.",
	}, nil
}
