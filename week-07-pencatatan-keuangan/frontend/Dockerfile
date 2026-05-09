# Tahap Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Tahap Produksi
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
# Konfigurasi nginx akan disuntikkan via docker-compose
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]