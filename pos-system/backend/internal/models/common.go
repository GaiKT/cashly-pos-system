package models

import (
	"time"

	"github.com/google/uuid"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents detailed error information
type ErrorInfo struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Error      *ErrorInfo  `json:"error,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

// PaginationQuery represents pagination query parameters
type PaginationQuery struct {
	Page   int    `json:"page" form:"page" binding:"omitempty,min=1"`
	Limit  int    `json:"limit" form:"limit" binding:"omitempty,min=1,max=100"`
	Sort   string `json:"sort" form:"sort"`
	Order  string `json:"order" form:"order" binding:"omitempty,oneof=asc desc"`
	Search string `json:"search" form:"search"`
}

// GetPage returns the page number, defaulting to 1
func (pq *PaginationQuery) GetPage() int {
	if pq.Page < 1 {
		return 1
	}
	return pq.Page
}

// GetLimit returns the limit, defaulting to 20
func (pq *PaginationQuery) GetLimit() int {
	if pq.Limit < 1 {
		return 20
	}
	if pq.Limit > 100 {
		return 100
	}
	return pq.Limit
}

// GetSkip returns the number of documents to skip
func (pq *PaginationQuery) GetSkip() int {
	return (pq.GetPage() - 1) * pq.GetLimit()
}

// GetSort returns the sort field, defaulting to "created_at"
func (pq *PaginationQuery) GetSort() string {
	if pq.Sort == "" {
		return "created_at"
	}
	return pq.Sort
}

// GetOrder returns the sort order for SQL (GORM), defaulting to "desc"
func (pq *PaginationQuery) GetOrder() string {
	if pq.Order == "asc" {
		return "asc"
	}
	return "desc" // Default to descending
}

// DateRange represents a date range query
type DateRange struct {
	StartDate *time.Time `json:"startDate" form:"startDate" time_format:"2006-01-02"`
	EndDate   *time.Time `json:"endDate" form:"endDate" time_format:"2006-01-02"`
}

// IsValid checks if the date range is valid
func (dr *DateRange) IsValid() bool {
	if dr.StartDate == nil || dr.EndDate == nil {
		return true // Optional range
	}
	return dr.StartDate.Before(*dr.EndDate) || dr.StartDate.Equal(*dr.EndDate)
}

// GetStartOfDay returns start date at 00:00:00
func (dr *DateRange) GetStartOfDay() *time.Time {
	if dr.StartDate == nil {
		return nil
	}
	start := time.Date(dr.StartDate.Year(), dr.StartDate.Month(), dr.StartDate.Day(),
		0, 0, 0, 0, dr.StartDate.Location())
	return &start
}

// GetEndOfDay returns end date at 23:59:59
func (dr *DateRange) GetEndOfDay() *time.Time {
	if dr.EndDate == nil {
		return nil
	}
	end := time.Date(dr.EndDate.Year(), dr.EndDate.Month(), dr.EndDate.Day(),
		23, 59, 59, 999999999, dr.EndDate.Location())
	return &end
}

// Response helpers

// SuccessResponse creates a successful response
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(message string, code string, details map[string]interface{}) Response {
	return Response{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// PaginatedSuccessResponse creates a successful paginated response
func PaginatedSuccessResponse(message string, data interface{}, pagination Pagination) PaginatedResponse {
	return PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}

// PaginatedErrorResponse creates an error paginated response
func PaginatedErrorResponse(message string, code string, details map[string]interface{}) PaginatedResponse {
	return PaginatedResponse{
		Success: false,
		Message: message,
		Data:    []interface{}{},
		Pagination: Pagination{
			Page:       1,
			Limit:      20,
			Total:      0,
			TotalPages: 0,
			HasNext:    false,
			HasPrev:    false,
		},
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, limit int, total int64) Pagination {
	if limit == 0 {
		limit = 20
	}
	if page == 0 {
		page = 1
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit)) // Ceiling division
	hasNext := page < totalPages
	hasPrev := page > 1

	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}

// IDRequest represents a request with an ID parameter
type IDRequest struct {
	ID uuid.UUID `json:"id" uri:"id" binding:"required"`
}

// HealthCheck represents health check response
type HealthCheck struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	Database  string    `json:"database"`
	DBName    string    `json:"dbName,omitempty"`
	Uptime    string    `json:"uptime,omitempty"`
}

// ValidationError represents field validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// Stats represents general statistics
type Stats struct {
	Users        int64     `json:"users"`
	Products     int64     `json:"products"`
	Categories   int64     `json:"categories"`
	Transactions int64     `json:"transactions"`
	Revenue      float64   `json:"revenue"`
	LastUpdated  time.Time `json:"lastUpdated"`
}

// FileUploadResponse represents file upload response
type FileUploadResponse struct {
	FileName string `json:"fileName"`
	FileURL  string `json:"fileUrl"`
	FileSize int64  `json:"fileSize"`
	MimeType string `json:"mimeType"`
}

// BulkOperation represents a bulk operation request
type BulkOperation struct {
	Operation string      `json:"operation" binding:"required,oneof=create update delete"`
	IDs       []uuid.UUID `json:"ids,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// BulkOperationResult represents the result of a bulk operation
type BulkOperationResult struct {
	Operation      string               `json:"operation"`
	TotalRequested int                  `json:"totalRequested"`
	Successful     int                  `json:"successful"`
	Failed         int                  `json:"failed"`
	Errors         []BulkOperationError `json:"errors,omitempty"`
}

// BulkOperationError represents an error in bulk operation
type BulkOperationError struct {
	Index   int    `json:"index"`
	ID      string `json:"id,omitempty"`
	Message string `json:"message"`
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query   string                 `json:"query" form:"query" binding:"required,min=1"`
	Fields  []string               `json:"fields" form:"fields"`
	Filters map[string]interface{} `json:"filters" form:"filters"`
	Page    int                    `json:"page" form:"page" binding:"omitempty,min=1"`
	Limit   int                    `json:"limit" form:"limit" binding:"omitempty,min=1,max=100"`
}

// SearchResult represents search results
type SearchResult struct {
	Query      string      `json:"query"`
	Results    interface{} `json:"results"`
	Total      int64       `json:"total"`
	Pagination Pagination  `json:"pagination"`
	TimeTaken  string      `json:"timeTaken"`
}

// ExportRequest represents data export request
type ExportRequest struct {
	Format    string                 `json:"format" binding:"required,oneof=csv xlsx json"`
	DateRange *DateRange             `json:"dateRange,omitempty"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	Fields    []string               `json:"fields,omitempty"`
}

// ExportResponse represents data export response
type ExportResponse struct {
	FileName    string    `json:"fileName"`
	DownloadURL string    `json:"downloadUrl"`
	Format      string    `json:"format"`
	RecordCount int       `json:"recordCount"`
	GeneratedAt time.Time `json:"generatedAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// EmailRequest represents email sending request
type EmailRequest struct {
	To          []string `json:"to" binding:"required,min=1"`
	CC          []string `json:"cc,omitempty"`
	BCC         []string `json:"bcc,omitempty"`
	Subject     string   `json:"subject" binding:"required,min=1"`
	Body        string   `json:"body" binding:"required,min=1"`
	IsHTML      bool     `json:"isHtml"`
	Attachments []string `json:"attachments,omitempty"`
}

// NotificationRequest represents notification request
type NotificationRequest struct {
	UserIDs []uuid.UUID            `json:"userIds" binding:"required,min=1"`
	Title   string                 `json:"title" binding:"required,min=1"`
	Message string                 `json:"message" binding:"required,min=1"`
	Type    string                 `json:"type" binding:"required,oneof=info warning error success"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// Constants for common error codes
const (
	ErrorCodeValidation        = "VALIDATION_ERROR"
	ErrorCodeNotFound          = "NOT_FOUND"
	ErrorCodeUnauthorized      = "UNAUTHORIZED"
	ErrorCodeForbidden         = "FORBIDDEN"
	ErrorCodeConflict          = "CONFLICT"
	ErrorCodeInternalError     = "INTERNAL_ERROR"
	ErrorCodeBadRequest        = "BAD_REQUEST"
	ErrorCodeInsufficientStock = "INSUFFICIENT_STOCK"
	ErrorCodePaymentFailed     = "PAYMENT_FAILED"
	ErrorCodeEmailExists       = "EMAIL_EXISTS"
	ErrorCodeInvalidToken      = "INVALID_TOKEN"
	ErrorCodeExpiredToken      = "EXPIRED_TOKEN"
)

// Constants for success messages
const (
	MessageCreatedSuccessfully   = "Created successfully"
	MessageUpdatedSuccessfully   = "Updated successfully"
	MessageDeletedSuccessfully   = "Deleted successfully"
	MessageRetrievedSuccessfully = "Retrieved successfully"
	MessageOperationSuccessful   = "Operation completed successfully"
)
