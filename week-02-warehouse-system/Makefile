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