# Product Requirements Document (PRD) – Aplikasi Pencatatan Keuangan Sederhana

**Versi:** 2.0  
**Tanggal:** 4 Mei 2026  
**Status:** Final untuk Development Kickoff  
**Last Updated:** 4 Mei 2026 | Alignment dengan BA Document v2.0

---

## 1. Ringkasan Eksekutif

### 1.1 Definisi Produk
Aplikasi Pencatatan Keuangan Sederhana adalah solusi mobile **Android-first** yang dirancang dengan **mindset "simplicity & ease-to-use"** untuk membantu individu (Bu Siti, Pak Budi, Andi) dan UMKM mencatat pemasukan/pengeluaran dengan cara paling **mudah, cepat, dan jelas**.

**Filosofi Desain**: Lebih sedikit fitur = lebih banyak pengguna. Menggunakan **"Napkin Test"** – jika fitur tidak bisa dijelaskan di napkin, tidak include di MVP.

### 1.2 Value Proposition
| Aspek | Benefit |
|-------|---------|
| **⚡ Speed** | Input hanya 3 klik, hemat 30 detik per transaksi vs mencatat manual |
| **🎯 Simplicity** | Tidak ada setup rumit, langsung pakai di hari pertama (zero learning curve) |
| **📊 Clarity** | Dashboard visual menunjukkan financial health dalam 1 layar |
| **🔐 Trust** | Data lokal, tidak ada internet requirement, sync otomatis v1.1+ |

### 1.3 Target User Segmentation
- 🏠 **Ibu Rumah Tangga (Bu Siti, 35 thn)**: Catat belanja harian, track pengeluaran keluarga (2-3x/hari)
- 🏪 **Pemilik UMKM (Pak Budi, 45 thn)**: Pisah cash flow bisnis vs personal (5-10x/hari)
- 💼 **Freelancer (Andi, 28 thn)**: Monitor income project vs expense operasional (1-2x/hari)
- 💰 **Karyawan Gaji Tetap (Sinta, 42 thn)**: Hemat 10% tiap bulan dengan budget control (1x/hari)

### 1.4 Key Principles
✅ Input transaksi < 5 detik (3 klik saja)  
✅ Dashboard one-screen yang comprehensive  
✅ Zero learning curve – tidak butuh training  
✅ Data aman & offline-first (SQLite lokal)  
✅ Strict MVP scope (defer everything to v1.1+)

---

## 2. Tujuan Bisnis dan Keberhasilan Produk

### 2.1 Tujuan Bisnis (Business Objectives)
1. **Speed**: Pengguna mencatat transaksi dalam **< 5 detik** (3 klik maksimal)
2. **Simplicity**: UI minimalis, kategori preset, auto-save draft → zero setup time
3. **Clarity**: Dashboard one-screen menunjukkan balance, pie chart, recent transactions
4. **Trust**: Data lokal (SQLite), offline-first, optional cloud backup v1.1+
5. **Adoption**: Zero learning curve, tidak perlu tutorial/onboarding panjang

**Go/No-Go Criteria**: Jika KPI tidak tercapai di MVP → fix before launch → retest

### 2.2 Key Performance Indicators (KPI - Tri-dimensional)

#### **Engagement Metrics** (Learning if app is sticky)
| Metric | Target | Definition |
|--------|--------|-----------|
| **DAU/MAU** | DAU > 40% of MAU | Daily vs Monthly active users |
| **Retention D7** | > 60% | % users active on day 7 after install |
| **Retention D30** | > 50% | % users active on day 30 after install |
| **Session Duration** | < 2 min/day | Users shouldn't overspend time (simplicity win) |
| **Input Frequency** | > 3 txn/day | Average daily transaction volume |

#### **Quality Metrics** (Learning if app is stable & reliable)
| Metric | Target | Definition |
|--------|--------|-----------|
| **Crash Rate** | < 0.5% | % of sessions with crash/error |
| **Performance** | Input <500ms, Dashboard <1s | Speed measurements |
| **Data Integrity** | 0% data loss | No transaction corruption |

