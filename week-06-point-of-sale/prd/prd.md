# Product Requirements Document (PRD)
## Aplikasi Point of Sale (POS) untuk Toko Kelontong

**Versi:** 1.0  
**Tanggal:** 17 April 2026  
**Status:** Final

---

## 1. Pendahuluan

### 1.1 Latar Belakang
Toko kelontong sebagai bagian dari usaha mikro, kecil, dan menengah (UMKM) masih banyak mengandalkan pencatatan manual dalam operasional sehari‑hari. Hal ini menyebabkan proses transaksi lambat, kesalahan perhitungan stok, serta kesulitan dalam memantau pendapatan dan profitabilitas. Oleh karena itu, diperlukan sebuah aplikasi *Point of Sale* (POS) yang sederhana, cepat, dan dapat berjalan secara *offline* untuk meningkatkan efisiensi serta akurasi pengelolaan toko.

### 1.2 Visi Produk
Menjadi solusi POS yang terjangkau, mudah digunakan, dan andal bagi pemilik toko kelontong, dengan fokus utama pada kecepatan transaksi (<10 detik), akurasi stok *real‑time*, serta kemampuan beroperasi tanpa ketergantungan penuh pada koneksi internet.

### 1.3 Target Pengguna
- Pemilik toko kelontong (skala mikro dan kecil).
- Kasir yang menjalankan transaksi harian.
- Admin (opsional) untuk pengelolaan pengguna dan pencadangan data.

---

## 2. Tujuan Bisnis

1. Mempercepat proses transaksi penjualan hingga kurang dari **10 detik** untuk setiap transaksi sederhana.
2. Menjaga akurasi data stok secara *real‑time* sehingga mengurangi kejadian kehabisan barang.
3. Menyediakan laporan operasional harian dan mingguan yang mudah dipahami untuk mendukung pengambilan keputusan cepat.
4. Mendukung metode pembayaran tunai dan non‑tunai (QRIS, *e‑wallet*) sesuai kebutuhan di lapangan.
5. Menekan total biaya kepemilikan (*total cost of ownership*) agar tetap terjangkau bagi UMKM.

---

## 3. Ruang Lingkup

### 3.1 Termasuk dalam Lingkup (In Scope)

- Manajemen produk (tambah, ubah, hapus, cari, lihat stok).
- Pencatatan transaksi penjualan tunai dan non‑tunai.
- Pemindaian *barcode* melalui kamera/scanner, atau input kode secara manual.
- Penghitungan diskon, total akhir, dan kembalian secara otomatis.
- Pencetakan struk melalui printer termal *Bluetooth*.
- Pencatatan barang masuk (pembelian dari pemasok) yang memperbarui stok.
- Notifikasi stok minimum berdasarkan ambang batas yang ditentukan.
- Pelaporan penjualan, stok, dan profitabilitas (harian/mingguan).
- Mode *offline* penuh (aplikasi tetap berfungsi tanpa internet).
- Ekspor laporan ke format PDF atau Excel.

### 3.2 Tidak Termasuk dalam Lingkup (Out of Scope)

- Manajemen gudang multi‑cabang (fitur ini tidak tersedia pada rilis awal).
- Integrasi penuh dengan sistem akuntansi enterprise.
- Aplikasi terpisah untuk pelanggan (*customer‑facing app*).
- Fitur *inventory forecasting* berbasis kecerdasan buatan.

---

## 4. Pengguna dan Hak Akses

| Peran | Tanggung Jawab Utama | Hak Akses |
|-------|----------------------|------------|
| **Kasir** | Menjalankan transaksi, menerima pembayaran, mencetak struk | Membuat transaksi, melihat stok, melihat riwayat transaksi harian |
| **Pemilik** | Mengelola barang, stok, pemasok, dan mengevaluasi kinerja toko | Seluruh akses operasional dan laporan |
| **Admin (opsional)** | Mengelola akun pengguna dan melakukan pencadangan data | Manajemen pengguna, *backup* dan *restore* data |

> **Catatan:** Setiap pengguna harus masuk (*login*) menggunakan PIN atau kata sandi untuk menjamin keamanan data.

---

## 5. Fitur Fungsional

### 5.1 Manajemen Data Master

