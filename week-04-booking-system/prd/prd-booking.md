# Padel Court Booking System – Complete System Design

**Version:** 1.0  
**Date:** 2026‑03‑21  
**Author:** Product Manager & Software Architect

---

## 1. Executive Summary

The Padel Court Booking System enables customers to book padel courts online in under 2 minutes, prevents double bookings through concurrency control, implements dynamic pricing, and provides an analytics dashboard for business owners. The system is built using **Go** (Fiber, GORM) on the backend, **Svelte** on the frontend, and **PostgreSQL** as the persistence layer, adhering to **Clean Architecture** principles where all business logic resides in the service layer and the database is used only for CRUD operations.

**Key Objectives:**
- Seamless booking experience (< 2 minutes end‑to‑end).
- Zero double bookings via temporary locks and optimistic locking.
- Dynamic pricing based on day type (weekday, weekend) and time slots.
- Integrated payment gateway (Midtrans/Xendit) with automatic state transitions.
- Admin dashboard for revenue, utilization, and VIP customer insights.

---

## 2. System Architecture (Clean Architecture)

The system follows **Clean Architecture** to separate concerns, improve testability, and maintain flexibility. All business logic is implemented in the **service layer**; the database is used solely for data persistence.

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Presentation Layer                            │
│                         (Svelte Frontend)                            │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │ HTTP/JSON
┌─────────────────────────────────▼───────────────────────────────────┐
│                         Delivery Layer (Handlers)                    │
│                         (Fiber, Middlewares)                         │
├─────────────────────────────────────────────────────────────────────┤
│                         Usecase / Service Layer                      │
│    (Business logic, pricing, locking, state transitions, validations)│
├─────────────────────────────────────────────────────────────────────┤
│                         Repository Layer                             │
│              (Database operations via GORM, raw SQL)                 │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │ SQL
┌─────────────────────────────────▼───────────────────────────────────┐
│                         Database (PostgreSQL)                        │
│                   (Only CRUD – no business logic)                    │
└─────────────────────────────────────────────────────────────────────┘

