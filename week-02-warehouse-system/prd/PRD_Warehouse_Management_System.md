# Product Requirements Document (PRD)
## Warehouse Management System – Complete System Design

**Versi:** 1.0  

---

## 1. Product Overview

### Executive Summary
Warehouse Management System (WMS) is a web‑based application designed for startups and SMEs (1‑5 users) to efficiently manage warehouse operations. It enables user management, real‑time stock tracking, inbound/outbound transactions, and actionable reports.

### Primary Goals:
- Automate inventory processes.
- Reduce manual entry errors.
- Provide real‑time stock visibility.
- Deliver actionable analytics for decision making.

### Target Market
- SMEs in retail, distribution, and trading.
- E‑commerce startups with small warehouse operations.
- Family businesses seeking inventory digitization.

### Problem Statement
- Difficulty in accurate stock tracking.
- No real‑time visibility of goods movement.
- Manual reporting is time‑consuming and error‑prone.
- Lack of access control for different user levels.

---

## 2. System Architecture (Clean Architecture)

The system follows Clean Architecture principles to separate concerns, improve testability, and maintain flexibility. All business logic resides in the service layer, while the database, API, and external agents are treated as implementation details.

```text
┌─────────────────────────────────────────────────────────────┐
│                     Presentation Layer                       │
│                    (React.js Frontend)                       │
└───────────────────────────┬─────────────────────────────────┘
                            │ HTTP/JSON
┌───────────────────────────▼─────────────────────────────────┐
│                     Delivery Layer (Handlers)                │
│                      (Fiber, Middlewares)                     │
├───────────────────────────────────────────────────────────────┤
│                     Usecase / Service Layer                   │
│          (Business logic, validation, orchestration)          │
├───────────────────────────────────────────────────────────────┤
│                     Repository Layer                           │
│              (Database operations, GORM / raw SQL)             │
└───────────────────────────┬─────────────────────────────────┘
                            │ SQL
┌───────────────────────────▼─────────────────────────────────┐
│                     Database (PostgreSQL)                    │
└───────────────────────────────────────────────────────────────┘

Additional Components:
- Background Jobs (Go channels/goroutines for async exports)
- Object Storage (S3-compatible for attachments)
```

---

### Layer Descriptions
##### 2.1 Models / Entities
- Plain Go structs with GORM tags (e.g., User, Product, StockTransaction).
- Define business domain objects.
##### 2.2 Repository Layer
- Interface‑based data access.
- Implements CRUD and complex queries using GORM or raw SQL (for performance‑critical reports).
- Example: ProductRepository, TransactionRepository.

##### 2.3 Usecase / Service Layer
- Contains all business logic.
- Validates inputs, enforces rules (e.g., stock sufficiency for outbound).
- Orchestrates repositories, background jobs, and external services.
- Example: TransactionService handles draft creation, approval, and ledger updates.

##### 2.4 Delivery Layer (Handlers)
- HTTP handlers using Fiber.
- Parse requests, bind to DTOs, call usecase methods, format responses.
- Middleware for authentication (JWT), RBAC, rate limiting, logging, and idempotency.

##### 2.5 Presentation Layer
- React.js single‑page application.
- Communicates with API via REST.
- State management with Zustand/Context, forms with React Hook Form + Zod.

##### 2.6 Background Jobs
- Lightweight concurrency using Go channels and goroutines.
- Handle report generation, email notifications, and long‑running exports.
- Job status tracked in a jobs table (optional).

---

## 3. Technology Stack
| Component | Technology |
|------|----------|
| **Backend** | Go (1.22+) with Fiber v2, GORM, github.com/golang-jwt/jwt |
| **Database** | PostgreSQL 15+ with uuid-ossp extension |
| **Frontend** | React 18+ (Vite), Tailwind CSS, shadcn/ui components, TanStack Table |
| **Auth** | JWT (access + refresh tokens), RBAC middleware |
| **Background** | Native Go goroutines + channels |
| **File Storage** | Supabase Storage (S3‑compatible) or any S3 provider |
| **Caching** | Redis (optional, for idempotency keys and rate limiting) |
| **Deployment** | Docker, Railway (API), Vercel (frontend) |
| **Monitoring** | Structured logs, Sentry, Go pprof |

