# Enterprise School Integration System (SMA)

Sistem Informasi Manajemen Sekolah berbasis fullstack dengan pendekatan Clean Architecture. Sistem ini dirancang untuk skalabilitas, keamanan data, dan pengalaman pengguna yang reaktif.

Fokus utama mencakup:
- Manajemen akademik
- Kedisiplinan dan absensi
- Penggajian guru (payroll)
- Sistem pendukung keputusan (rekomendasi siswa berprestasi)

## Daftar Isi
- [Ringkasan Sistem](#ringkasan-sistem)
- [Tech Stack dan Arsitektur](#tech-stack-dan-arsitektur)
- [Fase Pengembangan dan Fitur Utama](#fase-pengembangan-dan-fitur-utama)
- [Alur Kerja Sistem](#alur-kerja-sistem)
- [Nilai Jual Enterprise](#nilai-jual-enterprise)

## Ringkasan Sistem

Enterprise School Integration System membantu operasional sekolah dari hulu ke hilir:
- Setup tahun ajaran dan master data
- Proses belajar mengajar harian
- Evaluasi kinerja akademik
- Penggajian dan pelaporan
- Audit aktivitas penting

## Tech Stack dan Arsitektur

| Komponen | Teknologi |
| --- | --- |
| Backend | Golang (Go Fiber) |
| Database | PostgreSQL + GORM |
| Frontend | SvelteKit + Tailwind CSS |
| Arsitektur | Clean Architecture (Repository -> Service -> Handler) |
| Keamanan | JWT, RBAC, Bcrypt |

## Fase Pengembangan dan Fitur Utama

### Phase 1 - Foundation and Security
Fase ini membangun pondasi sistem agar kokoh dan aman dari akses tidak sah.

- JWT and RBAC middleware untuk membatasi akses endpoint berdasarkan role: ADMIN, GURU, TU, KEPSEK, BK, MURID.
- Centralized database migration dengan hierarki relasi (Tier 1 hingga Tier 5) untuk mencegah foreign key error saat deploy.
- Master data management (single active year): CRUD Tahun Ajaran, Kelas, Siswa, dan Mata Pelajaran, dengan aturan hanya satu Tahun Ajaran aktif pada satu waktu.

### Phase 2 - Core Academic and Attendance
Fase ini memfasilitasi kegiatan belajar mengajar sehari-hari.

- Batch attendance (bulk upsert): absensi massal satu klik (Tandai Semua Hadir), memangkas sekitar 40 kueri menjadi 1 kueri atomik.
- Teacher assignments validation: guru hanya bisa memasukkan nilai untuk kelas dan mata pelajaran yang ditugaskan.
- Grading system: input nilai Ulangan, UTS, UAS dengan validasi nilai 0-100 di frontend dan backend.

### Phase 3 - Payroll and Decision Support System (DSS)
Fase analitik dan operasional tingkat lanjut.

- Dynamic payroll system (JSONB): komponen gaji (tunjangan dan potongan) dihitung dinamis tanpa mengubah skema database.
- Student recommendation engine: algoritma Simple Additive Weighting (SAW) dengan bobot 70% rata-rata akademik + 30% attendance rate.

### Phase 4 - Announcements and Asynchronous Reporting
Fase penyempurnaan fitur komersial dan optimasi performa server.

- Targeted announcements: pengumuman berbasis role dengan JSONB `target_roles`.
- Asynchronous report generation (goroutines): ekspor rapor PDF (gofpdf) melalui background jobs agar API tetap responsif.
- Audit logging: pelacakan aksi, user, dan timestamp pada perubahan krusial.

## Alur Kerja Sistem

### 1. Tahap Persiapan (Admin/TU)
- Admin mengaktifkan Tahun Ajaran baru (tahun ajaran lama otomatis nonaktif).
- Admin membuat Kelas, Mata Pelajaran, dan registrasi akun siswa ke kelas.
- Admin menetapkan teacher assignments.

### 2. Tahap KBM dan Operasional Harian (Guru)
- Guru piket atau wali kelas melakukan absensi massal setiap pagi.
- Guru mata pelajaran memasukkan nilai siswa pada kelas yang diampu.

### 3. Tahap Evaluasi dan Keuangan (BK/Keuangan)
- Tim keuangan mencetak slip gaji bulanan berdasarkan base salary dan potongan absensi.
- Guru BK menjalankan kalkulasi rekomendasi siswa berprestasi untuk melihat peringkat otomatis.

### 4. Tahap Pelaporan Akhir (Kepala Sekolah/Admin)
- Kepala sekolah memantau audit log.
- Admin memicu pembuatan rapor PDF di background hingga siap diunduh staf atau siswa terkait.

## Nilai Jual Enterprise

- Zero N+1 query problem: optimasi dengan GORM preload dan `clause.OnConflict` (true bulk upsert) pada absensi massal.
- Business logic centric di service layer: logika bisnis tidak dipindahkan ke trigger database.
- Non-blocking I/O UX: pemanggilan data paralel dengan `Promise.all` di SvelteKit untuk mempercepat waktu muat layar.