#### **Satisfaction Metrics** (Learning if app is loved)
| Metric | Target | Definition |
|--------|--------|-----------|
| **NPS Score** | > 8/10 | "Would you recommend to friend?" |
| **Store Rating** | > 4.5/5 ⭐ | Play Store app review rating |
| **Support Tickets** | < 10% of users | Sign of UX problems if high |

---

## 3. Lingkup Produk (MVP - STRICT SCOPE)

### 3.1 MVP Features (MUST HAVE - v1.0 ONLY)

| Kode | Fitur | Description | Priority |
|------|-------|-------------|----------|
| **FR-001** | **Input Transaksi 3-Klik** | Jenis (Keluar/Masuk) → Kategori (icon grid) → Nominal → Simpan. Total waktu <15s. | 🔴 P0 |
| **FR-002** | **Dashboard One-Screen** | Balance bulan ini, ↓/↑ vs bulan lalu, Pie chart (top 5 categories), CTA buttons [+Keluar] [+Masuk] [Riwayat]. | 🔴 P0 |
| **FR-003** | **History & Edit** | List transaksi (newest first), swipe left = edit/delete, undo 5 min. Tab filter (Semua/Keluar/Masuk). | 🔴 P0 |
| **FR-004** | **Monthly Report** | Navigasi bulan (<Mei 2026>), pie chart, summary table (Category\|Amount\|%), vs last month comparison. | 🔴 P0 |
| **FR-005** | **Local SQLite Storage** | Offline-first: 100% berfungsi tanpa internet, auto-save draft, fast queries. | 🔴 P0 |
| **FR-006** | **Default Categories** | 10 expense + 5 income preset (fixed, not editable in MVP). Icons: 🍚🛒🚗⚡🎬💊✈️🏠🛠️📚 | 🟢 P1 |
| **FR-007** | **Single Wallet** | Default "Dompet Saya" auto-created. No multi-wallet until v1.1. | 🟢 P1 |

### 3.2 Nice-to-Have (DEFER to v1.1+)

| Feature | Reason | Timeline |
|---------|--------|----------|
| Multiple wallets | Adds complexity, MVP users can manage single wallet | v1.1 |
| Cloud sync (Firebase) | Local SQLite sufficient, sync optional | v1.1+ |
| Edit category names | Preset categories clear enough for MVP | v1.1+ |
| Search/filter advanced | Tab filter sufficient for MVP scope | v1.1+ |
| Export PDF/Excel | Manual backup optional, defer | v1.1+ |
| Budget alerts | Not MVP priority, defer | v1.2+ |
| Recurring transactions | Can input manually, defer | v1.1+ |
| PIN/Biometric lock | Security features, defer | v1.1+ |

### 3.3 Out of Scope (NOT for ANY version 1.x)

❌ **Corporate accounting** – Use ERPNext, SAP instead  
❌ **Tax compliance** – Too complex for personal app  
❌ **Investment portfolio** – Different use case  
❌ **Currency exchange** – Indonesia only, IDR only for MVP  
❌ **Bill payment integration** – Not personal finance tracker  
❌ **Multi-user family account** – Single user app for MVP  
❌ **Bank auto-import** – OCR integration too complex  
❌ **Cryptocurrency** – Beyond scope  

**Strict Rule**: Any feature request beyond FR-001 to FR-007 = DEFER, do not add to MVP

---

## 4. User Personas & Journey

### 4.1 Primary Personas (Research-backed)

| Persona | Profile | Main Goal | Pain Point | Frequency | Tech Comfort |
|---------|---------|-----------|-----------|-----------|-------------|
| 🏠 **Bu Siti** | Ibu RT, 35y, 2 anak | Track family daily spending | Buku catatan mudah hilang | 2-3x/day | Low (old Android) |
| 🏪 **Pak Budi** | Warung owner, 45y, cash flow | Separate business vs personal | Excel susah di HP | 5-10x/day | Low-Medium |
| 💼 **Andi** | Freelancer, 28y, multi-project | Monitor income vs expense | Multi-client tracking hard | 1-2x/day | Medium-High |
| 💰 **Sinta** | Employee, 42y, fixed salary | Budget discipline (save 10%) | Lupa-lupa pengeluaran | 1x/day | Medium |

