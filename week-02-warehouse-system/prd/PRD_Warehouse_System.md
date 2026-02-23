# Product Requirements Document (PRD)
## Sistem Manajemen Gudang (Warehouse Management System)

**Versi:** 1.0  
**Target User:** Startup & UMKM (1-5 pengguna)  

---

## 1. Executive Summary

Sistem Manajemen Gudang adalah aplikasi berbasis web/mobile yang dirancang untuk membantu UMKM dan startup mengelola operasional gudang dengan efisien. Sistem ini memungkinkan pengguna untuk mengelola data pengguna, memantau transaksi barang masuk/keluar, mengelola stok barang, dan menghasilkan laporan analytics real-time.

### Tujuan Utama
- Mengotomasi proses inventory management
- Mengurangi kesalahan input manual
- Memberikan visibility real-time terhadap stok barang
- Menyediakan laporan actionable untuk decision making

---

## 2. Deskripsi Produk

### 2.1 Visi Produk
Menjadi solusi gudang yang user-friendly, terjangkau, dan dapat dengan mudah diadopsi oleh UMKM lokal tanpa memerlukan pelatihan teknis yang rumit.

### 2.2 Target Market
- UMKM retail, distribusi, dan trading
- Startup e-commerce dengan operasional gudang kecil
- Usaha keluarga yang ingin digitalisasi inventory

### 2.3 Problem Statement
- Kesulitan tracking stok barang secara akurat
- Tidak ada visibility real-time terhadap pergerakan barang
- Laporan manual yang time-consuming dan error-prone
- Tidak ada kontrol akses untuk berbagai level pengguna

---

## 3. Fitur-Fitur Utama

### 3.1 Module User Management

#### 3.1.1 User Registration & Authentication
- **Register Account Baru**
  - Form registrasi dengan input: email, password, nama lengkap, nomor telepon
  - Validasi email unique dan password strength (minimal 8 karakter, kombinasi huruf & angka)
  - Email verification untuk konfirmasi akun
  - Terms & Conditions acceptance

- **Login & Session Management**
  - Login dengan email dan password
  - Remember me functionality (max 30 hari)
  - Session timeout setelah 30 menit inaktif
  - Logout dengan konfirmasi

- **Password Management**
  - Forgot password dengan email reset link
  - Change password (memerlukan password lama)
  - Password expiry notification (3 bulan sekali)

#### 3.1.2 User Profile Management
- View dan edit profil pengguna (nama, email, telepon, foto)
- Update informasi warehouse (nama gudang, alamat, lokasi)
- Preference settings (bahasa, timezone, notifikasi)

#### 3.1.3 Role & Permission Management
- **Tiga level role:**
  1. **Admin** - Full access ke semua fitur, dapat mengelola user, konfigurasi sistem
  2. **Manager** - Akses laporan, membuat transaksi, view stok
  3. **Operator** - Hanya input transaksi barang dan lihat stok dasar

- **Permission Matrix:**
  - Create, Read, Update, Delete untuk setiap fitur berdasarkan role
  - Activity logging untuk setiap action pengguna

---

### 3.2 Module Transaksi Barang

#### 3.2.1 Tipe Transaksi
1. **Barang Masuk (Inbound)**
   - Purchase Order dari supplier
   - Return dari customer
   - Retur/kompensasi dari pihak lain

2. **Barang Keluar (Outbound)**
   - Penjualan ke customer
   - Transfer ke warehouse lain
   - Rusak/waste barang

#### 3.2.2 Create Transaksi Barang Masuk
- Input form dengan field:
  - Tanggal transaksi (default: hari ini)
  - Tipe transaksi (dropdown: PO, Return, Lainnya)
  - Nomor referensi (PO number, invoice, dll)
  - Supplier/Pengirim
  - List barang (SKU, quantity, harga satuan)
  - Catatan/notes
  - Foto bukti (optional)

