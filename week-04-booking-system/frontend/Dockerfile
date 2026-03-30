# Tahap 1: Pembangun (Builder)
FROM node:18-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Tahap 2: Pelaksana (Runner)
FROM node:18-alpine

WORKDIR /app

# Menyalin hasil build dan dependensi produksi
COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY package.json .

# Porta bawaan adapter-node adalah 3000
EXPOSE 3000

# Menjalankan peladen Node
CMD ["node", "build/index.js"]