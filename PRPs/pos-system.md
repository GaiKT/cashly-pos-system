name: "Point of Sale (POS) System - Full-Stack Web Application"
description: |

## Purpose
Complete POS system implementation for small stores with inventory management, sales tracking, expenses, dashboard analytics, and intelligent stock recommendations using Next.js frontend, Go backend, and MongoDB database.

## Core Principles
1. **Context is King**: Include ALL necessary documentation, examples, and caveats
2. **Validation Loops**: Provide executable tests/lints the AI can run and fix
3. **Information Dense**: Use keywords and patterns from the codebase
4. **Progressive Success**: Start simple, validate, then enhance
5. **Full-Stack Thinking**: Consider both frontend and backend implications
6. **Global rules**: Be sure to follow all rules in copilot-instructions.md

---

## Goal
Build a comprehensive Point of Sale (POS) web application for small stores that enables:
- Secure authentication system with Google, Facebook, and email/password options
- Complete sales transaction processing with receipt generation
- Real-time inventory management with stock level monitoring
- Income and expense tracking with categorization
- Analytics dashboard with sales metrics and trends
- Top-selling products identification and reporting
- Intelligent stock recommendation system for reordering
- User-friendly interface optimized for retail environments
- Role-based access control (admin, manager, cashier)
- Docker-compose development environment setup

## Why
- **Business Value**: Eliminates manual sales tracking and reduces inventory management overhead
- **Cost Effective**: Provides enterprise POS features without expensive licensing
- **Data-Driven**: Enables informed business decisions through analytics and recommendations
- **Scalable**: Modern architecture supports growth from single store to multiple locations
- **Integrated**: Unified system for sales, inventory, and financial tracking

## What
A modern web-based POS system with real-time capabilities for retail operations

### Success Criteria
- [ ] **Authentication**: Multi-provider authentication with Google, Facebook, and email/password
- [ ] **Authorization**: Role-based access control with admin, manager, and cashier roles
- [ ] **Frontend**: Responsive POS interface with transaction processing, inventory views, and analytics dashboard
- [ ] **Backend**: RESTful API with real-time data processing, business logic, and recommendation engine
- [ ] **Database**: Optimized MongoDB schema with proper indexing for fast queries
- [ ] **Integration**: Complete end-to-end transaction flow from sale to inventory update
- [ ] **Analytics**: Dashboard with sales trends, top products, and stock recommendations
- [ ] **Environment**: Docker-compose setup for development and deployment

## All Needed Context

### Documentation & References
```yaml
# MUST READ - Include these in your context window
- url: https://nextjs.org/docs/getting-started
  why: App Router, API routes, SSR patterns for dashboard

- url: https://nextjs.org/docs/app/building-your-application/data-fetching
  why: Data fetching patterns for real-time POS data

- url: https://ui.shadcn.com/docs/installation/next
  why: UI component library setup and patterns

- url: https://ui.shadcn.com/docs/components/data-table
  why: Table components for inventory and sales lists

- url: https://ui.shadcn.com/docs/components/form
  why: Form handling for product entry and transactions

- url: https://golang.org/doc/effective_go
  why: Go best practices for backend services

- url: https://gin-gonic.com/docs/
  why: HTTP router and middleware patterns

- url: https://www.prisma.io/docs/concepts/database-connectors/mongodb
  why: MongoDB schema design and query optimization

- url: https://docs.mongodb.com/manual/indexes/
  why: Index strategies for POS query performance

- url: https://docs.docker.com/compose/
  why: Multi-service development environment setup

- url: https://recharts.org/en-US/guide
  why: Chart components for analytics dashboard

- url: https://next-auth.js.org/getting-started/introduction
  why: Authentication library for Next.js with multiple providers

- url: https://next-auth.js.org/providers/google
  why: Google OAuth provider setup and configuration

- url: https://next-auth.js.org/providers/facebook
  why: Facebook OAuth provider setup and configuration

- url: https://next-auth.js.org/providers/credentials
  why: Email/password authentication provider

- url: https://jwt.io/introduction/
  why: JWT token structure and best practices for API authentication

- url: https://github.com/golang-jwt/jwt
  why: Go JWT library for backend token validation

- url: https://bcrypt-generator.com/
  why: Password hashing best practices and implementation
```

### Current Codebase Structure (New Project)
```bash
# Project will be created with this structure
pos-system/
├── docker-compose.yml          # Development environment
├── frontend/                   # Next.js application
│   ├── app/                   # App Router structure
│   ├── components/            # Reusable UI components
│   ├── lib/                   # Utilities and configurations
│   ├── hooks/                 # Custom React hooks
│   ├── types/                 # TypeScript definitions
│   └── styles/                # Global styles
├── backend/                   # Go application
│   ├── cmd/                   # Application entry points
│   ├── internal/              # Private application code
│   │   ├── handlers/          # HTTP request handlers
│   │   ├── models/            # Data models
│   │   ├── services/          # Business logic layer
│   │   ├── middleware/        # HTTP middleware
│   │   └── utils/             # Utility functions
│   ├── pkg/                   # Public packages
│   └── scripts/               # Build and deployment scripts
├── database/                  # Database related files
│   ├── prisma/                # Prisma schema and migrations
│   └── seed/                  # Initial data scripts
└── docs/                      # Project documentation
```

