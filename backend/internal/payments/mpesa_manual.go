package payments

import (
	"context"
	"fmt"
	"strings"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

// MpesaManualStrategy handles Lipa Na M-Pesa Offline (Buy Goods Till or Business Paybill).
// The customer sends funds using standard SIM toolkit / M-Pesa app and supplies the confirmation code.
type MpesaManualStrategy struct{}

func NewMpesaManualStrategy() *MpesaManualStrategy {
	return &MpesaManualStrategy{}
}

func (s *MpesaManualStrategy) ID() string {
	return "mpesa_manual"
}

func (s *MpesaManualStrategy) Name() string {
	return "M-Pesa Offline (Till & Paybill)"
}

func (s *MpesaManualStrategy) Type() string {
	return "manual"
}

func (s *MpesaManualStrategy) Description() string {
	return "Send payment directly via M-Pesa Buy Goods Till or Business Paybill, then verify with your M-Pesa confirmation code."
}

func (s *MpesaManualStrategy) IsEnabled(cfg models.PaymentConfig) bool {
	return cfg.MpesaManualEnabled && (strings.TrimSpace(cfg.TillNumber) != "" || strings.TrimSpace(cfg.PaybillNumber) != "")
}

func (s *MpesaManualStrategy) GetPublicInfo(cfg models.PaymentConfig) map[string]any {
	return map[string]any{
		"tillNumber":     cfg.TillNumber,
		"paybillNumber":  cfg.PaybillNumber,
		"paybillAccount": cfg.PaybillAccount,
		"instructions":   cfg.MpesaInstructions,
	}
}

func (s *MpesaManualStrategy) Process(ctx context.Context, req models.PaymentInitiateRequest, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error) {
	ref := strings.TrimSpace(req.OrderID)
	channel := "till"
	number := cfg.TillNumber
	if cfg.PaybillNumber != "" && (cfg.TillNumber == "" || req.Metadata["offlineChannel"] == "paybill") {
		channel = "paybill"
		number = cfg.PaybillNumber
	}

	msg := fmt.Sprintf("Please send KES %d via M-Pesa %s (%s). Account: %s. Enter your M-Pesa confirmation code upon completion.",
		req.AmountKES, strings.ToUpper(channel), number, ref)

	return &models.PaymentInitiateResponse{
		Success:        true,
		Status:         "pending",
		TransactionID:  ref,
		PromptRequired: false,
		Message:        msg,
		Details: map[string]string{
			"channel":        channel,
			"number":         number,
			"paybillAccount": cfg.PaybillAccount,
			"orderId":        req.OrderID,
		},
	}, nil
}

func (s *MpesaManualStrategy) Verify(ctx context.Context, reference string, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error) {
	ref := strings.TrimSpace(strings.ToUpper(reference))
	if len(ref) < 8 {
		return nil, fmt.Errorf("invalid M-Pesa transaction reference code (minimum 8 characters required)")
	}
	return &models.PaymentInitiateResponse{
		Success:       true,
		Status:        "completed",
		TransactionID: ref,
		Message:       fmt.Sprintf("M-Pesa transaction %s submitted for verification.", ref),
	}, nil
}
