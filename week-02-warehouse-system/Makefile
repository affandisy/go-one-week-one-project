# 1. Inisialisasi Go Module
go mod init github.com/affandisy/go-one-week-one-project/week-02-warehouse-system

# 2. Buat struktur folder 3-Tier Architecture
mkdir config models repositories services handlers routes middlewares

# 3. Install dependencies yang dibutuhkan untuk MVP
go get github.com/gofiber/fiber/v2
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/joho/godotenv
go get golang.org/x/crypto/bcrypt
go get github.com/golang-jwt/jwt/v5

# 4. Jalankan docker compose
docker-compose up -d


# 5. Uji coba POST Produk
http://localhost:3000/api/v1/products

{
    "sku": "BKS-001",
    "name": "Beras Premium 5Kg",
    "description": "Beras putih pulen kualitas super",
    "category": "Sembako",
    "unit": "Sak",
    "price": 75000,
    "min_stock": 10,
    "max_stock": 100
}

# 6. Uji coba GET Produk
http://localhost:3000/api/v1/products

# 7. Uji coba POST User
http://localhost:3000/api/v1/auth/register

{
    "name": "Budi Gudang",
    "email": "budi@gudang.com",
    "password": "rahasia_budi",
    "role": "admin"
}

# 8. Uji coba POST Login
http://localhost:3000/api/v1/auth/login

{
    "email": "budi@gudang.com",
    "password": "rahasia_budi",
}

# 8. Uji coba POST Stock In
http://localhost:3000/api/v1/transactions

{
    "reference_no": "PO-2026-001",
    "type": "INBOUND",
    "entity_name": "PT Supplier Beras",
    "notes": "Restock bulanan",
    "items": [
    { "product_id": 1, "quantity": 50, "unit_price": 70000 }
    ]
}

# 9. Uji coba GET dashboard
http://localhost:3000/api/v1/dashboard

{
    "success": true,
    "message": "Dashboard WMS Metrics",
    "data": {
        "total_product_items": 1,
        "total_stock_value_idr": 2800000,
        "critical_stock_items": 0,
        "today_transactions": 2
    }
}

# 10. Uji Coba POST Adjustment
http://localhost:3000/api/v1/adjustments
{
    "product_id": 1,
    "qty": -2,
    "reason": "Kemasan rusak saat dipindahkan ke rak"
}

# 11. Uji Coba Approval Workflow
POST /api/v1/transactions (role operator) 
STATUS: draft

PUT /api/v1/transactions/1/approve (role operator)
STATUS: "Akses Ditolak"

PUT /api/v1/transactions/1/approve (role manager)
STATUS: "Transaksi Disetujui"

# 12. Uji Coba GET Analyitics
GET http://localhost:3000/api/v1/analytics/best-sellers?days=30&limit=3

{
    "success": true,
    "message": "Berhasil memuat Top Best Seller",
    "data": {
        "limit": 3,
        "periode_hari": 30,
        "items": [
            {
                "product_name": "Beras Premium 5Kg",
                "total_out": 150
            },
            {
                "product_name": "Minyak Goreng 2L",
                "total_out": 85
            }
        ]
    }
}

# 13. Uji Coba POST Data Supplier
POST http://localhost:3000/api/v1/partners
{
    "name": "PT Indofood Makmur",
    "type": "SUPPLIER",
    "phone": "08123456789",
    "address": "Kawasan Industri MM2100"
}

# 14. Uji Coba POST Data Customer
POST http://localhost:3000/api/v1/partners
{
    "name": "Toko Kelontong Berkah",
    "type": "CUSTOMER",
    "phone": "08987654321",
    "address": "Jl. Sudirman No. 10"
}

# 15. Uji Coba Transaksi Inbound dengan Validasi Silang
POST http://localhost:3000/api/v1/transactions
{
    "reference_no": "PO-TEST-001",
    "type": "INBOUND",
    "partner_id": 2, 
    "notes": "Tes Inbound dengan ID Customer",
    "items": [
    { "product_id": 1, "quantity": 10, "unit_price": 50000 }
    ]
}