- Validasi:
  - SKU harus ada di master barang
  - Quantity harus positif
  - Field mandatory harus terisi

- Action:
  - Save sebagai draft
  - Submit/approve untuk finalisasi
  - Auto-update stok barang saat transaksi approve

#### 3.2.3 Create Transaksi Barang Keluar
- Input form serupa dengan barang masuk
- Field tambahan:
  - Tujuan/customer
  - Nomor invoice penjualan (jika ada)
  
- Validasi tambahan:
  - Quantity keluar tidak boleh melebihi stok available
  - Warning jika stok mendekati minimum level

- Action:
  - Save draft
  - Submit/confirm
  - Bisa dibatalkan jika masih draft
  - Auto-update stok saat approve

#### 3.2.4 Transaction History & Management
- List semua transaksi dengan filter:
  - Range tanggal
  - Tipe transaksi
  - Status (draft, approved, rejected, cancelled)
  - Barang tertentu
  
- View detail transaksi:
  - Semua informasi transaksi
  - Approval history
  - User yang melakukan action
  
- Edit transaksi (hanya draft):
  - Update field-field
  - Delete draft transaksi

- Approval Workflow (opsional untuk startup):
  - Operator: Submit transaksi
  - Manager/Admin: Review & approve/reject
  - Email notification untuk approval

---

### 3.3 Module Stok Barang

#### 3.3.1 Master Barang Management
- **Create Barang Baru:**
  - Input: SKU (unique), nama barang, deskripsi, kategori
  - Unit of measurement (pcs, box, kg, liter, dll)
  - Harga beli, harga jual
  - Supplier default
  - Minimum stock level
  - Maksimum stock level
  - Photo barang (optional)
  
- **Edit Barang:**
  - Update informasi barang
  - Disable barang jika tidak aktif (soft delete)
  
- **List Barang:**
  - Tabel dengan column: SKU, nama, kategori, stok sekarang, min level, max level
  - Search by SKU atau nama
  - Filter by kategori
  - Sort by berbagai column

#### 3.3.2 Real-time Stock Tracking
- **Stok Overview:**
  - Total barang dalam gudang
  - Barang dengan stok rendah (di bawah minimum)
  - Barang dengan stok tinggi (di atas maksimum)
  - Barang expired/mendekati expired (jika ada date tracking)

- **Stok by Item:**
  - Detail stok setiap barang
  - Jumlah sekarang
  - Mutasi 7 hari terakhir (in/out)
  - Location/rack (jika ada)
  - Unit satuan
  - Tanggal input pertama

#### 3.3.3 Stock Adjustment & Opname
- **Stok Adjustment (untuk koreksi):**
  - Input: SKU, quantity adjustment (+ atau -), alasan
  - Catatan/notes
  - Attachment (foto, bukti)
  - Approval workflow

- **Opname Barang:**
  - Create opname session
  - List semua barang untuk di-count
  - Input physical count per barang
  - Auto-compare dengan sistem stok
  - Highlight discrepancies
  - Generate opname report
  - Approve untuk update sistem stok

#### 3.3.4 Stok Movement History
- Lihat riwayat pergerakan stok per barang
- Filter by tanggal, tipe transaksi
- Detail: tanggal, tipe, quantity, user, reference

---

### 3.4 Module Laporan & Analytics

#### 3.4.1 Laporan Stok
- **Stok Summary:**
  - Total item dalam gudang
  - Total nilai stok (qty × harga)
  - Barang dengan stok minimal
  - Barang dengan stok berlebih

- **Stok Detail Report:**
  - Table: SKU, nama barang, kategori, qty, harga satuan, total nilai
  - Export ke PDF/Excel
  - Print-friendly format

#### 3.4.2 Laporan Transaksi
- **Rangkuman Transaksi:**
  - Total barang masuk (qty, nilai)
  - Total barang keluar (qty, nilai)
  - Perbandingan periode (hari ini vs minggu lalu, vs bulan lalu)