#### 5.1.1 Produk
- **Atribut:** ID produk, nama, *barcode*, harga beli, harga jual, satuan, stok saat ini, stok minimum.
- **Fungsi:**
  - Tambah, ubah, hapus produk.
  - Cari produk berdasarkan nama atau *barcode*.
  - Lihat daftar produk dengan informasi stok dan harga.

#### 5.1.2 Pemasok
- **Atribut:** Nama pemasok, nomor telepon, catatan piutang/hutang (opsional).
- **Fungsi:** Tambah, ubah, hapus data pemasok.

### 5.2 Proses Transaksi Penjualan

**Aktor utama:** Kasir

**Alur normal:**
1. Kasir memilih menu **Transaksi Baru**.
2. Kasir menambahkan produk ke keranjang melalui:
   - Pemindaian *barcode* (kamera atau scanner USB/Bluetooth).
   - Pencarian produk dari katalog.
   - Input manual nama produk.
3. Kasir memasukkan jumlah item (kuantitas).
4. Sistem menampilkan subtotal.
5. Kasir memberikan diskon (jika ada promo) dan sistem menghitung total akhir.
6. Kasir memilih metode pembayaran:
   - Tunai → masukkan nominal uang, sistem hitung kembalian.
   - QRIS / *e‑wallet* → sistem menampilkan kode QR atau *deep link*.
7. Kasir menekan tombol **Selesai**.
8. Sistem menyimpan transaksi, mengurangi stok produk, dan mencetak struk.

**Aturan bisnis:**
- Stok tidak boleh negatif. Jika stok tidak mencukupi, sistem menolak penambahan produk.
- Diskon dapat berupa persen atau nominal, diterapkan ke seluruh transaksi atau per produk (opsional).

### 5.3 Pencatatan Barang Masuk

**Aktor utama:** Pemilik

**Alur normal:**
1. Buka menu **Stok** → **Barang Masuk**.
2. Pilih produk yang akan ditambah stoknya.
3. Pilih pemasok (opsional).
4. Masukkan kuantitas dan harga beli per unit.
5. Sistem menghitung total nilai pembelian.
6. Simpan transaksi → stok produk bertambah.

### 5.4 Manajemen Stok dan Notifikasi

- Stok produk berkurang secara otomatis setiap transaksi penjualan selesai.
- Stok bertambah melalui pencatatan barang masuk.
- Jika stok saat ini kurang dari **stok minimum**, sistem menampilkan notifikasi di halaman utama dan daftar produk.

### 5.5 Laporan Operasional

| Jenis Laporan | Isi Laporan |
|---------------|--------------|
| **Penjualan harian** | Omzet, jumlah transaksi, rincian metode pembayaran, daftar produk terjual beserta kuantitasnya. |
| **Stok** | Daftar barang dengan stok menipis, riwayat pergerakan stok (masuk/keluar), barang dengan pergerakan lambat (*slow moving*). |
| **Profitabilitas** | Margin kotor per produk (harga jual - harga beli), total margin untuk periode yang dipilih. |

> Semua laporan dapat diekspor ke PDF atau Excel.

### 5.6 Mode Offline dan Sinkronisasi (Opsional)

- Aplikasi berfungsi penuh tanpa koneksi internet.
- Data disimpan secara lokal (SQLite).
- Jika diperlukan sinkronisasi antar‑perangkat atau *backup* cloud, fitur ini bersifat opsional dan dapat diaktifkan kemudian.

---

## 6. Kebutuhan Non‑Fungsional

| Aspek | Target / Kriteria |
|-------|-------------------|
| **Ketersediaan** | Aplikasi tetap dapat digunakan secara penuh dalam mode *offline*. |
| **Kinerja** | Waktu transaksi < 10 detik; pemuatan daftar produk < 2 detik. |
| **Kegunaan** | Antarmuka sederhana dengan tombol besar, ramah untuk perangkat Android kelas bawah. |
| **Keamanan** | *Login* menggunakan PIN, enkripsi data lokal, opsi *backup* manual. |
| **Kompatibilitas Perangkat** | Mendukung pemindai *barcode* USB/Bluetooth serta printer termal dengan protokol ESC/POS. |

---

## 7. Alur Pengguna (Use Case)

