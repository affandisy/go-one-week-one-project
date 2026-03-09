# Dokumentasi Bisnis – Sistem Integrasi Sekolah (Jenjang SMA)

## 1. Latar belakang
Sistem ini dirancang untuk jenjang SMA yang mengalami overload administrasi karena data terpisah (nilai, absensi, gaji, dan rekap mutu).  
Tujuan utama: mengintegrasikan data dan proses untuk guru, murid, kepala sekolah, dan bagian keuangan dalam satu sistem.

## 2. Tujuan bisnis
- Memudahkan input dan pemantauan nilai, absensi, dan gaji.
- Menghasilkan laporan otomatis untuk kepala sekolah (rapor, kehadiran, keuangan).
- Mendukung rekomendasi murid terbaik berbasis data.

## 3. Stakeholder dan user role
| Role            | Deskripsi singkat                                                                 |
|-----------------|------------------------------------------------------------------------------------|
| Kepala sekolah  | Pengawas sistem, penerima laporan, pemberi kebijakan.                             |
| Wakil / TU      | Membuat master data, jadwal, kelas, dan pengumuman.                              |
| Guru            | Input nilai, presensi mengajar, melihat rekap murid.                             |
| Murid           | Terima data rapor, nilai, absensi, dan rekomendasi.                              |
| Bagian keuangan | Mengatur struktur gaji, memvalidasi perhitungan, cetak slip gaji.               |
| Admin sistem    | Mengatur user, hak akses, dan backup data.                                       |

## 4. Scope sistem
### Termasuk
- Modul akademik (nilai, rapor, kelas, jadwal).
- Modul presensi (guru dan murid).
- Modul keuangan (gaji guru).
- Modul rekomendasi murid terbaik.
- Modul pengumuman & dashboard.

### Tidak termasuk (out of scope)
- ERP lengkap (inventory, perpustakaan, catering).
- Aplikasi mobile khusus orang tua (jika tidak diinginkan).

## 5. Modul utama
### 5.1 Modul akademik
- Master kurikulum SMA (peminatan IPA/IPS, tahun ajaran).
- Manajemen kelas, rombel, jadwal mengajar.
- Input nilai harian, UTS, UAS.
- Cetak rapor digital (PDF).

### 5.2 Modul absensi
- Presensi harian guru (online/QR).
- Presensi harian murid (per kelas).
- Rekap kehadiran per bulan, per semester.

### 5.3 Modul gaji guru
- Struktur gaji: pokok, tunjangan, honor, lembur, potongan.
- Integrasi dengan data absensi (terlambat, izin, sakit).
- Perhitungan gaji otomatis dan slip gaji digital.

### 5.4 Modul rekomendasi murid terbaik
- Kriteria: nilai akademik, kehadiran, perilaku, prestasi.
- Bobot dinamis (configurable oleh sekolah).
- List nominasi per kelas, per semester.

### 5.5 Modul pengumuman & dashboard
- Pengumuman sekolah (PPDB, UAS, event).
- Dashboard kepala sekolah: kehadiran, keuangan, nilai.
- Dashboard guru: kelas mengajar, nilai rata‑rata, murid bermasalah (low performance).

## 6. Non‑functional requirements
- Role‑based access control (RBAC).
- Autentikasi login per user.
- Backup otomatis harian.
- Laporan PDF/Excel.

## 7. Use case utama

### 7.1 Use Case 1: Input Nilai Murid
- **Actor** : Guru  
- **Pre‑condition** : Jadwal mengajar sudah ada, kelas dan murid sudah terdaftar.  
- **Flow** :
  1. Guru login ke sistem.
  2. Pilih menu “Akademik” → “Input Nilai”.
  3. Pilih kelas, mata pelajaran, dan jenis nilai (harian/UTS/UAS).
  4. Sistem menampilkan daftar murid kelas tersebut.
  5. Guru mengisi nilai per murid.
  6. Sistem menyimpan nilai dan menghitung rata‑rata kelas.
  7. Sistem memberi notifikasi penyimpanan berhasil.

### 7.2 Use Case 2: Presensi Harian Guru
- **Actor** : Guru  
- **Pre‑condition** : Jadwal mengajar harian sudah terdaftar.  
- **Flow** :
  1. Guru login atau scan QR di kelas.
  2. Sistem menampilkan daftar kelas yang perlu diisi.
  3. Guru menandai status: hadir, izin, sakit, atau belum hadir.
  4. Sistem mencatat waktu dan status.
  5. Sistem update rekap kehadiran guru hari tersebut.

### 7.3 Use Case 3: Presensi Harian Murid
- **Actor** : Guru Wali Kelas / Guru Piket  
- **Pre‑condition** : Kelas sudah dibuka di sistem.  
- **Flow** :
  1. Guru login dan memilih kelas yang akan diisi presensi.
  2. Sistem menampilkan daftar murid.
  3. Guru mengisi status: hadir, izin, sakit, atau alfa.
  4. Sistem menyimpan data dan update rekap kehadiran per kelas dan per murid.

