# Product Requirements Document (PRD)

## Sistem Manajemen & Gizi Perawatan Hewan (Pet Care & Nutrition System)

**Versi:** 1.1
**Fokus Utama:** MVP – Rekam Gizi, Penjadwalan *Grooming*, dan Penagihan Terpisah per Individu.
**Tech Stack:** Backend (Golang, Fiber, GORM, PostgreSQL), Frontend (ReactJS, TailwindCSS).

---

## 1. Ringkasan Eksekutif

### 1.1 Definisi Produk

Sistem Manajemen & Gizi Perawatan Hewan adalah aplikasi berbasis *web* yang dirancang untuk fasilitas perawatan hewan peliharaan (klinik, *pet hotel*, atau *grooming center*). Aplikasi ini tidak hanya mencatat transaksi dasar, tetapi berfokus pada dua pilar utama: **Pemantauan kualitas diet/nutrisi** (khususnya untuk ras kucing lokal maupun eksotis) dan **Sistem penagihan layanan yang diisolasi per individu hewan**.

### 1.2 Value Proposition

* **Akurasi Finansial:** Kalkulasi harga layanan dihitung secara terpisah untuk masing-masing individu hewan. Tidak ada lagi kesalahan akibat formula gabungan saat satu pemilik membawa banyak hewan peliharaan.
* **Fokus Nutrisi:** Pencatatan spesifik mengenai merek makanan, jadwal diet berkualitas tinggi, dan transformasi fisik hewan peliharaan.
* **Performa Tinggi:** *Backend* Golang memastikan kalkulasi layanan dan pemanggilan riwayat medis terjadi dalam hitungan milidetik.

---

## 2. Asumsi, Peran, dan Batasan

### 2.1 Asumsi Operasional

- Sistem digunakan oleh staf internal (resepsionis, kasir, perawat hewan, manajer).
- Pembayaran bersifat langsung (kas/tunai, transfer, QRIS) tanpa termin.
- Satu lokasi operasional (multi-cabang di luar ruang lingkup MVP).

### 2.2 Peran Pengguna

- **Resepsionis:** registrasi pemilik dan hewan, jadwal layanan, check-in.
- **Perawat/Groomer:** update status layanan, input catatan gizi.
- **Kasir:** membuat invoice, menerima pembayaran, cetak bukti.
- **Manajer:** melihat dashboard operasional dan laporan ringkas.

### 2.3 Batasan MVP

- Tidak ada pemesanan mandiri oleh pelanggan.
- Tidak ada piutang atau termin pembayaran.
- Tidak ada marketplace/inventaris ritel.

## 3. Ruang Lingkup Produk (Fase 1 - MVP)

### 2.1 Fitur Wajib (MUST HAVE)

| Kode | Nama Fitur | Deskripsi | Prioritas |
| --- | --- | --- | --- |
| **FR-001** | **Manajemen Pemilik & Hewan** | CRUD data Pemilik (Nama, Kontak) dan Hewan Peliharaan (Nama, Spesies, Berat Badan, Kondisi Fisik). Satu Pemilik bisa memiliki banyak Hewan. | 🔴 P0 |
| **FR-002** | **Rekam Gizi & Transformasi** | Modul pencatatan harian/mingguan untuk melacak asupan kalori, jenis/merek makanan (diet berkualitas tinggi), dan catatan transformasi fisik (bulu, berat badan). | 🔴 P0 |
| **FR-003** | **Manajemen Layanan (*Grooming*/Inap)** | Pembuatan tiket layanan perawatan untuk hewan peliharaan yang mencakup status pengerjaan (Antre, Proses, Selesai). | 🔴 P0 |
| **FR-004** | **Penagihan Dinamis Individu** | Sistem kalkulasi tagihan yang **secara ketat memisahkan perhitungan harga per individu hewan**, bukan menggunakan formula gabungan per pemilik. | 🔴 P0 |
| **FR-005** | **Dasbor Operasional** | Menampilkan antrean *grooming* hari ini, ringkasan hewan yang sedang rawat inap, dan total estimasi pendapatan harian. | 🟢 P1 |
| **FR-006** | **Otentikasi & Role** | Login berbasis peran (Resepsionis, Kasir, Groomer, Manajer) dengan pembatasan akses halaman dan API. | 🔴 P0 |
| **FR-007** | **Penjadwalan Layanan** | Jadwal *grooming*/rawat inap per hewan dengan kapasitas harian dan status check-in/check-out. | 🔴 P0 |
| **FR-008** | **Penerimaan Pembayaran** | Pencatatan pembayaran per invoice (tunai/transfer/QRIS), status Paid/Unpaid, dan cetak bukti transaksi. | 🔴 P0 |