### UC‑01: Proses Transaksi Penjualan
| Elemen | Deskripsi |
|--------|------------|
| **Aktor** | Kasir |
| **Prasyarat** | Produk sudah terdaftar di data master. |
| **Alur utama** | 1. Buka transaksi baru.<br>2. Tambah produk (scan/cari/input).<br>3. Input jumlah.<br>4. Terapkan diskon (bila ada).<br>5. Pilih metode pembayaran.<br>6. Selesaikan transaksi → cetak struk. |
| **Alur alternatif** | - Stok tidak cukup → sistem tolak penambahan.<br>- Pembayaran tunai dengan nominal kurang → sistem minta ulang. |
| **Hasil** | Transaksi tersimpan, stok berkurang, struk tercetak. |

### UC‑02: Pencatatan Barang Masuk
| Elemen | Deskripsi |
|--------|------------|
| **Aktor** | Pemilik |
| **Prasyarat** | Produk sudah terdaftar; pemasok bersifat opsional. |
| **Alur utama** | 1. Buka menu Barang Masuk.<br>2. Pilih produk.<br>3. Pilih pemasok (opsional).<br>4. Isi kuantitas dan harga beli.<br>5. Simpan. |
| **Hasil** | Stok produk bertambah, tercatat riwayat pembelian. |

### UC‑03: Pemantauan Laporan Harian
| Elemen | Deskripsi |
|--------|------------|
| **Aktor** | Pemilik |
| **Prasyarat** | Data transaksi sudah ada. |
| **Alur utama** | 1. Buka menu Laporan.<br>2. Pilih tanggal periode.<br>3. Tinjau omzet, margin, produk terlaris.<br>4. Ekspor ke PDF/Excel jika diperlukan. |
| **Hasil** | Laporan tersedia untuk analisis. |

---

## 8. Indikator Keberhasilan (KPI)

| Indikator | Target |
|-----------|--------|
| Jumlah transaksi harian (pada jam operasional normal) | Minimal 200 transaksi |
| Akurasi stok (kesesuaian antara sistem dan stok fisik) | ≥ 99% |
| Rata‑rata waktu transaksi | < 10 detik |
| Frekuensi pemilik mengakses laporan | Minimal 1 kali per minggu |

---

## 9. Risiko dan Mitigasi

| Risiko | Dampak | Mitigasi |
|--------|--------|----------|
| Perangkat (scanner/printer) tidak kompatibel | Transaksi melambat | Uji coba dengan merek perangkat umum; sediakan opsi input manual. |
| Kehilangan data karena kerusakan perangkat | Kerugian pencatatan | Sediakan fitur *backup* lokal dan opsi sinkronisasi cloud (opsional). |
| Kasir kesulitan menggunakan antarmuka | Penolakan pengguna | Desain antarmuka dengan tombol besar, ikon jelas, dan panduan singkat. |
| Daya saing dengan aplikasi gratis | Kehilangan pangsa pasar | Fokus pada keandalan *offline* dan kecepatan yang tidak dimiliki aplikasi gratis berbasis cloud. |

---

## 10. Rekomendasi Implementasi (MVP)

**Fase 1 (4 minggu) – Minimum Viable Product (MVP)**
- Manajemen produk (CRUD).
- Transaksi penjualan tunai dengan pencetakan struk.
- Pencatatan barang masuk.
- Laporan penjualan dan stok sederhana.
- Mode *offline* penuh.

**Fase 2 (2 minggu) – Uji Pilot**
- Uji coba di satu toko kelontong.
- Validasi alur kasir dan kecepatan transaksi.
- Perbaikan kecil berdasarkan umpan balik.

**Fase 3 (Opsional) – Ekspansi Fitur**
- Dukungan QRIS/*e‑wallet*.
- Sinkronisasi multi‑perangkat dan *backup* cloud.
- Laporan profitabilitas yang lebih detail.

---

## 11. Glosarium

| Istilah | Penjelasan |
|---------|-------------|
| **POS** | *Point of Sale* – tempat dan sistem untuk melakukan transaksi penjualan. |
| **QRIS** | Standar kode QR pembayaran di Indonesia yang menyatukan berbagai *e‑wallet* dan bank. |
| **ESC/POS** | Protokol komunikasi standar untuk printer termal. |
| **Stok minimum** | Ambang batas jumlah barang yang memicu peringatan agar pemilik segera memesan ulang. |

---