**Persona Selection Criteria**:
- All are price-sensitive (prefer free/cheap)
- Low tech literacy (need simple UI)
- Cash-heavy (mostly Rp transactions)
- Daily financial concerns

### 4.2 Personas NOT in Scope
❌ Corporate accountants  
❌ Crypto traders  
❌ High-net-worth individuals  
❌ Teenagers

### 4.3 User Role (v1.0 MVP)
- **Single User Role**: All users have same permissions (add/view/edit/delete own transactions)
- **No access control** in MVP (personal device, personal data)
- **Multi-user support** deferred to v2.0+ (family account feature)

---

## 5. Deskripsi Fungsional Fitur

### 5.1 Input Transaksi Cepat (FR-001)

**Alur pengguna:**  
1. Dari dasbor, tekan tombol **+** (tambah).  
2. Pilih jenis transaksi: **Keluar** (pengeluaran) atau **Masuk** (pemasukan).  
3. Pilih kategori dari kisi ikon yang telah ditentukan.  
4. Masukkan nominal (keyboard numerik).  
5. (Opsional) tambahkan catatan.  
6. Tekan **Simpan** → kembali ke dasbor dengan saldo terbaru.

**Aturan bisnis:**  
- Nominal harus lebih dari 0.  
- Jika pengguna lupa memilih kategori, sistem akan menggunakan kategori **Lainnya** secara otomatis.  
- Waktu input diukur dari layar dasbor hingga transaksi tersimpan.

### 5.2 Dasbor Satu Layar (FR-002)

**Elemen yang ditampilkan:**  
- Saldo bulan berjalan (besar dan mencolok).  
- Perubahan dibanding bulan lalu (misal: `↓ Rp245.000`).  
- Total pemasukan (hijau) dan total pengeluaran (merah).  
- Diagram lingkaran 5 kategori pengeluaran teratas.  
- Tombol aksi: **+ Keluar**, **+ Masuk**, **Riwayat**.

**Interaksi:**  
- Sentuh diagram lingkaran → masuk ke rincian kategori.  
- Geser kiri/kanan pada area bulan → pindah ke bulan sebelumnya/sesudahnya.

### 5.3 Riwayat dan Pengeditan Transaksi (FR-003)

**Fitur:**  
- Daftar transaksi bulan berjalan, diurutkan dari yang terbaru.  
- Sentuh item → lihat detail dan tombol **Edit**.  
- Geser kiri/kanan pada item → opsi **Hapus** (dengan konfirmasi).  
- Tab **Semua | Pemasukan | Pengeluaran** untuk menyaring.

**Aturan bisnis:**  
- Penghapusan dapat dibatalkan dalam waktu 5 menit melalui pesan *undo*.  
- Setiap perubahan akan langsung memperbarui saldo dan diagram.

### 5.4 Laporan Bulanan dan Perbandingan (FR-004)

**Fitur:**  
- Tombol navigasi bulan (`<` Maret 2026 `>`).  
- Ringkasan: total masuk, total keluar, bersih.  
- Diagram lingkaran per kategori (pengeluaran & pemasukan terpisah).  
- Tabel sederhana: `Kategori | Jumlah | Persentase`.  
- Indikator dibanding bulan lalu (misal: `↑ 12%`).

**Aturan bisnis:**  
- Perhitungan persentase kategori:  
  $$\text{% Kategori} = \frac{\text{Jumlah kategori}}{\text{Total periode}} \times 100$$  
- Perbandingan bulan ke bulan:  
  $$\Delta\% = \frac{\text{Nilai bulan ini} - \text{Nilai bulan lalu}}{\text{Nilai bulan lalu}} \times 100$$