### Desired Codebase Implementation
```bash
# Files to be created with specific responsibilities
frontend/
├── app/
│   ├── layout.tsx             # Root layout with navigation
│   ├── page.tsx               # Dashboard landing page
│   ├── auth/
│   │   ├── signin/
│   │   │   └── page.tsx       # Sign in page with multiple providers
│   │   ├── signup/
│   │   │   └── page.tsx       # Sign up page for email/password
│   │   └── error/
│   │       └── page.tsx       # Authentication error handling
│   ├── pos/
│   │   └── page.tsx           # POS transaction interface (protected)
│   ├── inventory/
│   │   ├── page.tsx           # Inventory management (manager+)
│   │   └── [id]/page.tsx      # Product details/edit (manager+)
│   ├── sales/
│   │   └── page.tsx           # Sales history and reports (manager+)
│   ├── expenses/
│   │   └── page.tsx           # Expense tracking (admin only)
│   ├── analytics/
│   │   └── page.tsx           # Dashboard with insights (manager+)
│   └── admin/
│       ├── users/
│       │   └── page.tsx       # User management (admin only)
│       └── settings/
│           └── page.tsx       # System settings (admin only)
├── components/
│   ├── auth/
│   │   ├── SignInForm.tsx     # Email/password sign in form
│   │   ├── SignUpForm.tsx     # Email/password registration
│   │   ├── SocialAuth.tsx     # Google/Facebook auth buttons
│   │   ├── ProtectedRoute.tsx # Route protection wrapper
│   │   └── RoleGuard.tsx      # Role-based access control
│   ├── pos/
│   │   ├── TransactionInterface.tsx
│   │   ├── ProductSelector.tsx
│   │   ├── Cart.tsx
│   │   └── PaymentModal.tsx
│   ├── inventory/
│   │   ├── ProductForm.tsx
│   │   ├── StockAlert.tsx
│   │   └── ReorderSuggestions.tsx
│   ├── dashboard/
│   │   ├── SalesChart.tsx
│   │   ├── TopProducts.tsx
│   │   ├── RevenueMetrics.tsx
│   │   └── InventoryStatus.tsx
│   ├── admin/
│   │   ├── UserManagement.tsx # User role management
│   │   └── SystemSettings.tsx # System configuration
│   └── ui/                    # Shadcn components
├── lib/
│   ├── api.ts                 # API client with auth headers
│   ├── auth.ts                # NextAuth configuration
│   ├── utils.ts               # Helper functions
│   └── validations.ts         # Form validation schemas
├── middleware.ts              # Route protection middleware
└── hooks/
    ├── use-auth.ts            # Authentication state management
    ├── use-pos.ts             # POS transaction logic
    ├── use-inventory.ts       # Inventory management
    ├── use-analytics.ts       # Dashboard data
    └── use-recommendations.ts  # Stock recommendations

backend/
├── cmd/
│   └── server/
│       └── main.go            # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── auth.go            # Authentication endpoints
│   │   ├── users.go           # User management endpoints
│   │   ├── pos.go             # POS transaction endpoints
│   │   ├── inventory.go       # Inventory management
│   │   ├── sales.go           # Sales reporting
│   │   ├── expenses.go        # Expense tracking
│   │   └── analytics.go       # Dashboard analytics
│   ├── services/
│   │   ├── auth.go            # Authentication business logic
│   │   ├── users.go           # User management logic
│   │   ├── pos.go             # Transaction business logic
│   │   ├── inventory.go       # Stock management logic
│   │   ├── analytics.go       # Data analysis and insights
│   │   └── recommendations.go # Stock reorder algorithms
│   ├── models/
│   │   ├── user.go            # User and authentication models
│   │   ├── product.go         # Product data structures
│   │   ├── transaction.go     # Sale transaction models
│   │   ├── expense.go         # Expense tracking models
│   │   └── analytics.go       # Analytics data models
│   └── middleware/
│       ├── cors.go            # CORS configuration
│       ├── auth.go            # JWT authentication middleware
│       ├── rbac.go            # Role-based access control
│       └── logger.go          # Request logging
└── pkg/
    ├── database/
    │   └── connection.go      # MongoDB connection
    ├── auth/
    │   ├── jwt.go             # JWT token management
    │   ├── oauth.go           # OAuth provider handlers
    │   └── password.go        # Password hashing utilities
    └── config/
        └── config.go          # Application configuration

database/
├── prisma/
│   └── schema.prisma          # Complete data model
└── seed/
    └── initial_data.go        # Sample products and categories
```

### Known Gotchas & Library Quirks
```typescript
// FRONTEND CRITICAL PATTERNS
// Next.js App Router: Use 'use client' for interactive components
// NextAuth: Configure providers in [...nextauth]/route.ts
// Shadcn: Components need proper className merging with cn() utility
// Real-time updates: Consider WebSocket for live inventory updates
// Form handling: Use react-hook-form with Zod validation
// State management: Use Zustand for POS cart state management
// Protected routes: Use middleware.ts for route protection
// Role-based UI: Conditionally render based on user roles

// BACKEND CRITICAL PATTERNS  
// Go: Use context.Context for request scoping and cancellation
// JWT: Validate tokens in middleware before protected endpoints
// OAuth: Store provider tokens securely for API access
// MongoDB: ObjectId fields need proper bson tags for serialization
// Gin: Middleware order matters - CORS before authentication before RBAC
// Prisma: Use proper MongoDB query optimization for large datasets
// Error handling: Use structured error responses with status codes
// Password security: Use bcrypt with proper salt rounds (12+)

// POS SPECIFIC PATTERNS
// Transactions: Ensure atomicity for inventory updates during sales
// Pricing: Use decimal.Decimal for precise monetary calculations
// Stock: Implement concurrent-safe stock level modifications
// Analytics: Pre-aggregate data for faster dashboard loading
// Recommendations: Cache algorithm results to avoid recalculation
// Audit logging: Track all user actions for compliance
```

