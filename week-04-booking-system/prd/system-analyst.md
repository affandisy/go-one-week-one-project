# Dokumentasi Bisnis - Sistem Booking Lapangan Padel (SMA Business Analyst)

Versi 1.0  
Tanggal: 21 Maret 2026  
Author: Business & System Analyst

Ringkasan: Sistem booking lapangan padel yang seamless (booking <2 menit), anti-double booking, dynamic pricing, dan dashboard analitik untuk admin. Dirancang berdasarkan business requirement: zero double-booking, automated state management, user journey lengkap, dan admin control.

## 1. Tujuan Bisnis

- Memungkinkan pelanggan booking lapangan padel secara online dalam <2 menit.
- Mencegah kerugian bisnis seperti double booking dengan concurrency control dan temporary lock.
- Mendukung dynamic pricing (weekdays murah, peak hours mahal, weekend premium).
- Memberikan dashboard analitik untuk pemilik: revenue, utilization rate, dan data pelanggan VIP.
- Integrasi payment gateway dan notifikasi WhatsApp/Email untuk konversi tinggi.

## 2. Stakeholder & User Role

| Role | Deskripsi & Akses Utama |
|---|---|
| Pelanggan | Registrasi/login, cari jadwal, booking, bayar, lihat riwayat. |
| Admin/Manager | CRUD master data (lapangan, harga), override booking, laporan revenue/utilization. |
| Pemilik | Dashboard read-only: analitik pendapatan, pelanggan VIP, metrik utilisasi. |
| Sistem | Automated state management, payment callback, notifikasi. |

## 3. Core Business Requirements (Kriteria Wajib)

### 3.1 Zero Double-Booking (Concurrency Control)

- Temporary lock slot saat checkout (10 menit hold).
- Optimistic locking: cek versi timestamp sebelum commit booking.

### 3.2 Dynamic Pricing

- Harga bervariasi: weekdays siang (murah), malam peak (mahal), weekend (premium).
- Formula: Harga = Base Price × Multiplier (hari/jam).

### 3.3 Automated State Management

| Status Slot | Deskripsi |
|---|---|
| Available | Hijau, bisa dibooking. |
| Locked | Kuning, sedang checkout (10 menit). |
| Pending Payment | Menunggu callback payment gateway. |
| Paid/Success | Abu-abu, terpakai. |
| Expired/Failed | Kembali Available otomatis. |

### 3.4 User Authentication

- Registrasi: Nama, Email/WA, Password + OTP verifikasi.
- Login: Session token untuk multi-device.

## 4. Workflow User Journey (Pelanggan)

### 4.1 High-Level Flowchart (Text)

```text
1. Registrasi/Login -> 2. Pilih Tanggal -> 3. Lihat Lapangan & Slot (Warna: Hijau/Abu/Kuning)
	|
4. Pilih Slot(s) -> 5. Checkout (Lock 10 menit) -> 6. Bayar (QRIS/VA/Kartu)
	|
7. Payment Success -> 8. E-ticket + WA/Email -> 9. Riwayat Dashboard
```

### 4.2 Detail Use Case

UC-01: Booking Lapangan  
Actor: Pelanggan  
Pre-condition: Login, pilih tanggal.

Flow:
1. Tampilkan lapangan (Court A Indoor, Court B Semi-outdoor).
2. Slot per 1/1.5 jam dengan indikator warna.
3. Pilih slot -> Hitung total (dynamic price + pajak).
4. Checkout: Lock slot -> Redirect payment gateway.
5. Callback success: Generate e-ticket, kirim WA/Email.

UC-02: Riwayat Booking  
Flow: Lihat upcoming/past bookings, total spending untuk loyalitas (main 10x gratis 1 jam).

## 5. Workflow Admin Journey

### 5.1 High-Level Flowchart (Text)

```text
1. Login Admin -> 2. CRUD Master (Lapangan/Harga) -> 3. Manajemen Jadwal (Block/Override)
	|
4. Monitor Transaksi (Pending/Refund) -> 5. Generate Laporan -> 6. Analitik Dashboard
```

### 5.2 Detail Use Case

UC-A1: CRUD Master Data
1. Lapangan: Tambah/edit (nama, tipe indoor/outdoor), nonaktifkan maintenance.
2. Pricing: Base price + multiplier per hari/jam.

UC-A2: Operational Control
1. Manual booking (VIP/turnamen).
2. Block time (maintenance/coaching).

UC-A3: Manajemen Transaksi
1. Lihat status semua booking.
2. Refund/reschedule (force majeure).

UC-A4: Laporan & Analitik
1. Revenue harian/mingguan/bulanan.
2. Utilization rate: % Lapangan Terpakai = (Booked Slots / Total Slots) × 100.
3. Data pelanggan: Riwayat, spending untuk marketing VIP.

## 6. Sistem Design (High-Level)

### 6.1 Architecture Overview

```text
Frontend (Web/Mobile): React/Vue -> API Layer
Backend Service: Golang (business logic sepenuhnya di service layer) + GORM
Database: PostgreSQL (hanya untuk operasi CRUD data persisten)
Payment: Midtrans/Xendit (QRIS/VA/Kartu)
Notif: WhatsApp Business API / Email (SendGrid)
Cache: Redis (untuk lock slot concurrency dan cache data yang sering diakses)
```

Prinsip arsitektur:
1. PostgreSQL diposisikan sebagai data store transaksional (CRUD only).
2. Seluruh aturan bisnis (pricing, locking, state transition, validasi booking) dieksekusi di Golang.
3. GORM digunakan sebagai ORM untuk akses data repository, dengan transaction management di layer service/repository.

### 6.2 Database Schema (Core Tables)

| Table | Fields Utama |
|---|---|
| courts | id, name (Court A), type (indoor), status (active/inactive) |
| pricing | id, court_id, day_type (weekday/weekend), time_slot (08:00-09:00), base_price, multiplier |
| bookings | id, user_id, court_id, date, start_time, end_time, status (available/locked/paid), lock_expiry, total_price |
| users | id, name, email, whatsapp, password_hash |
| transactions | id, booking_id, payment_method, status, amount |

### 6.3 Key Algorithms

Dynamic Price Calculation:

$$
	ext{Total Price} = \sum (\text{Base Price} \times \text{Multiplier}) + \text{Pajak (10%)}
$$

Concurrency Lock Logic (implemented di Golang service dengan GORM transaction):

```text
1. User checkout: buat/ubah booking status='locked' dengan lock_expiry=now+10 menit (via GORM).
2. Validasi ketersediaan slot: query booking aktif pada slot yang sama; jika bentrok, tolak request.
3. Payment callback sukses: update booking status='paid/success' jika lock belum expired.
4. Scheduler/worker: ubah status='expired/failed' lalu release slot ketika lock_expiry terlewati.
```

## 7. Non-Functional Requirements

| Aspek | Requirement |
|---|---|
| Performance | Booking <2 menit end-to-end, handle 100 concurrent users peak hour. |
| Scalability | Horizontal scale backend, CDN untuk frontend. |
| Security | JWT auth, rate limiting, SQL injection prevention, PCI-DSS compliant payment gateway. |
| Reliability | 99.9% uptime, auto-backup DB, payment idempotency. |

## 8. Success Metrics (KPI)

- Conversion rate: >70% dari checkout ke paid.
- Utilization rate: >80% lapangan terpakai.
- Average booking time: <2 menit.
- Monthly revenue growth dari online bookings.