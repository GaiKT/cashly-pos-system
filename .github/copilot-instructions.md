### 🔄 Project Awareness & Context
- **Always read `PLANNING.md`** at the start of a new conversation to understand the project's architecture, goals, style, and constraints.
- **Check `TASK.md`** before starting a new task. If the task isn't listed, add it with a brief description and today's date.
- **Use consistent naming conventions, file structure, and architecture patterns** as described in `PLANNING.md`.
- **Use proper environment setup** for both frontend (Node.js/npm) and backend (Go modules) development.

### 🧱 Code Structure & Modularity
- **Never create a file longer than 500 lines of code.** If a file approaches this limit, refactor by splitting it into modules or helper files.
- **Organize code into clearly separated modules**, grouped by feature or responsibility.
  - **Frontend structure** (`/frontend` or `/client`):
    - `/components` - Reusable UI components
    - `/pages` - Next.js page components
    - `/lib` - Utility functions and configurations
    - `/hooks` - Custom React hooks
    - `/types` - TypeScript type definitions
  - **Backend structure** (`/backend` or `/server`):
    - `/handlers` - HTTP request handlers
    - `/models` - Data models and Prisma schema
    - `/services` - Business logic layer
    - `/middleware` - HTTP middleware functions
    - `/utils` - Utility functions
- **Use clear, consistent imports** (prefer relative imports within packages).
- **Use environment variables** with `.env` files for configuration (frontend: `NEXT_PUBLIC_` prefix for client-side vars).

### 🧪 Testing & Reliability
- **Always create Jest unit tests for new features** (functions, components, handlers, etc).
- **After updating any logic**, check whether existing unit tests need to be updated. If so, do it.
- **Tests should live in appropriate test folders**:
  - Frontend: `__tests__` folders or `.test.ts/.test.tsx` files alongside components
  - Backend: `/tests` folder mirroring the main app structure
  - Include at least:
    - 1 test for expected use
    - 1 edge case
    - 1 failure case

### ✅ Task Completion
- **Mark completed tasks in `TASK.md`** immediately after finishing them.
- Add new sub-tasks or TODOs discovered during development to `TASK.md` under a "Discovered During Work" section.

### 📎 Style & Conventions
- **Frontend (Next.js + TypeScript)**:
  - Use **TypeScript** with strict type checking
  - Follow **React/Next.js best practices**
  - Use **Tailwind CSS** for styling with consistent design system
  - Use **Shadcn/ui components** for UI elements
  - Format with **Prettier** and lint with **ESLint**
  - Use **functional components with hooks**
- **Backend (Go)**:
  - Follow **Go conventions** (gofmt, golint)
  - Use **Prisma** for database operations with MongoDB
  - Implement proper **error handling** with Go idioms
  - Use **structured logging**
  - Follow **clean architecture patterns**
- **Database**:
  - Use **Prisma schema** for MongoDB modeling
  - Implement proper **indexing strategies**
  - Use **connection pooling**
- Write **comprehensive documentation** for all functions:
  ```typescript
  /**
   * Brief summary.
   * 
   * @param param1 - Description
   * @returns Description
   */
  ```
  ```go
  // FunctionName does something specific.
  // It takes param1 and returns result with potential error.
  func FunctionName(param1 string) (result string, err error) {
  ```

### 📚 Documentation & Explainability
- **Update `README.md`** when new features are added, dependencies change, or setup steps are modified.
- **Comment non-obvious code** and ensure everything is understandable to a mid-level developer.
- When writing complex logic, **add an inline comment** explaining the why, not just the what.
- **Document API endpoints** with proper OpenAPI/Swagger documentation.
- **Document component props** and usage examples.

### 🧠 AI Behavior Rules
- **Never assume missing context. Ask questions if uncertain.**
- **Never hallucinate libraries or functions** – only use known, verified packages from npm (frontend) or Go modules (backend).
- **Always confirm file paths and module names** exist before referencing them in code or tests.
- **Never delete or overwrite existing code** unless explicitly instructed to or if part of a task from `TASK.md`.
- **Consider both frontend and backend implications** when implementing features.
- **Ensure proper error handling** across the full stack.