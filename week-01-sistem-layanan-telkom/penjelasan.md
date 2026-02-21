1. Menerapkan pola Asynchronous Event-Driven Architecture yang dipadukan dengan Cache-Aside Pattern

#### Komponen Utama Sistem
Sitem terdiri dari 4 container Docker yang saling bekerja sama secara harmonis:
* API Service (Golang + Fiber): Bertindak sebagai "Pintu Gerbang" utama. Tugasnya hanya menerima request HTTP, melakukan validasi dasar, dan memberikan response secepat mungkin
* Message Broker (RabbitMQ): Bertindak sebagai "Ruang Tunggu" atau buffer. Ia menampung ribuan tiket yang masuk secara bersamaan agar database tidak panik atau mati karena kelebihan beban
* Primary Database (PostgreSQL): Bertindak sebagai "Buku Besar" permanen (Source of Truth). Ia menyimpan data mentah tiket dan statusnya dengan aman
* Cache Server (Redis): Bertindak sebagai "Buku Catatan Kecil" di atas meja kasir. Ia menyimpan data yang sering dibaca dalam memori RAM agar bisa diambil dalam hitungan milidetik tanpa perlu mengganggu PostgreSQL

#### Alur Perjalanan Data (End-to-End)
1. Alur membuat Tiket (POST /tickets) - Asynchronous Flow
Tujuan alur ini adalah memastikan pelanggan (atau simulator yang memborbardir sistem dengan 100 request per detik) tidak perlu menunggu lama
    1. Client: Mengirim komplain dalam format JSON
    2. Fiber API: Menerima JSON, men-generate TicketID (UUID) dan CreatedAt secara otomatis
    3. Producer: Alih-alih menyimpannya langsung ke database, API membungkus data tersebut dan melemparnya ke antrean RabbitMQ bernama customer.tickets
    4. Respons Cepat: Fiber langsung mengembalikkan status HTTP 201 Created ke pengguna beserta ID ticketnya. Pengguna merasa aplikasinya sangat cepat (respons dalam <10ms)
2. Alur Pekerja Latar Belakang (Consumer / Worker)
Proses ini terjadi di belakang layar, tanpa disadari oleh pelanggan
    1. Consumer Goroutine: Fungsi StartTicketConsumer() yang berjalan di background secara terus-menerus memantau antrean customer.tickets
    2. Pemrosesan: Begitu ada tiket masuk di RabbitMQ, consumer mengambilnya, mengubahnya kembali menjadi struct Golang, lalu menggunakan GORM untuk melakukan operasi INSERT ke dalam PostgreSQL
3. Alur Membaca Data (GET /tickets) - Cache-Aside Pattern
Alur ini dirancang agar agen layanan pelanggan bisa memuat daftar tiket mereka secara instan
    1. Client: Agen melakukan request data dengan parameter agent_id
    2. Cek Redis (Hit/Miss): Golang pertama-tama mencari data ke Redis menggunakan kunci rahasia tickets:active:{agent_id}
        * Cache Hit: Jika data ditemukan di Redis, Golang langsung mengembalikannya ke agen. Proses ini selesai dalam -2 milidetik
        * Cache Miss: Jika data belum ada, Golang terpaksa melakukan SELECT ke PostgreSQL. Setelah PostgreSQL memberikan datanya, Golang akan menyimpannya (fotokopi) ke Redis dengan batas waktu (Time to Live) selama 5 menit, baru kemudian mengembalikannya ke agen
4. Alur Update Tiket (PUT /tickets/:id) - Cache Invalidation
Ketika agen mengambil tiket atau mengubah statusnya menjadi selesai.
    1. Client: Agen mengirim pembaruan status (misalnya menjadi "closed").
    2. Update Database: Golang langsung melakukan perintah UPDATE ke PostgreSQL via GORM.
    3. Invalidasi Cache: Ini adalah kunci keberhasilannya! Agar agen tidak melihat data lama (karena cache Redis belum kedaluwarsa), Golang secara proaktif menghapus (delete) kunci cache tickets:active:{agent_id} di Redis.
    4. Efek Domino: Saat agen tersebut me- refresh halaman (memanggil GET lagi), Golang tidak akan menemukannya di Redis, memaksanya menarik data terbaru dari PostgreSQL, dan siklus kembali berulang.