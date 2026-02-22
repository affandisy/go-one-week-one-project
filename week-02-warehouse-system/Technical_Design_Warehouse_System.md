## 1. Arsitektur & Komponen
- Client: Web responsive (React.js)
- API: REST API (monolith modular) + background job (untuk export report & email schedule).
- DB: PostgreSQL (relasional) + indeks untuk SKU, waktu transaksi, dan foreign key untuk integritas.
- Storage: Object storage (S3-compatible) untuk attachment (foto bukti, file export).
- Observability: structured log + audit log table, metrics latency/error untuk endpoint inventory/transaksi.

## 2. Desain Data (Database Schema)
Struktur ini mengikuti kebutuhan umum inventory: entity barang + transaksi + mutasi stok (ledger) agar stok bisa dihitung konsisten dari pergerakan
**Tabel inti**
**users**
- id (uuid, pk), email (unique), password_hash, full_name, phone, status, created_at, updated_at, last_login_at

**roles**
- id, code: ADMIN|MANAGER|OPERATOR, name

**user_roles**
- user_id, role_id (untuk fleksibilitas; jika 1 user 1 role tetap bisa dipakai)

**warehouses**

- id, name, address, timezone (untuk future multi-warehouse)

**products**

- id (uuid), sku (unique), name, description, category_id (nullable), uom, buy_price, sell_price, min_stock, max_stock, is_active, photo_url, created_at

**inventory_balances (opsional cache)**

- product_id, warehouse_id, on_hand_qty, reserved_qty (opsional), updated_at
Catatan: on_hand juga bisa dihitung dari ledger; tabel ini dipakai untuk performa.

**stock_transactions (header)**

- id, warehouse_id, type: IN|OUT|ADJUSTMENT|OPNAME, subtype (PO/SALE/RETURN/WASTE/etc), reference_no, partner_name (supplier/customer), status: DRAFT|APPROVED|REJECTED|CANCELLED, notes, attachment_url, created_by, approved_by, approved_at, created_at

**stock_transaction_lines (detail)**

- id, transaction_id, product_id, qty, unit_price (nullable), notes

**stock_ledger (mutasi / source of truth)**

- id, warehouse_id, product_id, ref_transaction_id, ref_line_id, direction: IN|OUT, qty, occurred_at, created_at
Ledger memastikan setiap pergerakan tercatat (mirip konsep “Stock In/Out Records” yang umum pada skema inventory).

**stock_opname_sessions**

- id, warehouse_id, code, status: DRAFT|COUNTING|SUBMITTED|APPROVED, started_at, submitted_at, approved_at, created_by

**stock_opname_lines**

- id, session_id, product_id, system_qty, physical_qty, diff_qty

**audit_logs**

- id, actor_user_id, action, entity_type, entity_id, before_json, after_json, ip, user_agent, created_at
Praktik audit: log semua event security-relevant (login sukses/gagal, perubahan role/permission, keputusan auth) dan jangan log data sensitif.

## 3. API Spec (REST)
Base URL: /api/v1
Content-Type: application/json
Auth Header: Authorization: Bearer <access_token>
Idempotency: untuk endpoint “commit/approve” transaksi gunakan Idempotency-Key: <uuid> dan simpan key+hasil (dedup window) agar retry tidak menggandakan transaksi.

Konvensi respons
```json
// sukses
{ "data": {...}, "meta": {...} }

// error
{ "error": { "code": "VALIDATION_ERROR", "message": "qty must be > 0", "details": [...] } }
```
1. Auth & User
POST /auth/register (public)
Request:
```json
{ "email":"a@b.com", "password":"Secret123", "fullName":"Andi", "phone":"08..." }
```
Response: user + token (opsional), atau “verify email” (opsional).
POST /auth/login (public, rate-limited)
Request:
```json
{ "email":"a@b.com", "password":"Secret123" }
```
Response:
```json
{ "data": { "accessToken":"...", "refreshToken":"...", "user": { "id":"...", "roles":["ADMIN"] } } }
```
POST /auth/refresh
POST /auth/logout
GET /me
PATCH /me (update profil)

