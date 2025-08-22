name: "Full-Stack Web App PRP Template v2 - Next.js + Go + MongoDB"
description: |

## Purpose
Template optimized for AI agents to implement full-stack web application features with Next.js frontend, Go backend, and MongoDB database, with sufficient context and self-validation capabilities.

## Core Principles
1. **Context is King**: Include ALL necessary documentation, examples, and caveats
2. **Validation Loops**: Provide executable tests/lints the AI can run and fix
3. **Information Dense**: Use keywords and patterns from the codebase
4. **Progressive Success**: Start simple, validate, then enhance
5. **Full-Stack Thinking**: Consider both frontend and backend implications
6. **Global rules**: Be sure to follow all rules in copilot-instructions.md

---

## Goal
[What needs to be built - be specific about the end state and user experience]

## Why
- [Business value and user impact]
- [Integration with existing features]
- [Problems this solves and for whom]

## What
[User-visible behavior and technical requirements for both frontend and backend]

### Success Criteria
- [ ] [Frontend: Specific UI/UX outcomes]
- [ ] [Backend: API endpoints and data handling]
- [ ] [Integration: End-to-end functionality]

## All Needed Context

### Documentation & References (list all context needed to implement the feature)
```yaml
# MUST READ - Include these in your context window
- url: [Next.js API docs URL]
  why: [Specific patterns for pages/components/API routes]
  
- url: [Go documentation URL]
  why: [HTTP handlers, middleware patterns]

- url: [Prisma MongoDB docs]
  why: [Schema design, query patterns]
  
- file: [frontend/components/example.tsx]
  why: [UI patterns to follow, Shadcn usage]
  
- file: [backend/handlers/example.go]
  why: [Handler patterns, error handling]

- file: [prisma/schema.prisma]
  why: [Current data models and relationships]

- docfile: [PRPs/ai_docs/file.md]
  why: [Project-specific documentation]
```

### Current Codebase Structure
```bash
# Frontend Structure
frontend/
├── components/
├── pages/
├── lib/
├── hooks/
├── types/
└── styles/

# Backend Structure  
backend/
├── handlers/
├── models/
├── services/
├── middleware/
├── utils/
└── main.go

# Database
prisma/
├── schema.prisma
└── migrations/
```

### Desired Codebase Changes
```bash
# Files to be added/modified with responsibilities
frontend/
├── components/[FeatureName]/
│   ├── index.tsx          # Main component
│   ├── [FeatureName].tsx  # Core logic
│   └── types.ts           # TypeScript definitions
├── pages/api/             # Next.js API routes (if needed)
└── hooks/use[Feature].ts  # Custom hooks

backend/
├── handlers/[feature].go  # HTTP handlers
├── models/[feature].go    # Data models
└── services/[feature].go  # Business logic

prisma/
└── schema.prisma          # Updated with new models
```

### Known Gotchas & Library Quirks
```typescript
// FRONTEND CRITICAL PATTERNS
// Next.js: Use getServerSideProps for dynamic data
// Tailwind: Use design system classes consistently
// Shadcn: Import components properly from @/components/ui

// BACKEND CRITICAL PATTERNS  
// Go: Always handle errors explicitly
// Prisma: Use proper MongoDB field types (@db.ObjectId)
// CORS: Configure properly for frontend domain
```

## Implementation Blueprint

### Data Models and Schema Design
```prisma
// Update prisma/schema.prisma
model [FeatureName] {
  id        String   @id @default(auto()) @map("_id") @db.ObjectId
  // Define fields with proper MongoDB types
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
}
```

### API Design
```yaml
# Define REST endpoints
GET    /api/[feature]      # List items
POST   /api/[feature]      # Create item  
GET    /api/[feature]/[id] # Get specific item
PUT    /api/[feature]/[id] # Update item
DELETE /api/[feature]/[id] # Delete item
```

### Implementation Tasks (in order)

