# Business Requirement Document (BRD) – Aplikasi Pencatatan Keuangan Sederhana

**Versi 2.0** | **Tanggal**: 4 Mei 2026 | **Author**: Business & System Analyst  
**Last Updated**: 4 Mei 2026 | **Status**: Draft for Development

---

## 📋 Ringkasan Eksekutif
Aplikasi mobile/web untuk pencatatan keuangan pribadi/UMKM dengan mindset **Simplicity & Easy-to-Use**. 

**Filosofi Design**: Lebih sedikit fitur = lebih banyak pengguna. 80% kebutuhan user tercover dengan 20% fitur core.

**Target Utama**:
- ⚡ Input transaksi dalam <5 detik (3 klik saja)
- 📊 Dashboard one-screen yang comprehensive
- 🚀 Zero learning curve - tidak butuh training/tutorial
- 🔐 Data aman & offline-first

## 1. Executive Summary & Tujuan Bisnis

### 1.1 Problem Statement
Pengguna Indonesia menghadapi masalah pencatatan keuangan manual:
- 📱 Tidak ada tracking real-time pengeluaran
- 📓 Buku/Excel mudah hilang dan tidak praktis untuk mobile
- ❌ Aplikasi existing terlalu kompleks (MMEX, Money Manager) dengan 100+ fitur yang overwhelming
- 🎯 Butuh solusi lokal, bahasa Indonesia, tanpa learning curve

### 1.2 Value Proposition
| Aspek | Benefit |
|-------|---------|
| **Speed** | Input hanya 3 klik, hemat 30 detik per transaksi vs mencatat manual |
| **Simplicity** | Tidak ada setup rumit, langsung pakai di hari pertama |
| **Clarity** | Dashboard visual menunjukkan financial health dalam 1 layar |
| **Trust** | Data lokal, tidak ada internet requirement, sync otomatis |

### 1.3 Target User & Use Cases
1. **Ibu Rumah Tangga**: Catat belanja harian, track pengeluaran keluarga
2. **Pemilik UMKM**: Pisah cash flow bisnis vs personal
3. **Freelancer**: Monitor income project vs expense operasional

### 1.4 Key Success Factors
✅ Input transaksi < 5 detik  
✅ 90% user aktif > 30 hari pertama  
✅ NPS Score > 8/10 (very likely to recommend)  
✅ Zero customer support untuk basic usage

## 2. Stakeholder & User Personas

### 2.1 User Persona Primer
| Persona | Profile | Pain Point | Goal | Frequency |
|---------|---------|-----------|------|-----------|
| **Bu Siti** (35, Ibu Rumah Tangga) | Catat belanja harian dengan anak 2 | Buku catatan selalu hilang/lupa | Track pengeluaran bulanan | 2-3x sehari |
| **Pak Budi** (45, Pemilik Warung) | Jual kebutuhan sehari-hari, cash flow tinggi | Sulit pisah kas bisnis vs pribadi, customer hutang | Lihat profit mingguan | 5-10x sehari |
| **Andi** (28, Freelancer Web) | Kerja project variasi income | Multi-client, income tidak pasti, butuh invoice track | Monitor monthly revenue vs expense | 1-2x sehari |
| **Sinta** (42, Karyawan Gaji Tetap) | Income regular dari kantor | Lupa-lupa pengeluaran, budget yang tidak terstruktur | Hemat 10% tiap bulan | 1x sehari |

### 2.2 User Roles (untuk v1.1+)
- **User Biasa**: Input/lihat transaksi, view balance
- **Power User** (future): Multiple wallet, advanced reporting, export

### 2.3 User Segment NOT dalam Scope
❌ Corporate accounting (use ERPNext, SAP)  
❌ Tax compliance deep dive  
❌ Investment portfolio management

## 3. Core Functional Requirements (MVP SCOPE - Simplicity First)

### 3.1 MVP Principle: "Napkin Test"
Fitur hanya included jika bisa dideskripsikan di napkin (simple). Jika butuh diagram kompleks = tidak include di MVP.

### 3.2 MUST HAVE (v1.0 - MVP)
#### FR-001: Input Transaksi Sederhana
**User Story**: "Sebagai pengguna, saya ingin input pengeluaran dalam 3 klik tanpa ribet"

**Input Types**:
- **Pengeluaran**: Uang yang keluar (belanja, makan, transport, tagihan)
- **Pemasukan**: Uang yang masuk (gaji, bonus, jualan)
- **Transfer**: Antar dompet (optional untuk MVP, v1.1)