# 16. Uji Coba Opname Bulanan
POST http://localhost:3000/api/v1/opname
{
    "notes": "Opname Rutin Akhir Bulan Februari 2026",
    "items": [
    {
        "product_id": 1,
        "actual_qty": 48
    },
    {
        "product_id": 2,
        "actual_qty": 10
    }
    ]
}

# 17. Uji Coba Unduh File Laporan Stok
GET http://localhost:3000/api/v1/reports/transactions/csv?month=2&year=2026

# 18. Uji Coba POST Master Lokasi
POST http://localhost:3000/api/v1/locations
{
    "code": "rak a1",
    "name": "Rak Makanan Kering Baris A",
    "description": "Kapasitas maksimal 500kg"
}

# 19. Uji Coba Batch Proses
INBOUND Batch 1: Masukkan 50 sak Beras, Batch: BRS-01, Expiry: 2026-12-01.
INBOUND Batch 2: Masukkan 50 sak Beras, Batch: BRS-02, Expiry: 2026-06-01 (Catatan: Batch 2 ini kedaluwarsanya lebih cepat!).
OUTBOUND: Lakukan permintaan pengeluaran barang sebanyak 60 sak Beras.

# 20. Uji Coba Upload Gambar Produk & User
POST http://localhost:3000/api/v1/products/1/image
Masuk ke tab Body, pilih form-data.
Di kolom Key, ketik image. Di sebelah kanan Key (tersembunyi), ada dropdown teks/file, ubah dari Text menjadi File.
Di kolom Value, klik tombol Select Files dan pilih gambar berukuran di bawah 2MB.
Masukkan Token JWT di tab Authorization.
Klik Send.

# 21. Uji Coba Movement Analytics
GET http://localhost:3000/api/v1/reports/movement-analytics?month=3&year=2026

# 22. Uji Coba Multi-Warehouse Support
1. POST http://localhost:3000/api/v1/warehouses
{
    "code": "JKT-01", 
    "name": "Gudang Utama Jakarta"
}
2. POST http://localhost:3000/api/v1/warehouses
{
    "code": "SBY-01", 
    "name": "Gudang Cabang Surabaya"
}
INBOUND (Memasukkan Barang ke Jakarta)
3. POST http://localhost:3000/api/v1/transactions
{
    "reference_no": "PO-IN-001",
    "type": "INBOUND",
    "warehouse_id": 1, 
    "partner_id": 1,
    "notes": "Barang masuk ke Jakarta",
    "items": [
        {
            "product_id": 1,
            "quantity": 50,
            "unit_price": 10000,
            "batch_no": "BATCH-JKT-A",
            "expiry_date": "2026-12-31"
        }
    ]
}
4. POST http://localhost:3000/api/v1/transactions/1/approve
5. OUTBOUND (Memindahkan Barang ke Surabaya)
POST http://localhost:3000/api/v1/transactions
{
    "reference_no": "TRF-001",
    "type": "TRANSFER",
    "warehouse_id": 1, 
    "target_warehouse_id": 2, 
    "notes": "Pindah stok ke Surabaya untuk persiapan promo",
    "items": [
        {
            "product_id": 1,
            "quantity": 20,
            "unit_price": 10000
        }
    ]
}
6. POST http://localhost:3000/api/v1/transactions/2/approve

## 23. Uji Coba WebSocket untuk Notifikasi Approve
1. Buka Postman, klik tombol New -> pilih WebSocket Request.
2. Masukkan URL: ws://localhost:3000/ws/notifications (Gunakan ws://, BUKAN http://).
3. Klik Connect. Anda akan melihat status "Connected" berwarna hijau. Biarkan tab ini terbuka (Ini mensimulasikan layar Manager yang sedang aktif).
4. Buka tab Postman biasa (HTTP). Buat transaksi Inbound baru ke POST http://localhost:3000/api/v1/transactions.
5. Klik Send.

## 24. Uji Coba Seeding Data Dummy
1. Karena skrip ini akan memasukkan data dalam jumlah masif, sangat disarankan Anda mengosongkan database terlebih dahulu 
(Anda bisa melakukan Drop Table di DBeaver/DataGrip, lalu jalankan ulang main.go agar tabelnya kosong dan bersih).
2. Buka terminal, jalankan: go run seeder.go