package payments

import (
	"context"
	"testing"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

func TestPaymentManagerStrategies(t *testing.T) {
	pm := NewPaymentManager()

	cfg := models.PaymentConfig{
		MpesaExpressEnabled: true,
		MpesaManualEnabled:  true,
		CardEnabled:         true,
		CashEnabled:         true,
		BankEnabled:         false, // disabled
		TillNumber:          "808500",
		PaybillNumber:       "247247",
		PaybillAccount:      "TestStore",
		ShortCode:           "174379",
	}

	methods := pm.ListActiveMethods(cfg)
	if len(methods) != 5 {
		t.Fatalf("expected 5 registered methods, got %d", len(methods))
	}

	enabledCount := 0
	for _, m := range methods {
		if m.Enabled {
			enabledCount++
		}
	}
	if enabledCount != 4 {
		t.Fatalf("expected 4 enabled methods, got %d", enabledCount)
	}

	ctx := context.Background()

	// 1. Test Mpesa Express Strategy
	res, err := pm.Process(ctx, models.PaymentInitiateRequest{
		OrderID:       "ord-001",
		CustomerPhone: "0712345678",
		AmountKES:     5000,
		Method:        "mpesa_express",
	}, cfg)
	if err != nil {
		t.Fatalf("mpesa_express process failed: %v", err)
	}
	if !res.Success || !res.PromptRequired {
		t.Fatalf("expected success and prompt required for mpesa_express")
	}

	// 2. Test Mpesa Manual Strategy
	res, err = pm.Process(ctx, models.PaymentInitiateRequest{
		OrderID:   "ord-002",
		AmountKES: 3500,
		Method:    "mpesa_manual",
	}, cfg)
	if err != nil {
		t.Fatalf("mpesa_manual process failed: %v", err)
	}
	if !res.Success {
		t.Fatalf("expected success for mpesa_manual")
	}

	// 3. Test Cash Strategy
	res, err = pm.Process(ctx, models.PaymentInitiateRequest{
		OrderID:   "ord-003",
		AmountKES: 2000,
		Method:    "cash",
	}, cfg)
	if err != nil {
		t.Fatalf("cash process failed: %v", err)
	}
	if !res.Success {
		t.Fatalf("expected success for cash")
	}

	// 4. Test Disabled Strategy (Bank)
	_, err = pm.Process(ctx, models.PaymentInitiateRequest{
		OrderID:   "ord-004",
		AmountKES: 10000,
		Method:    "bank_transfer",
	}, cfg)
	if err == nil {
		t.Fatalf("expected error when processing disabled bank_transfer strategy")
	}
}