**UI Flow**:
```
Home Screen
    ↓
    Tap [+] Button
    ↓
    Choose Type (❌ Keluar | ✅ Masuk)
    ↓
    Pick Category (Icon grid: 🍚 Makan, 🛒 Belanja, 🚗 Transport, etc)
    ↓
    Input Amount (default currency Rp, keyboard numerik)
    ↓
    [Optional] Add Note ("Belanja di Indomaret")
    ↓
    Tap [✓ Simpan]
    ↓
    Back to Dashboard (auto-update balance)
```

**Kategori Preset (Fixed, tidak bisa edit di MVP)**:
| Pengeluaran | Icon | Pemasukan | Icon |
|-------------|------|-----------|------|
| 🍚 Makan | food | 💰 Gaji | salary |
| 🛒 Belanja | shopping | 💼 Bonus | bonus |
| 🚗 Transport | car | 📱 Jualan | sales |
| ⚡ Tagihan | bill | 🎁 Lainnya | gift |
| 🎬 Hiburan | entertainment | | |
| 💊 Kesehatan | health | | |
| ✈️ Liburan | travel | | |
| 🏠 Rumah | home | | |
| 🛠️ Perbaikan | maintenance | | |
| 📚 Pendidikan | education | | |

#### FR-002: Dashboard One-Screen Overview
**User Story**: "Saya ingin tahu kondisi uang saya sekilas tanpa scroll banyak-banyak"

**Dashboard Elements**:
```
┌─────────────────────────────────────┐
│  📊 BULAN INI                       │
├─────────────────────────────────────┤
│  Balance: Rp 5.245.000              │  ← Big, highlight
│  ↓ vs Bulan Lalu: -Rp 245.000       │  ← Trend indicator
├─────────────────────────────────────┤
│  Pemasukan: Rp 10.000.000           │  ← Green
│  Pengeluaran: -Rp 4.755.000         │  ← Red
├─────────────────────────────────────┤
│  Pie Chart (top categories)         │
│  ┌──────────────────────────────┐   │
│  │  🍚 40% Makan                │   │
│  │  🛒 25% Belanja              │   │
│  │  🚗 15% Transport            │   │
│  │  ⚡ 10% Tagihan              │   │
│  │  🎬 10% Lainnya              │   │
│  └──────────────────────────────┘   │
├─────────────────────────────────────┤
│  [+ Keluar] [+ Masuk] [📋 History]  │ ← CTA buttons
└─────────────────────────────────────┘
```

**Interaction**: 
- Tap pie chart = drill-down ke category details
- Swipe kiri/kanan = navigate bulan sebelumnya/sesudahnya
- Pull-to-refresh = sync data

#### FR-003: Transaction History & Edit
**User Story**: "Saya perlu lihat riwayat transaksi dan edit kalau ada kesalahan input"

**Features**:
- Scroll list transaksi bulan berjalan (newest first)
- Tap item = expand untuk lihat detail + edit option
- Swipe left/right = quick delete dengan konfirmasi
- Search by category/date (future: v1.1)
- Filter: "Semua | Pemasukan | Pengeluaran" (tabs)

#### FR-004: Monthly View & Basic Report
**User Story**: "Saya ingin tahu pengeluaran bulan kemarin untuk budgeting bulan ini"

**Features**:
- Month navigation (prev/next button di dashboard)
- Summary: Total In/Out/Net per bulan
- Pie chart breakdown by category
- Simple table: "Category | Amount | %" (sortable by amount desc)
- Compare indicator: ↑/↓ vs bulan lalu

### 3.3 NICE-TO-HAVE (Defer ke v1.1+)
- ⏭️ Multiple wallet (Dompet Harian, Tabungan, Bisnis)
- ⏭️ Edit category name & icon (user custom)
- ⏭️ Export PDF/Excel report
- ⏭️ Search/filter advanced
- ⏭️ Budget limit alert
- ⏭️ Recurring transaction

### 3.4 OUT OF SCOPE (untuk versi future)
- ❌ Bill payment integration
- ❌ Cryptocurrency
- ❌ Stock portfolio
- ❌ Loan calculator
- ❌ Currency conversion
- ❌ Multi-user family account

## 4. User Journey & Workflow

### 4.1 First-Time User Experience (Onboarding - CRITICAL)
**Principle**: Masuk dan langsung bisa pakai, tanpa wizard panjang

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

**NO Tutorial Requirement** (Learning curve = 0 minutes, bukan 30 minutes)

