
### 1. Visi Produk (Simplified MVP)
Aplikasi web yang memungkinkan pengguna pemula untuk menguasai dasar bahasa asing melalui modul linier yang tidak bisa dilewati (*locked progression*), memastikan fondasi bahasa terbentuk dengan kuat.

---

### 2. Struktur Modul Pembelajaran
Kita akan membagi kurikulum menjadi 3 tahap utama:

| Tahap | Fokus Utama | Output Pengguna |
| :--- | :--- | :--- |
| **Level 1: Alfabet** | Pengenalan karakter & suara. | Hafal bentuk dan bunyi huruf. |
| **Level 2: Kosakata** | Kata benda & kata kerja dasar. | Mengetahui 50-100 kata benda umum. |
| **Level 3: Kalimat** | Struktur S-P-O sederhana. | Bisa menyusun "Saya makan apel". |

---

### 3. User Stories (Kebutuhan Pengguna)

* **Sebagai Pengguna Baru**, saya ingin belajar alfabet melalui kartu gambar dan audio agar saya tahu cara pengucapan yang benar.
* **Sebagai Pelajar**, saya ingin mencocokkan kata dengan gambar (kuis visual) agar lebih mudah menghafal kosakata tanpa terjemahan yang rumit.
* **Sebagai Pengguna**, saya ingin menyusun kata-kata yang diacak menjadi satu kalimat utuh untuk melatih logika tata bahasa.

---

### 4. Alur Logika Aplikasi (User Flow)

1.  **Halaman Utama:** Menampilkan daftar modul yang terkunci (Level 1 terbuka, Level 2 & 3 terkunci).
2.  **Sesi Belajar:** * **Belajar:** Menampilkan informasi (Huruf A = Ayam).
    * **Latihan:** User menjawab soal (Pilih huruf A dari 3 pilihan).
3.  **Validasi:** Jika jawaban benar 80% dari total soal, modul berikutnya terbuka.
4.  **Simpan Progres:** Menggunakan *Local Storage* atau database sederhana agar saat user kembali, mereka tidak mengulang dari awal.

---

### 5. Skema Database Sederhana (ERD Dasar)

Untuk mendukung aplikasi ini, kita hanya butuh 3 tabel utama:

1.  **Users:** `id, username, password, current_level`
2.  **Modules:** `id, title, description, level_order` (Contoh: Level 1 - Alfabet)
3.  **Materials:** `id, module_id, content_type, question, answer, image_url, audio_url`

---

### 6. Acceptance Criteria (Kriteria Keberhasilan)

* Satu sesi belajar tidak boleh lebih dari 10 soal agar tidak melelahkan.
* Tombol "Next" hanya aktif jika user sudah memberikan jawaban (benar/salah).
* Data progres user harus tetap tersimpan meskipun browser ditutup.

---
