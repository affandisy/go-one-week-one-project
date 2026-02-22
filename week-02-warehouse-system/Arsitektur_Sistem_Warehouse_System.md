### 1. Arsitektur Sistem

| Komponen | Teknologi |
|------|----------|
| **Backend** | Go (Fiber + GORM + Native SQL) |
| **Database** | PostgreSQL |
| **Frontend** | React + Tailwind CSS |
| **Auth** | JWT (github.com/golang-jwt) + RBAC middleware |
| **Background Jobs** | Go channels + goroutines |
| **Deployment** | Docker |

### 2. Database Schema
```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users & Auth
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL,
    full_name VARCHAR NOT NULL,
    phone VARCHAR,
    status VARCHAR DEFAULT 'active',
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

-- Warehouses
CREATE TABLE warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    address TEXT,
    timezone VARCHAR DEFAULT 'Asia/Jakarta'
);

-- Products (Master Barang)
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
    photo_url TEXT, -- Supabase Storage URL
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Stock Transactions (Header)
CREATE TABLE stock_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    warehouse_id UUID REFERENCES warehouses(id),
    type VARCHAR NOT NULL, -- IN, OUT, ADJUSTMENT
    subtype VARCHAR, -- PO, SALE, RETURN, WASTE
    reference_no VARCHAR NOT NULL,
    partner_name VARCHAR,
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

-- Stock Ledger (Source of Truth untuk Stok)
CREATE TABLE stock_ledger (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    product_id UUID NOT NULL REFERENCES products(id),
    ref_transaction_id UUID REFERENCES stock_transactions(id),
    ref_line_id UUID REFERENCES stock_transaction_lines(id),
    direction VARCHAR NOT NULL CHECK (direction IN ('IN', 'OUT')), -- IN: +qty, OUT: -qty
    qty INTEGER NOT NULL,
    occurred_at TIMESTAMPTZ DEFAULT NOW(),
    INDEX idx_ledger_product_warehouse (product_id, warehouse_id),
    INDEX idx_ledger_time (occurred_at)
);

-- Inventory Balance (Cache untuk performa)
CREATE TABLE inventory_balances (
    product_id UUID PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    on_hand_qty INTEGER DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Audit Logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_user_id UUID REFERENCES users(id),
    action VARCHAR NOT NULL,
    entity_type VARCHAR, -- USER, TRANSACTION, PRODUCT
    entity_id UUID,
    before_data JSONB,
    after_data JSONB,
    ip INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes untuk performa
CREATE INDEX idx_transactions_status_created ON stock_transactions(status, created_at);
CREATE INDEX idx_products_sku ON products(sku);
```

* GORM Model: Buat struct Go dengan gorm:"column:..." tags. Gunakan raw SQL untuk complex query (stock summary, report).

### 3. API Specification (Fiber)
Base: /api/v1 | CORS enabled untuk React | Middleware: Auth, RBAC, RateLimit, Logging. Idempotency untuk approve endpoints.

Core Endpoints

Auth
```text
POST /api/v1/auth/register  (public)
POST /api/v1/auth/login     (rate-limited) → {accessToken, refreshToken, user}
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET  /api/v1/me
PATCH /api/v1/me
```

Users (Admin only)

```text
GET   /api/v1/users?page=1&limit=20&q=andi
POST  /api/v1/users
PATCH /api/v1/users/{id}
POST  /api/v1/users/{id}/roles
```

Products

```text
GET   /api/v1/products?page=1&q=sku001&category=elektronik
POST  /api/v1/products
GET   /api/v1/products/{id}
PATCH /api/v1/products/{id}
DELETE /api/v1/products/{id}  (soft delete)
```

Transactions

```text
POST  /api/v1/transactions  → status=DRAFT
PATCH /api/v1/transactions/{id}  (draft only)
GET   /api/v1/transactions?status=APPROVED&from=2026-01-01&to=2026-01-31
GET   /api/v1/transactions/{id}
POST  /api/v1/transactions/{id}/approve  (Idempotency-Key header, Manager+)
POST  /api/v1/transactions/{id}/cancel   (draft only)
```

Inventory

```text
GET  /api/v1/inventory/summary?warehouseId=xxx
GET  /api/v1/inventory/items?warehouseId=xxx
GET  /api/v1/inventory/{productId}/movements?from=...&to=...
```

Reports (Manager+)
```text
GET  /api/v1/reports/stock?warehouseId=xxx&asOf=2026-01-16
GET  /api/v1/reports/transactions?type=IN&from=...&to=...
POST /api/v1/reports/export/stock  → {jobId}
GET  /api/v1/jobs/{jobId}  → downloadUrl ketika ready
```

File Upload
```text
POST /api/v1/files/upload-url  → signed PUT URL untuk client upload langsung
```

### 4. Frontend (React + Tailwind)
* Routing: React Router (Dashboard, Transactions/In, Transactions/Out, Products, Reports)
* State: Zustand atau Context untuk global (user, warehouse)
* Forms: React Hook Form + Zod validation
* Tables/Charts: TanStack Table + Recharts (untuk report dashboard)
* Upload: Direct to Supabase Storage dengan signed URL
* UI: Tailwind components konsisten (shadcn/ui atau Headless UI)

### 5. Implementation Roadmap (Go + React)
Week 1: Setup & Core
* Supabase project + schema migrations (GORM.AutoMigrate)
* Go API skeleton (Fiber, CORS, middleware auth)
* React setup (Vite + Tailwind + shadcn)
* Auth endpoints + login page

Week 2: Products & Inventory
* CRUD products + list dengan search/paginate
* Stock summary + movements (raw SQL: SELECT SUM(CASE WHEN direction='IN' THEN qty ELSE -qty END) FROM stock_ledger...)
* React tables + charts

Week 3: Transactions
* Transaction CRUD + approve (idempotency dengan Redis)
* Ledger update transactional (DB tx di GORM)

Week 4: Reports & Polish
* Report endpoints + async export (goroutine + email via Resend/Brevo)
* Opname module
* Audit logs + error monitoring (Sentry)
* Deploy: Dockerize Go API → Railway, React → Vercel

### 6. Deployment & Ops
* Supabase: Free tier cukup untuk startup (500MB DB, 1GB storage)
* Go API: Dockerfile → Railway/Fly.io ($5/bln)
* React: Vercel free tier
* Env Vars: SUPABASE_URL, SUPABASE_ANON_KEY, JWT_SECRET
* Monitoring: Supabase dashboard + Go pprof + Sentry

### 7. Success Metrics & Testing
* E2E Tests: Playwright untuk React flows (login → create transaksi → approve → check stok)
* API Tests: Go test + Postman collection
* Load Test: 100 concurrent transaksi approve (Go handle mudah)
* Metrics: Stok accuracy 99.9% (ledger vs physical), API latency <200ms