### 4.2 Daily User Flow (Ideal)
```
1. Open App
2. Home = Dashboard (lihat balance, pie chart, recent transaction)
3. Tap [+] → Input 1-3 transaksi harian
4. Swipe/tap history jika mau edit
5. Close app

Total: <2 menit
```

### 4.3 Weekly/Monthly Flow
```
1. Tap [📋 History] atau swipe prev month
2. Check breakdown pie chart
3. Compare dengan bulan lalu (automation)
```

### 4.4 Detailed Use Cases

#### **UC-001: Input Transaksi Cepat** (Happy Path)
**Aktor**: User  
**Precondition**: App terbuka di Dashboard  
**Goal**: Input makan siang 50rb dalam <15 detik

| Step | Action | System Response |
|------|--------|-----------------|
| 1 | Tap [+] button | Show input screen with type selector (Keluar/Masuk) |
| 2 | Select [Keluar] | Default, show category grid |
| 3 | Tap [🍚 Makan] icon | Icon highlight + text "Kategori: Makan" |
| 4 | Type "50000" | Show "Rp 50.000" preview |
| 5 | Tap [✓ Simpan] | Transaction saved → Back to dashboard |
| 6 | See Dashboard update | Balance updated real-time, transaction in history |

**Alternative**: User leave app mid-input → Auto-save draft → Resume on reopen

---

#### **UC-002: Edit/Delete Transaksi Salah**
**Goal**: Fix input amount yang salah

| Step | Action |
|------|--------|
| 1 | Swipe left on transaction in history |
| 2 | Show options: [✏️ Edit] [🗑️ Delete] |
| 3 | Tap [✏️ Edit] → popup edit form |
| 4 | Change amount, tap [✓ Update] |
| 5 | Balance auto-recalculate |

---

#### **UC-003: Check Monthly Report**
**Goal**: Lihat pengeluaran bulan ini vs bulan lalu

| Step | Action |
|------|--------|
| 1 | On Dashboard |
| 2 | Tap [< | Mei 2026 | >] month navigation |
| 3 | Select April 2026 (prev) |
| 4 | Dashboard refresh menunjukkan April data |
| 5 | Compare: Mei vs April (δ percentage displayed) |

---

#### **UC-004: Add Pemasukan**
**Goal**: Input gaji atau bonus

Same as UC-001, but:
- Select [✅ Masuk] instead of [Keluar]
- Categories: Gaji, Bonus, Jualan, Lainnya

---

### 4.5 Error Handling & Recovery (Simplicity Mindset)

| Scenario | Simple Solution |
|----------|-----------------|
| **User input nominal 0** | Show toast: "Nominal harus lebih dari 0" → Stay in form |
| **User forget category** | Default to "Lainnya" if they just tap Save |
| **Offline** | Show "📡 Sync pending" badge, but keep app working |
| **Phone low battery** | Auto-save every 10 detik, no warning needed |
| **User delete wrong transaction** | Undo for 5 minutes, then confirm delete |

## 5. Non-Functional Requirements (Simplicity & Reliability First)

### 5.1 Performance Requirements
| Requirement | Target | Rationale |
|-------------|--------|-----------|
| **Input save latency** | < 500ms | User feels instant response |
| **Dashboard load** | < 1 second | First screen must be fast |
| **List scroll** | 60 FPS | Smooth scrolling, not laggy |
| **Cold start** | < 3 seconds | App launch time |
| **Memory usage** | < 100MB | Work on low-end Android |

### 5.2 UI/UX Requirements (Accessibility for All)
| Requirement | Detail |
|-------------|--------|
| **Font Size** | Min 14pt untuk readability (para lansia) |
| **Button Size** | Min 44x44dp (thumb-friendly) |
| **Color Contrast** | WCAG AA standard minimum |
| **Dark Mode** | Support untuk mata tidak pedih malam hari |
| **Bahasa** | Bahasa Indonesia primary (no English UI text di MVP) |
| **Text** | Clear, simple Bahasa Indonesia (hindari jargon finance) |
| **Icons** | Clear, recognizable (Makan=🍚, bukan generic icons) |

### 5.3 Platform & Compatibility
| Aspect | Requirement |
|--------|-------------|
| **Primary Platform** | Android 10+ (cover 95% user Android Indonesia) |
| **Secondary** | iOS 13+ (future: v1.1) |
| **PWA** | Optional untuk web, defer ke v2.0 |
| **Offline-first** | Harus work 100% tanpa internet |
| **Sync** | Auto-sync when online (no manual action) |

