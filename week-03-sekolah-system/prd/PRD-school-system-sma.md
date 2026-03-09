# Technical Design Document – School Integration System (SMA)

---

## 1. Introduction

This document outlines the technical architecture and design for the **School Integration System (Jenjang SMA)** based on the provided Business Requirements Document (BRD). The system aims to unify academic, attendance, payroll, and recommendation processes into a single platform, reducing administrative overhead and providing real‑time insights for stakeholders.

### Primary Goals

- Centralise student and teacher data management
- Automate grading, attendance tracking, and salary calculations
- Generate comprehensive reports for school leadership
- Provide data‑driven recommendations for top students

---

## 2. System Architecture (Clean Architecture)

The system follows **Clean Architecture** principles to ensure separation of concerns, testability, and maintainability. All business logic resides in the service layer; the database is used solely for persistence.

### Architecture Diagram

```text
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                         │
│                     (Svelte Frontend)                         │
└───────────────────────────┬─────────────────────────────────┘
                            │ HTTP/JSON
┌───────────────────────────▼─────────────────────────────────┐
│                     Delivery Layer (Handlers)                 │
│                      (Fiber, Middlewares)                      │
├───────────────────────────────────────────────────────────────┤
│                     Usecase / Service Layer                    │
│          (Business logic, validation, orchestration)           │
├───────────────────────────────────────────────────────────────┤
│                     Repository Layer                           │
│              (Database operations, GORM / raw SQL)             │
└───────────────────────────┬─────────────────────────────────┘
                            │ SQL
┌───────────────────────────▼─────────────────────────────────┐
│                     Database (PostgreSQL)                    │
└───────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Description |
|-------|-------------|
| **Entities (Models)** | Plain Go structs representing core business objects (e.g., User, Student, Class, Grade). These are independent of frameworks. |
| **Repository Layer** | Interface‑based data access. Implements CRUD and complex queries using GORM or raw SQL. Each aggregate has its own repository. |
| **Service Layer (Usecase)** | Contains all business rules, validations, and workflows. Orchestrates repositories and external services. No database logic here. |
| **Delivery Layer (Handlers)** | HTTP handlers using Fiber. Parse requests, bind to DTOs, call service methods, and format responses. Includes middleware for authentication, RBAC, logging, and rate limiting. |
| **Presentation Layer** | Svelte single‑page application consuming the REST API. Handles UI/UX, state management, and client‑side validation. |

---

## 3. Technology Stack

| Component | Technology |
|-----------|------------|
| **Backend** | Go 1.22+ with Fiber v2, GORM, github.com/golang-jwt/jwt |
| **Database** | PostgreSQL 15+ with uuid-ossp extension |
| **Frontend** | Svelte (with SvelteKit), Tailwind CSS, Skeleton UI (or custom components) |
| **Auth** | JWT (access + refresh tokens), RBAC middleware |
| **Background** | Native Go goroutines + channels for async tasks (report generation, etc.) |
| **File Storage** | MinIO / S3‑compatible (for student photos, documents) |
| **Caching** | Redis (optional, for rate limiting, idempotency keys) |
| **Deployment** | Docker, Railway (backend), Vercel (frontend) |
| **Monitoring** | Structured logs, Sentry, Prometheus (optional) |

---

## 4. Database Schema

The schema is designed to support all modules while maintaining referential integrity and performance.

### SQL Schema

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==================== Users & Roles ====================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL,
    full_name VARCHAR NOT NULL,
    phone VARCHAR,
    status VARCHAR DEFAULT 'active',            -- active, inactive, suspended
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR UNIQUE NOT NULL,                -- KEPSEK, WAKIL, GURU, MURID, KEUANGAN, ADMIN
    name VARCHAR NOT NULL
);

CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

-- ==================== Core School Data ====================

CREATE TABLE academic_years (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,                        -- e.g., "2025/2026"
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT false
);

CREATE TABLE study_programs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR UNIQUE NOT NULL,                 -- IPA, IPS
    name VARCHAR NOT NULL
);

CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,                         -- e.g., "XII IPA 1"
    study_program_id UUID REFERENCES study_programs(id),
    academic_year_id UUID REFERENCES academic_years(id),
    homeroom_teacher_id UUID REFERENCES users(id)  -- guru wali kelas
);

CREATE TABLE subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    study_program_id UUID REFERENCES study_programs(id)  -- null if common subject
);

-- Teachers assigned to subjects in a specific class
CREATE TABLE teacher_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id UUID REFERENCES users(id),          -- must have role GURU
    class_id UUID REFERENCES classes(id),
    subject_id UUID REFERENCES subjects(id),
    academic_year_id UUID REFERENCES academic_years(id)
);

CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) UNIQUE,      -- each student is a user
    nis VARCHAR UNIQUE NOT NULL,                   -- student ID number
    class_id UUID REFERENCES classes(id),
    enrollment_year_id UUID REFERENCES academic_years(id)
);

-- ==================== Academic (Grades) ====================

CREATE TABLE grade_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,                          -- Harian, UTS, UAS
    weight DECIMAL(3,2) NOT NULL,                   -- e.g., 0.20, 0.30, 0.50
    academic_year_id UUID REFERENCES academic_years(id)
);

CREATE TABLE grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id),
    subject_id UUID REFERENCES subjects(id),
    component_id UUID REFERENCES grade_components(id),
    score DECIMAL(5,2) NOT NULL,
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(student_id, subject_id, component_id)    -- one grade per component per subject per student
);

-- Semester final report (snapshot)
CREATE TABLE report_cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id),
    semester INTEGER NOT NULL,                      -- 1 or 2
    academic_year_id UUID REFERENCES academic_years(id),
    final_score DECIMAL(5,2),                       -- computed average
    pdf_url TEXT,                                    -- generated PDF
    approved_by UUID REFERENCES users(id),           -- kepala sekolah
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ==================== Attendance ====================

CREATE TABLE attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),               -- either teacher or student
    attendance_date DATE NOT NULL,
    status VARCHAR NOT NULL,                          -- hadir, izin, sakit, alfa, (for teachers: terlambat, etc.)
    check_in_time TIMESTAMPTZ,                        -- optional
    notes TEXT,
    recorded_by UUID REFERENCES users(id),            -- who recorded (for student attendance)
    created_at TIMESTAMPTZ DEFAULT NOW(),
    INDEX idx_attendance_user_date (user_id, attendance_date)
);

-- ==================== Payroll ====================

CREATE TABLE salary_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,                             -- Gaji Pokok, Tunjangan, Honor, Potongan
    type VARCHAR NOT NULL CHECK (type IN ('earning', 'deduction'))
);

CREATE TABLE teacher_salary_config (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id UUID REFERENCES users(id) UNIQUE,
    base_salary DECIMAL(12,2) NOT NULL,               -- Gaji pokok
    config JSONB NOT NULL DEFAULT '{}'                  -- Dynamic allowances/deductions rules
);

CREATE TABLE salary_payslips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id UUID REFERENCES users(id),
    month INTEGER NOT NULL,
    year INTEGER NOT NULL,
    total_earnings DECIMAL(12,2) NOT NULL,
    total_deductions DECIMAL(12,2) NOT NULL,
    net_salary DECIMAL(12,2) NOT NULL,
    details JSONB,                                      -- breakdown of components
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ==================== Recommendations ====================

CREATE TABLE recommendation_criteria (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,                              -- Nilai, Kehadiran, Prestasi, Perilaku
    weight DECIMAL(3,2) NOT NULL                        -- sum of weights = 1
);

CREATE TABLE student_recommendations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id),
    semester INTEGER NOT NULL,
    academic_year_id UUID REFERENCES academic_years(id),
    score DECIMAL(5,4) NOT NULL,                         -- computed total score
    rank INTEGER,
    is_selected BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ==================== Announcements ====================

CREATE TABLE announcements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    content TEXT NOT NULL,
    target_roles VARCHAR[] DEFAULT '{}',                  -- e.g., {GURU, MURID}
    target_class_ids UUID[],                              -- specific classes
    target_program_ids UUID[],                            -- specific study programs
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ==================== Audit Logs ====================

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_user_id UUID REFERENCES users(id),
    action VARCHAR NOT NULL,                               -- CREATE, UPDATE, DELETE, LOGIN, etc.
    entity_type VARCHAR NOT NULL,
    entity_id UUID,
    before_data JSONB,
    after_data JSONB,
    ip INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_grades_student ON grades(student_id);
CREATE INDEX idx_attendance_user_date ON attendance(user_id, attendance_date);
CREATE INDEX idx_teacher_assignments ON teacher_assignments(teacher_id, academic_year_id);
```

