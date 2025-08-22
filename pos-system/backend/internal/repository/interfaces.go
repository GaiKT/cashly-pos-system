package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/pos-system/backend/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]interface{}, pagination *models.PaginationQuery) ([]models.User, int64, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	SetActiveStatus(ctx context.Context, id uuid.UUID, isActive bool) error
	UpdateRole(ctx context.Context, id uuid.UUID, role models.Role) error
}

// AccountRepository defines the interface for OAuth account operations
type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	GetByProviderAndAccountID(ctx context.Context, provider, providerAccountID string) (*models.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// SessionRepository defines the interface for session operations
type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByToken(ctx context.Context, token string) (*models.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
}

// PasswordRepository defines the interface for password operations
type PasswordRepository interface {
	Create(ctx context.Context, password *models.Password) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Password, error)
	Update(ctx context.Context, password *models.Password) error
	Delete(ctx context.Context, userID uuid.UUID) error
	SetResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	ValidateResetToken(ctx context.Context, token string) (*models.Password, error)
	ClearResetToken(ctx context.Context, userID uuid.UUID) error
}

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetBySKU(ctx context.Context, sku string) (*models.Product, error)
	GetByBarcode(ctx context.Context, barcode string) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters *models.ProductFilters, pagination *models.PaginationQuery) ([]models.Product, int64, error)
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int, reason string, userID uuid.UUID) error
	GetLowStock(ctx context.Context, threshold int) ([]models.Product, error)
	GetOutOfStock(ctx context.Context) ([]models.Product, error)
	BulkUpdateStock(ctx context.Context, updates []models.BulkStockUpdate, userID uuid.UUID) error
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetByName(ctx context.Context, name string) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, pagination *models.PaginationQuery) ([]models.Category, int64, error)
	GetWithProducts(ctx context.Context, id uuid.UUID) (*models.CategoryWithProducts, error)
	GetTree(ctx context.Context) ([]models.Category, error)
}

// TransactionRepository defines the interface for transaction data operations
type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	GetByReceiptID(ctx context.Context, receiptID string) (*models.Transaction, error)
	Update(ctx context.Context, transaction *models.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters *models.TransactionFilters, pagination *models.PaginationQuery) ([]models.Transaction, int64, error)
	GetDailySales(ctx context.Context, date time.Time) (*models.DailySales, error)
	GetSalesReport(ctx context.Context, startDate, endDate time.Time) (*models.SalesReport, error)
	GetTopProducts(ctx context.Context, startDate, endDate time.Time, limit int) ([]models.ProductSales, error)
	GetCashierPerformance(ctx context.Context, startDate, endDate time.Time) ([]models.CashierPerformance, error)
}

// StockMovementRepository defines the interface for stock movement operations
type StockMovementRepository interface {
	Create(ctx context.Context, movement *models.StockMovement) error
	GetByProductID(ctx context.Context, productID uuid.UUID, pagination *models.PaginationQuery) ([]models.StockMovement, int64, error)
	List(ctx context.Context, filters map[string]interface{}, pagination *models.PaginationQuery) ([]models.StockMovement, int64, error)
}

// ExpenseRepository defines the interface for expense operations
type ExpenseRepository interface {
	Create(ctx context.Context, expense *models.Expense) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Expense, error)
	Update(ctx context.Context, expense *models.Expense) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]interface{}, pagination *models.PaginationQuery) ([]models.Expense, int64, error)
	GetByCategory(ctx context.Context, category models.ExpenseCategory, startDate, endDate time.Time) ([]models.Expense, error)
	GetTotalByPeriod(ctx context.Context, startDate, endDate time.Time) (float64, error)
	Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error
}

// StockRecommendationRepository defines the interface for stock recommendation operations
type StockRecommendationRepository interface {
	Create(ctx context.Context, recommendation *models.StockRecommendation) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.StockRecommendation, error)
	Update(ctx context.Context, recommendation *models.StockRecommendation) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]interface{}, pagination *models.PaginationQuery) ([]models.StockRecommendation, int64, error)
	GetPending(ctx context.Context) ([]models.StockRecommendation, error)
	TakeAction(ctx context.Context, id uuid.UUID, action string, notes *string, userID uuid.UUID) error
}

// AuditLogRepository defines the interface for audit log operations
type AuditLogRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	List(ctx context.Context, filters map[string]interface{}, pagination *models.PaginationQuery) ([]models.AuditLog, int64, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, pagination *models.PaginationQuery) ([]models.AuditLog, int64, error)
	GetByResource(ctx context.Context, resource string, resourceID string, pagination *models.PaginationQuery) ([]models.AuditLog, int64, error)
	DeleteOldLogs(ctx context.Context, beforeDate time.Time) error
}

// SystemConfigRepository defines the interface for system configuration operations
type SystemConfigRepository interface {
	Get(ctx context.Context) (*models.SystemConfig, error)
	Update(ctx context.Context, config *models.SystemConfig, updatedBy uuid.UUID) error
	GetCompanyInfo(ctx context.Context) (*models.CompanyInfo, error)
}

// CartRepository defines the interface for shopping cart operations
type CartRepository interface {
	Create(ctx context.Context, cart *models.Cart) error
	GetByCashierID(ctx context.Context, cashierID uuid.UUID) (*models.Cart, error)
	Update(ctx context.Context, cart *models.Cart) error
	Delete(ctx context.Context, id uuid.UUID) error
	Clear(ctx context.Context, cashierID uuid.UUID) error
	AddItem(ctx context.Context, cartID uuid.UUID, item *models.CartItem) error
	UpdateItem(ctx context.Context, cartID uuid.UUID, item *models.CartItem) error
	RemoveItem(ctx context.Context, cartID uuid.UUID, productID uuid.UUID) error
}

// Repositories represents all repository interfaces
type Repositories struct {
	User                UserRepository
	Account             AccountRepository
	Session             SessionRepository
	Password            PasswordRepository
	Product             ProductRepository
	Category            CategoryRepository
	Transaction         TransactionRepository
	StockMovement       StockMovementRepository
	Expense             ExpenseRepository
	StockRecommendation StockRecommendationRepository
	AuditLog            AuditLogRepository
	SystemConfig        SystemConfigRepository
	Cart                CartRepository
	DB                  *gorm.DB
}

// NewRepositories creates new repository instances
// TODO: Implement repository constructors before enabling this function
/*
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:                NewUserRepository(db),
		Account:             NewAccountRepository(db),
		Session:             NewSessionRepository(db),
		Password:            NewPasswordRepository(db),
		Product:             NewProductRepository(db),
		Category:            NewCategoryRepository(db),
		Transaction:         NewTransactionRepository(db),
		StockMovement:       NewStockMovementRepository(db),
		Expense:             NewExpenseRepository(db),
		StockRecommendation: NewStockRecommendationRepository(db),
		AuditLog:            NewAuditLogRepository(db),
		SystemConfig:        NewSystemConfigRepository(db),
		Cart:                NewCartRepository(db),
		DB:                  db,
	}
}
*/