### 5.4 Data & Security (Trust is Everything)
| Requirement | Implementation |
|-------------|-----------------|
| **Data Storage** | SQLite lokal (tidak cloud di MVP) |
| **Backup** | Optional Cloud Backup (Firebase, v1.1) |
| **Encryption** | Optional PIN/Biometric lock (v1.1) |
| **Data Export** | Export as CSV/PDF (future, v1.1) |
| **Privacy** | No telemetry, zero data collection untuk MVP |
| **No Ads** | Clean experience, monetization v1.2+ |

### 5.5 Testing Requirements (Quality Assurance)
| Type | Target |
|------|--------|
| **Functional Testing** | 100% coverage untuk 4 main features |
| **UAT** | 10-15 user dari target persona (Bu Siti, Pak Budi, Andi) |
| **Crash Rate** | < 0.5% di production |
| **Bug Priority** | P0 (crash) fixed dalam 24h, P1 (data loss) dalam 48h |

## 6. Technical Architecture (Simplicity-First Design)

### 6.1 Design Principle: Minimal & Lean
- **Frontend Only** (MVP): UI + Local Storage, tidak butuh backend
- **Database**: SQLite (built-in di Flutter/React Native)
- **Sync**: Optional Firebase (future, v1.1)
- **No complex APIs**: Keep it simple

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

### 6.2 Database Schema (Minimal & Efficient)

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
  name TEXT NOT NULL,               -- "Dompet Saya" (v1.0)
  balance REAL DEFAULT 0,           -- Denormalized for speed
  currency TEXT DEFAULT 'IDR',      -- Only IDR in MVP
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Table: `categories` (Seed Data)
```sql
CREATE TABLE categories (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,               -- "Makan", "Belanja", etc
  icon TEXT NOT NULL,               -- "🍚", "🛒", etc (emoji)
  type TEXT NOT NULL,               -- 'income' or 'expense'
  color TEXT,                       -- HEX color for pie chart
  is_active BOOLEAN DEFAULT TRUE
);
-- Seed 10 expense + 5 income categories
```

### 6.3 Key Calculation Formulas

#### Monthly Balance
$$\text{Balance} = \sum_{\text{type='income'}} \text{amount} - \sum_{\text{type='expense'}} \text{amount}$$

#### By Category Breakdown
$$\text{Category\%} = \frac{\sum_{\text{category}} \text{amount}}{\sum_{\text{total}} \text{amount}} \times 100$$

#### Month-over-Month Comparison
$$\Delta = \text{Current Month} - \text{Previous Month}$$
$$\Delta\% = \frac{\Delta}{\text{Previous Month}} \times 100$$

### 6.4 Data Flow (Offline-First)
```
User Input → Validate (amount > 0, category valid) 
    → Save to SQLite (instant)
    → Update in-memory balance
    → Refresh UI (real-time)
    
Optional (Future): Auto-sync to Firebase when online
```

### 6.5 Technology Stack (Recommended)
| Layer | Technology | Why |
|-------|-----------|-----|
| **Frontend** | Flutter | Fast dev, Android/iOS, built-in SQLite |
| **Database** | SQLite | No server needed, offline-first, proven |
| **Sync** | Firebase (v1.1) | Free tier generous, simple API |
| **Testing** | Flutter Test + Mockito | Good coverage, easy to run |
| **CI/CD** | GitHub Actions | Free for open-source |

## 7. Success Metrics & OKR (Objectives & Key Results)

### 7.1 Key Performance Indicators (KPI)

#### **Engagement Metrics** (Learning if app is sticky)
| Metric | Target | Definition |
|--------|--------|-----------|
| **DAU/MAU** | DAU > 40% of MAU | Daily Active / Monthly Active users |
| **Retention D7** | > 60% | % users active on day 7 after install |
| **Retention D30** | > 50% | % users active on day 30 after install |
| **Session Duration** | < 2 min/day | Users shouldn't spend too long (simplicity win) |
| **Input/Day** | > 3 transactions | Average user habit |

#### **Quality Metrics** (Learning if app is stable)
| Metric | Target | Definition |
|--------|--------|-----------|
| **Crash Rate** | < 0.5% | % of sessions with crash |
| **App Stability** | > 99% | Uptime (for sync, when online) |
| **Bug Density** | < 1 bug per 1000 LOC | Code quality check |

