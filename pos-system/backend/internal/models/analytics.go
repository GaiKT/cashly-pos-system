package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExpenseCategory represents expense categories
type ExpenseCategory string

const (
	ExpenseCategorySupplies    ExpenseCategory = "SUPPLIES"
	ExpenseCategoryUtilities   ExpenseCategory = "UTILITIES"
	ExpenseCategoryRent        ExpenseCategory = "RENT"
	ExpenseCategorySalaries    ExpenseCategory = "SALARIES"
	ExpenseCategoryMarketing   ExpenseCategory = "MARKETING"
	ExpenseCategoryMaintenance ExpenseCategory = "MAINTENANCE"
	ExpenseCategoryOther       ExpenseCategory = "OTHER"
)

// Expense represents a business expense
type Expense struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title       string          `json:"title" gorm:"not null"`
	Description *string         `json:"description,omitempty" gorm:"type:text"`
	Amount      float64         `json:"amount" gorm:"not null;check:amount > 0"`
	Category    ExpenseCategory `json:"category" gorm:"type:expense_category;not null"`
	Date        time.Time       `json:"date" gorm:"not null;index"`
	Receipt     *string         `json:"receipt,omitempty"` // File URL
	CreatedBy   uuid.UUID       `json:"createdBy" gorm:"type:uuid;not null;index"`
	ApprovedBy  *uuid.UUID      `json:"approvedBy,omitempty" gorm:"type:uuid"`
	ApprovedAt  *time.Time      `json:"approvedAt,omitempty"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships
	CreatedByUser  User  `json:"createdByUser,omitempty" gorm:"foreignKey:CreatedBy"`
	ApprovedByUser *User `json:"approvedByUser,omitempty" gorm:"foreignKey:ApprovedBy"`
}

// TableName specifies the table name for GORM
func (Expense) TableName() string {
	return "expenses"
}

// StockRecommendationStatus represents the status of a stock recommendation
type StockRecommendationStatus string

const (
	RecommendationStatusPending   StockRecommendationStatus = "PENDING"
	RecommendationStatusAccepted  StockRecommendationStatus = "ACCEPTED"
	RecommendationStatusRejected  StockRecommendationStatus = "REJECTED"
	RecommendationStatusProcessed StockRecommendationStatus = "PROCESSED"
)

// StockRecommendationPriority represents the priority of a recommendation
type StockRecommendationPriority string

const (
	RecommendationPriorityLow    StockRecommendationPriority = "LOW"
	RecommendationPriorityMedium StockRecommendationPriority = "MEDIUM"
	RecommendationPriorityHigh   StockRecommendationPriority = "HIGH"
	RecommendationPriorityUrgent StockRecommendationPriority = "URGENT"
)

// StockRecommendation represents an automated stock recommendation
type StockRecommendation struct {
	ID                  uuid.UUID                   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ProductID           uuid.UUID                   `json:"productId" gorm:"type:uuid;not null;index"`
	ProductName         string                      `json:"productName" gorm:"not null"`
	ProductSKU          string                      `json:"productSku" gorm:"not null"`
	CurrentStock        int                         `json:"currentStock" gorm:"not null"`
	MinStock            int                         `json:"minStock" gorm:"not null"`
	RecommendedQuantity int                         `json:"recommendedQuantity" gorm:"not null;check:recommended_quantity > 0"`
	Priority            StockRecommendationPriority `json:"priority" gorm:"type:recommendation_priority;not null;default:'MEDIUM'"`
	Reason              string                      `json:"reason" gorm:"not null"`
	EstimatedCost       float64                     `json:"estimatedCost" gorm:"not null;check:estimated_cost >= 0"`
	SalesVelocity       float64                     `json:"salesVelocity" gorm:"not null;default:0"` // units per day
	DaysUntilStockout   *int                        `json:"daysUntilStockout,omitempty"`
	Status              StockRecommendationStatus   `json:"status" gorm:"type:recommendation_status;not null;default:'PENDING'"`
	ActionTakenBy       *uuid.UUID                  `json:"actionTakenBy,omitempty" gorm:"type:uuid"`
	ActionTakenAt       *time.Time                  `json:"actionTakenAt,omitempty"`
	Notes               *string                     `json:"notes,omitempty" gorm:"type:text"`
	CreatedAt           time.Time                   `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt           time.Time                   `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships (removed due to circular dependency)
	ActionTakenByUser *User `json:"actionTakenByUser,omitempty" gorm:"foreignKey:ActionTakenBy"`
}

// TableName specifies the table name for GORM
func (StockRecommendation) TableName() string {
	return "stock_recommendations"
}

// AuditLogAction represents different audit log actions
type AuditLogAction string

const (
	AuditActionLogin             AuditLogAction = "LOGIN"
	AuditActionLogout            AuditLogAction = "LOGOUT"
	AuditActionCreateUser        AuditLogAction = "CREATE_USER"
	AuditActionUpdateUser        AuditLogAction = "UPDATE_USER"
	AuditActionDeleteUser        AuditLogAction = "DELETE_USER"
	AuditActionCreateProduct     AuditLogAction = "CREATE_PRODUCT"
	AuditActionUpdateProduct     AuditLogAction = "UPDATE_PRODUCT"
	AuditActionDeleteProduct     AuditLogAction = "DELETE_PRODUCT"
	AuditActionUpdateStock       AuditLogAction = "UPDATE_STOCK"
	AuditActionCreateTransaction AuditLogAction = "CREATE_TRANSACTION"
	AuditActionRefundTransaction AuditLogAction = "REFUND_TRANSACTION"
	AuditActionCreateExpense     AuditLogAction = "CREATE_EXPENSE"
	AuditActionUpdateExpense     AuditLogAction = "UPDATE_EXPENSE"
	AuditActionDeleteExpense     AuditLogAction = "DELETE_EXPENSE"
	AuditActionSystemConfig      AuditLogAction = "SYSTEM_CONFIG"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID     uuid.UUID              `json:"userId" gorm:"type:uuid;not null;index"`
	UserName   string                 `json:"userName" gorm:"not null"`
	UserRole   Role                   `json:"userRole" gorm:"type:user_role;not null"`
	Action     AuditLogAction         `json:"action" gorm:"type:audit_action;not null;index"`
	Resource   string                 `json:"resource" gorm:"not null;index"` // "user", "product", "transaction", etc.
	ResourceID *string                `json:"resourceId,omitempty"`
	OldValues  map[string]interface{} `json:"oldValues,omitempty" gorm:"type:jsonb"`
	NewValues  map[string]interface{} `json:"newValues,omitempty" gorm:"type:jsonb"`
	IPAddress  string                 `json:"ipAddress" gorm:"type:inet;not null"`
	UserAgent  string                 `json:"userAgent" gorm:"type:text;not null"`
	Timestamp  time.Time              `json:"timestamp" gorm:"not null;default:now();index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (AuditLog) TableName() string {
	return "audit_logs"
}

