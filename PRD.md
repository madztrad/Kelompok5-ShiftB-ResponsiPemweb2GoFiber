# Product Requirements Document (PRD)
## SimInventaris — Sistem Manajemen Inventaris & Peminjaman Barang Lab

**Mata Kuliah:** Responsi Praktikum Pemrograman Web 2 (Teknik Komputer, UNSOED)  
**Kelompok:** Kelompok 5 — Shift B (3 Anggota)  
**Tenggat Pengumpulan:** Minggu, 29 November 2026, pukul 23:59 WIB  

---

## 1. Ringkasan Eksekutif & Rekomendasi Paling Mudah

### Mengapa Memilih Tema & Stack Ini?
Dengan waktu pengerjaan yang singkat hingga **29 November 2026 23:59 WIB**, strategi terbaik adalah memilih domain dengan **relasi data paling datar (flat)**, **bebas dari logika bisnis rumit** (seperti keranjang belanja atau kalkulasi pembayaran), dan **frontend tanpa build step yang rawan error 404 saat di-refresh**.

1. **Tema Paling Mudah: Inventaris & Peminjaman Barang Lab (`SimInventaris`)**
   - Memiliki tepat **3 resource utama** yang sangat natural: `Kategori`, `Barang`, dan `Peminjaman`.
   - Penamaan resource (`/api/barang`) selaras langsung dengan contoh di soal responsi, sehingga minim kebingungan saat implementasi.
   - Pembagian 2 role (`admin` dan `member`) sangat jelas untuk kebutuhan uji otorisasi `403 Forbidden`.
2. **Frontend Paling Mudah: Multi-Page HTML/CSS/JavaScript (Vanilla JS + Tailwind CDN)**
   - **Tanpa konfigurasi build/bundler** (`npm run build` tidak diperlukan) dan minim risiko konflik merge antar 3 orang.
   - Setiap halaman adalah file `.html` terpisah (`login.html`, `barang.html`, dll.) sehingga **me-refresh halaman selain beranda di production tidak akan pernah menghasilkan 404** (memenuhi syarat wajib deployment).
   - *(Opsional)* Jika tim sudah terbiasa dengan React/Vite, struktur endpoint dan PRD ini tetap 100% kompatibel; cukup pastikan menambahkan `vercel.json` rewrite untuk SPA.
3. **Database & Deployment Paling Mudah:**
   - **Frontend:** Vercel / Cloudflare Pages / Netlify (gratis, HTTPS otomatis, deploy hitungan detik).
   - **Backend (Go Fiber v3):** Koyeb / Render / Railway / VPS Modul 10 (HTTPS otomatis, mendukung environment variable `PORT`).
   - **Database:** MySQL (Aiven / TiDB Cloud / Railway) atau PostgreSQL (Neon.tech — gratis & stabil, diperbolehkan oleh soal jika hosting menyediakan PostgreSQL).

---

## 2. Tujuan & Cakupan Wajib (Compliance Matrix)

Seluruh fitur di bawah ini dirancang untuk memenuhi **100% syarat kelulusan** pada dokumen soal responsi:

| Syarat Soal Responsi | Implementasi pada SimInventaris |
| :--- | :--- |
| Backend Go + Fiber v3 + GORM | Go >= 1.23/1.24, `github.com/gofiber/fiber/v3`, `gorm.io/gorm` |
| Validasi & Auth | `go-playground/validator/v10`, `golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt` |
| Minimal 3 CRUD (di luar `users`/auth) | 1. `Kategori` (`/api/kategori`)<br>2. `Barang` (`/api/barang`)<br>3. `Peminjaman` (`/api/peminjaman`) |
| Minimal 1 relasi dengan `Preload` | `Barang` -> `Preload("Kategori")`; `Peminjaman` -> `Preload("Barang")` & `Preload("User")` |
| Minimal 1 endpoint dengan pagination & search/filter | `GET /api/barang?page=1&limit=10&search=laptop&kategori_id=1&kondisi=baik` |
| Auth (Register, Login, Profil, Logout) | `/api/auth/register`, `/api/auth/login`, `/api/auth/profile`, hapus JWT di frontend (+ opsional `/api/auth/logout`) |
| Minimal 2 Role + Uji 403 | Role `admin` dan `member`. Registrasi publik dipaksa `member` (tidak bisa pilih `admin`). |
| Standar API & Middleware | Struktur `cmd/api` + `internal/{config,model,repository,service,handler,middleware,router}`, `recover`, `logger`, `CORS`, global error handler, validasi `422` |
| Deployment & Integrasi | HTTPS aktif untuk frontend & backend, CORS mengizinkan domain frontend, healthcheck `GET /api/status` -> `200`, refresh non-beranda tidak 404 |