#### **User Satisfaction** (Learning if app is loved)
| Metric | Target | Definition |
|--------|--------|-----------|
| **NPS Score** | > 8/10 | "Would you recommend to friend?" |
| **Store Rating** | > 4.5/5 stars | Play Store/App Store rating |
| **Support Tickets** | < 10% of users | Sign of UX problems if high |

---

### 7.2 MVP Success Criteria (Go/No-Go Decision)
Before launching v1.1, MVP must hit:

✅ **Functional**: All 4 core features working (input, dashboard, history, report)  
✅ **Performance**: Dashboard <1s, input <500ms on mid-range Android  
✅ **UAT**: 10 users dari target persona test 3 hari, NPS > 7  
✅ **Stability**: 0 crashes dalam 1000 user sessions  
✅ **Localization**: Semua UI bahasa Indonesia, no English words  

If ANY criteria tidak tercapai → Fix before launch → Retest

---

### 7.3 Roadmap & OKR by Release

#### **v1.0 (MVP) - 4 weeks**
**OKR**: "Simplify personal finance tracking"
- ✓ Input transaksi 3-klik
- ✓ Dashboard overview one-screen
- ✓ Local SQLite storage
- ✓ Basic reporting (pie chart, monthly)
- **Success**: Retention D30 > 50%

---

#### **v1.1 (Stabilization) - 2 weeks after v1.0**
**OKR**: "Empower power users + improve retention"
- Multiple wallet (Dompet Harian, Tabungan, Bisnis)
- Edit category (user-defined names)
- Search/filter transactions
- Basic cloud backup (Firebase optional)
- **Success**: NPS score 8.5+

---

#### **v2.0 (Growth) - 4 weeks after v1.1**
**OKR**: "Scale to 10K+ users"
- Cloud sync (Firebase)
- Premium tier (subscription Rp25rb/year)
- Advanced reports (export PDF, compare year-over-year)
- Web PWA version
- **Success**: 5K+ DAU, retention D30 > 70%

---

## 8. Risk Assessment & Mitigation

### 8.1 Product Risks

| # | Risk | Likelihood | Impact | Mitigation |
|---|------|-----------|--------|-----------|
| R1 | User overwhelmed by features (scope creep) | Medium | High | Strict MVP scope, defer all "nice-to-have" |
| R2 | Input too slow, not meeting <5s target | Medium | High | UAT early with real users, measure & optimize |
| R3 | Data loss/corruption in SQLite | Low | Critical | Backup regularly, test recovery flow, auto-backup v1.1 |
| R4 | Offline sync conflict when going online | Medium | Medium | Simple conflict resolution (last-write-wins), v1.1 |
| R5 | User forgot/lost PIN/Biometric | Low | Medium | Defer to v1.1, use device lock instead |

**Mitigation Strategy**:
- Early UAT dengan target user (week 2 of dev) untuk validate speed/UX
- Code review + unit testing untuk data integrity
- Keep MVP scope lean (kill scope creep immediately)
- Build telemetry untuk track issues di production (v1.1+)

---

### 8.2 Market Risks

| Risk | Competitor | Our Strategy |
|------|-----------|--------------|
| Money Manager EX, MMEX too complex | Kompleks, banyak fitur overwhelming | **Our edge**: Simple, lokal, zero learning curve |
| Google Sheets (free alternative) | User pakai Excel untuk budget | **Our edge**: Visual (pie chart), mobile-first, dedicated |
| Manual notebook masih populer | Simplicity vs adoption | **Our edge**: Cloud backup, searchable, exportable |

**Go-to-Market**: TikTok/Instagram ads → Ibu rumah tangga, UMKM → Word-of-mouth → Scale

---

## 9. Business Model & Monetization

### 9.1 Pricing Strategy: Freemium Model
| Feature | Free | Premium (Rp25rb/tahun) |
|---------|------|----------------------|
| **Input transaksi** | ✅ Unlimited | ✅ Unlimited |
| **Dashboard & Reports** | ✅ Basic (pie chart, monthly) | ✅ Advanced (compare YoY) |
| **Multiple wallet** | ❌ 1 wallet only | ✅ Unlimited |
| **Cloud backup** | ❌ Manual export | ✅ Auto-sync |
| **Export PDF/Excel** | ❌ | ✅ |
| **Custom categories** | ❌ | ✅ |
| **No ads** | ❌ | ✅ |
| **Budget alert** | ❌ | ✅ v1.2 |

**Rationale**: Free version covers 80% user needs. Premium untuk power users willing to pay.

