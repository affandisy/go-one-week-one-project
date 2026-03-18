# Improvisasi Arsitektur dan Implementasi

Dokumen ini merangkum penyesuaian dari PRD agar sistem sekolah lebih stabil, efisien, dan nyaman digunakan.

## Ringkasan Cepat

| No | Topik | Permintaan PRD | Improvisasi | Dampak Utama |
|---|---|---|---|---|
| 1 | Urutan Migrasi Terpusat | GORM untuk operasi database | Migrasi diatur berurutan di `config/database.go` | Deploy di server kosong lebih aman |
| 2 | Input Absensi Massal | Endpoint `POST /attendance/student` | Absensi batch (array siswa) dalam satu request | UX guru lebih cepat, request lebih hemat |
| 3 | Single Active Year | `academic_years.is_active` | Logika service + transaction untuk 1 tahun aktif | Cegah anomali data akademik |
| 4 | `target_roles` JSONB | `target_roles` bertipe `VARCHAR[]` | Ubah ke `JSONB` (`gorm:"type:jsonb"`) | Integrasi Golang-Svelte lebih mulus |
| 5 | UX Laporan Asinkron | Generate PDF asynchronous | Auto-refresh data via `setTimeout(fetchAllData, 3000)` | Tombol download muncul otomatis |

## 1) Urutan Migrasi Terpusat (Centralized Migration)

### Kondisi di PRD
PRD hanya menyebut penggunaan GORM untuk operasi database.

### Improvisasi
Migrasi tidak dijalankan bebas lewat `db.AutoMigrate` di `main.go`. Urutan migrasi diatur terpusat pada `config/database.go`, dari tabel independen (contoh: `User`) ke tabel dependen (contoh: `Grade`).

### Alasan dan Manfaat
- Menghindari error foreign key saat tabel dependen dibuat lebih dulu.
- Menjaga startup aplikasi tetap konsisten.
- Memastikan proses deploy awal di server kosong berjalan mulus.

## 2) Input Absensi Massal (Batch Attendance)

### Kondisi di PRD
Disediakan endpoint `POST /attendance/student`.

### Improvisasi
API dan UI Svelte dirancang menerima array data siswa dalam satu request, bukan input satu per satu.

### Alasan dan Manfaat
- Guru dapat mengabsen satu kelas sekaligus dalam satu halaman.
- Interaksi lebih cepat: pilih status hadir/sakit/izin/alfa lalu simpan sekali.
- Trafik request HTTP berkurang drastis.

## 3) Logika Tahun Ajaran Tunggal (Single Active Year)

### Kondisi di PRD
Tabel `academic_years` memiliki kolom `is_active`.

### Improvisasi
Di service layer Golang ditambahkan aturan: ketika admin mengaktifkan tahun ajaran baru, sistem otomatis menonaktifkan tahun ajaran aktif sebelumnya dalam satu database transaction.

### Alasan dan Manfaat
- Mencegah kondisi dua tahun ajaran aktif bersamaan.
- Menjaga konsistensi data akademik lintas modul.
- Menurunkan risiko bug pada laporan dan proses belajar mengajar.

## 4) Perubahan `target_roles` menjadi JSONB

### Kondisi di PRD
Kolom `announcements.target_roles` bertipe PostgreSQL array (`VARCHAR[]`).

### Improvisasi
Tipe data diubah ke `JSONB` pada model GORM (`gorm:"type:jsonb"`).

### Alasan dan Manfaat
- Parsing array string dari frontend Svelte menjadi lebih sederhana.
- Proses marshal/unmarshal di Golang lebih stabil.
- Mengurangi potensi error kompatibilitas driver PostgreSQL.

## 5) UX Laporan Asinkron (Simulasi Real-time)

### Kondisi di PRD
Generate laporan PDF berjalan asynchronous agar tidak blocking.

### Improvisasi
Setelah tombol generate diklik, frontend tidak berhenti di notifikasi sukses. UI menjadwalkan pengambilan ulang data dengan `setTimeout(fetchAllData, 3000)`.

### Alasan dan Manfaat
- Selaras dengan proses backend yang memproses laporan di goroutine.
- User tidak perlu refresh manual untuk melihat status terbaru.
- Tombol download muncul otomatis ketika file sudah siap.

## Kesimpulan

Lima improvisasi ini berfokus pada tiga nilai utama:

- Stabilitas teknis saat deploy dan migrasi.
- Konsistensi data pada proses akademik kritikal.
- Pengalaman pengguna yang lebih cepat dan minim langkah manual.