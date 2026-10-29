# TJUS Backend (PajangAset)

Backend RESTful API untuk platform konsinyasi (titip jual) aset properti dan kendaraan, dibangun menggunakan **Go (Golang)**, **Gin Framework**, dan **PostgreSQL**.

---

## 🚀 Tech Stack
- **Language**: Go (Golang) 1.23+
- **Web Framework**: Gin
- **Database**: PostgreSQL dengan connection pooling (`pgx/v5`)
- **Authentication**: Stateful JWT dengan instant revocation
- **API Documentation**: OpenAPI 3.0 / Swagger UI

---

## 📋 Fitur Utama

1. **Autentikasi & Sesi Admin**
   - Login admin dengan hashing kata sandi (Bcrypt).
   - Stateful JWT dengan pengecekan sesi aktif di database dan fitur *True Logout*.

2. **Katalog Publik & Pencarian Aset**
   - Pencarian dinamis berbasis kata kunci (judul, lokasi, kode aset).
   - Filter berdasarkan kategori serta batas harga minimum dan maksimum.
   - Auto-increment penghitung tampilan (`view_count`) saat detail aset dibuka via slug.
   - Proteksi privasi: data kontak dan catatan pemilik aset disembunyikan dari publik.

3. **Manajemen Aset (Admin)**
   - Pembuatan aset lengkap dengan atribut kustom dan data pemilik.
   - Pembaruan modular berbasis flag (data umum, pemilik, atribut, status).
   - *Soft delete* aset secara aman.

4. **Manajemen Media (Foto Aset)**
   - Upload banyak foto sekaligus (*batch upload*) dengan penamaan unik berbasis timestamp.
   - Validasi keamanan file: batas ukuran, format gambar, dan *MIME sniffing* (`http.DetectContentType`).
   - Penghapusan foto fisik dari storage disk lokal sekaligus dari database.

5. **Kategori & Atribut Dinamis**
   - Manajemen kategori aset serta atribut kustom yang fleksibel untuk tiap jenis kategori.

6. **Pengaturan Sistem (Setting)**
   - Konfigurasi identitas situs dan integrasi WhatsApp CS resmi (dilindungi autentikasi admin).

---

## 📖 Dokumentasi API
Dokumentasi API interaktif dapat diakses saat server berjalan:
- **Swagger UI**: `http://localhost:8080/swagger`
- **OpenAPI Spec**: `http://localhost:8080/swagger.yaml`

---

## ▶️ Cara Menjalankan
Jalankan aplikasi secara lokal:
```bash
go run cmd/tjus/main.go
```
Aplikasi akan aktif di port `8080`.