---

## 3. Arsitektur Monorepo Modular

Pendekatan **1 repositori (monorepo)** berisi dua folder terpisah (`backend/` dan `frontend/`) dengan **pemisahan 1 file per modul/resource** agar 3 anggota bisa bekerja paralel tanpa konflik merge, sekaligus tetap mematuhi struktur folder wajib dari modul.

### 3.1 Struktur Root Monorepo
```text
Kelompok5-ShiftB-ResponsiPemweb2GoFiber/
├── backend/                     # Root directory deploy Backend (Koyeb/Render/Railway/VPS)
├── frontend/                    # Root directory deploy Frontend (Vercel/Netlify/CF Pages)
├── PRD.md
└── README.md                    # Dokumentasi utama (job desk, akun demo, API table, link video)
```

### 3.2 Struktur Modular Backend (`backend/`)
Tetap memakai folder wajib modul (`cmd/api`, `internal/{config,model,handler,middleware,router}`), ditambah `repository` & `service` versi **ringan (struct langsung tanpa interface ribet)** agar kode tetap simpel dan mudah dijelaskan di video:

```text
backend/
├── cmd/
│   └── api/
│       └── main.go              # Entry point, load env, init DB, setup Fiber & middleware
├── internal/
│   ├── config/
│   │   ├── env.go               # Loader PORT / APP_PORT, DB_DSN, JWT_SECRET, CORS_ORIGIN
│   │   └── database.go          # Koneksi GORM, AutoMigrate, dan Seeder akun demo
│   ├── model/                   # 1 file = 1 entitas
│   │   ├── user.go              # Model User (password json:"-")
│   │   ├── kategori.go          # Model Kategori
│   │   ├── barang.go            # Model Barang (belongsTo Kategori)
│   │   └── peminjaman.go        # Model Peminjaman (belongsTo User & Barang)
│   ├── repository/              # (Opsional/dianjurkan) Query GORM per modul, tanpa interface
│   │   ├── user_repo.go
│   │   ├── kategori_repo.go
│   │   ├── barang_repo.go
│   │   └── peminjaman_repo.go
│   ├── service/                 # Logika bisnis & validasi aturan domain per modul
│   │   ├── auth_service.go
│   │   ├── kategori_service.go
│   │   ├── barang_service.go
│   │   └── peminjaman_service.go
│   ├── handler/                 # 1 file = 1 modul handler
│   │   ├── status_handler.go    # GET /api/status
│   │   ├── auth_handler.go
│   │   ├── kategori_handler.go
│   │   ├── barang_handler.go
│   │   └── peminjaman_handler.go
│   ├── middleware/
│   │   ├── auth.go              # WajibLogin (verifikasi Bearer JWT)
│   │   └── role.go              # HanyaPeran("admin") — diletakkan SEBELUM handler!
│   ├── router/                  # Dipisah per modul agar bebas konflik merge
│   │   ├── router.go            # Setup group /api + panggil Register*Routes
│   │   ├── auth_routes.go       # Route /api/auth/*
│   │   ├── kategori_routes.go   # Route /api/kategori/*
│   │   ├── barang_routes.go     # Route /api/barang/*
│   │   └── peminjaman_routes.go # Route /api/peminjaman/*
│   └── utils/
│       ├── response.go          # Format konsisten {sukses, pesan, data}
│       └── validator.go         # Wrapper go-playground/validator -> format error 422
├── .env.example                 # Template env (WAJIB ada di repo)
├── .gitignore                   # Wajib mengabaikan .env
└── go.mod
```