### 7.4 Use Case 4: Hitung Gaji Guru
- **Actor** : Sistem (otomatis) + Operator Keuangan  
- **Pre‑condition** : Struktur gaji & tunjangan sudah dikonfigurasi.  
- **Flow** :
  1. Operator keuangan memilih bulan dan periode gaji.
  2. Sistem mengambil data:
     - Absensi guru (terlambat, izin, sakit, tanpa keterangan).
     - Jam mengajar ekstra dan honor.
  3. Sistem menghitung komponen gaji berdasarkan formula:
     \[
     \text{Gaji akhir} = \text{Gaji pokok} + \text{Tunjangan} + \text{Honor} - \text{Potongan}
     \]
  4. Sistem menyimpan hasil perhitungan.
  5. Operator dapat melihat dan mengonfirmasi slip gaji.
  6. Slip gaji dapat di‑export PDF.

### 7.5 Use Case 5: Rekomendasi Murid Terbaik
- **Actor** : Sistem (otomatis) + Guru BK / Kepala Sekolah  
- **Pre‑condition** : Kriteria dan bobot rekomendasi sudah dikonfigurasi.  
- **Flow** :
  1. Sistem mengumpulkan data:
     - Nilai rata‑rata per murid.
     - Persentase kehadiran.
     - Data prestasi & perilaku (opsional).
  2. Sistem menghitung skor rekomendasi:
     \[
     \text{Skor total} = w_1 \cdot \text{Nilai} + w_2 \cdot \text{Kehadiran} + w_3 \cdot \text{Prestasi}
     \]
     dengan \(w_1+w_2+w_3 = 1\).
  3. Sistem membuat daftar 10–20 murid teratas per kelas.
  4. Guru BK / Kepala sekolah dapat menyetujui atau mengubah nominasi.
  5. Sistem menampilkan daftar resmi “Murid Terbaik Semester”.

### 7.6 Use Case 6: Lihat Dashboard Kepala Sekolah
- **Actor** : Kepala Sekolah  
- **Pre‑condition** : Data sekolah sudah terisi beberapa bulan.  
- **Flow** :
  1. Kepala sekolah login.
  2. Sistem menampilkan dashboard:
     - Kehadiran rata‑rata guru & murid.
     - Rata‑rata nilai per kelas/mapel.
     - Total gaji guru bulan ini.
     - Murid yang sering absen atau nilai rendah.
  3. Kepala sekolah dapat klik laporan detail per modul.

### 7.7 Use Case 7: Publikasi Pengumuman
- **Actor** : Wakil / TU / Admin  
- **Pre‑condition** : User memiliki role Administrator.  
- **Flow** :
  1. Admin login.
  2. Pilih menu “Pengumuman” → “Tambah Pengumuman”.
  3. Mengisi judul, isi, tipe, dan target (semua guru, semua murid, kelas tertentu, peminatan tertentu).
  4. Sistem menyimpan pengumuman dan mengirim notifikasi (in‑app).
  5. Penerima dapat melihat di halaman pengumuman.

## 8. Tabel ringkas use case

| Use Case ID | Nama Use Case                 | Actor            | Tujuan Utama |
|-------------|-------------------------------|------------------|--------------|
| UC‑01       | Input Nilai Murid            | Guru            | Mencatat data akademik murid. |
| UC‑02       | Presensi Harian Guru         | Guru            | Mencatat kehadiran guru. |
| UC‑03       | Presensi Harian Murid        | Guru Wali Kelas | Mencatat kehadiran murid. |
| UC‑04       | Hitung Gaji Guru             | Sistem + Keuangan | Menghasilkan slip gaji bulanan. |
| UC‑05       | Rekomendasi Murid Terbaik    | Sistem + BK/Kepsek | Memilih murid dengan prestasi terbaik. |
| UC‑06       | Lihat Dashboard Kepala Sekolah | Kepala Sekolah  | Memantau kinerja akademik & keuangan. |
| UC‑07       | Publikasi Pengumuman         | Admin/TU        | Menyebarkan informasi sekolah. |

## 9. High‑level workflow

### 9.1 Workflow proses akademik semester
1. Pendataan awal: murid, kelas, jadwal & kurikulum dimasukkan oleh TU.  
2. Guru mengajar → input nilai harian, UTS, UAS → sistem menyimpan.  
3. Sistem menghitung rata‑rata Semester & rapor.  
4. Kepala sekolah menyetujui rapor.  
5. Sistem memunculkan laporan rekap nilai & murid terbaik.

### 9.2 Workflow proses gaji guru bulanan
1. Konfigurasi struktur gaji & tunjangan (oleh Keuangan).  
2. Setiap harinya: guru melakukan presensi.  
3. Akhir bulan: sistem menarik data absensi & jam mengajar.  
4. Sistem menghitung gaji otomatis → generate slip.  
5. Operator keuangan meninjau → setujui → export PDF.
