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