---

## 4. Database Schema

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL,
    full_name VARCHAR NOT NULL,
    phone VARCHAR,
    status VARCHAR DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Roles
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR UNIQUE NOT NULL, -- ADMIN, MANAGER, OPERATOR
    name VARCHAR NOT NULL
);

-- User Roles (many-to-many)
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

-- Warehouses (supports future multi-warehouse)
CREATE TABLE warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    address TEXT,
    timezone VARCHAR DEFAULT 'Asia/Jakarta'
);

-- Products (master barang)
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sku VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    description TEXT,
    category VARCHAR,
    uom VARCHAR NOT NULL, -- pcs, kg, box
    buy_price DECIMAL(12,2),
    sell_price DECIMAL(12,2),
    min_stock INTEGER DEFAULT 0,
    max_stock INTEGER DEFAULT 999999,
    is_active BOOLEAN DEFAULT true,
    photo_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Stock Transactions (Header)
CREATE TABLE stock_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    warehouse_id UUID REFERENCES warehouses(id),
    type VARCHAR NOT NULL, -- IN, OUT, ADJUSTMENT, OPNAME
    subtype VARCHAR, -- PO, SALE, RETURN, WASTE, etc.
    reference_no VARCHAR NOT NULL,
    partner_name VARCHAR, -- supplier / customer
    status VARCHAR DEFAULT 'DRAFT', -- DRAFT, APPROVED, REJECTED, CANCELLED
    notes TEXT,
    attachment_url TEXT,
    created_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Transaction Lines
CREATE TABLE stock_transaction_lines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID REFERENCES stock_transactions(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id),
    qty INTEGER NOT NULL,
    unit_price DECIMAL(12,2),
    notes TEXT
);

-- Stock Ledger (source of truth for stock movements)
CREATE TABLE stock_ledger (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    product_id UUID NOT NULL REFERENCES products(id),
    ref_transaction_id UUID REFERENCES stock_transactions(id),
    ref_line_id UUID REFERENCES stock_transaction_lines(id),
    direction VARCHAR NOT NULL CHECK (direction IN ('IN', 'OUT')),
    qty INTEGER NOT NULL,
    occurred_at TIMESTAMPTZ DEFAULT NOW(),
    INDEX idx_ledger_product_warehouse (product_id, warehouse_id),
    INDEX idx_ledger_time (occurred_at)
);

-- Inventory Balances (cache for performance)
CREATE TABLE inventory_balances (
    product_id UUID PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    on_hand_qty INTEGER DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Stock Opname Sessions
CREATE TABLE stock_opname_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    warehouse_id UUID REFERENCES warehouses(id),
    code VARCHAR UNIQUE,
    status VARCHAR DEFAULT 'DRAFT', -- DRAFT, COUNTING, SUBMITTED, APPROVED
    started_at TIMESTAMPTZ,
    submitted_at TIMESTAMPTZ,
    approved_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id)
);

-- Opname Lines
CREATE TABLE stock_opname_lines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID REFERENCES stock_opname_sessions(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id),
    system_qty INTEGER,
    physical_qty INTEGER,
    diff_qty INTEGER
);

-- Audit Logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_user_id UUID REFERENCES users(id),
    action VARCHAR NOT NULL,
    entity_type VARCHAR, -- USER, TRANSACTION, PRODUCT, etc.
    entity_id UUID,
    before_data JSONB,
    after_data JSONB,
    ip INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_transactions_status_created ON stock_transactions(status, created_at);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
```

Note: Use GORM models with gorm:"column:..." tags. For complex reports (stock summary, movements), use raw SQL for optimal performance.

---

## 5. API Specifications (Fiber)
**Base URL**: /api/v1
**Content-Type**: application/json
**Auth Header**: Authorization: Bearer <access_token>
**Idempotency**: Endpoints that change state (e.g., approve transaction) require Idempotency-Key: <uuid> header to prevent duplicate processing.

###### Response Conventions
```json
// Success
{ "data": {...}, "meta": { "page": 1, "limit": 20, "total": 100 } }