## Implementation Blueprint

### Data Models and Schema Design
```prisma
// database/prisma/schema.prisma
generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "mongodb"
  url      = env("DATABASE_URL")
}

model User {
  id            String        @id @default(auto()) @map("_id") @db.ObjectId
  email         String        @unique
  name          String
  avatar        String?       // Profile picture URL
  role          Role          @default(CASHIER)
  isActive      Boolean       @default(true)
  lastLoginAt   DateTime?
  createdAt     DateTime      @default(now())
  updatedAt     DateTime      @updatedAt
  
  // Authentication methods
  accounts      Account[]     // OAuth accounts (Google, Facebook)
  sessions      Session[]     // Active sessions
  password      Password?     // Email/password auth
  
  // POS relationships
  transactions  Transaction[] // Transactions created by this user
  expenses      Expense[]     // Expenses recorded by this user

  @@map("users")
}

model Account {
  id                String  @id @default(auto()) @map("_id") @db.ObjectId
  userId            String  @db.ObjectId
  type              String  // oauth, email
  provider          String  // google, facebook, credentials
  providerAccountId String
  refresh_token     String?
  access_token      String?
  expires_at        Int?
  token_type        String?
  scope             String?
  id_token          String?
  session_state     String?
  
  user User @relation(fields: [userId], references: [id], onDelete: Cascade)

  @@unique([provider, providerAccountId])
  @@map("accounts")
}

model Session {
  id           String   @id @default(auto()) @map("_id") @db.ObjectId
  sessionToken String   @unique
  userId       String   @db.ObjectId
  expires      DateTime
  user         User     @relation(fields: [userId], references: [id], onDelete: Cascade)

  @@map("sessions")
}

model Password {
  id       String @id @default(auto()) @map("_id") @db.ObjectId
  userId   String @unique @db.ObjectId
  hash     String // bcrypt hashed password
  user     User   @relation(fields: [userId], references: [id], onDelete: Cascade)

  @@map("passwords")
}

enum Role {
  ADMIN     // Full system access, user management, expenses
  MANAGER   // Inventory, sales reports, analytics
  CASHIER   // POS transactions only
}

model Category {
  id          String    @id @default(auto()) @map("_id") @db.ObjectId
  name        String    @unique
  description String?
  products    Product[]
  createdAt   DateTime  @default(now())
  updatedAt   DateTime  @updatedAt

  @@map("categories")
}

model Product {
  id           String            @id @default(auto()) @map("_id") @db.ObjectId
  name         String
  description  String?
  sku          String            @unique
  barcode      String?           @unique
  price        Float
  cost         Float             // Cost price for profit calculation
  stock        Int               @default(0)
  minStock     Int               @default(5) // Minimum stock alert level
  maxStock     Int               @default(100) // Maximum stock capacity
  categoryId   String            @db.ObjectId
  category     Category          @relation(fields: [categoryId], references: [id])
  transactions TransactionItem[]
  isActive     Boolean           @default(true)
  createdAt    DateTime          @default(now())
  updatedAt    DateTime          @updatedAt

  @@map("products")
  @@index([categoryId])
  @@index([sku])
  @@index([stock])
}

model Transaction {
  id          String            @id @default(auto()) @map("_id") @db.ObjectId
  items       TransactionItem[]
  subtotal    Float
  tax         Float             @default(0)
  discount    Float             @default(0)
  total       Float
  paymentMethod String          // cash, card, digital
  status      String            @default("completed") // completed, refunded, cancelled
  customerId  String?           @db.ObjectId
  customer    Customer?         @relation(fields: [customerId], references: [id])
  userId      String            @db.ObjectId // User who processed the transaction
  user        User              @relation(fields: [userId], references: [id])
  receiptId   String            @unique // For receipt lookup
  createdAt   DateTime          @default(now())
  updatedAt   DateTime          @updatedAt

  @@map("transactions")
  @@index([createdAt])
  @@index([receiptId])
  @@index([userId])
}

model TransactionItem {
  id            String      @id @default(auto()) @map("_id") @db.ObjectId
  transactionId String      @db.ObjectId
  transaction   Transaction @relation(fields: [transactionId], references: [id])
  productId     String      @db.ObjectId
  product       Product     @relation(fields: [productId], references: [id])
  quantity      Int
  unitPrice     Float       // Price at time of sale
  totalPrice    Float       // quantity * unitPrice
  createdAt     DateTime    @default(now())

  @@map("transaction_items")
  @@index([transactionId])
  @@index([productId])
}

model Customer {
  id           String        @id @default(auto()) @map("_id") @db.ObjectId
  name         String
  email        String?       @unique
  phone        String?
  address      String?
  transactions Transaction[]
  createdAt    DateTime      @default(now())
  updatedAt    DateTime      @updatedAt

  @@map("customers")
}

model Expense {
  id          String   @id @default(auto()) @map("_id") @db.ObjectId
  description String
  amount      Float
  category    String   // rent, utilities, supplies, etc.
  vendor      String?
  date        DateTime @default(now())
  isRecurring Boolean  @default(false)
  receiptUrl  String?  // For receipt image storage
  userId      String   @db.ObjectId // User who recorded the expense
  user        User     @relation(fields: [userId], references: [id])
  createdAt   DateTime @default(now())
  updatedAt   DateTime @updatedAt

  @@map("expenses")
  @@index([date])
  @@index([category])
  @@index([userId])
}

model StockRecommendation {
  id             String   @id @default(auto()) @map("_id") @db.ObjectId
  productId      String   @db.ObjectId
  currentStock   Int
  recommendedQty Int
  reason         String   // "low_stock", "high_demand", "seasonal_trend"
  confidence     Float    // 0.0 to 1.0 confidence score
  isActioned     Boolean  @default(false)
  createdAt      DateTime @default(now())
  updatedAt      DateTime @updatedAt

  @@map("stock_recommendations")
  @@index([productId])
  @@index([createdAt])
}

model AuditLog {
  id        String   @id @default(auto()) @map("_id") @db.ObjectId
  userId    String   @db.ObjectId
  action    String   // "sale", "inventory_update", "user_create", etc.
  resource  String   // "transaction", "product", "user", etc.
  resourceId String? // ID of the affected resource
  details   Json?    // Additional action details
  ipAddress String?
  userAgent String?
  createdAt DateTime @default(now())

  @@map("audit_logs")
  @@index([userId])
  @@index([action])
  @@index([createdAt])
}
```

