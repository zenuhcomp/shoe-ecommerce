package models

import "time"

type AdminUser struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	User  AdminUser `json:"user"`
}

type BookingRecord struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email,omitempty"`
	ServiceSlug string    `json:"serviceSlug,omitempty"`
	PreferredAt string    `json:"preferredAt,omitempty"`
	Message     string    `json:"message,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ReorderRequest struct {
	IDs []string `json:"ids"`
}

type ReorderIntRequest struct {
	IDs []int `json:"ids"`
}

type OrderRecord struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email,omitempty"`
	ServiceSlug      string    `json:"serviceSlug"`
	ServiceName      string    `json:"serviceName"`
	Size             string    `json:"size,omitempty"`
	Quantity         int       `json:"quantity"`
	AmountKES        int       `json:"amountKes"`
	PaymentMethod    string    `json:"paymentMethod"`
	MpesaPhone       string    `json:"mpesaPhone,omitempty"`
	OfflineChannel   string    `json:"offlineChannel,omitempty"`
	PaymentReference string    `json:"paymentReference,omitempty"`
	Status           string    `json:"status"`
	Notes            string    `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}