### Schema Design Notes

- **UUID Primary Keys**: Ensure global uniqueness across distributed systems
- **Foreign Keys**: Enforce referential integrity at database level
- **JSONB Columns**: Store flexible configuration (e.g., salary rules) without extra tables
- **Audit Logs**: Capture all critical changes for compliance and troubleshooting
- **Indexes**: Optimized for common query patterns (student grades, attendance lookups, teacher assignments)

---

## 5. API Specifications (Fiber)

### API Configuration

- **Base URL**: `/api/v1`
- **Content-Type**: `application/json`
- **Auth Header**: `Authorization: Bearer <access_token>`

### Response Conventions

**Success Response:**
```json
{
  "data": { ... },
  "meta": { "page": 1, "limit": 20, "total": 100 }
}
```

**Error Response:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": [ ... ]
  }
}
```

### 5.1 Authentication & User Management

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| POST | `/auth/login` | Login (rate‑limited) | Public |
| POST | `/auth/refresh` | Refresh access token | Public |
| POST | `/auth/logout` | Logout | Auth |
| GET | `/me` | Get current user profile | Auth |
| PATCH | `/me` | Update own profile | Auth |
| GET | `/users` | List users (filter by role, status) | Admin |
| POST | `/users` | Create new user (teacher, student, etc.) | Admin |
| PATCH | `/users/{id}` | Update user | Admin |
| POST | `/users/{id}/roles` | Assign roles | Admin |

### 5.2 Master Data (Academic)

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/academic-years` | List academic years | Admin/TU |
| POST | `/academic-years` | Create academic year | Admin/TU |
| GET | `/study-programs` | List study programs (IPA, IPS) | Admin/TU |
| GET | `/classes` | List classes with filters | All (role‑based) |
| POST | `/classes` | Create new class | Admin/TU |
| GET | `/subjects` | List subjects | All |
| POST | `/subjects` | Create subject | Admin/TU |
| GET | `/teacher-assignments` | Get teacher assignments | All |
| POST | `/teacher-assignments` | Assign teacher to class + subject | Admin/TU |
| GET | `/students` | List students (by class, etc.) | Guru, TU |
| POST | `/students` | Register new student | Admin/TU |

