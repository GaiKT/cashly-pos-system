package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusActive       ProductStatus = "active"
	ProductStatusInactive     ProductStatus = "inactive"
	ProductStatusDiscontinued ProductStatus = "discontinued"
)

// StockMovementType represents the type of stock movement
type StockMovementType string

const (
	StockMovementIn     StockMovementType = "in"
	StockMovementOut    StockMovementType = "out"
	StockMovementAdjust StockMovementType = "adjust"
)

// Product represents a product in the system
type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	SKU         string         `gorm:"type:varchar(100);unique;not null;index" json:"sku"`
	Barcode     string         `gorm:"type:varchar(255);unique;index" json:"barcode"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT" json:"category"`
	Price       float64        `gorm:"type:decimal(10,2);not null;check:price >= 0" json:"price"`
	Cost        float64        `gorm:"type:decimal(10,2);not null;check:cost >= 0" json:"cost"`
	Stock       int            `gorm:"not null;default:0;check:stock >= 0" json:"stock"`
	MinStock    int            `gorm:"not null;default:0;check:min_stock >= 0" json:"min_stock"`
	MaxStock    int            `gorm:"not null;default:0;check:max_stock >= min_stock" json:"max_stock"`
	Status      ProductStatus  `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active','inactive','discontinued')" json:"status"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"image_url"`
	Weight      float64        `gorm:"type:decimal(8,3);check:weight >= 0" json:"weight"`
	Dimensions  string         `gorm:"type:varchar(100)" json:"dimensions"`
	Supplier    string         `gorm:"type:varchar(255)" json:"supplier"`
	Notes       string         `gorm:"type:text" json:"notes"`
	IsActive    bool           `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Associations
	StockMovements []StockMovement `gorm:"foreignKey:ProductID" json:"stock_movements,omitempty"`
}

// Category represents a product category
type Category struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null;unique;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	ParentID    *uuid.UUID     `gorm:"type:uuid;index" json:"parent_id"`
	Parent      *Category      `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"parent,omitempty"`
	IsActive    bool           `gorm:"not null;default:true;index" json:"is_active"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Associations
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Products []Product  `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// StockMovement represents stock movements for products
type StockMovement struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID   uuid.UUID         `gorm:"type:uuid;not null;index" json:"product_id"`
	Product     Product           `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product"`
	Type        StockMovementType `gorm:"type:varchar(10);not null;check:type IN ('in','out','adjust')" json:"type"`
	Quantity    int               `gorm:"not null" json:"quantity"`
	Reason      string            `gorm:"type:varchar(500);not null" json:"reason"`
	Reference   string            `gorm:"type:varchar(255)" json:"reference"`
	Notes       string            `gorm:"type:text" json:"notes"`
	PerformedBy uuid.UUID         `gorm:"type:uuid;not null;index" json:"performed_by"`
	CreatedAt   time.Time         `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName returns the table name for Product model
func (Product) TableName() string {
	return "products"
}

// TableName returns the table name for Category model
func (Category) TableName() string {
	return "categories"
}

// TableName returns the table name for StockMovement model
func (StockMovement) TableName() string {
	return "stock_movements"
}

// BeforeCreate hook for Product
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for Category
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for StockMovement
func (sm *StockMovement) BeforeCreate(tx *gorm.DB) error {
	if sm.ID == uuid.Nil {
		sm.ID = uuid.New()
	}
	return nil
}

// ProductWithRelations represents a product with its relations loaded
type ProductWithRelations struct {
	Product
	CategoryName string  `json:"category_name"`
	TotalSold    int     `json:"total_sold"`
	Revenue      float64 `json:"revenue"`
}

// CreateProductRequest represents the request to create a new product
type CreateProductRequest struct {
	Name        string        `json:"name" binding:"required,min=1,max=255"`
	Description string        `json:"description"`
	SKU         string        `json:"sku" binding:"required,min=1,max=100"`
	Barcode     string        `json:"barcode"`
	CategoryID  uuid.UUID     `json:"category_id" binding:"required"`
	Price       float64       `json:"price" binding:"required,min=0"`
	Cost        float64       `json:"cost" binding:"required,min=0"`
	Stock       int           `json:"stock" binding:"min=0"`
	MinStock    int           `json:"min_stock" binding:"min=0"`
	MaxStock    int           `json:"max_stock" binding:"min=0"`
	Status      ProductStatus `json:"status"`
	ImageURL    string        `json:"image_url"`
	Weight      float64       `json:"weight" binding:"min=0"`
	Dimensions  string        `json:"dimensions"`
	Supplier    string        `json:"supplier"`
	Notes       string        `json:"notes"`
	IsActive    bool          `json:"is_active"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        *string        `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string        `json:"description,omitempty"`
	SKU         *string        `json:"sku,omitempty" binding:"omitempty,min=1,max=100"`
	Barcode     *string        `json:"barcode,omitempty"`
	CategoryID  *uuid.UUID     `json:"category_id,omitempty"`
	Price       *float64       `json:"price,omitempty" binding:"omitempty,min=0"`
	Cost        *float64       `json:"cost,omitempty" binding:"omitempty,min=0"`
	Stock       *int           `json:"stock,omitempty" binding:"omitempty,min=0"`
	MinStock    *int           `json:"min_stock,omitempty" binding:"omitempty,min=0"`
	MaxStock    *int           `json:"max_stock,omitempty" binding:"omitempty,min=0"`
	Status      *ProductStatus `json:"status,omitempty"`
	ImageURL    *string        `json:"image_url,omitempty"`
	Weight      *float64       `json:"weight,omitempty" binding:"omitempty,min=0"`
	Dimensions  *string        `json:"dimensions,omitempty"`
	Supplier    *string        `json:"supplier,omitempty"`
	Notes       *string        `json:"notes,omitempty"`
	IsActive    *bool          `json:"is_active,omitempty"`
}

// ProductFilters represents filters for product queries
type ProductFilters struct {
	CategoryID *uuid.UUID     `json:"category_id,omitempty"`
	Status     *ProductStatus `json:"status,omitempty"`
	IsActive   *bool          `json:"is_active,omitempty"`
	MinPrice   *float64       `json:"min_price,omitempty"`
	MaxPrice   *float64       `json:"max_price,omitempty"`
	LowStock   *bool          `json:"low_stock,omitempty"`
	SearchTerm string         `json:"search_term,omitempty"`
	Supplier   string         `json:"supplier,omitempty"`
}

// BulkStockUpdate represents a bulk stock update request
type BulkStockUpdate struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
	Reason    string    `json:"reason" binding:"required"`
}

// CreateCategoryRequest represents the request to create a new category
type CreateCategoryRequest struct {
	Name        string     `json:"name" binding:"required,min=1,max=255"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	IsActive    bool       `json:"is_active"`
	SortOrder   int        `json:"sort_order"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        *string    `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string    `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
	SortOrder   *int       `json:"sort_order,omitempty"`
}

// CategoryWithProductCount represents a category with product count
type CategoryWithProductCount struct {
	Category
	ProductCount int `json:"product_count"`
}

// CategoryWithProducts represents a category with its products loaded
type CategoryWithProducts struct {
	Category
	Products []Product `json:"products"`
}

// StockAdjustmentRequest represents a stock adjustment request
type StockAdjustmentRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
	Reason    string    `json:"reason" binding:"required,min=1,max=500"`
	Reference string    `json:"reference"`
	Notes     string    `json:"notes"`
}

// ProductSummary represents a summary of product statistics
type ProductSummary struct {
	TotalProducts    int     `json:"total_products"`
	ActiveProducts   int     `json:"active_products"`
	InactiveProducts int     `json:"inactive_products"`
	LowStockProducts int     `json:"low_stock_products"`
	TotalValue       float64 `json:"total_value"`
	TotalCost        float64 `json:"total_cost"`
}