### API Design
```yaml
# RESTful API Endpoints

# Authentication & Authorization
POST   /api/auth/signin           # Email/password sign in
POST   /api/auth/signup           # Email/password registration
POST   /api/auth/signout          # Sign out (clear session)
GET    /api/auth/session          # Get current user session
POST   /api/auth/refresh          # Refresh JWT token
GET    /api/auth/providers        # Get available OAuth providers
POST   /api/auth/verify-email     # Email verification
POST   /api/auth/reset-password   # Password reset request
POST   /api/auth/change-password  # Change password (authenticated)

# User Management (Admin only)
GET    /api/users                 # List users with roles
POST   /api/users                 # Create new user
GET    /api/users/{id}            # Get user details
PUT    /api/users/{id}            # Update user
PUT    /api/users/{id}/role       # Update user role
PUT    /api/users/{id}/status     # Activate/deactivate user
DELETE /api/users/{id}            # Delete user

# Products & Inventory (Manager+ access)
GET    /api/products              # List products with filtering
POST   /api/products              # Create new product (Manager+)
GET    /api/products/{id}         # Get product details
PUT    /api/products/{id}         # Update product (Manager+)
DELETE /api/products/{id}         # Soft delete product (Manager+)
POST   /api/products/{id}/stock   # Update stock levels (Manager+)
GET    /api/products/low-stock    # Get products below minimum stock

# Categories (Manager+ access)
GET    /api/categories            # List categories
POST   /api/categories            # Create category (Manager+)
PUT    /api/categories/{id}       # Update category (Manager+)
DELETE /api/categories/{id}       # Delete category (Manager+)

# POS Transactions (All authenticated users)
POST   /api/transactions          # Create new sale
GET    /api/transactions          # List transactions with filters
GET    /api/transactions/{id}     # Get transaction details
POST   /api/transactions/{id}/refund # Process refund (Manager+)
GET    /api/transactions/receipt/{receiptId} # Get receipt data

# Analytics & Reporting (Manager+ access)
GET    /api/analytics/dashboard   # Main dashboard metrics
GET    /api/analytics/sales       # Sales trends and charts
GET    /api/analytics/products/top # Top selling products
GET    /api/analytics/revenue     # Revenue breakdowns
GET    /api/analytics/inventory   # Inventory status and alerts
GET    /api/analytics/users       # User performance metrics (Admin)

# Expenses (Admin only)
GET    /api/expenses              # List expenses
POST   /api/expenses              # Create expense record
PUT    /api/expenses/{id}         # Update expense
DELETE /api/expenses/{id}         # Delete expense

# Stock Recommendations (Manager+ access)
GET    /api/recommendations       # Get stock reorder suggestions
POST   /api/recommendations/generate # Generate new recommendations
PUT    /api/recommendations/{id}/action # Mark recommendation as actioned

# Audit & Compliance (Admin only)
GET    /api/audit-logs            # Get audit trail
GET    /api/audit-logs/user/{id}  # Get user-specific audit logs
GET    /api/audit-logs/export     # Export audit logs

# System
GET    /api/health                # Health check endpoint
GET    /api/config                # System configuration (Admin)
PUT    /api/config                # Update system configuration (Admin)
```

### Implementation Tasks (in order)