// SystemConfig represents system configuration
type SystemConfig struct {
	ID                          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CompanyName                 string    `json:"companyName" gorm:"not null"`
	CompanyAddress              string    `json:"companyAddress" gorm:"not null"`
	CompanyPhone                string    `json:"companyPhone" gorm:"not null"`
	CompanyEmail                string    `json:"companyEmail" gorm:"not null"`
	CompanyWebsite              *string   `json:"companyWebsite,omitempty"`
	CompanyTaxID                *string   `json:"companyTaxId,omitempty"`
	DefaultCurrency             string    `json:"defaultCurrency" gorm:"not null;default:'USD'"`
	TaxRate                     float64   `json:"taxRate" gorm:"not null;default:0;check:tax_rate >= 0 AND tax_rate <= 1"`
	ReceiptHeader               *string   `json:"receiptHeader,omitempty" gorm:"type:text"`
	ReceiptFooter               *string   `json:"receiptFooter,omitempty" gorm:"type:text"`
	LowStockThreshold           int       `json:"lowStockThreshold" gorm:"not null;default:10;check:low_stock_threshold >= 0"`
	AutoGenerateRecommendations bool      `json:"autoGenerateRecommendations" gorm:"not null;default:true"`
	UpdatedBy                   uuid.UUID `json:"updatedBy" gorm:"type:uuid;not null"`
	UpdatedAt                   time.Time `json:"updatedAt" gorm:"not null;default:now()"`

	// Relationships
	UpdatedByUser User `json:"updatedByUser,omitempty" gorm:"foreignKey:UpdatedBy"`
}

