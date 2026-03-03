A. Mengatasi Race Condition (Masalah Concurrency)
Skenario: Karena PostgreSQL tidak melakukan kalkulasi (UPDATE stock = stock - 1), apa yang terjadi jika Operator A dan Operator B menekan tombol Approve untuk barang yang sama di milidetik yang sama persis?

Golang A membaca stok: 10

Golang B membaca stok: 10

Golang A menghitung sisa: 10 - 5 = 5. Disimpan ke DB.

Golang B menghitung sisa: 10 - 2 = 8. Disimpan ke DB.
Hasil Akhir di DB: 8 (Padahal seharusnya 10 - 5 - 2 = 3).
Improvement: PRD perlu ditambahkan standar Concurrency Control. Di Golang, kita bisa menggunakan Optimistic Locking (menambahkan field Version di tabel Product) atau menggunakan Distributed Lock menggunakan Redis di Service Layer sebelum Golang membaca stok.

B. Optimasi Export Laporan untuk Jutaan Data
Skenario: Di PRD tertulis target Scalability adalah 1.000.000+ transaksi. Saat ini, Golang menarik semua data transaksi selama sebulan ke dalam RAM (bytes.Buffer) untuk membuat CSV. Jika datanya 1 juta baris, RAM server akan penuh (Out of Memory) dan request HTTP akan timeout (karena PRD meminta response di bawah 10 detik).
Improvement: Ubah alur Export CSV menjadi Asynchronous. Saat Manager klik Export, API langsung merespons "Laporan sedang diproses". Golang akan menggunakan Goroutine untuk mencicil query (chunking 10.000 baris), menulisnya ke GCS (Google Cloud Storage), dan mengirimkan link download-nya via Email atau Notifikasi Web.

C. Menangani Anomali Soft-Delete
Skenario: Fitur Soft Delete (DeletedAt) sudah ada. Jika "Barang X" dinonaktifkan (di-soft delete), lalu ternyata masih ada "Transaksi Draf" yang menggunakan barang tersebut, apa yang terjadi saat Manager menekan Approve?
Improvement: Service Layer di ApproveTransaction harus divalidasi ulang: "Apakah barang ini IsActive == false atau sudah terhapus?" Jika ya, transaksi harus ditolak agar stok riwayat tidak berantakan.

D. Persiapan Multi-Warehouse (Phase 4 PRD)
Skenario: PRD menyebutkan kesiapan untuk "Multi-warehouse (future)". Desain database kita saat ini mengikat field CurrentStock langsung di tabel Product. Jika nanti ada Gudang Jakarta dan Gudang Bali, struktur ini akan rusak.
Improvement: Mulai Phase 4, PRD harus merombak skema. Tabel Product hanya berisi nama dan SKU. Kita butuh tabel Pivot bernama WarehouseStock (relasi antara WarehouseID dan ProductID) untuk menyimpan kuantitas stok masing-masing cabang.