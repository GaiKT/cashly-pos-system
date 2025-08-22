# Create Full-Stack Web App PRP

## Feature file: $ARGUMENTS

Generate a complete PRP for full-stack web application feature implementation with thorough research. Ensure context is passed to the AI agent to enable self-validation and iterative refinement across both frontend (Next.js + Tailwind + Shadcn) and backend (Go + Prisma + MongoDB).

Read the feature file first to understand what needs to be created, how the examples provided help, and any other considerations for both UI/UX and API requirements.

The AI agent only gets the context you are appending to the PRP and training data. Assume the AI agent has access to the codebase and the same knowledge cutoff as you, so it's important that your research findings are included or referenced in the PRP. The Agent has web search capabilities, so pass URLs to documentation and examples.

## Research Process

1. **Codebase Analysis**
   - Search for similar features/patterns in both frontend and backend
   - Identify component patterns in frontend/components/
   - Check API handler patterns in backend/handlers/
   - Review Prisma schema for data model patterns
   - Note existing conventions to follow (naming, structure, error handling)
   - Check test patterns for both Jest (frontend) and Go testing (backend)

2. **Frontend Research**
   - Next.js patterns (pages, API routes, SSR/ISR)
   - Tailwind design system usage
   - Shadcn component implementations
   - React hooks patterns for data fetching
   - State management approaches
   - Error boundary implementations

3. **Backend Research** 
   - Go HTTP handler patterns
   - Prisma with MongoDB best practices
   - Error handling in Go
   - Middleware patterns
   - Database schema design
   - API design conventions (RESTful)

4. **Integration Research**
   - Frontend-backend communication patterns
   - Authentication/authorization flows
   - Error handling across the stack
   - Environment variable management
   - CORS configuration

5. **External Research**
   - Library documentation (include specific URLs)
   - Implementation examples (GitHub/StackOverflow/blogs)
   - Best practices and common pitfalls for full-stack development
   - Performance optimization techniques

6. **User Clarification** (if needed)
   - Specific UI/UX requirements and design patterns?
   - API design preferences and data flow?
   - Integration requirements with existing features?

## PRP Generation

Using PRPs/templates/prp_base.md as template:

### Critical Context to Include
- **Documentation**: URLs with specific sections for each technology
- **Code Examples**: Real snippets from both frontend and backend
- **Gotchas**: Library quirks, version issues, common integration problems
- **Patterns**: Existing approaches to follow for components, handlers, and data models
- **Database Schema**: Current Prisma models and relationships
- **API Design**: Existing endpoint patterns and conventions

### Implementation Blueprint
- Start with database schema design
- Define API endpoints and data flow
- Plan frontend component hierarchy
- Reference real files for patterns
- Include error handling strategy for both frontend and backend
- List tasks in implementation order (usually: schema → backend → frontend → integration)

### Validation Gates (Must be Executable)
```bash
# Frontend validation
cd frontend
npm run lint && npm run type-check && npm test

# Backend validation  
cd backend
go fmt ./... && go vet ./... && golangci-lint run && go test ./...

# Database validation
npx prisma validate && npx prisma db push --accept-data-loss
```

*** CRITICAL: AFTER RESEARCHING AND EXPLORING THE CODEBASE ***

*** THINK DEEPLY ABOUT THE FULL-STACK IMPLEMENTATION AND PLAN YOUR APPROACH ***
- Consider data flow from database → backend → frontend
- Think about error states and loading states
- Plan for responsive design and user experience
- Consider authentication and authorization needs
- Think about performance implications

*** THEN START WRITING THE PRP ***

## Output
Save as: `PRPs/{feature-name}.md`

## Quality Checklist
- [ ] All necessary context for both frontend and backend included
- [ ] Database schema design considerations documented
- [ ] API endpoints clearly defined
- [ ] Frontend component hierarchy planned
- [ ] Validation gates are executable by AI
- [ ] References existing patterns from codebase
- [ ] Clear implementation path with task ordering
- [ ] Error handling documented for both sides
- [ ] Integration points clearly identified
- [ ] Performance and UX considerations included

Score the PRP on a scale of 1-10 (confidence level to succeed in one-pass implementation using Claude with full-stack context)

Remember: The goal is one-pass implementation success through comprehensive full-stack context that enables the AI to build working features that integrate seamlessly with existing patterns.