### 2.2 Out of Scope (Ditunda ke Fase Selanjutnya)

❌ *Booking* mandiri (*Self-service*) oleh pelanggan dari luar.
❌ Sistem kasbon/piutang pelanggan (fokus MVP pada kas tuntas).
❌ Manajemen inventaris barang eceran (fokus pada layanan dan gizi).

---

## 4. Desain Arsitektur & Skema Database (Golang + GORM)

Untuk memastikan logika penagihan (*billing*) beroperasi dengan benar (FR-004), kita menggunakan desain *database* relasional di mana sebuah tagihan (*invoice*) dipecah menjadi *line items* yang terikat kuat pada ID Hewan spesifik.

### 4.1 Entitas & Model GORM

**A. Pemilik & Hewan (`owners` & `pets`)**

```go
type Owner struct {
	ID        string `gorm:"type:varchar(36);primaryKey"`
	Name      string `gorm:"not null"`
	Phone     string
	Pets      []Pet  `gorm:"foreignKey:OwnerID"`
}

type Pet struct {
	ID        string `gorm:"type:varchar(36);primaryKey"`
	OwnerID   string `gorm:"type:varchar(36);index"`
	Name      string `gorm:"not null"`
	Species   string // Kucing, Anjing, dll.
	Breed     string
	Weight    float64
	DietNotes string // Preferensi atau pantangan makanan
}

```

**B. Rekam Gizi (`nutrition_logs`)**

```go
type NutritionLog struct {
	ID          string    `gorm:"type:varchar(36);primaryKey"`
	PetID       string    `gorm:"type:varchar(36);index"`
	LogDate     time.Time `gorm:"index"`
	FoodBrand   string    // Merek makanan berkualitas tinggi yang diberikan
	Calories    int       // Asupan kalori
	HealthNotes string    // Catatan transformasi bulu, pencernaan, dll.
}

```

**C. Layanan & Penagihan (`services`, `invoices`, `invoice_items`)**

```go
type Service struct {
	ID          string  `gorm:"type:varchar(36);primaryKey"`
	Name        string  // e.g., "Premium Grooming Kucing", "Konsultasi Diet"
	BasePrice   float64
}

type Invoice struct {
	ID          string    `gorm:"type:varchar(36);primaryKey"`
	OwnerID     string    `gorm:"type:varchar(36);index"`
	TotalAmount float64
	Status      string    // "Unpaid", "Paid"
	CreatedAt   time.Time
	Items       []InvoiceItem `gorm:"foreignKey:InvoiceID"`
}

// InvoiceItem: Kunci dari FR-004. Kalkulasi wajib dipisah per hewan.
type InvoiceItem struct {
	ID          string  `gorm:"type:varchar(36);primaryKey"`
	InvoiceID   string  `gorm:"type:varchar(36);index"`
	PetID       string  `gorm:"type:varchar(36);index"` // Item ini untuk hewan yang mana?
	ServiceID   string  `gorm:"type:varchar(36)"`
	Price       float64 // Harga aktual (bisa disesuaikan berdasarkan ukuran hewan)
}

```

**D. Penjadwalan Layanan (`appointments`, `service_orders`)**

```go
type Appointment struct {
	ID        string    `gorm:"type:varchar(36);primaryKey"`
	PetID     string    `gorm:"type:varchar(36);index"`
	ServiceID string    `gorm:"type:varchar(36);index"`
	StartAt   time.Time `gorm:"index"`
	EndAt     time.Time
	Status    string    // "Scheduled", "CheckedIn", "InProgress", "Completed", "Cancelled"
}

type ServiceOrder struct {
	ID           string    `gorm:"type:varchar(36);primaryKey"`
	AppointmentID string   `gorm:"type:varchar(36);index"`
	PetID        string    `gorm:"type:varchar(36);index"`
	Notes        string
	CreatedAt    time.Time
}

```

**E. Pembayaran (`payments`)**