- **Detail Transaksi:**
  - List lengkap transaksi dengan filter
  - Export ke PDF/Excel
  - Breakdown by tipe transaksi

#### 3.4.3 Laporan Pergerakan Stok
- **Stock Movement Report:**
  - Top 10 barang dengan pergerakan tertinggi
  - Chart: trend pergerakan barang per hari/minggu
  - Barang paling banyak keluar (best seller)
  - Barang paling banyak masuk

#### 3.4.4 Dashboard Analytics
- **Key Metrics:**
  - Total stok value (IDR)
  - Jumlah transaksi hari ini
  - Barang dengan stok kritis (baru ditampilkan)
  - Avg. turnover rate per barang

- **Charts & Visualizations:**
  - Bar chart: Top 5 barang by movement
  - Pie chart: Distribusi stok by kategori
  - Line chart: Trend stok value per minggu
  - Line chart: In/Out transactions per hari

- **Period Comparison:**
  - Filter by periode: hari, minggu, bulan
  - Bandingkan dengan periode sebelumnya
  - Hitung growth/decline percentage

#### 3.4.5 Export & Sharing
- Export report ke format PDF, Excel
- Schedule automatic report (email setiap minggu/bulan)
- Print report
- Share report link (read-only)

---

## 4. User Workflows

### 4.1 Workflow Barang Masuk
1. Operator menerima barang dari supplier
2. Operator buka aplikasi → Transaksi Barang → Barang Masuk
3. Isi form: tanggal, supplier, list barang (SKU, qty, harga)
4. Submit transaksi (auto-approve untuk startup kecil, atau perlu approval)
5. Sistem auto-update stok barang
6. Manager bisa view laporan stok terbaru

### 4.2 Workflow Barang Keluar
1. Customer/order masuk
2. Operator cek stok available
3. Operator input transaksi barang keluar (SKU, qty, tujuan)
4. Sistem validasi stok cukup
5. Submit transaksi
6. Stok otomatis berkurang
7. Laporan penjualan ter-update

### 4.3 Workflow Opname Stok
1. Manager inisiasi opname bulanan
2. Cetak list barang atau buka di mobile
3. Operator/team melakukan physical count
4. Input hasil count di aplikasi
5. Sistem compare dengan stok sistem
6. Manager review discrepancies
7. Approve adjustment untuk update stok sistem

### 4.4 Workflow Melihat Laporan
1. Manager login
2. Buka menu Laporan → pilih tipe laporan
3. Filter periode/kategori (optional)
4. View dashboard atau export ke PDF/Excel
5. Share atau print report

---

## 5. Non-Functional Requirements

### 5.1 Performance
- Halaman dashboard load dalam < 3 detik
- Transaksi submit dalam < 5 detik
- Report generate dalam < 10 detik

### 5.2 Security
- HTTPS untuk semua komunikasi
- Password encryption (bcrypt/argon2)
- Session token dengan expiry
- Audit log untuk setiap action (create, update, delete)
- Data backup otomatis (daily)

### 5.3 Scalability
- Dapat menampung 100,000+ barang
- Dapat menampung 1,000,000+ transaksi
- Ready untuk multi-warehouse (future)

### 5.4 Availability
- 99.5% uptime SLA
- Support offline mode untuk input transaksi (sync saat online)
- Graceful degradation jika ada fitur yang down

### 5.5 Usability
- Mobile responsive design
- Clear navigation
- Minimal 3 clicks untuk setiap action
- Tooltip & help text untuk setiap form field
- Dark/light mode toggle

---

## 6. User Interface Specifications

### 6.1 Layout Structure
- **Header:** Logo, user info, logout, settings
- **Sidebar:** Navigation menu (collapse-able untuk mobile)
- **Main Content:** Dynamic berdasarkan menu yang dipilih
- **Footer:** Copyright, version info (optional)

