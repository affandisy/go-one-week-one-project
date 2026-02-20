### Tes Teknikal: Sistem Layanan Pelanggan PT Telekom

#### Deskripsi Proyek
Anda diminta untuk mengembangkan layanan backend untuk sistem layanan pelanggan PT Telekom yang dapat:
1. Menerima tiket pelanggan melalui RabbitMQ.
2. Menyimpan tiket dan statusnya ke PostgreSQL.
3. Menyediakan API untuk membuat tiket baru, mendapatkan daftar tiket, dan update status tiket.
4. Menggunakan Redis untuk caching daftar tiket aktif per agen.
5. Menggunakan Docker untuk menjalankan seluruh sistem.

#### Teknologi yang Digunakan
* **Golang** untuk pengembangan layanan backend.
* **MQTT (Eclipse Mosquitto)** untuk menerima tiket pelanggan secara asinkron.
* **PostgreSQL** untuk menyimpan data tiket dan riwayat.
* **Redis** untuk caching daftar tiket aktif.
* **Docker** untuk menjalankan seluruh sistem dalam container.

#### Spesifikasi Teknis
**1. Menerima Tiket via RabbitMQ**
Buat consumer RabbitMQ di Golang yang mendengarkan queue customer.tickets.
Data tiket dalam format JSON:
```json
{
  "ticket_id": "TKT-123", 
  "customer_id": "CUST-456", 
  "issue": "Koneksi lambat", 
  "priority": "high", 
  "timestamp": 1715003456
}
```
* Validasi data sebelum proses lebih lanjut.

**2. Menyimpan ke PostgreSQL**
Simpan ke tabel tickets dengan fields:
Field: ticket_id (PK), customer_id, issue, priority, status (open/closed), assigned_agent, created_at
* Implementasikan **service di Golang** yang memasukkan data ke PostgreSQL.

**3. API Endpoints (gunakan Gin atau Fiber)**
Buat **REST API** menggunakan **Fiber** dengan endpoint berikut:
1. POST /tickets untuk buat tiket baru (kirim ke RabbitMQ). Response: 
```json
{"ticket_id": "TKT-123"}.
```
2. GET /tickets?agent_id=AGENT-001&status=open untuk daftar tiket.
3. PUT /tickets/:id untuk update status.

**4. Redis Caching**
* Cache daftar tiket aktif per agen di key tickets:active:AGENT-001 (TTL 5 menit).
Update cache saat assign atau close tiket.

**5. Docker & Deployment**
Buat docker-compose.yml dengan: Golang app, PostgreSQL, RabbitMQ, Redis.
Include Dockerfile untuk Golang build multi-stage.

#### Persyaratan Tambahan
**Kandidat diwajibkan untuk membuat script yang mempublikasikan data lokasi kendaraan ke MQTT.**
* Buat producer script Golang untuk simulasi kirim tiket ke RabbitMQ.
* Logging dengan structured logs di stdout untuk Docker.

#### Cara Pengujian
| Pengujian | Hasil yang Diharapkan
| -------- | -------- 
| Konsumsi RabbitMQ | Tiket diterima dan diproses. 
| Penyimpanan PostgreSQL   | Data yang diterima disimpan di database dengan benar.   
| API Create Ticket | POST sukses, kirim ke queue.
| API List (Cache) | GET hit cache Redis cepat.
| Docker Compose | Semua layanan berjalan dengan lancar.