// Error
{ "error": { "code": "VALIDATION_ERROR", "message": "qty must be > 0", "details": [...] } }
```

### 5.1 Auth & User
| Method | Endpoint | Description | Access |
|------|----------|--------------|------|
| **POST** | /auth/register | Register new user | Public |
| **POST** | /auth/login | Login (rate-limited) | Public |
| **POST** | /auth/logout | Logout (invalidate refresh token) |Auth |
| **GET** | /me | Get current user profile | Auth |
| **PATCH** | /me | Update own profile | Auth |
| **GET** | /users | List users (filter: q, role) | Admin |
| **POST** | /users | Create user (operator/manager)| Admin |
| **PATCH** | /users/{id} | Update user | Admin |
| **POST** | /users/{id}/roles | Assign roles to user |Admin |

### 5.3 Inventory
| Method | Endpoint | Description | Access |
|------|----------|--------------|------|
| **GET** | /inventory/summary?warehouseId=... | Stock summary (totals, low/over stock) | Operator+ |
| **GET** | /inventory/items?warehouseId=...&q=... | Paginated stock items with balances | Operator+ |
| **GET** | /inventory/items/{productId}/movements?from=...&to=... | Stock movement history for a product |Operator+ |

### 5.4 Transactions
| Method | Endpoint | Description | Access |
|------|----------|--------------|------|
| **POST** | /transactions | Create draft transaction | Operator+ |
| **PATCH** | /transactions/{id} | Update draft transaction (only if DRAFT) | Operator+ |
| **GET** | /transactions | List transactions (filter by date, type, status) | Operator+ |
| **GET** | /transactions/{id} | Get transaction details |Operator+ |
| **POST** | /transactions/{id}/approve | Approve transaction (idempotent) | Manager+ |
| **POST** | /transactions/{id}/cancel | Cancel draft transaction | Operator+ |

**Request body for POST /transactions:**
```json
{
  "warehouseId": "uuid",
  "type": "IN",              // IN, OUT, ADJUSTMENT
  "subtype": "PO",           // PO, SALE, RETURN, WASTE, etc.
  "referenceNo": "PO-2026-001",
  "partnerName": "PT Supplier",
  "notes": "Arrived in the morning",
  "lines": [
    { "productId": "uuid", "qty": 10, "unitPrice": 5000 },
    { "productId": "uuid", "qty": 5,  "unitPrice": 12000 }
  ]
}
```
**Approve Logic (in service layer):**
- Validate stock sufficiency for OUT transactions.
- Insert rows into stock_ledger for each line (IN adds, OUT subtracts).
- Update inventory_balances atomically.
- Change transaction status to APPROVED and record approved_by, approved_at.
- Idempotency: store processed idempotency keys (in Redis or DB) to reject duplicate requests.

### 5.5 Opname
| Method | Endpoint | Description | Access |
|------|----------|--------------|------|
| **POST** | /opname-sessions | Create a new opname session | Manager+ |
| **POST** | /opname-sessions/{id}/start | Start counting | Manager+ |
| **POST** | /opname-sessions/{id}/submit | Submit physical counts | Operator+ |
| **POST** | /opname-sessions/{id}/approve | Approve opname (auto‑create adjustment transaction) | Manager+ |

### 5.6 Reports (Manager+)
| Method | Endpoint | Description | Access |
|------|----------|--------------|------|
| **GET** | /reports/stock?warehouseId=...&asOf=2026-01-13 | Stock report as of date | Manager+ |
| **GET** | /reports/transactions?from=...&to=...&type=... | Transaction report with filters | Manager+ |
| **GET** | /reports/movement/top?from=...&to=...&limit=10 | Top moving products | Manager+ |
| **POST** | /reports/export | Request async export (returns jobId) | Manager+ |
| **GET** | /jobs/{jobId} | Get job status + download URL when ready | Manager+ |

### 5.7 File Upload
| Method | Endpoint | Description |
|------|----------|--------------|
| **POST** | /files/upload-url | Request signed URL for direct upload to S3 |

---

## 6. Frontend (React + Tailwind)
* Framework: React with Vite.
* Routing: React Router v6.
* State Management: Zustand (or Context) for global state (user, warehouse, theme).
* Forms: React Hook Form + Zod validation.
* UI Components: shadcn/ui (built on Radix UI) for consistent, accessible components.
* Tables: TanStack Table (React Table) for advanced data grids with sorting, filtering, pagination.
* Charts: Recharts for dashboard visualizations.
* HTTP Client: Axios with interceptors for token refresh.
* File Upload: Direct to Supabase Storage using signed URLs.

**Key Pages:**
- Dashboard (summary charts, low stock alerts)
- Transactions (Inbound, Outbound, History)
- Products (list, create, edit)
- Inventory (stock summary, item movements)
- Reports (stock, transaction, analytics)
- User Management (admin only)
- Settings (profile, warehouse)

---

## 7. Non‑Functional Requirements

* Performance:
    * Dashboard load < 3s.
    * Transaction submit < 5s.
    * Report generation < 10s (async for longer exports).
* Security:
    * HTTPS only.
    * Passwords hashed with bcrypt/argon2.
    * JWT short‑lived access tokens + refresh tokens.
    * Audit logs for all create/update/delete actions.
    * Daily automated backups.
* Scalability:
    * Designed for 100k+ products, 1M+ transactions.
    * Multi‑warehouse ready (future).
* Availability:
    * 99.5% uptime target.
    * Graceful degradation if a service fails.
* Usability:
    * Mobile responsive design.
    * Clear navigation, tooltips, helpful error messages.
    * Maximum 3 clicks to reach any feature.

---

## 8. Implementation Roadmap
#### Phase 1: MVP (Weeks 1–4)
- Project setup: Go (Fiber), PostgreSQL, React (Vite + Tailwind)
- Auth endpoints (register, login, refresh, logout)
- User management (admin only)
- Product CRUD
- Transaction draft creation (inbound/outbound)
- Basic stock tracking (inventory_balances updated on approve)
- Dashboard with simple metrics

#### Phase 2: Core Features (Weeks 5–8)
- Approval workflow with idempotency
- Stock ledger as source of truth
- Stock movement history
- Report endpoints (stock, transactions)
- Async export jobs
- Audit logging

#### Phase 3: Enhancement (Weeks 9–12)
- Stock opname module
- Advanced analytics dashboard (charts, trends)
- Email notifications (Resend/Brevo)
- Multi‑warehouse support (basic)
- Performance optimization (indexes, query tuning)

#### Phase 4: Polish & Scale (Weeks 13–16)
- Mobile‑friendly refinements
- Integration tests (Playwright, Go test)
- Load testing (100 concurrent approvals)
- Dockerize and deploy to production (Railway + Vercel)
- Monitoring (Sentry, structured logs)

---

## 9. Deployment & Operations
- **Database**: Supabase PostgreSQL (free tier sufficient for startup: 500 MB).
- **Storage**: Supabase Storage (S3‑compatible) for product photos and attachments.
- **Backend API**: Dockerized Go app → deployed on Railway or Fly.io ($5/month).
- **Frontend**: Static hosting on Vercel (free tier).
- **Environment Variables**: SUPABASE_URL, SUPABASE_ANON_KEY, JWT_SECRET, DATABASE_URL, etc.
- **Monitoring**: Supabase dashboard, Go pprof, Sentry for errors.
- **Backup**: Automated daily via Supabase.

---

## 10. Success Metrics
- **Adoption**: Weekly active users, 30‑day retention.
- **Operational**: Number of transactions/day, stock accuracy (ledger vs physical), time to process a transaction.
- **Business**: User satisfaction (NPS), feature adoption, cost savings (manual work reduction).
- **Technical**: API latency <200ms, 99.9% ledger accuracy, zero duplicate transactions due to idempotency.

---

## 11. Glossary
| Term | Definisi |
|------|----------|
| SKU |	Stock Keeping Unit – unique identifier for a product. |
| Inbound | Goods received into the warehouse. |
| Outbound | Goods shipped out of the warehouse. |
| Opname | Physical stock counting process. |
| Ledger | Immutable record of every stock movement (source of truth). | 
| Draft | Transaction not yet finalised. |
| Approved | Transaction that has been confirmed and updated stock. |
| Min/Max Stock | Thresholds for low and over stock alerts. |

---

**Document Version:** 1.0
**Last Updated:** 2026-02-22