```yaml
Task 1 - Project Setup & Environment:
CREATE docker-compose.yml:
  - SETUP MongoDB service with persistent volume
  - CONFIGURE Go backend with hot reload
  - CONFIGURE Next.js frontend with development server
  - SETUP environment variables and networking
  - ADD OAuth provider credentials (Google, Facebook)

Task 2 - Database Schema & Connection:
CREATE database/prisma/schema.prisma:
  - DEFINE complete data model with relationships
  - ADD user, authentication, and authorization models
  - SETUP proper indexing for query optimization
  - CONFIGURE MongoDB connection settings
RUN: npx prisma db push

Task 3 - Backend Foundation:
CREATE backend/cmd/server/main.go:
  - SETUP Gin HTTP server with middleware
  - CONFIGURE CORS for frontend communication
  - IMPLEMENT health check endpoint
  - SETUP graceful shutdown

Task 4 - Authentication Infrastructure:
CREATE backend/pkg/auth/:
  - IMPLEMENT JWT token management
  - ADD OAuth provider handlers (Google, Facebook)
  - IMPLEMENT password hashing utilities with bcrypt
  - ADD token validation and refresh logic

Task 5 - Backend Models & Database Layer:
CREATE backend/internal/models/:
  - IMPLEMENT Go structs matching Prisma schema
  - ADD JSON serialization tags
  - INCLUDE validation tags for input validation
  - ADD user and authentication models
CREATE backend/pkg/database/connection.go:
  - SETUP MongoDB connection with Prisma
  - IMPLEMENT connection pooling
  - ADD error handling and retry logic

Task 6 - Authentication Services & Middleware:
CREATE backend/internal/services/auth.go:
  - IMPLEMENT user registration and login
  - ADD OAuth flow handling
  - IMPLEMENT password reset functionality
  - ADD session management
CREATE backend/internal/middleware/auth.go:
  - IMPLEMENT JWT validation middleware
  - ADD role-based access control (RBAC)
  - IMPLEMENT audit logging middleware

Task 7 - User Management Services:
CREATE backend/internal/services/users.go:
  - IMPLEMENT user CRUD operations
  - ADD role management functionality
  - IMPLEMENT user activation/deactivation
  - ADD user profile management

Task 8 - Product & Inventory Services:
CREATE backend/internal/services/inventory.go:
  - IMPLEMENT product CRUD operations
  - ADD stock level management with atomic updates
  - IMPLEMENT low stock detection logic
  - ADD category management functionality
  - INCLUDE user tracking for audit logs

Task 9 - POS Transaction Services:
CREATE backend/internal/services/pos.go:
  - IMPLEMENT sale processing with inventory updates
  - ADD transaction validation and atomicity
  - IMPLEMENT receipt generation logic
  - ADD refund processing capabilities
  - INCLUDE user association for transactions

Task 10 - Analytics & Recommendation Services:
CREATE backend/internal/services/analytics.go:
  - IMPLEMENT sales trend analysis
  - ADD top products calculation
  - IMPLEMENT revenue reporting
  - ADD user performance analytics
CREATE backend/internal/services/recommendations.go:
  - IMPLEMENT stock reorder algorithm
  - ADD demand forecasting logic
  - IMPLEMENT confidence scoring

Task 11 - Backend HTTP Handlers:
CREATE backend/internal/handlers/:
  - IMPLEMENT authentication endpoints
  - ADD user management endpoints
  - IMPLEMENT all POS API endpoints
  - ADD proper error handling and status codes
  - IMPLEMENT request validation
  - ADD response serialization
  - INCLUDE role-based endpoint protection

Task 12 - Frontend Authentication Setup:
CREATE frontend/lib/auth.ts:
  - SETUP NextAuth configuration
  - CONFIGURE Google OAuth provider
  - CONFIGURE Facebook OAuth provider
  - CONFIGURE email/password provider
  - ADD JWT session strategy
CREATE frontend/middleware.ts:
  - IMPLEMENT route protection middleware
  - ADD role-based route access

Task 13 - Frontend Authentication Components:
CREATE frontend/components/auth/:
  - IMPLEMENT sign in form with email/password
  - ADD social authentication buttons
  - IMPLEMENT registration form
  - ADD password reset functionality
  - IMPLEMENT protected route wrapper
  - ADD role-based component guards

Task 14 - Frontend Foundation & Layout:
CREATE frontend/app/layout.tsx:
  - SETUP root layout with navigation
  - CONFIGURE Tailwind CSS and Shadcn
  - IMPLEMENT responsive sidebar navigation
  - ADD global error boundary
  - INCLUDE authentication provider
  - ADD user profile dropdown

Task 15 - Frontend API Client:
CREATE frontend/lib/api.ts:
  - IMPLEMENT typed API client with error handling
  - ADD request/response interceptors
  - CONFIGURE base URL and timeout settings
  - IMPLEMENT retry logic for failed requests
  - ADD automatic JWT token attachment
  - IMPLEMENT token refresh logic

Task 16 - POS Interface Components:
CREATE frontend/components/pos/:
  - IMPLEMENT transaction interface with cart
  - ADD product search and selection
  - IMPLEMENT payment processing modal
  - ADD receipt printing functionality
  - INCLUDE cashier identification in transactions

Task 17 - Inventory Management Interface:
CREATE frontend/components/inventory/:
  - IMPLEMENT product listing with search/filter
  - ADD product form for CRUD operations
  - IMPLEMENT stock level indicators
  - ADD bulk stock update functionality
  - INCLUDE role-based access controls

Task 18 - Admin & User Management Interface:
CREATE frontend/components/admin/:
  - IMPLEMENT user management interface
  - ADD role assignment functionality
  - IMPLEMENT system settings management
  - ADD audit log viewer
  - INCLUDE user performance analytics

Task 19 - Analytics Dashboard:
CREATE frontend/components/dashboard/:
  - IMPLEMENT sales charts and metrics
  - ADD top products display
  - IMPLEMENT inventory status overview
  - ADD expense tracking integration
  - INCLUDE role-based data filtering

Task 20 - Frontend Pages & Routing:
CREATE frontend/app/*/page.tsx:
  - IMPLEMENT all main application pages
  - ADD proper loading and error states
  - IMPLEMENT responsive design for mobile
  - ADD navigation and breadcrumbs
  - INCLUDE authentication redirects
  - ADD role-based page access

Task 21 - Testing & Validation:
CREATE tests for both frontend and backend:
  - IMPLEMENT unit tests for business logic
  - ADD integration tests for API endpoints
  - IMPLEMENT component tests for UI
  - ADD end-to-end transaction flow tests
  - INCLUDE authentication flow testing
  - ADD role-based access testing
```