```yaml
Task 1 - Database Schema:
MODIFY prisma/schema.prisma:
  - ADD new model definition
  - DEFINE relationships with existing models
  - RUN: npx prisma db push

Task 2 - Backend Models:
CREATE backend/models/[feature].go:
  - DEFINE Go structs matching Prisma schema
  - ADD validation tags
  - INCLUDE JSON serialization tags

Task 3 - Backend Services:
CREATE backend/services/[feature].go:
  - IMPLEMENT business logic
  - USE Prisma client for database operations
  - HANDLE errors with Go idioms

Task 4 - Backend Handlers:
CREATE backend/handlers/[feature].go:
  - IMPLEMENT HTTP handlers
  - USE existing middleware patterns
  - RETURN proper HTTP status codes

Task 5 - Frontend Types:
CREATE frontend/types/[feature].ts:
  - DEFINE TypeScript interfaces
  - MATCH backend model structure
  - EXPORT for component usage

Task 6 - Frontend Components:
CREATE frontend/components/[FeatureName]/:
  - IMPLEMENT UI with Shadcn components
  - USE Tailwind for styling
  - HANDLE loading and error states

Task 7 - Frontend Hooks:
CREATE frontend/hooks/use[Feature].ts:
  - IMPLEMENT data fetching logic
  - USE React Query or SWR for caching
  - HANDLE optimistic updates

Task 8 - Integration:
MODIFY frontend/pages/[page].tsx:
  - INTEGRATE new component
  - HANDLE routing if needed
  - TEST full user flow
```

### Per Task Implementation Details
```typescript
// Task 6 - Frontend Component Example
export function FeatureComponent() {
  const { data, loading, error } = useFeature();
  
  if (loading) return <Skeleton className="w-full h-20" />;
  if (error) return <Alert variant="destructive">Error loading data</Alert>;
  
  return (
    <Card className="p-6">
      {/* Use Shadcn components consistently */}
    </Card>
  );
}
```

```go
// Task 4 - Backend Handler Example
func CreateFeatureHandler(w http.ResponseWriter, r *http.Request) {
    // PATTERN: Always validate input first
    var req CreateFeatureRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // PATTERN: Use service layer for business logic
    result, err := services.CreateFeature(req)
    if err != nil {
        // PATTERN: Proper error handling
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // PATTERN: Consistent JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

## Validation Loop

### Level 1: Syntax & Style
```bash
# Frontend
cd frontend
npm run lint          # ESLint check
npm run type-check     # TypeScript check
npm run format         # Prettier formatting

# Backend  
cd backend
go fmt ./...           # Go formatting
go vet ./...           # Go static analysis
golangci-lint run      # Extended linting
```

### Level 2: Unit Tests
```typescript
// Frontend - Jest + React Testing Library
// CREATE __tests__/[FeatureName].test.tsx
describe('FeatureComponent', () => {
  it('renders loading state', () => {
    render(<FeatureComponent />);
    expect(screen.getByTestId('skeleton')).toBeInTheDocument();
  });

  it('handles error state', () => {
    // Mock error state
    render(<FeatureComponent />);
    expect(screen.getByRole('alert')).toBeInTheDocument();
  });

  it('displays data correctly', () => {
    // Mock successful data
    render(<FeatureComponent />);
    expect(screen.getByText('Expected Content')).toBeInTheDocument();
  });
});
```

```go
// Backend - Go testing
// CREATE handlers/[feature]_test.go
func TestCreateFeatureHandler(t *testing.T) {
    // Test happy path
    // Test validation errors  
    // Test service errors
}
```

```bash
# Run tests
cd frontend && npm test
cd backend && go test ./...
```

### Level 3: Integration Test
```bash
# Start backend
cd backend && go run main.go

# Start frontend  
cd frontend && npm run dev

# Test API endpoints
curl -X POST http://localhost:8080/api/feature \
  -H "Content-Type: application/json" \
  -d '{"name": "test"}'

# Test frontend at http://localhost:3000
```

## Final Validation Checklist
- [ ] Database schema applied: `npx prisma db push`
- [ ] Backend tests pass: `go test ./...`
- [ ] Frontend tests pass: `npm test`
- [ ] No linting errors: `npm run lint && golangci-lint run`
- [ ] Type checking passes: `npm run type-check`
- [ ] API endpoints respond correctly
- [ ] Frontend renders without errors
- [ ] End-to-end user flow works
- [ ] Error states handled gracefully
- [ ] Loading states implemented
- [ ] Responsive design works on mobile

---

## Anti-Patterns to Avoid
- ❌ Don't mix server and client code in Next.js components
- ❌ Don't ignore MongoDB ObjectId requirements in Prisma
- ❌ Don't skip error handling in Go handlers
- ❌ Don't use inline styles instead of Tailwind classes
- ❌ Don't create new UI patterns when Shadcn components exist
- ❌ Don't skip TypeScript strict mode checks
- ❌ Don't hardcode API URLs - use environment variables
- ❌ Don't ignore CORS issues in development