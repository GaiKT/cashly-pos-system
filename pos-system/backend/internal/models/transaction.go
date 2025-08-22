package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
	TransactionStatusRefunded  TransactionStatus = "REFUNDED"
)

// PaymentMethod represents the payment method used
type PaymentMethod string

const (
	PaymentMethodCash         PaymentMethod = "CASH"
	PaymentMethodCard         PaymentMethod = "CARD"
	PaymentMethodDigital      PaymentMethod = "DIGITAL"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentMethodCredit       PaymentMethod = "CREDIT"
)

// Transaction represents a POS transaction
type Transaction struct {
	ID             uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ReceiptID      string            `json:"receiptId" gorm:"uniqueIndex;not null"`
	CashierID      uuid.UUID         `json:"cashierId" gorm:"type:uuid;not null;index"`
	CustomerName   *string           `json:"customerName,omitempty"`
	CustomerEmail  *string           `json:"customerEmail,omitempty"`
	CustomerPhone  *string           `json:"customerPhone,omitempty"`
	Subtotal       float64           `json:"subtotal" gorm:"not null;check:subtotal >= 0"`
	TaxAmount      float64           `json:"taxAmount" gorm:"not null;default:0;check:tax_amount >= 0"`
	DiscountAmount float64           `json:"discountAmount" gorm:"not null;default:0;check:discount_amount >= 0"`
	Total          float64           `json:"total" gorm:"not null;check:total >= 0"`
	AmountPaid     float64           `json:"amountPaid" gorm:"not null;check:amount_paid >= 0"`
	Change         float64           `json:"change" gorm:"not null;default:0;check:change >= 0"`
	PaymentMethod  PaymentMethod     `json:"paymentMethod" gorm:"type:payment_method;not null"`
	PaymentRef     *string           `json:"paymentRef,omitempty"`
	Status         TransactionStatus `json:"status" gorm:"type:transaction_status;not null;default:'PENDING'"`
	Notes          *string           `json:"notes,omitempty" gorm:"type:text"`
	RefundedAt     *time.Time        `json:"refundedAt,omitempty"`
	RefundedBy     *uuid.UUID        `json:"refundedBy,omitempty" gorm:"type:uuid"`
	RefundReason   *string           `json:"refundReason,omitempty" gorm:"type:text"`
	CreatedAt      time.Time         `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt      time.Time         `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships
	Cashier        User              `json:"cashier,omitempty" gorm:"foreignKey:CashierID"`
	RefundedByUser *User             `json:"refundedByUser,omitempty" gorm:"foreignKey:RefundedBy"`
	Items          []TransactionItem `json:"items,omitempty" gorm:"foreignKey:TransactionID"`
	Payments       []Payment         `json:"payments,omitempty" gorm:"foreignKey:TransactionID"`
}

// TableName specifies the table name for GORM
func (Transaction) TableName() string {
	return "transactions"
}

// TransactionItem represents an item in a transaction
type TransactionItem struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TransactionID uuid.UUID `json:"transactionId" gorm:"type:uuid;not null;index"`
	ProductID     uuid.UUID `json:"productId" gorm:"type:uuid;not null;index"`
	ProductName   string    `json:"productName" gorm:"not null"`
	ProductSKU    string    `json:"productSku" gorm:"not null"`
	Quantity      int       `json:"quantity" gorm:"not null;check:quantity > 0"`
	UnitPrice     float64   `json:"unitPrice" gorm:"not null;check:unit_price >= 0"`
	Discount      float64   `json:"discount" gorm:"not null;default:0;check:discount >= 0"`
	Subtotal      float64   `json:"subtotal" gorm:"not null;check:subtotal >= 0"`
	CreatedAt     time.Time `json:"createdAt" gorm:"not null;default:now()"`

	// Relationships
	Transaction Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
	// Product relationship removed to avoid circular dependency
}

// TableName specifies the table name for GORM
func (TransactionItem) TableName() string {
	return "transaction_items"
}

// Payment represents a payment record
type Payment struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TransactionID uuid.UUID     `json:"transactionId" gorm:"type:uuid;not null;index"`
	Amount        float64       `json:"amount" gorm:"not null;check:amount > 0"`
	Method        PaymentMethod `json:"method" gorm:"type:payment_method;not null"`
	Reference     *string       `json:"reference,omitempty"`
	Status        string        `json:"status" gorm:"not null;default:'COMPLETED'"`
	ProcessedAt   *time.Time    `json:"processedAt,omitempty"`
	CreatedAt     time.Time     `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt     time.Time     `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships
	Transaction Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
}

