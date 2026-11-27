# SimInventaris — Sistem Manajemen Inventaris & Peminjaman Barang Lab

Aplikasi web full-stack untuk manajemen inventaris dan peminjaman barang laboratorium, dibangun menggunakan **Go Fiber v3 + GORM** di sisi backend dan **HTML/Tailwind CSS/Vanilla JS** di sisi frontend.

---

## Anggota Kelompok 5 (Shift B) & Pembagian Job Desk

> **Catatan:** Pastikan NIM Dimas (`H1H024032`, sebelumnya tertulis `hih204032`) dan NIM Wendy sudah benar sebelum pengumpulan.

| Nama | NIM | Job Desk | Endpoint atau Halaman yang Dikerjakan |
| :--- | :--- | :--- | :--- |
| Muhammad Faiz Mubarok | H1H024051 | Autentikasi, Middleware JWT/Role, CRUD Kategori, Shared Fetch Helper | `/api/auth/*`, `/api/kategori/*`, `login.html`, `register.html`, `profile.html`, `kategori.html`, `js/api.js`, `js/auth.js`, `js/kategori.js` |
| Dimas Rafif Zaidam | H1H024032 | Koneksi Database, Seeder, Middleware Global (CORS/Logger/Recover), CRUD Barang (Preload, Pagination, Search/Filter) | `/api/status`, `/api/barang/*`, `index.html`, `barang.html`, `js/barang.js`, `js/navbar.js` |
| Wendy Virtus | H1H0240xx | CRUD Peminjaman (Preload Barang & User, Ownership/Role), Deployment aaPanel HTTPS & Dokumentasi | `/api/peminjaman/*`, `peminjaman.html`, `js/peminjaman.js`, `js/config.js`, Deploy aaPanel & `README.md` |

---

## Prinsip Monorepo & Pengerjaan Independen

1. **1 Repo, 2 Zona Utama:** `backend/` (RESTful API Go Fiber v3) dan `frontend/` (Multi-page HTML + Tailwind CDN + Vanilla JS).
2. **1 Modul = 1 File per Layer:** Setiap resource dipisah pada `model/`, `repository/`, `service/`, `handler/`, dan `router/*_routes.go` agar 3 anggota bisa bekerja paralel tanpa konflik merge.
3. **Bisa Dikerjakan Sendiri Tanpa Saling Tunggu:** Fondasi bersama (koneksi MySQL, `AutoMigrate`, seeder akun demo `admin@lab.test` & `member@lab.test`, middleware JWT/role, dan `js/api.js`) sudah tersedia di `main`. Masing-masing anggota bisa langsung menjalankan backend lokal dan mengerjakan/menguji modulnya secara mandiri.

---

## Teknologi yang Digunakan
- **Backend:** Go (`1.26`), Fiber v3, GORM, MySQL
- **Validasi & Keamanan:** `go-playground/validator/v10`, `golang-jwt/jwt/v5`, `bcrypt`
- **Frontend:** Multi-page HTML5, Tailwind CSS (CDN), Vanilla JavaScript (`Fetch API`)

---

## Cara Menjalankan Secara Lokal

### 1. Backend
```bash
cd backend
cp .env.example .env
# Edit .env sesuai konfigurasi MySQL lokal & JWT_SECRET Anda
go mod tidy
go run ./cmd/api
```
Server akan berjalan di `http://localhost:3000` (cek kesehatan API di `GET http://localhost:3000/api/status`).

### 2. Frontend
Jalankan folder `frontend/` menggunakan static server (misal VS Code **Live Server** di port `5500`, atau `python3 -m http.server 5500` di dalam folder `frontend/`):
```bash
cd frontend
python3 -m http.server 5500
```
Buka `http://localhost:5500/login.html`.

---

## Tutorial Deploy di aaPanel (Ringkas)

Panduan lengkap langkah demi langkah tersedia di [`DEPLOYMENT_AAPANEL.md`](DEPLOYMENT_AAPANEL.md).

