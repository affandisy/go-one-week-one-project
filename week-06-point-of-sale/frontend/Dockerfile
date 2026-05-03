# Tahap 1: Build SvelteKit
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Tahap 2: Server Nginx
FROM nginx:alpine
# Menyalin hasil build dari tahap 1
COPY --from=builder /app/build /usr/share/nginx/html
# Menyalin konfigurasi kustom Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]