### 5.5 Database Schema & Storage (FR-005)

#### Table: `transactions`
```sql
CREATE TABLE transactions (
  id TEXT PRIMARY KEY,              -- UUID
  wallet_id TEXT NOT NULL,          -- FK to wallets
  type TEXT NOT NULL,               -- 'income' or 'expense'
  category_id TEXT NOT NULL,        -- FK to categories
  amount REAL NOT NULL,             -- Currency amount (Rp)
  note TEXT,                        -- Optional description
  date_time DATETIME NOT NULL,      -- When transaction occurred
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  
  FOREIGN KEY(wallet_id) REFERENCES wallets(id),
  FOREIGN KEY(category_id) REFERENCES categories(id)
);
```

#### Table: `wallets`
```sql
CREATE TABLE wallets (
  id TEXT PRIMARY KEY,              -- UUID
  name TEXT NOT NULL,               -- "Dompet Saya" (v1.0 fixed)
  balance REAL DEFAULT 0,           -- Denormalized for speed
  currency TEXT DEFAULT 'IDR',      -- Only IDR in MVP
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Table: `categories` (Seed Data - Fixed in MVP)
```sql
CREATE TABLE categories (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,               -- "Makan", "Belanja", etc
  icon TEXT NOT NULL,               -- "🍚", "🛒", etc (emoji)
  type TEXT NOT NULL,               -- 'income' or 'expense'
  color TEXT,                       -- HEX color for pie chart
  is_active BOOLEAN DEFAULT TRUE
);

-- Seed expense categories (10):
-- Makan, Belanja, Transport, Tagihan, Hiburan, Kesehatan, Liburan, Rumah, Perbaikan, Pendidikan

-- Seed income categories (5):
-- Gaji, Bonus, Jualan, Lainnya
```

**Storage Architecture**:
- **Offline-first**: All data in SQLite local
- **No cloud**: Optional Firebase sync in v1.1+
- **Performance**: Queries optimize for month-based filtering
- **Backup**: Auto-draft save, manual export in v1.1+

---

## 6. Alur Pengguna dan Skenario Penting

### 6.1 First-Time User Experience (Onboarding - Critical)
**Goal**: Masuk dan langsung bisa pakai, tanpa wizard panjang

```
Step 1: Install app [5 detik]
   ↓
Step 2: Welcome screen + Create first wallet "Dompet Saya" [auto-filled, skip OK]
   ↓
Step 3: Demo input 1 transaksi (text "Cobalah: ketik 50000 untuk makan")
   ↓
Step 4: Dashboard dengan contoh data
   ↓