### 5.3 Academic (Grades & Reports)

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/grade-components` | List grading components (with weights) | All |
| POST | `/grade-components` | Create component | Admin/TU |
| GET | `/grades` | Get grades (query: classId, subjectId, componentId) | Guru |
| POST | `/grades` | Enter/update grade for a student | Guru |
| PATCH | `/grades/{id}` | Update grade | Guru |
| GET | `/report-cards` | Get report card (query: studentId, semester) | Murid, Wali, Kepsek |
| POST | `/report-cards/generate` | Trigger report card generation (async) | Kepsek |
| POST | `/report-cards/{id}/approve` | Approve final report | Kepsek |

### 5.4 Attendance

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| POST | `/attendance/teacher` | Record teacher attendance (self) | Guru |
| POST | `/attendance/student` | Record student attendance | Guru (wali/piket) |
| GET | `/attendance` | Get attendance records (query: userId, from, to) | All (own) / Guru |
| GET | `/attendance/summary` | Attendance summary (query: classId, month, year) | Guru, Wali, Kepsek |

### 5.5 Payroll

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/salary-components` | List salary components | Keuangan |
| POST | `/salary-components` | Create salary component | Keuangan |
| GET | `/teacher-salary-config` | Get salary config for teachers | Keuangan |
| POST | `/teacher-salary-config` | Set/update teacher salary config | Keuangan |
| POST | `/salary-payslips/calculate` | Trigger payslip calculation | Keuangan |
| GET | `/salary-payslips` | View payslips (query: teacherId, month, year) | Teacher (own), Keuangan |
| POST | `/salary-payslips/{id}/approve` | Approve payslip | Keuangan |
| GET | `/salary-payslips/{id}/pdf` | Download payslip as PDF | Teacher, Keuangan |

### 5.6 Student Recommendations

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/recommendation-criteria` | List criteria with weights | Kepsek, BK |
| POST | `/recommendation-criteria` | Set criteria weights | Kepsek |
| POST | `/student-recommendations/calculate` | Trigger recommendation calculation | BK |
| GET | `/student-recommendations` | Get ranked students (query: semester, classId) | Kepsek, BK |
| PATCH | `/student-recommendations/{id}/select` | Manually select/deselect a student | Kepsek |

### 5.7 Announcements

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/announcements` | List announcements (filtered by user role) | All |
| POST | `/announcements` | Create announcement | Admin/TU |
| PATCH | `/announcements/{id}` | Update announcement | Admin/TU |
| DELETE | `/announcements/{id}` | Delete announcement | Admin/TU |