### 9.2 Monetization Timeline
- **v1.0 - v1.1**: 100% free (build user base, gather feedback)
- **v2.0**: Introduce Premium tier
- **v2.5**: Optional: Affiliate links (e.g., insurance, investment products) - ethical placement

### 9.3 Launch & Acquisition Strategy
**Phase 1 (Week 1-2)**: Soft launch di Play Store (beta)
- Invite 100 power users (personal network, FB groups)
- Gather feedback, fix bugs quickly
- Iterate weekly

**Phase 2 (Week 3-4)**: Public launch
- Marketing: TikTok, Instagram (target: Ibu Rumah Tangga 25-45 years old)
- Messaging: "Catat uang dengan 3 klik, bukan ribet"
- Influencer collab: Finance educators, parenting channels
- Reddit: r/Indonesia, r/Finance communities

**Phase 3 (Post-launch)**: Growth
- ASO (App Store Optimization)
- Content marketing: Blog tips budgeting
- Referral program (future)

---

## 10. Constraints & Scope Boundaries

### 10.1 MVP Scope: What's EXCLUDED
| Feature | Status | Reason |
|---------|--------|--------|
| **Multiple wallet** | v1.1+ | Adds complexity for MVP |
| **Cloud sync** | v1.1+ | SQLite only, manual backup option |
| **Recurring transaction** | v2.0+ | Can be done manually |
| **Bill reminder/Alert** | v1.2+ | Out of scope for simplicity |
| **Income projection** | v2.0+ | Analytics heavy |
| **Tax report** | Out of scope | Compliance too complex |
| **Currency exchange** | Out of scope | Indonesia only, IDR only |
| **Multi-language** | Out of scope | Bahasa Indonesia only |
| **Social sharing** | v2.0+ | Growth feature |

### 10.2 Technical Constraints
- **Database limit**: SQLite OK untuk ~10,000 transactions/user (years of data)
- **App size**: Target < 50MB
- **Android version**: API 29+ (Android 10, cover ~90% devices in Indonesia)
- **Development time**: 4 weeks for MVP (1 person or 2-3 developers)
- **Server**: ZERO backend requirement for MVP (important for cost!)

### 10.3 Stakeholder Expectations
**Clear to Stakeholder**: This is MVP, not a complete finance app. We're deliberately minimal to:
- Launch faster (4 weeks vs 6 months)
- Get user feedback early
- Reduce risk (simple = fewer bugs)
- Make money sooner (launch before v2.0 only)

---

## 11. Assumptions & Dependencies

### 11.1 Key Assumptions
✓ Target user comfort dengan smartphone (Android/iOS)  
✓ Consistent internet (for cloud backup, not for basic usage)  
✓ Users willing to input transactions manually (no automatic bank sync in MVP)  
✓ Currency = IDR only (Indonesia market)  
✓ No family sharing in MVP (single-user app)  

### 11.2 Dependencies
| Dependency | Version | Risk |
|------------|---------|------|
| Flutter SDK | 3.10+ | Low (mature framework) |
| Android SDK | API 29+ | Low (widely available) |
| SQLite | Native | Low (built-in) |
| Firebase | Spark (free) | Low (only v1.1+) |

---

## 12. Approval & Sign-off

| Role | Name | Date | Status |
|------|------|------|--------|
| **Product Owner** | [Name] | May 4, 2026 | Pending |
| **Tech Lead** | [Name] | May 4, 2026 | Pending |
| **Design Lead** | [Name] | May 4, 2026 | Pending |
| **Business Analyst** | - | May 4, 2026 | Draft v2.0 |

---

## 13. Document History

| Version | Date | Author | Change |
|---------|------|--------|--------|
| v1.0 | May 4, 2026 | BA | Initial BRD (basic structure) |
| v2.0 | May 4, 2026 | BA | Enhanced for simplicity + detailed specs |
| - | TBD | - | Awaiting team feedback |

---

## 14. Appendix

### A. Glossary
- **DAU**: Daily Active Users
- **MAU**: Monthly Active Users
- **NPS**: Net Promoter Score
- **MVP**: Minimum Viable Product
- **UAT**: User Acceptance Testing
- **OKR**: Objectives & Key Results
- **KPI**: Key Performance Indicator

### B. Reference Links
- [Flutter Documentation](https://flutter.dev/docs)
- [SQLite vs Firebase](https://firebase.google.com/docs/database)
- [Material Design 3](https://m3.material.io/)
- [Indonesian UX Design Guidelines](https://www.uxpa.org/) (future)