### Per Task Implementation Details
```typescript
// Task 13 - Authentication Components Example
'use client';

import { signIn, signOut, useSession } from 'next-auth/react';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useState } from 'react';

export function SignInForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleEmailSignIn = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    
    try {
      const result = await signIn('credentials', {
        email,
        password,
        redirect: false,
      });
      
      if (result?.error) {
        toast({
          title: "Sign in failed",
          description: result.error,
          variant: "destructive",
        });
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="w-full max-w-md mx-auto p-6">
      <form onSubmit={handleEmailSignIn} className="space-y-4">
        <div>
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <Label htmlFor="password">Password</Label>
          <Input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? 'Signing in...' : 'Sign In'}
        </Button>
      </form>
      
      <div className="mt-4 space-y-2">
        <Button
          onClick={() => signIn('google')}
          variant="outline"
          className="w-full"
        >
          Sign in with Google
        </Button>
        <Button
          onClick={() => signIn('facebook')}
          variant="outline"
          className="w-full"
        >
          Sign in with Facebook
        </Button>
      </div>
    </Card>
  );
}

// Task 14 - Protected Route Component
export function ProtectedRoute({ 
  children, 
  requiredRole = 'CASHIER' 
}: { 
  children: React.ReactNode;
  requiredRole?: 'ADMIN' | 'MANAGER' | 'CASHIER';
}) {
  const { data: session, status } = useSession();

  if (status === 'loading') {
    return <div>Loading...</div>;
  }

  if (!session) {
    redirect('/auth/signin');
  }

  if (!hasRole(session.user.role, requiredRole)) {
    return <div>Access denied. Insufficient permissions.</div>;
  }

  return <>{children}</>;
}
```

```go
// Task 6 - Authentication Service Example
package services

import (
    "context"
    "errors"
    "time"
    
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "your-project/internal/models"
)

type AuthService struct {
    db        *DatabaseClient
    jwtSecret []byte
}

func NewAuthService(db *DatabaseClient, jwtSecret string) *AuthService {
    return &AuthService{
        db:        db,
        jwtSecret: []byte(jwtSecret),
    }
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
    // Check if user already exists
    existingUser, _ := s.db.GetUserByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, errors.New("user already exists")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
    if err != nil {
        return nil, err
    }
    
    // Create user
    user := &models.User{
        Email:     req.Email,
        Name:      req.Name,
        Role:      models.RoleCashier,
        IsActive:  true,
        CreatedAt: time.Now(),
    }
    
    if err := s.db.CreateUser(ctx, user); err != nil {
        return nil, err
    }
    
    // Create password record
    password := &models.Password{
        UserID: user.ID,
        Hash:   string(hashedPassword),
    }
    
    if err := s.db.CreatePassword(ctx, password); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
    // Get user by email
    user, err := s.db.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    if !user.IsActive {
        return nil, errors.New("account is deactivated")
    }
    
    // Get password hash
    userPassword, err := s.db.GetPasswordByUserID(ctx, user.ID)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(userPassword.Hash), []byte(password)); err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    // Generate JWT token
    token, err := s.generateJWTToken(user)
    if err != nil {
        return nil, err
    }
    
    // Update last login
    user.LastLoginAt = &time.Time{}
    *user.LastLoginAt = time.Now()
    s.db.UpdateUser(ctx, user)
    
    return &models.LoginResponse{
        User:        user,
        AccessToken: token,
        ExpiresIn:   3600, // 1 hour
    }, nil
}

func (s *AuthService) generateJWTToken(user *models.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 1).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}

// Task 6 - JWT Middleware Example
func (s *AuthService) ValidateJWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(401, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }
        
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("unexpected signing method")
            }
            return s.jwtSecret, nil
        })
        
        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(401, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }
        
        // Set user info in context
        c.Set("user_id", claims["user_id"])
        c.Set("user_email", claims["email"])
        c.Set("user_role", claims["role"])
        
        c.Next()
    }
}

// Role-based access control middleware
func RequireRole(minRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("user_role")
        if !exists {
            c.JSON(403, gin.H{"error": "User role not found"})
            c.Abort()
            return
        }
        
        if !hasPermission(userRole.(string), minRole) {
            c.JSON(403, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

func hasPermission(userRole, requiredRole string) bool {
    roleHierarchy := map[string]int{
        "CASHIER": 1,
        "MANAGER": 2,
        "ADMIN":   3,
    }
    
    return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}
```

## Validation Loop

### Level 1: Syntax & Style
```bash
# Frontend validation
cd frontend
npm run lint          # ESLint check
npm run type-check     # TypeScript check
npm run format         # Prettier formatting
npm run build          # Build check

# Backend validation  
cd backend
go fmt ./...           # Go formatting
go vet ./...           # Go static analysis
golangci-lint run      # Extended linting
go mod tidy            # Dependency cleanup
```

