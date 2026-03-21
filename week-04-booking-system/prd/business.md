# Business Requirements Sistem Booking Lapangan Padel

Saya menginginkan sebuah sistem yang seamless bagi pelanggan agar mereka bisa mem-booking dalam waktu kurang dari 2 menit, namun memiliki arsitektur di belakang layar yang sangat solid untuk mencegah kerugian bisnis (seperti double booking) dan memudahkan pelacakan data.

Berikut adalah Business Requirements dan alur logika wajib untuk sistem lapangan padel kita.

## 1. Kriteria Sistem Wajib (Core Business Requirements)

Sebelum masuk ke alur, sistem ini wajib memiliki kriteria fundamental berikut:

### 1.1 Zero Double-Booking (Concurrency Control)

Ini prioritas absolut. Saat jam sibuk (peak hours), banyak user akan merebutkan slot yang sama. Sistem harus memiliki mekanisme temporary lock (penguncian sementara) saat user masuk ke halaman pembayaran.

### 1.2 Dynamic Pricing

Harga sewa lapangan padel biasanya berbeda antara weekdays siang (murah), weekdays malam/jam pulang kerja (peak), dan akhir pekan (weekend). Master data harus mendukung variasi harga per slot jam.

### 1.3 Automated State Management

Status transaksi harus berjalan otomatis tanpa intervensi admin (misalnya: Available $\rightarrow$ Locked $\rightarrow$ Pending Payment $\rightarrow$ Paid/Success atau Expired/Failed).

## 2. Workflow & Logika Pelanggan (User Journey)

### A. Registrasi & Login

- Daftar: Pengguna memasukkan Nama Lengkap, Email, Nomor WhatsApp aktif, dan Password. Verifikasi OTP via WhatsApp atau Email wajib dilakukan untuk memastikan validitas kontak (penting untuk marketing dan notifikasi booking).
- Login: Menggunakan Email/Nomor WA dan Password. Sistem harus menyimpan sesi (session/token) agar pengguna tidak perlu login berulang kali di perangkat mereka.

### B. Mem-booking Jadwal Lapangan

- Pencarian: Pengguna memilih tanggal bermain. Sistem menampilkan daftar lapangan (misal: Court A - Indoor, Court B - Semi-outdoor) beserta slot waktu (per 1 jam atau 1.5 jam).
- Visualisasi Ketersediaan: Slot waktu harus memiliki indikator warna yang jelas: Hijau (Tersedia), Abu-abu (Penuh), atau Kuning (Sedang diproses pengguna lain).
- Proses Booking:
	1. Pengguna memilih satu atau beberapa slot waktu yang Available.
	2. Saat klik "Lanjut ke Pembayaran", sistem melakukan penguncian (booking hold) selama misalnya 10 menit.
	3. Jika pengguna membatalkan atau waktu 10 menit habis tanpa pembayaran, slot otomatis kembali menjadi Available.

### C. Pembayaran (Checkout)

- Rincian Transaksi: Menampilkan total biaya (harga lapangan + pajak/layanan), detail tanggal, dan jam.
- Metode Pembayaran: Integrasi langsung dengan Payment Gateway pihak ketiga untuk menerima QRIS, Virtual Account (VA), atau Kartu Kredit.
- Logika Notifikasi: Begitu Payment Gateway mengirimkan callback bahwa pembayaran sukses, sistem otomatis merilis e-ticket atau kode booking dan mengirimkannya ke WhatsApp/Email pengguna.

### D. Riwayat Historis (Dashboard User)

- Upcoming Bookings: Menampilkan jadwal bermain yang akan datang beserta kode booking untuk check-in di meja resepsionis fasilitas padel.
- Past Bookings: Riwayat transaksi dan permainan yang sudah selesai. Hal ini bisa kita gunakan untuk program loyalitas (misal: main 10 kali gratis 1 jam).

## 3. Workflow & Logika Admin (Admin Journey)

Sebagai pemilik, operasional harian akan dikelola oleh admin/manajer melalui dashboard khusus. Tugas Create, Read, Update, Delete (CRUD) wajib ada pada entitas berikut:

### A. CRUD Master Data Lapangan & Harga

- Lapangan: Admin bisa menambah lapangan baru (jika ada ekspansi), menonaktifkan lapangan sementara (jika ada perbaikan/bocor), dan mengatur deskripsi.
- Pricing Engine: Admin bisa mengatur base price dan multiplier/override harga berdasarkan hari dan jam.

### B. CRUD Manajemen Jadwal (Operational Control)

- Manual Override: Memungkinkan admin untuk mem-booking lapangan secara manual (misal untuk tamu VIP, turnamen, atau pelanggan offline yang datang langsung).
- Block Time: Admin bisa memblokir slot jam tertentu untuk maintenance atau coaching clinic.

### C. CRUD Manajemen Transaksi & Pengguna

- Data Pelanggan: Melihat daftar pelanggan, riwayat booking mereka, dan total spending (untuk analisis VIP).
- Manajemen Transaksi: Melihat seluruh status transaksi (Success, Pending, Expired). Admin memiliki tombol "Refund" atau "Reschedule" dengan otorisasi khusus jika terjadi force majeure (misal lampu lapangan mati).

### D. Laporan & Analitik (Read-Only)

- Laporan Pendapatan: Generate laporan harian, mingguan, bulanan berdasarkan transaksi yang sukses.
- Tingkat Utilisasi (Utilization Rate): Metrik yang menunjukkan persentase lapangan terpakai vs tidak terpakai untuk membantu pengambilan keputusan promo/diskon.