1. **Buat Database MySQL di aaPanel:** Menu **Databases** -> **Add DB** (`siminventaris_db`, Access: `127.0.0.1`).
2. **Build Binary Backend & Upload:**
   ```bash
   cd backend
   sh build-linux.sh
   ```
   Upload `backend/dist/siminventaris-api` ke `/www/wwwroot/api-inventaris/`, set permission `755`, buat `.env` (isi DB, `JWT_SECRET` acak panjang, dan `CORS_ORIGINS=https://inventaris.domainkamu.com`), lalu jalankan via **Supervisor Manager** / **Go Project** pada port `3000`.
3. **Domain + Reverse Proxy + SSL Backend:** Di **Website -> Add site** (`api-inventaris.domainkamu.com`, Pure static) -> **Reverse Proxy** ke `http://127.0.0.1:3000` -> aktifkan **SSL Let's Encrypt (Force HTTPS)**. Cek `https://api-inventaris.domainkamu.com/api/status`.
4. **Deploy Frontend Static + SSL:** Ubah URL produksi di `frontend/js/config.js`, buat site static baru di aaPanel (`inventaris.domainkamu.com`, Pure static), upload isi folder `frontend/`, lalu aktifkan **SSL Let's Encrypt (Force HTTPS)**.

---

## URL Deployment (Produksi)
- **Frontend:** *(isi setelah deploy, misal `https://inventaris.domainkamu.com`)*
- **Backend API:** *(isi setelah deploy, misal `https://api-inventaris.domainkamu.com`)*
- **Status API:** *(misal `https://api-inventaris.domainkamu.com/api/status`)*

---

## Akun Uji Demo (Khusus Akun Demo)
| Role | Email | Password |
| :--- | :--- | :--- |
| `admin` | `admin@lab.test` | `Admin123!` |
| `member` | `member@lab.test` | `Member123!` |

---

## Dokumentasi API

| Method | Endpoint | Keterangan | Auth | Role |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/status` | Cek status server & database | Tidak | Publik |
| `POST` | `/api/auth/register` | Registrasi akun baru (otomatis `member`) | Tidak | Publik |
| `POST` | `/api/auth/login` | Login & dapatkan JWT | Tidak | Publik |
| `GET` | `/api/auth/profile` | Lihat profil pengguna login | Ya | Semua |
| `POST` | `/api/auth/logout` | Konfirmasi logout (token dihapus di FE) | Ya | Semua |
| `GET` | `/api/kategori` | Daftar seluruh kategori | Ya | Semua |
| `GET` | `/api/kategori/:id` | Detail kategori | Ya | Semua |
| `POST` | `/api/kategori` | Tambah kategori | Ya | `admin` |
| `PUT` | `/api/kategori/:id` | Edit kategori | Ya | `admin` |
| `DELETE` | `/api/kategori/:id` | Hapus kategori | Ya | `admin` |
| `GET` | `/api/barang` | Daftar barang (`Preload("Kategori")`, pagination, search, filter) | Ya | Semua |
| `GET` | `/api/barang/:id` | Detail barang (`Preload("Kategori")`) | Ya | Semua |
| `POST` | `/api/barang` | Tambah barang | Ya | `admin` |
| `PUT` | `/api/barang/:id` | Edit barang | Ya | `admin` |
| `DELETE` | `/api/barang/:id` | Hapus barang | Ya | `admin` |
| `GET` | `/api/peminjaman` | Daftar peminjaman (`Preload("Barang")` & `Preload("User")`) | Ya | Semua (`member` milik sendiri) |
| `GET` | `/api/peminjaman/:id` | Detail peminjaman | Ya | Semua (`member` milik sendiri) |
| `POST` | `/api/peminjaman` | Ajukan peminjaman barang | Ya | Semua |
| `PUT` | `/api/peminjaman/:id` | Update peminjaman / status | Ya | `admin` (status) / `member` (catatan milik sendiri) |
| `DELETE` | `/api/peminjaman/:id` | Hapus peminjaman | Ya | `admin` |

---

## Video Presentasi YouTube
- **Tautan Video:** *(isi link YouTube Public/Unlisted sebelum pengumpulan)*
