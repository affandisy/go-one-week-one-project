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