```go
type Payment struct {
	ID         string    `gorm:"type:varchar(36);primaryKey"`
	InvoiceID  string    `gorm:"type:varchar(36);index"`
	Method     string    // "Cash", "Transfer", "QRIS"
	Amount     float64
	PaidAt     time.Time
	Reference  string    // nomor referensi transfer/QR
}

```

---

## 5. Alur Pengguna (Skenario Utama)

### UC-001: Menerima Pelanggan dengan Multi-Hewan

1. **Pelanggan datang** membawa 3 ekor kucing peliharaannya.
2. **Staf Frontend (ReactJS)** mendaftarkan Pemilik dan 3 profil Kucing secara terpisah di aplikasi.
3. **Staf membuat Order:**
* Kucing A: Layanan *Premium Grooming* & Terapi Nutrisi Bulu.
* Kucing B: Layanan Penitipan Harian (*Boarding*).
* Kucing C: Cek berat badan (tanpa biaya).


4. **Kalkulasi (Golang Backend):** Sistem *use case* secara ketat memproses perhitungan untuk Kucing A, Kucing B, dan Kucing C secara mandiri.
5. **Cetak Tagihan:** *Invoice* tercetak dengan struktur rincian (*breakdown*) yang dikelompokkan berdasarkan nama hewan peliharaan, bukan diakumulasi buta.

### UC-002: Penjadwalan Grooming dengan Kapasitas Harian

1. Resepsionis memilih tanggal layanan dan melihat kapasitas slot tersedia.
2. Sistem menolak jika slot penuh atau jadwal bentrok untuk hewan yang sama.
3. Saat hewan datang, status berubah ke "CheckedIn" lalu "InProgress".
4. Setelah selesai, sistem otomatis membuat item tagihan sesuai layanan.

---

## 6. API Contract (Ringkas)

Contoh endpoint utama untuk integrasi frontend React:

- `POST /auth/login`
- `GET /owners`, `POST /owners`, `PUT /owners/:id`
- `GET /pets?owner_id=...`, `POST /pets`
- `GET /services`, `POST /services`
- `POST /appointments`, `PUT /appointments/:id/status`
- `POST /invoices`, `GET /invoices/:id`
- `POST /payments`

Setiap respons mengembalikan format JSON konsisten:

```json
{
	"data": {},
	"message": "ok",
	"error": null
}
```

## 7. Persyaratan Non-Fungsional

* **Arsitektur:** Menggunakan **Hexagonal Architecture** di Golang untuk memisahkan domain billing dan scheduling dari Fiber/GORM.
* **Antarmuka UI:** React dengan *state management* terpusat (mis. Zustand/Redux) agar form multi-hewan tetap responsif.
* **Waktu Respons:** API pembuatan tagihan multi-hewan < 200ms pada beban 100 req/menit.
* **Keamanan:** Password di-hash (bcrypt/argon2), JWT dengan refresh token, role-based access control.
* **Audit Log:** Pencatatan perubahan status layanan, pembayaran, dan pembatalan.
* **Reliabilitas:** Backup database harian, strategi restore terdokumentasi.
* **Observabilitas:** Logging terstruktur (JSON), trace ID per request, dan metric dasar (latensi, error rate).

---

## 8. Kriteria Penerimaan (Acceptance Criteria)

- Invoice selalu memiliki `InvoiceItem` dengan `PetID` valid; sistem menolak jika kosong.
- Jadwal menolak *double booking* untuk hewan yang sama pada waktu yang tumpang tindih.
- Status pembayaran berubah otomatis ke "Paid" setelah pembayaran penuh tercatat.
- Kasir dapat mencetak bukti pembayaran yang memuat rincian per hewan.
- Role Resepsionis tidak bisa mengakses endpoint pembayaran.

## 9. Tahapan Pengembangan (Roadmap Eksekusi)

* **Langkah 1 (Backend):** Implementasi entitas Domain dan Skema GORM untuk `Owner`, `Pet`, `NutritionLog`, `Service`, `Appointment`, `Invoice`. Membangun API CRUD dasar.
* **Langkah 2 (Backend):** Implementasi modul `Billing` dan `Scheduling` dengan aturan bisnis (validasi `PetID`, cek konflik jadwal, status layanan).
* **Langkah 3 (Frontend):** Setup ReactJS + Tailwind. Pembuatan Dasbor utama, form registrasi, dan form rekam gizi.
* **Langkah 4 (Frontend & Integrasi):** Form "Kasir" dan "Penjadwalan" dengan pemilihan layanan per hewan dan *live calculation*.