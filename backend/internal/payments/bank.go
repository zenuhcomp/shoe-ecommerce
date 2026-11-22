package payments

import (
	"context"
	"fmt"
	"strings"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

// BankTransferStrategy handles Direct Bank Transfer, RTGS, EFT, and bank integrations like KCB Buni, Stanbic, and NCBA.
type BankTransferStrategy struct{}

func NewBankTransferStrategy() *BankTransferStrategy {
	return &BankTransferStrategy{}
}

func (s *BankTransferStrategy) ID() string {
	return "bank_transfer"
}

func (s *BankTransferStrategy) Name() string {
	return "Bank Wire / RTGS / KCB Buni"
}

func (s *BankTransferStrategy) Type() string {
	return "bank"
}

func (s *BankTransferStrategy) Description() string {
	return "Direct bank deposit, RTGS wire, or corporate bank transfer (e.g. KCB Buni, Stanbic, NCBA) with instant invoice generation."
}

func (s *BankTransferStrategy) IsEnabled(cfg models.PaymentConfig) bool {
	return cfg.BankEnabled
}

func (s *BankTransferStrategy) GetPublicInfo(cfg models.PaymentConfig) map[string]any {
	return map[string]any{
		"bankName":        cfg.BankName,
		"accountName":     cfg.BankAccountName,
		"accountNumber":   cfg.BankAccountNo,
		"branch":          cfg.BankBranch,
		"swiftCode":       cfg.BankSwiftCode,
		"instructions":    cfg.BankInstructions,
	}
}

func (s *BankTransferStrategy) Process(ctx context.Context, req models.PaymentInitiateRequest, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error) {
	bank := cfg.BankName
	if bank == "" {
		bank = "KCB Bank Kenya / Commercial Bank"
	}
	acc := cfg.BankAccountNo
	if acc == "" {
		acc = "1289004455"
	}
	accName := cfg.BankAccountName
	if accName == "" {
		accName = "Luna Collections & Co."
	}

	msg := fmt.Sprintf("Please transfer KES %d to %s (Account: %s, Name: %s). Use Reference: %s.",
		req.AmountKES, bank, acc, accName, req.OrderID)
	if cfg.BankInstructions != "" {
		msg = cfg.BankInstructions
	}

	return &models.PaymentInitiateResponse{
		Success:        true,
		Status:         "pending",
		TransactionID:  fmt.Sprintf("bank_%s", req.OrderID),
		PromptRequired: false,
		Message:        msg,
		Details: map[string]string{
			"bankName":      bank,
			"accountNumber": acc,
			"accountName":   accName,
			"branch":        cfg.BankBranch,
			"swiftCode":     cfg.BankSwiftCode,
			"reference":     req.OrderID,
			"amountKes":     fmt.Sprintf("%d", req.AmountKES),
		},
	}, nil
}

func (s *BankTransferStrategy) Verify(ctx context.Context, reference string, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error) {
	ref := strings.TrimSpace(reference)
	if ref == "" {
		return nil, fmt.Errorf("bank transaction reference / slip number is required")
	}
	return &models.PaymentInitiateResponse{
		Success:       true,
		Status:        "completed",
		TransactionID: ref,
		Message:       fmt.Sprintf("Bank transfer %s logged for clearing confirmation.", ref),
	}, nil
}