Step 5: Ready! Bisa start input transaksi atau import history (optional)
```
**NO Tutorial Required** → Learning curve = 0 minutes

### 6.2 Skenario: Input Pengeluaran Makan Siang (UC‑001 - Happy Path)

| Langkah | Aksi Pengguna | Respons Sistem |
|---------|---------------|----------------|
| 1 | Buka aplikasi | Menampilkan dasbor dengan saldo terkini |
| 2 | Tekan tombol **+** | Layar input transaksi, show Keluar/Masuk selector |
| 3 | Pilih **Keluar** | Menampilkan grid kategori pengeluaran |
| 4 | Pilih kategori **🍚 Makan** | Kategori terpilih, highlight icon, focus ke kolom nominal |
| 5 | Ketik `50000` | Pratinjau `Rp 50.000` |
| 6 | (Optional) Catatan: "Indomaret" | Text field updated |
| 7 | Tekan **Simpan** | Transaksi tersimpan, kembali ke dasbor, saldo updated real-time |

**Measurement**: Total waktu dari step 1-7 < 15 detik, ideally < 10 detik

### 6.3 Skenario: Mengedit Transaksi yang Salah (UC‑002)

| Langkah | Aksi Pengguna |
|---------|----------------|
| 1 | Pada Dashboard, tap **[📋 Riwayat]** button |
| 2 | Scroll ke transaksi yang salah |
| 3 | Tap item transaksi → expand show detail |
| 4 | Tap **[✏️ Edit]** button |
| 5 | Ubah nominal dari 50000 ke 75000 |
| 6 | Tap **[✓ Update]** |
| 7 | Otomatis kembali history, balance recalculate |

### 6.4 Skenario: Lihat Laporan Bulan Lalu (UC‑003)

| Langkah | Aksi Pengguna |
|---------|----------------|
| 1 | Pada dasbor, tap **[< April 2026 >]** navigation |
| 2 | Dashboard refresh menampilkan data April |
| 3 | Tap pie chart → modal detail kategori April |
| 4 | Close modal, lihat comparison: Mei vs April (Δ%) |
| 5 | Tap **[📊 Laporan]** untuk detail table |

### 6.5 Skenario: Input Pemasukan Gaji (UC‑004)

Sama seperti UC-001, tetapi:
- Step 3: Pilih **Masuk** sebagai jenis transaksi
- Step 4: Kategori yang tersedia: **Gaji**, **Bonus**, **Jualan**, **Lainnya**
- Step 5: Ketik `5000000` (gaji bulanan)

### 6.6 Error Handling & Recovery (Simplicity-First Approach)

| Scenario | Expected Behavior |
|----------|-------------------|
| **User input nominal 0** | Show toast: "Nominal harus lebih dari 0" → Stay in form |
| **User forget kategori** | Auto-default \"Lainnya\" jika tap Save |
| **Offline** | App tetap jalan 100%, show "📡 Offline mode" badge |
| **User delete transaksi** | Show 5-min undo toast before permanently delete |
| **Duplicate input** | Allow (user might have multiple same txn) |
| **App crash mid-input** | Auto-save draft, resume on reopen |

---

## 7. Persyaratan Non‑Fungsional (Comprehensive)

### 7.1 Kinerja & Responsiveness

| Aspek | Target | Rationale |
|-------|--------|-----------|
| **Waktu simpan transaksi** | < 500 ms | User feels instant response |
| **Waktu muat dasbor** | < 1 detik | First screen must be fast |
| **Pengguliran daftar** | 60 fps (halus) | Smooth scrolling, not jerky |
| **Waktu buka aplikasi (cold start)** | < 3 detik | App launch fast |
| **Penggunaan memori** | < 100 MB | Work on low-end Android (2GB RAM) |
| **Ukuran APK** | < 50 MB | Easy download for low-bandwidth users |

### 7.2 Antarmuka dan Pengalaman Pengguna (Accessibility for All)

- Ukuran font minimal **14pt** untuk keterbacaan.
- Tombol sentuh minimal **44x44dp** (thumb-friendly).
- Kontras warna memenuhi standar WCAG AA.
- Mendukung **mode gelap** (dark mode).
- Seluruh teks antarmuka dalam **Bahasa Indonesia** (tidak ada istilah asing di MVP).
- Ikon yang jelas dan recognizable: 🍚 untuk makan, 🛒 untuk belanja, dll.
- **Onboarding**: Max 30 seconds, skip option always available
- **Undo/Redo**: Support for accidental deletions (5-min window)

### 7.3 Platform dan Kompatibilitas

| Requirement | Detail |
|-------------|--------|
| **Target utama** | Android 10 (API 29) ke atas |
| **Cakupan pasar** | >90% perangkat di Indonesia |
| **iOS** | Tidak didukung di MVP (future: v1.1+) |
| **Web (PWA)** | Ditunda ke v2.0 |
| **Mode offline** | 100% berfungsi tanpa koneksi internet |
| **Screen sizes** | Support 4" - 7" smartphones (most common) |

### 7.4 Keamanan dan Privasi

| Requirement | Implementation |
|-------------|-----------------|
| **Data Storage** | SQLite lokal, tidak dikirim ke server |
| **Telemetry** | Zero data collection di MVP |
| **Ads** | Tidak ada di versi gratis MVP |
| **User tracking** | No analytics, no GA, no Firebase tracking |
| **Backup manual** | User dapat export/import (v1.1+) |
| **PIN/Biometric** | Defer to v1.1+ security features |
| **Data wipe** | Option untuk reset all data |

### 7.5 Testing & Quality Assurance

| Type | Criteria | Success |
|------|----------|---------|
| **Unit Testing** | Core calculations, date filtering | >80% code coverage |
| **Integration Testing** | DB CRUD, date navigation | All flows pass |
| **Functional Testing** | All FR-001 to FR-005 working | 100% pass rate |
| **UAT** | 10-15 users dari target persona test 3 hari | NPS > 7 |
| **Performance Testing** | Measure FPS, load time, memory | Hit all targets |
| **Stability** | 0 crashes dalam 1000 user sessions | <0.5% crash rate acceptable |

---

## 8. Batasan dan Asumsi

### 8.1 MVP Constraints (What's NOT in v1.0)
- ❌ Hanya satu dompet (tidak bisa membedakan uang pribadi dan bisnis)
- ❌ Kategori tidak dapat diubah atau ditambah oleh pengguna
- ❌ Tidak ada cadangan awan otomatis
- ❌ Tidak ada fitur pencarian teks
- ❌ Tidak ada notifikasi atau pengingat
- ❌ No multiple devices sync
- ❌ No custom categories
- ❌ No budget limits/alerts

### 8.2 Technical Constraints
- **SQLite Limit**: OK untuk ~10,000 transactions/user (years of data)
- **App Size**: Target < 50MB APK
- **Min Android Version**: API 29 (Android 10)
- **Dev Time**: 4 weeks for MVP (1-3 developers)
- **Backend**: ZERO backend requirement (important for cost!)

### 8.3 Key Assumptions
✓ Target user comfortable dengan smartphone (Android/iOS)  
✓ Consistent internet (for future cloud backup, not required for MVP)  
✓ Users willing to input transactions manually (no automatic bank sync in MVP)  
✓ Currency = IDR only (Indonesia market)  
✓ Single user per device (no family sharing yet)  
✓ Device storage sufficient (SQLite footprint ~1-5MB)

---

## 9. Technology Stack & Development Timeline

### 9.1 Recommended Technology Stack

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| **Frontend** | Flutter 3.10+ | Cross-platform, fast dev, Android/iOS, built-in SQLite |
| **Database** | SQLite | No server needed, offline-first, proven, built-in |
| **Language** | Dart | Flutter native, type-safe, well-optimized |
| **Sync** | Firebase (v1.1+) | Free tier generous, simple API (not for MVP) |
| **Testing** | Flutter Test + Mockito | Good coverage, fast, easy to run |
| **CI/CD** | GitHub Actions | Free for open-source, easy setup |

### 9.2 MVP Development Timeline (4 Minggu)

| Phase | Duration | Key Activities | Deliverable |
|-------|----------|-----------------|-------------|
| **Minggu 1** | 5 hari | Setup Flutter project, design SQLite schema, implement FR-001 (input), basic storage | Input feature working |
| **Minggu 2** | 5 hari | Build FR-002 (dashboard), pie charts, month navigation, edit/delete (FR-003) | Dashboard + History working |
| **Minggu 3** | 5 hari | Complete FR-004 (reports), category breakdowns, internal testing, bug fixes | Reports + all features working |
| **Minggu 4** | 5 hari | UAT dengan target users, polish UI, Play Store preparation | Ready for launch |

**Total**: 4 weeks | **Team**: 1-3 developers | **Owner**: Product Manager

### 9.3 Architecture Principle

```
┌─────────────────────────────┐
│  User Interface (Flutter)   │
│  (Input form, Dashboard,    │
│   History, Month view)      │
├─────────────────────────────┤
│  Business Logic             │
│  (Validation, Calculation,  │
│   Category management)      │
├─────────────────────────────┤
│  SQLite Local Database      │
│  (Transactions, wallets,    │
│   categories)               │
└─────────────────────────────┘
```

**Key Principle**: Frontend-only for MVP (no complex backend)

---

## 10. MVP Success Criteria & Go/No-Go Checklist

### Before Launch Decision:

| Criteria | Status | Evidence Required |
|----------|--------|-------------------|
| ✅ **All 4 features working** | Must Pass | FR-001 to FR-005 complete, zero P0 bugs |
| ✅ **Performance targets met** | Must Pass | Input <500ms, Dashboard <1s on mid-range Android |
| ✅ **UAT passed** | Must Pass | 10 users, 3 days test, NPS > 7, no data loss |
| ✅ **Zero crashes** | Must Pass | 1000+ user sessions, 0 crashes observed |
| ✅ **Indonesian UI** | Must Pass | All text in Bahasa Indonesia, no English |
| ✅ **Accessibility** | Must Pass | Font >14pt, buttons >44dp, dark mode works |
| ✅ **Offline working** | Must Pass | App 100% functional without internet |

**IF any criteria fail → FIX → RETEST before launch**

---

## 11. Risiko dan Mitigasi

| Risiko | Dampak | Mitigasi |
|--------|--------|----------|
| Input masih terasa lambat (>5 detik) | High | UAT awal (minggu ke-2) untuk mengukur waktu; iterasi UI design jika needed |
| Kehilangan data karena kerusakan SQLite | Critical | Test recovery flow, add backup v1.1 segera |
| Pengguna merasa fitur terlalu sedikit | Medium | Clear communication: ini MVP, v1.1 coming soon |
| Scope creep (dev ingin add features) | High | Strict enforcement: DEFER all to v1.1, no exceptions |
| Kompetitor sudah di market | Medium | Focus unique angle: simple UI lokal, zero learning curve |

---

## 12. Model Bisnis dan Monetisasi

### v1.0 - v1.1: 100% Free
- No ads, no premium tier
- Goal: Build user base, gather feedback

### v2.0: Freemium Model Introduction
| Feature | Free | Premium (Rp25rb/tahun) |
|---------|------|----------------------|
| Input transaksi | ✅ Unlimited | ✅ Unlimited |
| Dashboard & Reports | ✅ Basic | ✅ Advanced reports |
| Multiple wallets | ❌ 1 wallet | ✅ Unlimited |
| Cloud backup | ❌ Manual | ✅ Auto-sync |
| Export | ❌ | ✅ PDF/Excel |
| Custom categories | ❌ | ✅ |
| No ads | ✅ | ✅ |

---

## 13. Post-Launch Roadmap

### v1.1 (2 minggu setelah v1.0)
- Multiple wallets
- Cloud sync (Firebase)
- Export PDF/Excel
- Bug fixes dari UAT feedback

### v1.2 (1 bulan setelah v1.1)
- Budget limits & alerts
- Advanced search/filter
- Recurring transactions

### v2.0 (3 bulan setelah v1.0)
- iOS support
- Premium tier & monetization
- Web PWA version
- Family account (multi-user)

---

## 14. Persetujuan

| Peran | Nama | Tanggal | Status |
|-------|------|---------|--------|
| *Product Manager* | [diisi] | 4 Mei 2026 | Pending |
| *Tech Lead* | [diisi] | 4 Mei 2026 | Pending |
| *Business Analyst* | - | 4 Mei 2026 | Draft v2.0 |

---

## 15. Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| v1.0 | May 4, 2026 | PM | Initial PRD (basic structure) |
| v2.0 | May 4, 2026 | BA | Enhanced for clarity, tech stack, detailed specs |
| - | TBD | - | Awaiting team feedback & approval |

---

**Dokumen ini merupakan PRD v2.0 untuk pengembangan MVP Aplikasi Pencatatan Keuangan Sederhana. Alignment penuh dengan BA Document v2.0. Perubahan setelah tanda tangan harus melalui proses Change Management yang ketat.**