### Level 2: Unit Tests
```typescript
// Frontend - Jest + React Testing Library
// CREATE __tests__/components/auth/SignInForm.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { SignInForm } from '@/components/auth/SignInForm';
import { SessionProvider } from 'next-auth/react';

describe('SignInForm', () => {
  it('renders sign in form', () => {
    render(
      <SessionProvider session={null}>
        <SignInForm />
      </SessionProvider>
    );
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
    expect(screen.getByLabelText('Password')).toBeInTheDocument();
  });

  it('handles email/password submission', async () => {
    const mockSignIn = jest.fn();
    jest.mock('next-auth/react', () => ({
      signIn: mockSignIn,
    }));

    render(
      <SessionProvider session={null}>
        <SignInForm />
      </SessionProvider>
    );

    fireEvent.change(screen.getByLabelText('Email'), {
      target: { value: 'test@example.com' }
    });
    fireEvent.change(screen.getByLabelText('Password'), {
      target: { value: 'password123' }
    });
    fireEvent.click(screen.getByText('Sign In'));

    await waitFor(() => {
      expect(mockSignIn).toHaveBeenCalledWith('credentials', {
        email: 'test@example.com',
        password: 'password123',
        redirect: false,
      });
    });
  });

  it('handles OAuth sign in', () => {
    const mockSignIn = jest.fn();
    jest.mock('next-auth/react', () => ({
      signIn: mockSignIn,
    }));

    render(
      <SessionProvider session={null}>
        <SignInForm />
      </SessionProvider>
    );

    fireEvent.click(screen.getByText('Sign in with Google'));
    expect(mockSignIn).toHaveBeenCalledWith('google');
  });
});

// CREATE __tests__/components/pos/POSInterface.test.tsx
describe('POSInterface', () => {
  it('requires authentication', () => {
    render(<POSInterface />);
    // Should redirect to sign in or show auth error
  });

  it('adds products to cart when authenticated', () => {
    const mockSession = {
      user: { id: '1', email: 'cashier@test.com', role: 'CASHIER' }
    };
    
    render(
      <SessionProvider session={mockSession}>
        <POSInterface />
      </SessionProvider>
    );
    
    const addButton = screen.getByText('Add to Cart');
    fireEvent.click(addButton);
    expect(screen.getByText('1 item in cart')).toBeInTheDocument();
  });

  it('processes sale successfully with user tracking', async () => {
    const mockSession = {
      user: { id: '1', email: 'cashier@test.com', role: 'CASHIER' }
    };
    
    render(
      <SessionProvider session={mockSession}>
        <POSInterface />
      </SessionProvider>
    );
    
    // Add items and process sale
    await waitFor(() => {
      expect(screen.getByText('Sale completed')).toBeInTheDocument();
    });
  });
});
```

```go
// Backend - Go testing
// CREATE backend/internal/services/auth_test.go
package services

import (
    "context"
    "testing"
    "your-project/internal/models"
    "github.com/stretchr/testify/assert"
)

func TestAuthService_Register(t *testing.T) {
    ctx := context.Background()
    service := setupTestAuthService(t)
    
    t.Run("successful registration", func(t *testing.T) {
        req := &models.RegisterRequest{
            Email:    "test@example.com",
            Name:     "Test User",
            Password: "password123",
        }
        
        user, err := service.Register(ctx, req)
        
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, req.Email, user.Email)
        assert.Equal(t, req.Name, user.Name)
        assert.Equal(t, models.RoleCashier, user.Role)
    });
    
    t.Run("duplicate email registration", func(t *testing.T) {
        // Create initial user
        req1 := &models.RegisterRequest{
            Email:    "duplicate@test.com",
            Password: "password123",
        }
        service.Register(ctx, req1)
        
        // Try to register with same email
        req2 := &models.RegisterRequest{
            Email:    "duplicate@test.com",
            Password: "differentpass",
        }
        
        _, err := service.Register(ctx, req2)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "user already exists")
    });
}

func TestAuthService_Login(t *testing.T) {
    ctx := context.Background()
    service := setupTestAuthService(t)
    
    // Create test user
    req := &models.RegisterRequest{
        Email:    "login@test.com",
        Password: "password123",
    }
    user, _ := service.Register(ctx, req)
    
    t.Run("successful login", func(t *testing.T) {
        response, err := service.Login(ctx, "login@test.com", "password123")
        
        assert.NoError(t, err)
        assert.NotNil(t, response)
        assert.Equal(t, user.Email, response.User.Email)
        assert.NotEmpty(t, response.AccessToken)
        assert.Equal(t, 3600, response.ExpiresIn)
    });
    
    t.Run("invalid password", func(t *testing.T) {
        _, err := service.Login(ctx, "login@test.com", "wrongpassword")
        
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "invalid credentials")
    });
    
    t.Run("non-existent user", func(t *testing.T) {
        _, err := service.Login(ctx, "nonexistent@test.com", "password123")
        
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "invalid credentials")
    });
    
    t.Run("deactivated user", func(t *testing.T) {
        // Deactivate user
        user.IsActive = false
        service.db.UpdateUser(ctx, user)
        
        _, err := service.Login(ctx, "login@test.com", "password123")
        
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "account is deactivated")
    });
}

// CREATE backend/internal/middleware/auth_test.go
func TestJWTMiddleware(t *testing.T) {
    service := setupTestAuthService(t)
    middleware := service.ValidateJWTMiddleware()
    
    t.Run("valid token", func(t *testing.T) {
        // Create valid token
        user := &models.User{ID: "123", Email: "test@test.com", Role: "CASHIER"}
        token, _ := service.generateJWTToken(user)
        
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("GET", "/", nil)
        c.Request.Header.Set("Authorization", "Bearer "+token)
        
        middleware(c)
        
        assert.Equal(t, http.StatusOK, w.Code)
        assert.Equal(t, "123", c.GetString("user_id"))
        assert.Equal(t, "test@test.com", c.GetString("user_email"))
        assert.Equal(t, "CASHIER", c.GetString("user_role"))
    });
    
    t.Run("missing authorization header", func(t *testing.T) {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("GET", "/", nil)
        
        middleware(c)
        
        assert.Equal(t, http.StatusUnauthorized, w.Code)
    });
    
    t.Run("invalid token", func(t *testing.T) {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("GET", "/", nil)
        c.Request.Header.Set("Authorization", "Bearer invalid-token")
        
        middleware(c)
        
        assert.Equal(t, http.StatusUnauthorized, w.Code)
    });
}

func TestRoleMiddleware(t *testing.T) {
    t.Run("sufficient permissions", func(t *testing.T) {
        middleware := RequireRole("CASHIER")
        
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Set("user_role", "MANAGER") // Manager has sufficient permissions for Cashier
        
        middleware(c)
        
        assert.False(t, c.IsAborted())
    });
    
    t.Run("insufficient permissions", func(t *testing.T) {
        middleware := RequireRole("ADMIN")
        
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Set("user_role", "CASHIER") // Cashier doesn't have Admin permissions
        
        middleware(c)
        
        assert.Equal(t, http.StatusForbidden, w.Code)
        assert.True(t, c.IsAborted())
    });
}
```

