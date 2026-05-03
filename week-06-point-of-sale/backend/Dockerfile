# Tahap 1: Kompilasi Golang
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Menyalin dependensi
COPY go.mod go.sum ./
RUN go mod download

# Menyalin kode dan melakukan build
COPY . .
# CGO_ENABLED=0 penting agar binary bisa berjalan mandiri di Alpine Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pos-backend ./cmd/main.go

# Tahap 2: Lingkungan Produksi Minimal
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# Menyalin binary dari tahap 1
COPY --from=builder /app/pos-backend .

# Membuat direktori untuk menyimpan database SQLite
RUN mkdir -p data

EXPOSE 3000
CMD ["./pos-backend"]