package payments

import (
	"context"
	"fmt"
	"sync"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

// PaymentStrategy is the core Strategy interface for pluggable payment processing.
// In the future, additional providers like KCB Buni, Airtel Money, IntaSend, or custom gateways
// can be added simply by implementing this interface and registering it with the PaymentManager.
type PaymentStrategy interface {
	ID() string
	Name() string
	Type() string
	Description() string
	IsEnabled(cfg models.PaymentConfig) bool
	GetPublicInfo(cfg models.PaymentConfig) map[string]any
	Process(ctx context.Context, req models.PaymentInitiateRequest, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error)
	Verify(ctx context.Context, reference string, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error)
}

// PaymentManager manages and orchestrates payment strategies.
type PaymentManager struct {
	mu         sync.RWMutex
	strategies map[string]PaymentStrategy
	order      []string
}

// NewPaymentManager initializes a new PaymentManager with registered strategies.
func NewPaymentManager() *PaymentManager {
	pm := &PaymentManager{
		strategies: make(map[string]PaymentStrategy),
		order:      make([]string, 0),
	}
	// Register default standard strategies
	pm.Register(NewMpesaExpressStrategy())
	pm.Register(NewMpesaManualStrategy())
	pm.Register(NewCreditCardStrategy())
	pm.Register(NewCashStrategy())
	pm.Register(NewBankTransferStrategy())
	return pm
}

// Register registers a payment strategy into the manager.
// This allows plugging in new strategies like KcbBuniStrategy at runtime or startup.
func (pm *PaymentManager) Register(s PaymentStrategy) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	id := s.ID()
	if _, exists := pm.strategies[id]; !exists {
		pm.order = append(pm.order, id)
	}
	pm.strategies[id] = s
}

// GetStrategy retrieves a strategy by its ID.
func (pm *PaymentManager) GetStrategy(id string) (PaymentStrategy, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	s, ok := pm.strategies[id]
	if !ok {
		return nil, fmt.Errorf("payment strategy '%s' not found", id)
	}
	return s, nil
}

// ListActiveMethods returns all payment methods that are enabled in the current site configuration.
func (pm *PaymentManager) ListActiveMethods(cfg models.PaymentConfig) []models.PaymentMethodInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var active []models.PaymentMethodInfo
	for _, id := range pm.order {
		s := pm.strategies[id]
		enabled := s.IsEnabled(cfg)
		active = append(active, models.PaymentMethodInfo{
			ID:           s.ID(),
			Name:         s.Name(),
			Type:         s.Type(),
			Enabled:      enabled,
			Description:  s.Description(),
			PublicParams: s.GetPublicInfo(cfg),
		})
	}
	return active
}

// Process routes a payment request to the appropriate strategy.
func (pm *PaymentManager) Process(ctx context.Context, req models.PaymentInitiateRequest, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error) {
	s, err := pm.GetStrategy(req.Method)
	if err != nil {
		return nil, err
	}
	if !s.IsEnabled(cfg) {
		return nil, fmt.Errorf("payment method '%s' is currently disabled by administrator", s.Name())
	}
	return s.Process(ctx, req, cfg)
}

// Verify routes a verification request to the appropriate strategy.
func (pm *PaymentManager) Verify(ctx context.Context, method, reference string, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error) {
	s, err := pm.GetStrategy(method)
	if err != nil {
		return nil, err
	}
	return s.Verify(ctx, reference, cfg)
}