// TableName specifies the table name for GORM
func (Payment) TableName() string {
	return "payments"
}

// Cart represents a shopping cart (for draft transactions)
type Cart struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CashierID uuid.UUID `json:"cashierId" gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships
	Cashier User       `json:"cashier,omitempty" gorm:"foreignKey:CashierID"`
	Items   []CartItem `json:"items,omitempty" gorm:"foreignKey:CartID"`
}

// TableName specifies the table name for GORM
func (Cart) TableName() string {
	return "carts"
}

// CartItem represents an item in the cart
type CartItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CartID    uuid.UUID `json:"cartId" gorm:"type:uuid;not null;index"`
	ProductID uuid.UUID `json:"productId" gorm:"type:uuid;not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null;check:quantity > 0"`
	Discount  float64   `json:"discount" gorm:"not null;default:0;check:discount >= 0"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships
	Cart Cart `json:"cart,omitempty" gorm:"foreignKey:CartID"`
	// Product relationship removed to avoid circular dependency
}

// TableName specifies the table name for GORM
func (CartItem) TableName() string {
	return "cart_items"
}

// TransactionWithDetails represents a transaction with cashier details
type TransactionWithDetails struct {
	Transaction `gorm:"embedded"`
	Cashier     *User `json:"cashier,omitempty"`
}

// CreateTransactionRequest represents the request to create a new transaction
type CreateTransactionRequest struct {
	CustomerName   *string                 `json:"customerName,omitempty" binding:"omitempty,max=100"`
	CustomerEmail  *string                 `json:"customerEmail,omitempty" binding:"omitempty,email"`
	CustomerPhone  *string                 `json:"customerPhone,omitempty" binding:"omitempty,max=20"`
	Items          []CreateTransactionItem `json:"items" binding:"required,min=1,dive"`
	DiscountAmount *float64                `json:"discountAmount,omitempty" binding:"omitempty,gte=0"`
	PaymentMethod  PaymentMethod           `json:"paymentMethod" binding:"required"`
	AmountPaid     float64                 `json:"amountPaid" binding:"required,gt=0"`
	PaymentRef     *string                 `json:"paymentRef,omitempty" binding:"omitempty,max=100"`
	Notes          *string                 `json:"notes,omitempty" binding:"omitempty,max=500"`
}

// CreateTransactionItem represents an item in the create transaction request
type CreateTransactionItem struct {
	ProductID uuid.UUID `json:"productId" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,gt=0"`
	Discount  *float64  `json:"discount,omitempty" binding:"omitempty,gte=0"`
}

// RefundTransactionRequest represents the request to refund a transaction
type RefundTransactionRequest struct {
	Reason        string                  `json:"reason" binding:"required,min=1,max=500"`
	Items         []RefundTransactionItem `json:"items,omitempty" binding:"omitempty,dive"`
	PartialRefund bool                    `json:"partialRefund" binding:"omitempty"`
}

// RefundTransactionItem represents an item to be refunded
type RefundTransactionItem struct {
	ProductID uuid.UUID `json:"productId" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,gt=0"`
}

// TransactionFilters represents filters for transaction queries
type TransactionFilters struct {
	CashierID     *uuid.UUID         `json:"cashierId,omitempty"`
	UserID        *uuid.UUID         `json:"userId,omitempty"`
	Status        *TransactionStatus `json:"status,omitempty"`
	PaymentMethod *PaymentMethod     `json:"paymentMethod,omitempty"`
	MinTotal      *float64           `json:"minTotal,omitempty"`
	MaxTotal      *float64           `json:"maxTotal,omitempty"`
	MinAmount     *float64           `json:"minAmount,omitempty"`
	MaxAmount     *float64           `json:"maxAmount,omitempty"`
	StartDate     *time.Time         `json:"startDate,omitempty"`
	EndDate       *time.Time         `json:"endDate,omitempty"`
	ReceiptID     *string            `json:"receiptId,omitempty"`
	CustomerEmail *string            `json:"customerEmail,omitempty"`
	CustomerPhone *string            `json:"customerPhone,omitempty"`
}

// DailySales represents daily sales summary
type DailySales struct {
	Date               time.Time          `json:"date"`
	TransactionCount   int                `json:"transactionCount"`
	TotalSales         float64            `json:"totalSales"`
	TotalTax           float64            `json:"totalTax"`
	TotalDiscount      float64            `json:"totalDiscount"`
	AverageTransaction float64            `json:"averageTransaction"`
	PaymentMethods     map[string]float64 `json:"paymentMethods"`
}