### 5.8 Dashboard & Reports

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/dashboard/overview` | Get overview dashboard data | All |
| GET | `/dashboard/kepsek` | Kepala Sekolah dashboard metrics | Kepsek |
| GET | `/dashboard/teacher` | Teacher dashboard (classes, students) | Guru |
| GET | `/dashboard/student` | Student dashboard (grades, attendance) | Murid |
| GET | `/reports/attendance` | Generate attendance report | Admin, Kepsek |
| GET | `/reports/academic` | Generate academic report | Admin, Kepsek |
| GET | `/reports/payroll` | Generate payroll report | Keuangan, Kepsek |

---

## 6. Frontend Overview (Svelte)
- **Framework**: SvelteKit (for routing and server‑side rendering if needed, though primarily SPA).
- **UI Library**: Skeleton UI (or custom Tailwind components) for consistent design.
- **State Management**: Svelte stores (writable/readable) for global state (user, current academic year).
- **Forms**: Svelte forms with built‑in reactivity; validation with svelte-forms-lib or Zod.
- **HTTP Client**: fetch with interceptors (or axios) for JWT handling.
- **Charts**: svelte-chartjs or LayerChart for dashboards.
- **Tables**: svelte-headless-table or custom components for data grids.

**Key Pages & Components:**
- Login page
- Dashboard (role‑based)
- Master data management (classes, subjects, students, teachers)
- Grade entry screen (with class list, subject tabs)
- Attendance recording (for teachers and students)
- Payroll configuration and payslip viewer
- Recommendation ranking list
- Announcements board
- Report viewers (PDF embedded)

**Responsive Design**: Tailwind CSS ensures mobile‑friendly layouts.

---
## 7. Non‑Functional Requirements
- Security:
    - JWT tokens with short expiry (15 min) + refresh tokens.
    - RBAC enforced at both API and UI levels.
    - All passwords hashed with bcrypt.
    - Audit logs for all data modifications.
    - HTTPS only.
- Performance:
    - API response time < 200ms for 95% of requests
    - Database queries optimised with indexes.
    - Asynchronous report generation (background jobs) to avoid blocking.
- Scalability:
    - Designed to handle up to 2000 students, 200 teachers, 500k grades per year.
    - Stateless API allows horizontal scaling.
- Availability: 99.5% uptime target; graceful degradation.
- Maintainability: Clean Architecture ensures separation of concerns; unit tests for service layer.
---

## 8. Implementation Roadmap
#### Phase 1: Foundation (Weeks 1‑2)
- Set up Go project structure (Clean Architecture).
- Implement authentication (JWT) and user management.
- Create basic master data (academic years, study programs, classes, subjects).
- SvelteKit project setup with Tailwind and basic routing.
- Login page and dashboard skeleton.

#### Phase 2: Core Academic Features (Weeks 3‑4)
- Teacher assignment to classes.
- Grade components configuration.
- Grade entry and viewing.
- Attendance recording (teacher self, student by wali kelas).
- Basic reporting (grade summaries).

#### Phase 3: Payroll & Recommendations (Weeks 5‑6)
- Salary components and teacher salary config.
- Automated payslip calculation.
- Recommendation criteria and calculation engine.
- Headmaster dashboard with key metrics.

##### Phase 4: Announcements, Reports & Polish (Weeks 7‑8)
- Announcements module.
- Advanced reports (export to PDF/Excel).
- Audit logging.
- Asynchronous job handling (goroutines).
- Integration tests and load testing.
- Dockerisation and deployment.

---

## 9. Deployment & Operations
- Backend: Docker container deployed on Railway / Fly.io.
- Frontend: Static hosting on Vercel.
- Database: PostgreSQL on Supabase (or managed AWS RDS).
- File Storage: MinIO or Supabase Storage for PDFs and photos.
- Environment Variables: DATABASE_URL, JWT_SECRET, STORAGE_ENDPOINT, etc.
- Monitoring: Structured logs (JSON) → Grafana Loki; Sentry for errors.
- Backup: Automated daily backups via database provider.

---
## 10. Glossary

| Term | Definition |
|------|-----------|
| NIS | Student ID number (Nomor Induk Siswa).|
| Rombel | Classroom group (Rombongan Belajar).|
| Wali Kelas | Homeroom teacher responsible for a class.|
| BK | Guidance and counselling teacher (Bimbingan Konseling). |
| Payslip | Slip gaji – monthly salary statement for teachers.|

---

## 11. Summary

This technical design document provides a comprehensive blueprint for building the School Integration System (SMA) using Clean Architecture principles. The system is designed to be:

- **Scalable**: Clean separation of concerns allows independent scaling of components
- **Maintainable**: Clear layer boundaries make the codebase easy to understand and modify
- **Secure**: RBAC middleware and JWT authentication ensure proper access control
- **Performant**: Strategic indexing and optional Redis caching optimize query performance
- **Extensible**: Interface-based repositories and service layer make adding features straightforward

### Next Steps

1. **Phase 1**: Set up project structure, database schema, and authentication
2. **Phase 2**: Implement core academic and attendance modules
3. **Phase 3**: Develop payroll and recommendation systems
4. **Phase 4**: Build frontend with Svelte and integrate with backend API
5. **Phase 5**: Testing, optimization, and deployment

---

**Document Version**: 1.0  
**Last Updated**: 2026-03-09  
**Status**: Ready for Implementation