(Admin) GET /users (filter: q, role, status, page)
(Admin) POST /users (buat user operator/manager)
(Admin) PATCH /users/{id}
(Admin) POST /users/{id}/roles (set roles)
Audit log untuk login/logout, failed login, role change.

2. Master Barang (Products)
GET /products (q, category, isActive, sort, page)
POST /products (Admin/Manager)
GET /products/{id}
PATCH /products/{id}
DELETE /products/{id} (soft delete: is_active=false)
Request create:
```json
{
  "sku":"SKU-001",
  "name":"Kabel USB",
  "uom":"pcs",
  "buyPrice":5000,
  "sellPrice":8000,
  "minStock":10,
  "maxStock":200
}
```

3. Stok (Inventory)
GET /inventory/summary?warehouseId=...
Response:
```json
{
  "data":{
    "totalSkus": 120,
    "lowStockCount": 7,
    "overStockCount": 3,
    "stockValue": 125000000
  }
}
```
GET /inventory/items?warehouseId=...&q=SKU-001
Response per item: onHandQty, minStock, maxStock, lastMovementAt.

GET /inventory/items/{productId}/movements?warehouseId=...&from=...&to=...
Mengambil ledger/mutasi (paginate).

4. Transaksi Barang (Inbound/Outbound/Adjustment)
Draft & edit
POST /transactions (Operator+)
Request:
```json
{
  "warehouseId":"...",
  "type":"IN",
  "subtype":"PO",
  "referenceNo":"PO-2026-001",
  "partnerName":"PT Supplier",
  "notes":"Barang datang pagi",
  "lines":[
    { "productId":"...", "qty":10, "unitPrice":5000 },
    { "productId":"...", "qty":5, "unitPrice":12000 }
  ]
}
```
Response: status=DRAFT
PATCH /transactions/{id} (hanya DRAFT)
GET /transactions (filter by date range, type, status, productId)
GET /transactions/{id} (detail + lines)

Approve/Cancel (idempotent)
POST /transactions/{id}/approve (Manager/Admin)
Header: Idempotency-Key: uuid
Server behavior:
* Validasi stok cukup untuk OUT.
* Insert stock_ledger untuk setiap line (IN menambah, OUT mengurangi, ADJUSTMENT sesuai diff).
* Update inventory_balances secara atomik.
Response:
```json
{ "data": { "id":"...", "status":"APPROVED", "approvedAt":"2026-01-13T..." } }
```
POST /transactions/{id}/cancel (DRAFT only; atau policy tertentu)

5. Opname
POST /opname-sessions (Manager+)
POST /opname-sessions/{id}/start
POST /opname-sessions/{id}/submit (mengunci input)
POST /opname-sessions/{id}/approve (Manager/Admin, idempotent)
Saat approve: buat transaksi ADJUSTMENT otomatis berbasis diff dan hasilkan ledger, supaya histori stok tetap konsisten

6. Laporan & Analytics
GET /reports/stock?warehouseId=...&asOf=2026-01-13
GET /reports/transactions?warehouseId=...&from=...&to=...&type=IN|OUT
GET /reports/movement/top?warehouseId=...&from=...&to=...&limit=10
POST /reports/export (async job) → return jobId
GET /jobs/{jobId} → status + downloadUrl ketika selesai
Monitoring endpoint inventory/transaksi untuk latency & error (praktik ketahanan API).

### Error Codes & Validasi
* VALIDATION_ERROR (SKU duplikat, qty <= 0, field wajib kosong)
* INSUFFICIENT_STOCK (untuk OUT)
* FORBIDDEN (RBAC)
* CONFLICT_IDEMPOTENCY (Idempotency-Key pernah dipakai dengan payload berbeda) — mencegah perilaku ambigu saat retry.
​