Additional Components:
- Cache (Redis) – optional, for temporary locks and high‑frequency data.
- Background Jobs – Go goroutines for expired lock cleanup, scheduled reports.
- Payment Gateway – Midtrans/Xendit.
- Notification Service – WhatsApp Business API / SendGrid.
```

### Layer Responsibilities

- **Entities (Models)** – Plain Go structs representing core domain objects (e.g., `Court`, `Booking`, `User`). No database logic.
- **Repository Layer** – Interface‑based data access. Uses GORM for CRUD operations and simple queries. Complex reporting queries may use raw SQL.
- **Service Layer (Usecase)** – Contains all business rules: dynamic pricing calculation, slot availability validation, lock management, state transitions, payment callback handling. This layer orchestrates repositories, external services (payment, notifications), and background jobs.
- **Delivery Layer (Handlers)** – HTTP handlers using Fiber. Parse requests, bind to DTOs, call service methods, return JSON responses. Includes middleware for authentication (JWT), RBAC, rate limiting, and logging.
- **Presentation Layer** – Svelte single‑page application consuming the REST API. Handles UI, state management, and client‑side validation.

---

## 3. Technology Stack

| Component           | Technology                                                                 |
|---------------------|----------------------------------------------------------------------------|
| **Backend**         | Go 1.22+ with Fiber v2, GORM, `github.com/golang-jwt/jwt`                  |
| **Database**        | PostgreSQL 15+ (with `uuid-ossp`)                                          |
| **Cache (optional)**| Redis (for locks, session storage, rate limiting)                          |
| **Frontend**        | Svelte (SvelteKit) with Tailwind CSS, Skeleton UI components               |
| **Auth**            | JWT (access + refresh tokens), RBAC middleware                             |
| **Payment**         | Midtrans / Xendit API (QRIS, Virtual Account, Credit Card)                 |
| **Notifications**   | WhatsApp Business API (via Twilio) / SendGrid for email                    |
| **Background Jobs** | Native Go goroutines + channels                                            |
| **Deployment**      | Docker, Railway (backend), Vercel (frontend)                               |
| **Monitoring**      | Structured logs, Sentry, Prometheus (optional)                             |

---

## 4. Database Schema

The schema is designed for CRUD operations only. All business logic (pricing, availability, locking) is handled in the service layer.

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Courts
CREATE TABLE courts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,                  -- e.g., "Court A Indoor"
    type VARCHAR NOT NULL,                  -- indoor, outdoor
    status VARCHAR DEFAULT 'active',        -- active, inactive, maintenance
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Pricing Rules (dynamic price configuration)
CREATE TABLE pricing_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    court_id UUID REFERENCES courts(id) ON DELETE CASCADE,
    day_type VARCHAR NOT NULL,              -- weekday, weekend
    start_time TIME NOT NULL,               -- e.g., 08:00
    end_time TIME NOT NULL,                 -- e.g., 10:00
    base_price DECIMAL(10,2) NOT NULL,      -- price per hour
    multiplier DECIMAL(3,2) DEFAULT 1.0,    -- surge multiplier (1.0 = normal)
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR UNIQUE NOT NULL,
    whatsapp VARCHAR UNIQUE NOT NULL,        -- for notifications
    full_name VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL,
    role VARCHAR DEFAULT 'customer',         -- customer, admin, owner
    is_verified BOOLEAN DEFAULT false,       -- after OTP
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Bookings (state machine)
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_code VARCHAR UNIQUE NOT NULL,    -- human‑readable code (e.g., PDL-260321-001)
    user_id UUID REFERENCES users(id),
    court_id UUID REFERENCES courts(id),
    booking_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status VARCHAR NOT NULL,                 -- available, locked, pending, paid, cancelled, expired
    lock_expiry TIMESTAMPTZ,                 -- when lock expires (for locked status)
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    INDEX idx_booking_date_status (booking_date, status),
    INDEX idx_booking_user (user_id),
    INDEX idx_booking_court_date (court_id, booking_date)
);

-- Payments (linked to booking)
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID REFERENCES bookings(id) UNIQUE,
    payment_method VARCHAR,                  -- qris, va, card, etc.
    external_reference VARCHAR,              -- payment gateway transaction ID
    status VARCHAR NOT NULL,                 -- pending, success, failed, refunded
    amount DECIMAL(10,2) NOT NULL,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Notifications (optional, for audit)
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    type VARCHAR NOT NULL,                   -- email, whatsapp
    content TEXT,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Audit Logs (for admin actions)
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_user_id UUID REFERENCES users(id),
    action VARCHAR NOT NULL,
    entity_type VARCHAR NOT NULL,            -- booking, court, user
    entity_id UUID,
    before_data JSONB,
    after_data JSONB,
    ip INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

**GORM Model Example (simplified):**
```go
type Booking struct {
    ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
    BookingCode string     `gorm:"unique;not null"`
    UserID      uuid.UUID
    CourtID     uuid.UUID
    BookingDate time.Time  `gorm:"type:date"`
    StartTime   time.Time  `gorm:"type:time"`
    EndTime     time.Time  `gorm:"type:time"`
    Status      string     `gorm:"not null"`
    LockExpiry  *time.Time
    TotalPrice  float64    `gorm:"type:decimal(10,2)"`
    CreatedAt   time.Time
    UpdatedAt   time.Time

    // Relationships
    User   User   `gorm:"foreignKey:UserID"`
    Court  Court  `gorm:"foreignKey:CourtID"`
    Payment Payment `gorm:"foreignKey:BookingID"`
}
```

---

## 5. API Specifications (Fiber)

**Base URL:** `/api/v1`  
**Content-Type:** `application/json`  
**Auth Header:** `Authorization: Bearer <access_token>`

**Idempotency:** For payment callbacks and booking confirmations, use `Idempotency-Key` header to prevent duplicate processing.

### 5.1 Authentication

| Method | Endpoint          | Description                          | Access          |
|--------|-------------------|--------------------------------------|-----------------|
| POST   | `/auth/register`  | Register new user (with OTP)         | Public          |
| POST   | `/auth/verify-otp`| Verify OTP and activate account      | Public          |
| POST   | `/auth/login`     | Login (rate‑limited)                 | Public          |
| POST   | `/auth/refresh`   | Refresh access token                 | Public          |
| POST   | `/auth/logout`    | Logout                               | Auth            |
| GET    | `/me`             | Get current user profile             | Auth            |
| PATCH  | `/me`             | Update own profile                   | Auth            |

### 5.2 Courts & Pricing (Public)

| Method | Endpoint                | Description                          | Access          |
|--------|-------------------------|--------------------------------------|-----------------|
| GET    | `/courts`               | List all active courts               | Public          |
| GET    | `/courts/{id}`          | Get court details                    | Public          |
| GET    | `/courts/{id}/pricing`  | Get pricing rules for a court        | Public          |
| GET    | `/availability`         | Check available slots for a date     | Public          |

**Availability Request Example:**
```
GET /api/v1/availability?date=2026-03-25
```
Response:
```json
{
  "data": [
    {
      "courtId": "uuid",
      "courtName": "Court A Indoor",
      "slots": [
        { "start": "08:00", "end": "09:00", "price": 120000, "status": "available" },
        { "start": "09:00", "end": "10:00", "price": 150000, "status": "locked" },
        { "start": "10:00", "end": "11:00", "price": 150000, "status": "paid" }
      ]
    }
  ]
}
```

### 5.3 Bookings

| Method | Endpoint                      | Description                          | Access          |
|--------|-------------------------------|--------------------------------------|-----------------|
| POST   | `/bookings`                   | Create a booking (lock slot)         | Customer        |
| GET    | `/bookings`                   | List user's bookings                 | Customer        |
| GET    | `/bookings/{id}`              | Get booking details                  | Customer/Admin  |
| POST   | `/bookings/{id}/checkout`     | Initiate payment (after lock)        | Customer        |
| POST   | `/bookings/{id}/cancel`       | Cancel booking (only if not paid)    | Customer        |
| DELETE | `/bookings/{id}`              | Delete (admin only)                  | Admin           |

**POST `/bookings` Request Body:**
```json
{
  "courtId": "uuid",
  "bookingDate": "2026-03-25",
  "startTime": "08:00",
  "endTime": "09:00"
}
```
Response (with locked status and lock expiry):
```json
{
  "data": {
    "id": "uuid",
    "bookingCode": "PDL-260321-001",
    "status": "locked",
    "lockExpiry": "2026-03-21T10:30:00Z",
    "totalPrice": 120000,
    "paymentUrl": "https://payment-gateway.com/checkout/..."
  }
}
```

### 5.4 Payments (Callback)

| Method | Endpoint              | Description                          | Access          |
|--------|-----------------------|--------------------------------------|-----------------|
| POST   | `/payments/callback`  | Payment gateway webhook (idempotent) | Public (signed) |
| GET    | `/payments/{id}`      | Get payment status                   | Customer/Admin  |

**Webhook Payload (example from Midtrans):**
```json
{
  "transaction_id": "midtrans-123",
  "order_id": "PDL-260321-001",
  "status": "settlement",
  "amount": 120000
}
```
Service layer will:
- Verify signature.
- Update booking status to `paid`.
- Create payment record.
- Send notifications (WhatsApp/Email).

### 5.5 Admin Endpoints

| Method | Endpoint                      | Description                          | Access          |
|--------|-------------------------------|--------------------------------------|-----------------|
| GET    | `/admin/bookings`             | List all bookings (with filters)     | Admin/Owner     |
| PATCH  | `/admin/bookings/{id}`        | Override booking (e.g., refund)      | Admin           |
| POST   | `/admin/courts`               | Create court                         | Admin           |
| PATCH  | `/admin/courts/{id}`          | Update court                         | Admin           |
| POST   | `/admin/pricing-rules`        | Set pricing rules                    | Admin           |
| GET    | `/admin/dashboard`            | Revenue & utilization stats          | Admin/Owner     |
| GET    | `/admin/reports`              | Export reports (Excel/PDF)           | Admin/Owner     |

**Dashboard Response Example:**
```json
{
  "data": {
    "dailyRevenue": 1200000,
    "weeklyRevenue": 5600000,
    "utilizationRate": 78.5,
    "topCustomers": [
      { "name": "John Doe", "totalSpent": 1500000, "bookings": 12 }
    ]
  }
}
```

---

## 6. Frontend (Svelte)

- **Framework:** SvelteKit for routing and SSR/SPA flexibility.
- **UI Library:** Skeleton UI (or custom Tailwind components) for consistent design.
- **State Management:** Svelte stores (`writable` for user session, cart, etc.).
- **HTTP Client:** `fetch` with interceptors for JWT handling.
- **Forms:** Svelte forms with validation (e.g., `svelte-forms-lib`).
- **Real‑time Updates:** Optionally use Server‑Sent Events (SSE) for lock expiry notifications.

**Key Pages:**
- Home / Court Selection
- Availability Calendar (date picker + slot grid)
- Checkout / Payment
- User Dashboard (upcoming bookings, history)
- Admin Panel (court management, pricing, reports)

**Responsive Design:** Tailwind CSS ensures mobile‑friendly booking experience.

---

## 7. Business Logic Implementation (Service Layer)

All core algorithms reside in the service layer, using repositories for data access.

### 7.1 Dynamic Price Calculation

```go
// service/price.go
func (s *PriceService) CalculatePrice(courtID uuid.UUID, date time.Time, start, end time.Time) (float64, error) {
    // Determine day type (weekday/weekend)
    dayType := "weekday"
    if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
        dayType = "weekend"
    }

    // Fetch pricing rules from repository (GORM)
    rules, err := s.repo.GetPricingRules(courtID, dayType, start, end)
    if err != nil { ... }

    // Apply multipliers
    total := 0.0
    for _, rule := range rules {
        hours := rule.EndTime.Sub(rule.StartTime).Hours()
        total += rule.BasePrice * rule.Multiplier * hours
    }

    // Add tax (10%)
    return total * 1.10, nil
}
```

### 7.2 Concurrency Control (Locking)

**Approach:** Use database‑based locking with `lock_expiry` and optimistic versioning.

```go
// service/booking.go
func (s *BookingService) CreateBooking(ctx context.Context, req CreateBookingRequest) (*Booking, error) {
    // Begin transaction
    tx := s.repo.BeginTx(ctx)

    // 1. Check availability: no overlapping booking with status in ('locked','pending','paid')
    existing, err := s.repo.FindConflictingBookings(tx, req.CourtID, req.Date, req.Start, req.End)
    if err != nil { ... }
    if len(existing) > 0 {
        return nil, ErrSlotUnavailable
    }

    // 2. Calculate price
    price, err := s.priceService.CalculatePrice(req.CourtID, req.Date, req.Start, req.End)
    if err != nil { ... }

    // 3. Create booking with status 'locked' and lock_expiry = now + 10 min
    booking := &models.Booking{
        CourtID:     req.CourtID,
        UserID:      currentUser.ID,
        BookingDate: req.Date,
        StartTime:   req.Start,
        EndTime:     req.End,
        Status:      "locked",
        LockExpiry:  time.Now().Add(10 * time.Minute),
        TotalPrice:  price,
    }
    if err := s.repo.CreateBooking(tx, booking); err != nil { ... }

    // 4. Commit transaction
    tx.Commit()

    // 5. Optionally start a goroutine to release lock if not paid within 10 min
    go s.scheduleLockRelease(booking.ID, booking.LockExpiry)

    return booking, nil
}
```

**Lock Release Worker:**
```go
func (s *BookingService) scheduleLockRelease(bookingID uuid.UUID, expiry time.Time) {
    time.Sleep(time.Until(expiry))
    // Check if still locked; if yes, release it
    s.repo.UpdateBookingStatusIfLocked(bookingID, "expired")
}
```

### 7.3 Payment Callback Handling

```go
// service/payment.go
func (s *PaymentService) HandleCallback(ctx context.Context, payload PaymentCallback) error {
    // Idempotency: check if already processed (store callback ID)
    if s.idempotencyChecker.Exists(payload.TransactionID) {
        return nil
    }

    // Find booking by reference (order_id)
    booking, err := s.repo.FindBookingByCode(payload.OrderID)
    if err != nil { ... }

    // Validate that booking status is 'locked' or 'pending'
    if booking.Status != "locked" && booking.Status != "pending" {
        return ErrInvalidStatus
    }

    // Start transaction
    tx := s.repo.BeginTx(ctx)

    // Update booking status
    if payload.Status == "settlement" {
        booking.Status = "paid"
    } else {
        booking.Status = "failed"
    }
    if err := s.repo.UpdateBooking(tx, booking); err != nil { ... }

    // Create payment record
    payment := &models.Payment{
        BookingID:         booking.ID,
        PaymentMethod:     payload.PaymentMethod,
        ExternalReference: payload.TransactionID,
        Status:            payload.Status,
        Amount:            payload.Amount,
        PaidAt:            time.Now(),
    }
    if err := s.repo.CreatePayment(tx, payment); err != nil { ... }

    tx.Commit()

    // Send notifications
    if booking.Status == "paid" {
        s.notificationService.SendBookingConfirmation(booking.UserID, booking.BookingCode)
    }

    return nil
}
```

---

## 8. Non‑Functional Requirements

| Aspect          | Requirement                                                                 |
|-----------------|-----------------------------------------------------------------------------|
| **Performance** | Booking end‑to‑end < 2 minutes; API response < 200ms for 95% of requests.  |
| **Concurrency** | Support 100 concurrent users during peak hours with zero double bookings.   |
| **Security**    | JWT authentication, rate limiting, SQL injection prevention, PCI‑DSS compliant payment integration. |
| **Reliability** | 99.9% uptime; automatic retries for failed callbacks; idempotent webhooks. |
| **Scalability** | Horizontal scaling of backend; Redis for distributed locks if needed.       |
| **Backup**      | Daily automated PostgreSQL backups.                                         |

---

## 9. Implementation Roadmap

### Phase 1: Foundation (Week 1‑2)
- Set up Go project with Clean Architecture structure.
- Implement authentication (JWT, OTP verification).
- Create base repositories with GORM.
- Svelte project setup, Tailwind, basic routing.
- Court & pricing CRUD (admin).

### Phase 2: Booking Core (Week 3‑4)
- Availability calculation service.
- Booking creation with lock logic.
- Price calculation with dynamic rules.
- Admin override endpoints.

### Phase 3: Payment & Notifications (Week 5‑6)
- Integrate Midtrans/Xendit.
- Payment callback handler with idempotency.
- WhatsApp/Email notifications.
- Lock expiry worker.

### Phase 4: Admin Dashboard & Reports (Week 7‑8)
- Admin dashboard with charts.
- Revenue and utilization reports.
- Export functionality.
- Audit logs.

### Phase 5: Testing & Deployment (Week 9‑10)
- Unit tests (service layer).
- Integration tests (API + DB).
- Load testing.
- Dockerize and deploy to Railway + Vercel.
- Monitoring setup.

---

## 10. Success Metrics (KPIs)

| Metric                   | Target            |
|--------------------------|-------------------|
| Conversion rate          | >70%              |
| Utilization rate         | >80%              |
| Average booking time     | <2 minutes        |
| Monthly revenue growth   | Positive trend    |
| Zero double booking incidents | 100%           |

---

## 11. Glossary

| Term               | Definition                                                                 |
|--------------------|----------------------------------------------------------------------------|
| **Dynamic Pricing**| Price variation based on day type (weekday/weekend) and time of day.      |
| **Temporary Lock** | A 10‑minute hold on a slot during checkout to prevent double booking.     |
| **Idempotency**    | Ensuring that processing a request multiple times has the same effect as once. |
| **Utilization Rate**| Percentage of booked slots vs total available slots over a period.        |
| **Clean Architecture**| A design pattern that separates concerns into layers, keeping business logic independent of frameworks. |

---

**Document Version:** 1.0  
**Last Updated:** 2026‑03‑21