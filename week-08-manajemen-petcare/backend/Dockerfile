# Tahap 1: Build Binary (Menggunakan Image Golang Penuh)
FROM golang:1.21-alpine AS builder

# Setel direktori kerja
WORKDIR /app

# Salin file modul agar dapat di-*cache*
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode aplikasi
COPY . .

# Kompilasi aplikasi menjadi binary statis (CGO_ENABLED=0 agar aman jalan di Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o petcare-backend ./cmd/main.go

# Tahap 2: Runtime Image (Menggunakan Image Alpine Kosong yang Sangat Ringan)
FROM alpine:latest

# Instal sertifikat dan zona waktu untuk standar keamanan HTTP
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Salin hasil *build binary* dari tahap 1
COPY --from=builder /app/petcare-backend .

# Salin file .env agar variabel environment terbaca
COPY .env .

# Buka akses port (Sesuaikan dengan port aplikasi Anda)
EXPOSE 3000

# Eksekusi aplikasi saat container berjalan
CMD ["./petcare-backend"]