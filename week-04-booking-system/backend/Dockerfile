# Tahap 1: Pembangun (Builder)
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Menyalin berkas modul dan mengunduh dependensi
COPY go.mod go.sum ./
RUN go mod download

# Menyalin seluruh kode sumber
COPY . .

# Membangun berkas biner (tanpa CGO untuk Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o padel-backend ./cmd/main.go

# Tahap 2: Pelaksana (Runner)
FROM alpine:latest

WORKDIR /app

# Menyalin biner dari tahap pembangun
COPY --from=builder /app/padel-backend .

# Menyalin folder public untuk akses PDF
COPY --from=builder /app/public ./public

# Membuka porta 3000
EXPOSE 3000

# Menjalankan aplikasi
CMD ["./padel-backend"]