// ProductSales represents product sales summary
type ProductSales struct {
	ProductID        uuid.UUID `json:"productId"`
	ProductName      string    `json:"productName"`
	ProductSKU       string    `json:"productSku"`
	TotalQuantity    int       `json:"totalQuantity"`
	TotalRevenue     float64   `json:"totalRevenue"`
	TransactionCount int       `json:"transactionCount"`
}

// CashierPerformance represents cashier performance summary
type CashierPerformance struct {
	CashierID          uuid.UUID `json:"cashierId"`
	CashierName        string    `json:"cashierName"`
	TransactionCount   int       `json:"transactionCount"`
	TotalSales         float64   `json:"totalSales"`
	AverageTransaction float64   `json:"averageTransaction"`
	ItemsSold          int       `json:"itemsSold"`
}

// TransactionSummary represents transaction summary statistics
type TransactionSummary struct {
	TotalTransactions  int32            `json:"totalTransactions"`
	TotalRevenue       float64          `json:"totalRevenue"`
	TotalTax           float64          `json:"totalTax"`
	TotalDiscount      float64          `json:"totalDiscount"`
	AverageTransaction float64          `json:"averageTransaction"`
	PaymentMethods     map[string]int64 `json:"paymentMethods"`
}

// Receipt represents a formatted receipt for printing
type Receipt struct {
	Transaction   Transaction `json:"transaction"`
	CompanyInfo   CompanyInfo `json:"companyInfo"`
	FormattedDate string      `json:"formattedDate"`
	FormattedTime string      `json:"formattedTime"`
	QRCode        *string     `json:"qrCode,omitempty"`
}

// CompanyInfo represents company information for receipts
type CompanyInfo struct {
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Phone   string  `json:"phone"`
	Email   string  `json:"email"`
	Website *string `json:"website,omitempty"`
	TaxID   *string `json:"taxId,omitempty"`
}

// Helper methods for Transaction model

// IsRefundable checks if transaction can be refunded
func (t *Transaction) IsRefundable() bool {
	return t.Status == TransactionStatusCompleted && t.RefundedAt == nil
}

// IsCompleted checks if transaction is completed
func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

// IsPending checks if transaction is pending
func (t *Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// GetItemCount returns total number of items in transaction
func (t *Transaction) GetItemCount() int {
	total := 0
	for _, item := range t.Items {
		total += item.Quantity
	}
	return total
}

// GetTotalProfit calculates total profit (requires product cost prices map)
func (t *Transaction) GetTotalProfit(productCosts map[uuid.UUID]float64) float64 {
	totalProfit := 0.0
	for _, item := range t.Items {
		if costPrice, exists := productCosts[item.ProductID]; exists {
			profit := (item.UnitPrice - costPrice) * float64(item.Quantity)
			totalProfit += profit
		}
	}
	return totalProfit
}

// ValidateTransactionStatus checks if a transaction status string is valid
func ValidateTransactionStatus(status string) bool {
	switch TransactionStatus(status) {
	case TransactionStatusPending, TransactionStatusCompleted, TransactionStatusCancelled, TransactionStatusRefunded:
		return true
	default:
		return false
	}
}

// ValidatePaymentMethod checks if a payment method string is valid
func ValidatePaymentMethod(method string) bool {
	switch PaymentMethod(method) {
	case PaymentMethodCash, PaymentMethodCard, PaymentMethodDigital, PaymentMethodBankTransfer, PaymentMethodCredit:
		return true
	default:
		return false
	}
}

// Helper methods for TransactionItem model

// GetTotalWithDiscount calculates item total including discount
func (ti *TransactionItem) GetTotalWithDiscount() float64 {
	return ti.Subtotal - ti.Discount
}

// GetDiscountPercentage calculates discount percentage
func (ti *TransactionItem) GetDiscountPercentage() float64 {
	if ti.Subtotal == 0 {
		return 0
	}
	return (ti.Discount / ti.Subtotal) * 100
}

// BeforeCreate hook for Transaction model
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for Transaction model
func (t *Transaction) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate hook for TransactionItem model
func (ti *TransactionItem) BeforeCreate(tx *gorm.DB) error {
	if ti.ID == uuid.Nil {
		ti.ID = uuid.New()
	}
	ti.CreatedAt = time.Now()
	return nil
}

// BeforeCreate hook for Payment model
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for Payment model
func (p *Payment) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate hook for Cart model
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for Cart model
func (c *Cart) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate hook for CartItem model
func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	ci.CreatedAt = time.Now()
	ci.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for CartItem model
func (ci *CartItem) BeforeUpdate(tx *gorm.DB) error {
	ci.UpdatedAt = time.Now()
	return nil
}