```bash
# Run tests with authentication coverage
cd frontend && npm test -- --coverage
cd backend && go test ./... -cover

# Test authentication endpoints specifically
go test ./internal/handlers -run TestAuth
go test ./internal/services -run TestAuth
go test ./internal/middleware -run TestAuth
```

### Level 3: Integration Test
```bash
# Start services with Docker Compose
docker-compose up -d

# Wait for services to be ready
sleep 10

# Test API endpoints
curl -X GET http://localhost:8080/api/health
curl -X GET http://localhost:8080/api/products
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"items":[{"productId":"test","quantity":1}],"paymentMethod":"cash"}'

# Test frontend at http://localhost:3000
# - Navigate to POS interface
# - Add products to cart
# - Process a sale
# - Check inventory updates
# - View analytics dashboard
```

## Final Validation Checklist
- [ ] Database schema applied: `npx prisma db push`
- [ ] Docker services start: `docker-compose up -d`
- [ ] Backend tests pass: `go test ./...`
- [ ] Frontend tests pass: `npm test`
- [ ] No linting errors: `npm run lint && golangci-lint run`
- [ ] Type checking passes: `npm run type-check`
- [ ] All API endpoints respond correctly
- [ ] Authentication flows work (email/password, Google, Facebook)
- [ ] JWT token validation works correctly
- [ ] Role-based access control functions properly
- [ ] Protected routes redirect unauthenticated users
- [ ] Password hashing and verification works
- [ ] OAuth provider integration functions
- [ ] Frontend renders without errors on all pages
- [ ] POS transaction flow works end-to-end
- [ ] Inventory updates reflect sales in real-time
- [ ] Analytics dashboard displays correct data
- [ ] Stock recommendations generate properly
- [ ] User management interface works for admins
- [ ] Audit logging captures all user actions
- [ ] Error states handled gracefully across all interfaces
- [ ] Loading states implemented for all async operations
- [ ] Responsive design works on mobile and tablet
- [ ] Receipt generation and printing functionality works
- [ ] Session management and token refresh work
- [ ] Security headers and CORS configured properly

---

## Anti-Patterns to Avoid
- ❌ Don't process sales without atomic inventory updates
- ❌ Don't use floating point arithmetic for monetary calculations (use decimal)
- ❌ Don't skip input validation on financial data
- ❌ Don't allow negative stock levels without explicit business rules
- ❌ Don't cache sensitive transaction data in browser storage
- ❌ Don't ignore concurrent access to inventory operations
- ❌ Don't hardcode tax rates or business rules in frontend code
- ❌ Don't skip error boundaries for financial components
- ❌ Don't use weak typing for product prices or quantities
- ❌ Don't forget to handle network failures in POS operations
- ❌ Don't implement stock recommendations without proper data analysis
- ❌ Don't skip backup strategies for transaction data
- ❌ **Don't store JWT tokens in localStorage (use httpOnly cookies)**
- ❌ **Don't use weak passwords or skip password complexity requirements**
- ❌ **Don't expose sensitive user data in JWT payload**
- ❌ **Don't skip rate limiting on authentication endpoints**
- ❌ **Don't allow access to protected routes without proper authentication**
- ❌ **Don't trust client-side role validation alone**
- ❌ **Don't skip CSRF protection on state-changing operations**
- ❌ **Don't use predictable user IDs or session tokens**
- ❌ **Don't skip OAuth state parameter validation**
- ❌ **Don't log sensitive information (passwords, tokens)**
- ❌ **Don't skip input sanitization on user data**
- ❌ **Don't allow privilege escalation through API manipulation**

---

## PRP Quality Score: 9/10

**Confidence Level**: High confidence for one-pass implementation with comprehensive context provided.

**Reasoning**:
- Complete database schema designed for POS requirements
- Full API specification with all necessary endpoints
- Detailed component breakdown with error handling
- Business logic patterns specifically for retail operations
- Proper validation gates for both frontend and backend
- Docker environment setup for immediate development
- Real-world POS considerations (atomic transactions, stock management)
- Analytics and recommendation engine architecture planned

**Potential Risk Areas**:
- Stock recommendation algorithm may need refinement based on business requirements
- Performance optimization for high-volume sales environments may require additional work

**Mitigation**:
- Start with simple recommendation logic and iterate
- Implement proper indexing and caching strategies from the beginning
- Use validation loops to catch performance issues early
