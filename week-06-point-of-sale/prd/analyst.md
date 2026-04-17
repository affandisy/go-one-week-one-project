# Analisis Produk Point of Sale (POS) untuk Toko Kelontong

Versi: 1.1  
Tanggal: 17 April 2026  
Penulis: Business and System Analyst

## Ringkasan Eksekutif
Dokumen ini merangkum analisis kebutuhan dan rancangan tingkat tinggi untuk aplikasi Point of Sale (POS) skala UMKM, khususnya toko kelontong. Solusi difokuskan pada kecepatan transaksi, akurasi stok, kemudahan operasional harian, serta kemampuan berjalan saat koneksi internet terbatas. Target utama produk adalah meningkatkan efisiensi operasional pemilik toko tanpa kompleksitas sistem enterprise.

## 1. Tujuan Bisnis
1. Mempercepat proses transaksi penjualan hingga kurang dari 10 detik per transaksi sederhana.
2. Menjaga akurasi data stok secara real-time untuk mengurangi kejadian kehabisan barang.
3. Menyediakan laporan operasional harian dan mingguan untuk pengambilan keputusan cepat.
4. Mendukung metode pembayaran tunai dan non-tunai (QRIS/e-wallet) sesuai kebutuhan lapangan.
5. Menjaga total biaya kepemilikan tetap terjangkau bagi usaha mikro dan kecil.

## 2. Pemangku Kepentingan dan Peran Pengguna
| Peran | Tanggung Jawab Utama | Hak Akses |
| --- | --- | --- |
| Kasir | Menjalankan transaksi, menerima pembayaran, mencetak struk | Membuat transaksi, melihat stok, melihat riwayat transaksi harian |
| Pemilik | Mengelola barang, stok, pemasok, dan evaluasi kinerja toko | Seluruh akses operasional dan laporan |
| Admin (opsional) | Pengelolaan pengguna dan backup data | Manajemen pengguna, backup dan restore |

## 3. Ruang Lingkup Fitur Inti
### 3.1 Proses Transaksi Penjualan
1. Menambahkan item melalui pemindaian barcode, pencarian katalog, atau input manual.
2. Menghitung subtotal, diskon, total akhir, serta kembalian secara otomatis.
3. Merekam metode pembayaran: tunai, QRIS, atau e-wallet.
4. Mencetak struk melalui printer thermal Bluetooth.

### 3.2 Manajemen Stok
1. Mencatat barang masuk dengan kuantitas, harga beli, dan pemasok.
2. Mengurangi stok otomatis saat transaksi penjualan berhasil.
3. Menampilkan notifikasi stok minimum berdasarkan ambang yang ditentukan.

### 3.3 Pelaporan Operasional
| Jenis Laporan | Konten Utama |
| --- | --- |
| Penjualan harian | Omzet, jumlah transaksi, metode pembayaran, item terjual |
| Stok | Barang menipis, pergerakan stok masuk-keluar, slow moving item |
| Profitabilitas | Margin kotor per barang dan total margin periode |

### 3.4 Data Master
1. Produk: nama, barcode, harga beli, harga jual, satuan, stok, dan minimum stok.
2. Pemasok: nama, kontak, dan catatan piutang/hutang dagang (opsional).

## 4. Alur Pengguna
### 4.1 Alur Utama Kasir
1. Pengguna membuka aplikasi dan login menggunakan PIN.
2. Pengguna membuat transaksi baru dan menambahkan item.
3. Sistem menghitung total, pengguna memilih metode pembayaran.
4. Sistem menyimpan transaksi, memperbarui stok, dan mencetak struk.

### 4.2 Use Case Prioritas
#### UC-01 Proses Transaksi Penjualan
Actor: Kasir  
Alur:
1. Pilih menu Transaksi Baru.
2. Tambahkan produk (scan/cari/input manual).
3. Verifikasi jumlah item dan diskon bila ada promo.
4. Proses pembayaran dan validasi nominal.
5. Selesaikan transaksi, cetak struk, dan simpan data.

#### UC-02 Pencatatan Barang Masuk
Actor: Pemilik  
Alur:
1. Buka menu Stok, lalu pilih Barang Masuk.
2. Pilih produk dan pemasok.
3. Input kuantitas serta harga beli.
4. Simpan data dan perbarui total stok.