// TableName specifies the table name for GORM
func (SystemConfig) TableName() string {
	return "system_configs"
}

// Request/Response DTOs

// CreateExpenseRequest represents the request to create a new expense
type CreateExpenseRequest struct {
	Title       string          `json:"title" binding:"required,min=1,max=200"`
	Description *string         `json:"description,omitempty" binding:"omitempty,max=1000"`
	Amount      float64         `json:"amount" binding:"required,gt=0"`
	Category    ExpenseCategory `json:"category" binding:"required"`
	Date        time.Time       `json:"date" binding:"required"`
	Receipt     *string         `json:"receipt,omitempty"`
}

// UpdateExpenseRequest represents the request to update an expense
type UpdateExpenseRequest struct {
	Title       *string          `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Description *string          `json:"description,omitempty" binding:"omitempty,max=1000"`
	Amount      *float64         `json:"amount,omitempty" binding:"omitempty,gt=0"`
	Category    *ExpenseCategory `json:"category,omitempty"`
	Date        *time.Time       `json:"date,omitempty"`
	Receipt     *string          `json:"receipt,omitempty"`
}

// ActionRecommendationRequest represents the request to take action on a recommendation
type ActionRecommendationRequest struct {
	Action string  `json:"action" binding:"required,oneof=accept reject process"`
	Notes  *string `json:"notes,omitempty" binding:"omitempty,max=500"`
}

// UpdateSystemConfigRequest represents the request to update system config
type UpdateSystemConfigRequest struct {
	CompanyName                 *string  `json:"companyName,omitempty" binding:"omitempty,min=1,max=200"`
	CompanyAddress              *string  `json:"companyAddress,omitempty" binding:"omitempty,min=1,max=500"`
	CompanyPhone                *string  `json:"companyPhone,omitempty" binding:"omitempty,min=1,max=50"`
	CompanyEmail                *string  `json:"companyEmail,omitempty" binding:"omitempty,email"`
	CompanyWebsite              *string  `json:"companyWebsite,omitempty" binding:"omitempty,url"`
	CompanyTaxID                *string  `json:"companyTaxId,omitempty" binding:"omitempty,max=50"`
	DefaultCurrency             *string  `json:"defaultCurrency,omitempty" binding:"omitempty,len=3"`
	TaxRate                     *float64 `json:"taxRate,omitempty" binding:"omitempty,gte=0,lte=1"`
	ReceiptHeader               *string  `json:"receiptHeader,omitempty" binding:"omitempty,max=500"`
	ReceiptFooter               *string  `json:"receiptFooter,omitempty" binding:"omitempty,max=500"`
	LowStockThreshold           *int     `json:"lowStockThreshold,omitempty" binding:"omitempty,gte=0"`
	AutoGenerateRecommendations *bool    `json:"autoGenerateRecommendations,omitempty"`
}

// Analytics DTOs (using references to other models)

// Dashboard represents dashboard analytics data
type Dashboard struct {
	TodaySales         DashboardSales        `json:"todaySales"`
	WeeklySales        DashboardSales        `json:"weeklySales"`
	MonthlySales       DashboardSales        `json:"monthlySales"`
	TopProducts        []ProductSales        `json:"topProducts"`
	RecentTransactions []TransactionSummary  `json:"recentTransactions"`
	SalesChart         []ChartData           `json:"salesChart"`
	CategoryChart      []ChartData           `json:"categoryChart"`
	PaymentMethodChart []ChartData           `json:"paymentMethodChart"`
	Recommendations    []StockRecommendation `json:"recommendations"`
}

// DashboardSales represents sales summary for dashboard
type DashboardSales struct {
	TotalSales       float64 `json:"totalSales"`
	TransactionCount int     `json:"transactionCount"`
	ItemsSold        int     `json:"itemsSold"`
	AverageOrder     float64 `json:"averageOrder"`
	Growth           float64 `json:"growth"` // Percentage growth from previous period
}

// ChartData represents data for charts
type ChartData struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Count int     `json:"count,omitempty"`
}

// SalesReport represents detailed sales report
type SalesReport struct {
	Period             string               `json:"period"`
	StartDate          time.Time            `json:"startDate"`
	EndDate            time.Time            `json:"endDate"`
	TotalSales         float64              `json:"totalSales"`
	TotalTransactions  int                  `json:"totalTransactions"`
	TotalItems         int                  `json:"totalItems"`
	AverageOrder       float64              `json:"averageOrder"`
	TopProducts        []ProductSales       `json:"topProducts"`
	CategoryBreakdown  []ChartData          `json:"categoryBreakdown"`
	DailySales         []DailySales         `json:"dailySales"`
	CashierPerformance []CashierPerformance `json:"cashierPerformance"`
	PaymentMethods     map[string]float64   `json:"paymentMethods"`
}

// InventoryReport represents inventory status report
type InventoryReport struct {
	TotalProducts      int     `json:"totalProducts"`
	ActiveProducts     int     `json:"activeProducts"`
	LowStockProducts   int     `json:"lowStockProducts"`
	OutOfStockProducts int     `json:"outOfStockProducts"`
	TotalStockValue    float64 `json:"totalStockValue"`
}

// Helper methods

// ValidateExpenseCategory checks if an expense category is valid
func ValidateExpenseCategory(category string) bool {
	switch ExpenseCategory(category) {
	case ExpenseCategorySupplies, ExpenseCategoryUtilities, ExpenseCategoryRent,
		ExpenseCategorySalaries, ExpenseCategoryMarketing, ExpenseCategoryMaintenance,
		ExpenseCategoryOther:
		return true
	default:
		return false
	}
}

// ValidateRecommendationStatus checks if a recommendation status is valid
func ValidateRecommendationStatus(status string) bool {
	switch StockRecommendationStatus(status) {
	case RecommendationStatusPending, RecommendationStatusAccepted,
		RecommendationStatusRejected, RecommendationStatusProcessed:
		return true
	default:
		return false
	}
}

// ValidateRecommendationPriority checks if a recommendation priority is valid
func ValidateRecommendationPriority(priority string) bool {
	switch StockRecommendationPriority(priority) {
	case RecommendationPriorityLow, RecommendationPriorityMedium,
		RecommendationPriorityHigh, RecommendationPriorityUrgent:
		return true
	default:
		return false
	}
}

// Helper methods for StockRecommendation

// IsUrgent checks if recommendation is urgent
func (sr *StockRecommendation) IsUrgent() bool {
	return sr.Priority == RecommendationPriorityUrgent ||
		(sr.DaysUntilStockout != nil && *sr.DaysUntilStockout <= 3)
}

// IsPending checks if recommendation is pending action
func (sr *StockRecommendation) IsPending() bool {
	return sr.Status == RecommendationStatusPending
}

// CanTakeAction checks if action can be taken on recommendation
func (sr *StockRecommendation) CanTakeAction() bool {
	return sr.Status == RecommendationStatusPending
}

// GORM Hooks

// BeforeCreate hook for Expense model
func (e *Expense) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for Expense model
func (e *Expense) BeforeUpdate(tx *gorm.DB) error {
	e.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate hook for StockRecommendation model
func (sr *StockRecommendation) BeforeCreate(tx *gorm.DB) error {
	if sr.ID == uuid.Nil {
		sr.ID = uuid.New()
	}
	sr.CreatedAt = time.Now()
	sr.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for StockRecommendation model
func (sr *StockRecommendation) BeforeUpdate(tx *gorm.DB) error {
	sr.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate hook for AuditLog model
func (al *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	al.Timestamp = time.Now()
	return nil
}

// BeforeCreate hook for SystemConfig model
func (sc *SystemConfig) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	sc.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook for SystemConfig model
func (sc *SystemConfig) BeforeUpdate(tx *gorm.DB) error {
	sc.UpdatedAt = time.Now()
	return nil
}