### 3.3 Struktur Modular Frontend (`frontend/`)
Multi-page statis dengan komponen navbar bersama agar setiap anggota memegang file `.html` + `.js` modulnya masing-masing dan refresh halaman non-beranda tidak pernah 404:

```text
frontend/
├── index.html                   # Redirect / Dashboard ringkas (cek token)
├── login.html                   # Halaman login
├── register.html                # Halaman registrasi
├── profile.html                 # Halaman profil & tombol logout
├── kategori.html                # Halaman CRUD Kategori (Anggota 1)
├── barang.html                  # Halaman CRUD Barang (Anggota 2)
├── peminjaman.html              # Halaman CRUD Peminjaman (Anggota 3)
└── js/
    ├── config.js                # BASE_API_URL (URL produksi, bukan localhost)
    ├── api.js                   # Helper fetch() dengan header Authorization & auto-redirect 401
    ├── navbar.js                # Komponen navbar bersama (render menu + tombol logout sesuai role)
    ├── auth.js                  # Logika login, register, profile, logout
    ├── kategori.js              # Logika CRUD Kategori
    ├── barang.js                # Logika CRUD Barang + pagination + search/filter
    └── peminjaman.js            # Logika CRUD Peminjaman
```

---

## 4. Desain Basis Data (ERD & Skema Tabel)

### Diagram Relasi
- `users (1) ────< (N) peminjaman`
- `kategori (1) ────< (N) barang`
- `barang (1) ────< (N) peminjaman`

### Detail Tabel

1. **`users`** (Tidak dihitung sebagai 3 CRUD wajib)
   | Kolom | Tipe | Keterangan |
   | :--- | :--- | :--- |
   | `id` | `uint` (PK, AutoIncrement) | ID pengguna |
   | `nama` | `varchar(100)` (Not Null) | Nama lengkap |
   | `email` | `varchar(100)` (Unique, Not Null) | Email login |
   | `password` | `varchar(255)` (Not Null) | Hash bcrypt (`json:"-"` wajib!) |
   | `role` | `varchar(20)` (Default: `'member'`) | `'admin'` atau `'member'` |
   | `created_at`, `updated_at` | `datetime` | Timestamp GORM |