### 6.2 Main Menu Items
- Dashboard
- Transaksi Barang
  - Barang Masuk
  - Barang Keluar
  - History
- Inventory
  - Master Barang
  - Stok Real-time
  - Stok Adjustment
- Laporan
  - Stok Report
  - Transaksi Report
  - Analytics Dashboard
- User Management (Admin only)
- Settings

### 6.3 Design Principles
- Minimalist design
- Consistent color scheme (brand color + neutral)
- Typography: Clear, readable fonts (minimum 12px)
- Icons untuk quick recognition
- Empty states dengan call-to-action
- Error messages yang helpful

---

## 7. Technical Stack (Rekomendasi)

### Frontend
- **Framework:** React.js
- **State Management:** Redux
- **Styling:** Tailwind CSS
- **Charts:** Chart.js

### Backend
- **Framework:** Golang Fiber
- **Database:** PostgreSQL
- **Authentication:** JWT
- **File Storage:** Google Cloud Storage
- **Caching:** Redis (optional, untuk performance)

### Infrastructure
- **Hosting:** DigitalOcean
- **Containerization:** Docker
- **CI/CD:** GitHub Actions
- **Monitoring:** DataDog

---

## 8. Implementation Roadmap

### Phase 1: MVP (Bulan 1-2)
- [ ] User Management (registration, login, roles)
- [ ] Master Barang (create, read, update, list)
- [ ] Transaksi Barang Masuk
- [ ] Transaksi Barang Keluar
- [ ] Real-time Stock Tracking
- [ ] Basic Dashboard

### Phase 2: Enhancement (Bulan 3)
- [ ] Laporan & Analytics
- [ ] Stock Adjustment & Opname
- [ ] Approval Workflow
- [ ] Email Notifications

### Phase 3: Optimization (Bulan 4)
- [ ] Mobile App (React Native / Flutter)
- [ ] Performance Optimization
- [ ] Advanced Analytics
- [ ] Integrasi dengan sistem POS/Akuntansi (optional)

### Phase 4: Scale-Up (Bulan 5+)
- [ ] Multi-warehouse support
- [ ] Advanced forecasting (AI/ML)
- [ ] Supplier management
- [ ] API untuk integrasi pihak ketiga

---

## 9. Success Metrics

### Adoption Metrics
- Jumlah pengguna aktif per minggu
- User retention rate (30-hari, 90-hari)
- Feature adoption rate per fitur

### Operational Metrics
- Rata-rata transaksi per hari
- Stok accuracy vs physical count
- Time to process transaksi (reduction %)

### Business Metrics
- User satisfaction score (NPS)
- Feature usage distribution
- Error/discrepancy rate
- Cost savings (manual process reduction)

---

## 10. Constraints & Assumptions

### Constraints
- Budget terbatas untuk startup → Simple tech stack
- Team developer 1-3 orang → Focus pada MVP dulu
- Server resources minimal → Optimize database queries

### Assumptions
- User sudah familiar dengan basic computer/mobile
- Internet connection stabil (minimum 1Mbps)
- Data barang sudah tersedia (tidak perlu migrate dari sistem lama)
- Warehouse single location (expand later)

---

## 11. Glossary & Terminology

| Term | Definisi |
|------|----------|
| **SKU** | Stock Keeping Unit - kode unik untuk setiap barang |
| **Stok** | Jumlah barang dalam gudang |
| **Opname** | Physical count barang di gudang |
| **Inbound** | Barang masuk ke gudang |
| **Outbound** | Barang keluar dari gudang |
| **Draft** | Status transaksi yang belum final |
| **Approved** | Status transaksi yang sudah dikonfirmasi |
| **Minimum Level** | Batas minimum stok sebelum perlu re-order |
| **Maksimum Level** | Batas maksimum stok yang boleh disimpan |
| **Turnover** | Perputaran stok (berapa kali stok habis & diisi ulang) |

---

**End of Document**