#### UC-03 Pemantauan Laporan Harian
Actor: Pemilik  
Alur:
1. Buka menu Laporan dan tentukan tanggal periode.
2. Tinjau omzet, margin, dan produk terlaris.
3. Ekspor laporan ke PDF/Excel bila dibutuhkan.

## 5. Desain Sistem Tingkat Tinggi
### 5.1 Arsitektur Solusi
1. Frontend: aplikasi Android native atau Progressive Web App (PWA).
2. Data lokal: SQLite sebagai penyimpanan utama mode offline.
3. Sinkronisasi cloud: layanan opsional untuk backup dan multi-perangkat.
4. Integrasi perangkat: barcode scanner dan printer thermal ESC/POS.
5. Pembayaran digital: integrasi QRIS melalui mekanisme deep link.

### 5.2 Rancangan Data Utama (SQLite)
| Tabel | Kolom Kunci |
| --- | --- |
| products | id, barcode, name, sell_price, buy_price, stock_qty, min_stock |
| sales | id, date, total_amount, payment_method, items_json |
| suppliers | id, name, phone, debt_amount |
| stock_logs | id, product_id, type, qty, date |

### 5.3 Rumus Bisnis Kunci
Profit per transaksi:

$$
	ext{Profit} = \sum (\text{Harga Jual} - \text{Harga Beli}) \times \text{Kuantitas Terjual}
$$

Aturan peringatan stok:

Jika $stock\_qty < min\_stock$, sistem menampilkan notifikasi stok minimum untuk produk terkait.

## 6. Kebutuhan Non-Fungsional
| Aspek | Target |
| --- | --- |
| Ketersediaan | Tetap dapat dipakai penuh dalam mode offline |
| Kinerja | Waktu transaksi < 10 detik, pemuatan daftar produk < 2 detik |
| Usability | Antarmuka sederhana, tombol besar, ramah perangkat Android low-end |
| Keamanan | Login PIN, enkripsi data lokal, backup opsional |
| Kompatibilitas perangkat | Mendukung barcode scanner USB/Bluetooth dan printer ESC/POS |

## 7. Analisis SWOT
| Strengths | Weaknesses | Opportunities | Threats |
| --- | --- | --- | --- |
| Operasional cepat, sederhana, dan dapat berjalan offline | Fitur lanjutan analitik masih terbatas | Pertumbuhan digitalisasi UMKM dan adopsi pembayaran digital | Persaingan aplikasi gratis dan proses manual berbasis spreadsheet |
| Kontrol stok lebih akurat dibanding pencatatan manual | Dukungan multi-cabang belum menjadi fokus awal | Potensi perluasan ke modul akuntansi sederhana | Risiko operasional saat listrik/perangkat bermasalah |

## 8. Estimasi Biaya dan ROI
Dokumen ini menggunakan dua skenario agar estimasi realistis dan tidak kontradiktif:

| Skenario | Komponen | Estimasi Biaya |
| --- | --- | --- |
| Implementasi siap pakai (single store) | Lisensi/aplikasi ringan + perangkat dasar | Rp500.000 - Rp3.000.000 |
| Pengembangan kustom | Pengembangan aplikasi + perangkat | Rp12.000.000 - Rp23.000.000 |

Asumsi manfaat bulanan: penghematan waktu operasional 2 jam per hari x Rp50.000 per jam x 30 hari = Rp3.000.000 per bulan.

Dengan asumsi tersebut, estimasi break-even berada pada rentang sekitar 1 sampai 8 bulan, tergantung skenario implementasi.

## 9. Key Performance Indicators (KPI)
1. Jumlah transaksi harian minimal 200 transaksi pada jam operasional normal.
2. Akurasi stok minimal 99 persen.
3. Rata-rata waktu transaksi kurang dari 10 detik.
4. Laporan operasional digunakan pemilik minimal sekali per minggu.

## 10. Rekomendasi Implementasi
1. Fase MVP (4 minggu): transaksi, stok dasar, laporan harian.
2. Fase pilot (1-2 minggu): uji lapangan pada satu toko, validasi alur kasir.
3. Fase ekspansi: penambahan sinkronisasi multi-perangkat, QRIS lebih stabil, dan laporan lanjutan.