2. **`kategori`** (Resource CRUD #1)
   | Kolom | Tipe | Keterangan |
   | :--- | :--- | :--- |
   | `id` | `uint` (PK, AutoIncrement) | ID kategori |
   | `nama` | `varchar(100)` (Unique, Not Null) | Contoh: Elektronik, Jaringan, Mekanik |
   | `deskripsi` | `text` | Keterangan singkat kategori |
   | `created_at`, `updated_at` | `datetime` | Timestamp GORM |

3. **`barang`** (Resource CRUD #2 — dengan `Preload("Kategori")`)
   | Kolom | Tipe | Keterangan |
   | :--- | :--- | :--- |
   | `id` | `uint` (PK, AutoIncrement) | ID barang |
   | `kode_barang` | `varchar(50)` (Unique, Not Null) | Contoh: `LAB-001` |
   | `nama` | `varchar(150)` (Not Null) | Nama barang |
   | `stok` | `int` (Not Null, >= 0) | Jumlah unit tersedia |
   | `kondisi` | `varchar(30)` (Not Null) | `'baik'`, `'rusak_ringan'`, `'rusak_berat'` |
   | `kategori_id` | `uint` (FK -> `kategori.id`) | Relasi ke tabel kategori |
   | `created_at`, `updated_at` | `datetime` | Timestamp GORM |

4. **`peminjaman`** (Resource CRUD #3 — dengan `Preload("Barang")` & `Preload("User")`)
   | Kolom | Tipe | Keterangan |
   | :--- | :--- | :--- |
   | `id` | `uint` (PK, AutoIncrement) | ID peminjaman |
   | `user_id` | `uint` (FK -> `users.id`) | Peminjam (diambil otomatis dari JWT) |
   | `barang_id` | `uint` (FK -> `barang.id`) | Barang yang dipinjam |
   | `jumlah` | `int` (Not Null, >= 1) | Jumlah unit dipinjam |
   | `tanggal_pinjam` | `varchar(20)` / `date` | Format `YYYY-MM-DD` |
   | `tanggal_kembali` | `varchar(20)` / `date` (Nullable) | Format `YYYY-MM-DD` |
   | `status` | `varchar(30)` (Default: `'dipinjam'`) | `'dipinjam'`, `'dikembalikan'`, `'ditolak'` |
   | `catatan` | `text` | Keperluan peminjaman |
   | `created_at`, `updated_at` | `datetime` | Timestamp GORM |

---

## 5. Spesifikasi API (RESTful Endpoints)

### 5.1 Standar Format Respons
Semua endpoint **wajib** mengembalikan struktur JSON yang konsisten:

- **Sukses (`200` / `201`):**
  ```json
  {
    "sukses": true,
    "pesan": "Data berhasil diambil",
    "data": { ... }
  }
  ```
- **Sukses dengan Pagination (`GET /api/barang`):**
  ```json
  {
    "sukses": true,
    "pesan": "Daftar barang berhasil diambil",
    "data": {
      "items": [ ... ],
      "meta": {
        "page": 1,
        "limit": 10,
        "total_data": 25,
        "total_halaman": 3
      }
    }
  }
  ```
- **Error Validasi (`422`):**
  ```json
  {
    "sukses": false,
    "pesan": "Validasi gagal",
    "data": {
      "nama": "Field nama wajib diisi",
      "stok": "Field stok minimal 0"
    }
  }
  ```
- **Error Lain (`400`, `401`, `403`, `404`, `500`):**
  ```json
  {
    "sukses": false,
    "pesan": "Akses ditolak: hanya admin yang diizinkan"
  }
  ```
  > **Catatan Keamanan:** Untuk error `500`, pesan ke klien harus generik (misal `"Terjadi kesalahan internal server"`), jangan pernah mengirim detail error SQL/stack trace.

---

### 5.2 Daftar Lengkap Endpoint

#### A. Sistem & Healthcheck
| Method | Endpoint | Auth | Role | Deskripsi | Status Code |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/status` | Tidak | Publik | Cek status server & koneksi DB | `200`, `500` |

#### B. Autentikasi (`/api/auth`)
| Method | Endpoint | Auth | Role | Deskripsi | Status Code |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/auth/register` | Tidak | Publik | Registrasi user baru. Role **wajib di-hardcode `'member'`** di server. | `201`, `400`, `422` |
| `POST` | `/api/auth/login` | Tidak | Publik | Verifikasi email + bcrypt password, kembalikan JWT (berlaku mis. 24 jam). | `200`, `401`, `422` |
| `GET` | `/api/auth/profile` | Ya | Semua | Ambil data user yang sedang login (tanpa field `password`). | `200`, `401` |
| `POST` | `/api/auth/logout` | Ya | Semua | Opsional konfirmasi logout (frontend wajib hapus token dari storage). | `200`, `401` |

#### C. Resource 1: Kategori (`/api/kategori`)
| Method | Endpoint | Auth | Role | Deskripsi | Status Code |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/kategori` | Ya | `admin`, `member` | Daftar seluruh kategori. | `200`, `401` |
| `GET` | `/api/kategori/:id` | Ya | `admin`, `member` | Detail satu kategori. | `200`, `401`, `404` |
| `POST` | `/api/kategori` | Ya | **`admin`** | Tambah kategori baru. Non-admin dijawab `403`. | `201`, `401`, `403`, `422` |
| `PUT` | `/api/kategori/:id` | Ya | **`admin`** | Ubah kategori. Non-admin dijawab `403`. | `200`, `401`, `403`, `404`, `422` |
| `DELETE` | `/api/kategori/:id` | Ya | **`admin`** | Hapus kategori. Non-admin dijawab `403`. | `200`, `401`, `403`, `404` |

#### D. Resource 2: Barang (`/api/barang`)
| Method | Endpoint | Auth | Role | Deskripsi | Status Code |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/barang` | Ya | `admin`, `member` | Daftar barang + `Preload("Kategori")` + query `page`, `limit`, `search`, `kategori_id`, `kondisi`. | `200`, `401` |
| `GET` | `/api/barang/:id` | Ya | `admin`, `member` | Detail barang + `Preload("Kategori")`. | `200`, `401`, `404` |
| `POST` | `/api/barang` | Ya | **`admin`** | Tambah barang baru. Non-admin dijawab `403`. | `201`, `401`, `403`, `422` |
| `PUT` | `/api/barang/:id` | Ya | **`admin`** | Edit barang. Non-admin dijawab `403`. | `200`, `401`, `403`, `404`, `422` |
| `DELETE` | `/api/barang/:id` | Ya | **`admin`** | Hapus barang. Non-admin dijawab `403`. | `200`, `401`, `403`, `404` |

#### E. Resource 3: Peminjaman (`/api/peminjaman`)
| Method | Endpoint | Auth | Role | Deskripsi | Status Code |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/peminjaman` | Ya | `admin`, `member` | `Preload("Barang")` & `Preload("User")`. `admin` melihat semua; `member` hanya melihat miliknya. | `200`, `401` |
| `GET` | `/api/peminjaman/:id` | Ya | `admin`, `member` | Detail peminjaman (+ Preload). `member` hanya boleh melihat miliknya (`403` jika milik orang lain). | `200`, `401`, `403`, `404` |
| `POST` | `/api/peminjaman` | Ya | `admin`, `member` | Ajukan peminjaman (`user_id` diambil dari JWT, bukan dari body). Cek `barang_id` ada & `jumlah <= stok`. | `201`, `400`, `401`, `404`, `422` |
| `PUT` | `/api/peminjaman/:id` | Ya | `admin`, `member` | `admin` dapat ubah `status` (`dipinjam`/`dikembalikan`/`ditolak`) & tanggal kembali. `member` hanya boleh ubah `catatan`/`jumlah` miliknya (`403` jika coba ubah status atau milik orang lain). | `200`, `401`, `403`, `404`, `422` |
| `DELETE` | `/api/peminjaman/:id` | Ya | **`admin`** | Hapus data peminjaman. Non-admin dijawab `403`. | `200`, `401`, `403`, `404` |

---

## 6. Aturan Kritis Implementasi (Anti-Jebakan Soal)

1. **Urutan Middleware Fiber v3 (WAJIB):**
   Fiber v3 mengeksekusi handler sesuai urutan penulisan. Letakkan `HanyaPeran` **sebelum** handler akhir:
   ```go
   // BENAR
   kategori.Post("/", middleware.HanyaPeran("admin"), handler.TambahKategori)

   // SALAH (handler langsung merespons tanpa cek role!)
   // kategori.Post("/", handler.TambahKategori, middleware.HanyaPeran("admin"))
   ```
2. **Perbedaan Fiber v3 vs v2:**
   - Gunakan `c.Bind().Body(&req)` (bukan `c.BodyParser(&req)`).
   - Gunakan `c.Query("search")` atau `fiber.Query[int](c, "page", 1)`.
   - Pastikan versi Go pada `go.mod` sesuai dengan versi Go yang terpasang di platform hosting.
3. **Proteksi Registrasi Publik:**
   - Struct input register hanya menerima `nama`, `email`, `password`. Jangan pernah membaca field `role` dari body request register.
4. **Password Tidak Pernah Bocor:**
   - Model `User` wajib memiliki tag `json:"-"` pada field `Password`.
5. **Frontend Auth Guard & UX:**
   - Jika halaman terlindungi (`barang.html`, `kategori.html`, `peminjaman.html`, `profile.html`) tidak memiliki token di `localStorage` atau API membalas `401`, langsung redirect ke `login.html`.
   - Setiap aksi (simpan/edit/hapus) wajib menampilkan indikator **loading** (mis. tombol disable + teks "Memproses..."), serta banner pesan **berhasil** atau **gagal** (termasuk menampilkan detail pesan validasi `422` dari server).
   - Sembunyikan/nonaktifkan tombol CRUD khusus admin jika role user adalah `member`, **namun tetap uji endpoint backend dengan akun `member` via Postman/curl untuk memastikan server mengembalikan `403`**.

---

## 7. Pembagian Kerja 3 Anggota Tim (Job Desk & Commit Plan)

Sesuai ketentuan untuk 3 anggota: *setiap anggota memegang 1 CRUD penuh (backend + frontend), lalu autentikasi, deployment, dan integrasi dibagi merata.* Setiap anggota **wajib melakukan commit dari akun GitHub masing-masing** dan **menjelaskan kodenya sendiri di video**.

| Anggota | Tanggung Jawab Utama (Backend) | Tanggung Jawab Utama (Frontend) | Tanggung Jawab Bersama (Shared) | File yang Dikuasai & Dijelaskan di Video |
| :--- | :--- | :--- | :--- | :--- |
| **Anggota 1** *(Auth & CRUD 1)* | - Setup struktur `cmd/api`, `utils/response.go`, `utils/validator.go`, global error handler<br>- **Auth & Middleware**: `Register`, `Login` (JWT + bcrypt), `Profile`, `WajibLogin`, `HanyaPeran`<br>- **CRUD 1 (`Kategori`)**: model, repo, service, handler, route | - `js/config.js`, `js/api.js` (Fetch wrapper + auto-redirect 401)<br>- `login.html`, `register.html`, `profile.html`, `js/auth.js`<br>- `kategori.html`, `js/kategori.js` | - Inisialisasi repo & `.env.example` / `.gitignore`<br>- Review integrasi auth di frontend | `middleware/auth.go`, `middleware/role.go`, `handler/auth_handler.go`, `handler/kategori_handler.go`, `js/api.js`, `kategori.html` |
| **Anggota 2** *(Infra & CRUD 2)* | - Setup `config/env.go`, `config/database.go` (koneksi DB, AutoMigrate, Seeder akun demo `admin` & `member`)<br>- Middleware global (`recover`, `logger`, `CORS`) & `GET /api/status`<br>- **CRUD 2 (`Barang`)**: model relasi `Kategori`, `Preload("Kategori")`, **pagination + search + filter**, validasi, role `admin` | - `index.html` (Dashboard navigasi & ringkasan)<br>- `barang.html`, `js/barang.js` (tabel, search input, filter dropdown kategori/kondisi, tombol pagination, form CRUD) | - Pengujian validasi `422` & uji penolakan role `403`<br>- Sinkronisasi CORS origin produksi | `config/database.go`, `router/router.go`, `repository/barang_repo.go`, `handler/barang_handler.go`, `barang.html`, `js/barang.js` |
| **Anggota 3** *(Ops & CRUD 3)* | - **CRUD 3 (`Peminjaman`)**: model relasi `User` & `Barang`, `Preload("Barang")` & `Preload("User")`, validasi stok, aturan kepemilikan `member` vs `admin` | - `peminjaman.html`, `js/peminjaman.js` (form pengajuan pinjam, daftar peminjaman, update status oleh admin, hapus) | - **Deployment Backend, DB & Frontend ke Hosting HTTPS**<br>- Uji end-to-end URL produksi (incognito & refresh non-beranda)<br>- Penyusunan `README.md` lengkap (tabel job desk, akun demo, tabel API) | `model/peminjaman.go`, `service/peminjaman_service.go`, `handler/peminjaman_handler.go`, `peminjaman.html`, `js/peminjaman.js`, konfigurasi deploy & `README.md` |

> **Tips GitHub Anti-Konflik:**
> - Buat branch per fitur/anggota (mis. `feat/auth-kategori`, `feat/barang-db`, `feat/peminjaman-deploy`) lalu merge ke `main` via Pull Request atau merge teratur.
> - Pastikan seluruh 3 akun GitHub muncul di `git log` pada branch `main`!

---

## 8. Jadwal Eksekusi Cepat (27 – 29 November 2026)

| Hari / Waktu | Target Milestone | PIC |
| :--- | :--- | :--- |
| **Jumat, 27 Nov (Malam)** | - Init Go module (`fiber/v3`, `gorm`, `jwt/v5`, `validator`, `bcrypt`) & struktur folder<br>- DB & `.env.example` siap, `GET /api/status` jalan lokal<br>- Auth (`/register`, `/login`, `/profile`) & middleware `WajibLogin` + `HanyaPeran` selesai | Anggota 1 & 2 (setup awal bersama), Anggota 3 mulai desain model `Peminjaman` |
| **Sabtu, 28 Nov (Siang)** | - Ketiga CRUD backend (`Kategori`, `Barang` + pagination/search/preload, `Peminjaman` + preload) selesai & lolos uji Postman lokal (`200`, `201`, `401`, `403`, `404`, `422`) | Anggota 1 (`Kategori`), Anggota 2 (`Barang`), Anggota 3 (`Peminjaman`) |
| **Sabtu, 28 Nov (Malam)** | - Seluruh halaman frontend selesai & terintegrasi dengan backend lokal<br>- Deploy backend + DB + frontend ke hosting HTTPS | Semua anggota untuk halaman masing-masing; Anggota 3 memimpin deploy |
| **Minggu, 29 Nov (Siang)** | - Uji lengkap di **URL produksi** (Incognito mode: register, login, 3 CRUD, penolakan role `403`, refresh halaman selain beranda)<br>- Lengkapi `README.md` (job desk, akun demo, URL deploy, tabel API) | Semua anggota |
| **Minggu, 29 Nov (Sore – Malam)** | - Rekam & upload **Video Presentasi YouTube (5 menit)** (wajah & nama terlihat, tunjukkan kode masing-masing, **JANGAN tampilkan `.env` / password DB!**)<br>- Isi **Laporan Penilaian Anggota** masing-masing & submit Google Form sebelum **23:59 WIB** | Semua anggota |

---

## 9. Checklist Luaran Sebelum Pengumpulan

### 9.1 Checklist Teknis Produksi
- [ ] `GET <URL_BACKEND>/api/status` menjawab `200 OK`.
- [ ] `js/config.js` di frontend produksi memanggil URL API produksi (bukan `localhost`).
- [ ] CORS backend mengizinkan origin frontend produksi.
- [ ] Frontend & backend berjalan di atas **HTTPS**.
- [ ] `JWT_SECRET` di produksi menggunakan string acak panjang (bukan secret default modul).
- [ ] Aplikasi membaca port dari `PORT` (dengan fallback `APP_PORT`).
- [ ] Database produksi terhubung, tabel termigrasi, dan akun demo (`admin` + `member`) tersedia.
- [ ] Registrasi, login, ketiga CRUD, dan penolakan role (`403`) diuji di URL produksi dalam jendela **Incognito**.
- [ ] Membuka/me-refresh `/barang.html`, `/kategori.html`, `/peminjaman.html`, `/login.html` tidak menghasilkan `404`.
- [ ] File `.env` tidak ada di repositori GitHub.

### 9.2 Checklist `README.md`
- [ ] Nama dan deskripsi aplikasi.
- [ ] Tabel anggota, NIM, job desk, dan endpoint/halaman yang dikerjakan (sesuai format soal).
- [ ] Daftar teknologi yang digunakan.
- [ ] Cara menjalankan secara lokal.
- [ ] URL deployment aktif (frontend & backend).
- [ ] Akun uji demo untuk setiap role (`admin` dan `member`).
- [ ] Tabel dokumentasi API (method, endpoint, keterangan, auth, role) atau tautan koleksi Postman.
- [ ] Tautan video YouTube presentasi.

### 9.3 Checklist Video Presentasi YouTube (~5 Menit, Public/Unlisted)
- [ ] Pengenalan singkat aplikasi + demo langsung dari **URL produksi** (~1,5 menit).
- [ ] Penjelasan kode oleh Anggota 1 (~1 menit): tunjukkan kode auth/middleware/CRUD Kategori & cara kerjanya, wajah + nama terlihat.
- [ ] Penjelasan kode oleh Anggota 2 (~1 menit): tunjukkan kode database/CORS/CRUD Barang (pagination, search, Preload) & cara kerjanya, wajah + nama terlihat.
- [ ] Penjelasan kode oleh Anggota 3 (~1 menit): tunjukkan kode CRUD Peminjaman & setup integrasi/deployment, wajah + nama terlihat.
- [ ] **Pastikan `.env`, `JWT_SECRET`, dan password database TIDAK tampil di layar!**

### 9.4 Checklist Individu
- [ ] Setiap anggota mengisi **Laporan Penilaian Anggota** secara jujur & objektif (rencana vs realisasi job desk, keaktifan, kendala, nilai kontribusi 1–5).
- [ ] Pengumpulan ke Google Form asisten dilakukan sebelum **2026-11-29 23:59 WIB**.
