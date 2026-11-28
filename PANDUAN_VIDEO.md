# Panduan Video Presentasi & Pembagian Penjelasan Kode

Dokumen ini adalah naskah kerja untuk merekam **video presentasi SimInventaris (±5 menit)** yang diunggah ke YouTube. Isinya: susunan waktu, langkah demo, file yang ditunjukkan **setiap anggota beserta kalimat inti yang perlu diucapkan**, dan checklist sebelum menekan tombol rekam.

Aturan dari soal responsi yang **wajib** dipenuhi:

- Video berdurasi sekitar 5 menit, diunggah ke YouTube sebagai **Public** atau **Unlisted**.
- **Wajah dan nama setiap anggota terlihat** (webcam aktif + overlay nama/NIM).
- Setiap anggota **menjelaskan kode miliknya sendiri**, bukan diwakilkan.
- Demo dilakukan dari **URL produksi (HTTPS)**, bukan `localhost`.
- **Jangan pernah menampilkan** isi `.env`, `JWT_SECRET`, password database, atau token JWT di layar.

---

## 1. Susunan Waktu (5 menit)

| Waktu | Segmen | PIC | Yang tampil di layar |
| :--- | :--- | :--- | :--- |
| 0:00 – 0:20 | Pembukaan: nama kelompok, anggota + NIM, nama aplikasi | Faiz | Slide/README |
| 0:20 – 1:30 | Demo aplikasi end-to-end dari URL produksi | Dimas | Browser (dua jendela: admin & member) |
| 1:30 – 2:30 | Penjelasan kode: Auth, Middleware, CRUD Kategori, shared helper | **Faiz** | Editor + terminal |
| 2:30 – 3:30 | Penjelasan kode: Env, Database, Middleware global, CRUD Barang | **Dimas** | Editor |
| 3:30 – 4:30 | Penjelasan kode: CRUD Peminjaman + deployment | **Wendy** | Editor + panel aaPanel |
| 4:30 – 5:00 | Penutup: bukti `git log` 3 akun, URL frontend & backend | Wendy | Terminal + README |

Kalau waktu mepet, potong demo (0:20–1:30), **jangan** potong bagian penjelasan kode — bagian itu yang dinilai per anggota.

---

## 2. Persiapan Sebelum Merekam

Lakukan semuanya sebelum rekaman dimulai, supaya tidak ada waktu terbuang untuk menunggu.

1. **Dua sesi browser** agar tidak bolak-balik login:
   - Jendela normal → login `admin@lab.test`
   - Jendela Incognito / profil kedua → login `member@lab.test`
2. **Data contoh sudah ada** di database produksi: minimal 1 kategori, 2 barang, 1 peminjaman berstatus `dipinjam`.
3. **Editor** dengan font ≥ 16pt, zoom 110–125%, dan word-wrap mati agar nomor baris tidak berantakan.
4. **Tutup semua tab `.env`.** Bila perlu menunjukkan konfigurasi, buka `.env.example` (tanpa nilai rahasia).
5. **Terminal** siap di folder `backend/` untuk menampilkan `go test ./...` dan `git log --oneline`.
6. **Matikan notifikasi** (Focus Mode / Do Not Disturb) agar tidak ada pop-up masuk.
7. **Muat ulang halaman dengan hard reload** (`Cmd` + `Shift` + `R`) atau pakai jendela Incognito. Browser sering menyimpan `js/*.js` lama, dan versi lama itu masih menampilkan bug yang sudah diperbaiki (label kondisi huruf kecil, tanggal ISO, dropdown status di form tambah). Kalau temuan lama muncul lagi saat merekam, hampir pasti penyebabnya cache — bukan kode.
8. Catat urutan klik demo di kertas — jangan improvisasi saat merekam.

---

## 3. Demo Aplikasi (0:20 – 1:30) — langkah klik

Dibawakan **Dimas**, sementara Faiz & Wendy sudah siap di editor masing-masing.

1. Buka `https://<domain-frontend>/login.html` → login sebagai **admin**.
2. Dashboard → **Kategori**: tambahkan `Elektronik`. Tunjukkan banner hijau *"Kategori berhasil ditambahkan"*.
3. **Barang**: tambahkan satu barang (kode `LAB-001`, stok 5). Tunjukkan kolom **Kategori** terisi — sebutkan bahwa itu hasil `Preload("Kategori")`.
4. Masih di Barang: ketik kata kunci di kolom cari, ganti filter kondisi, lalu klik **Berikutnya** → tunjukkan teks *"Halaman X dari Y"*.
5. **Peminjaman**: ajukan pinjam 2 unit → kembali ke **Barang**, tunjukkan **stok turun dari 5 ke 3**.
6. Kembali ke Peminjaman sebagai admin → ubah status menjadi **Dikembalikan** → tunjukkan **stok kembali ke 5**.
7. Logout → login sebagai **member**:
   - Kategori: **panel form tambah hilang**, kolom aksi hanya berisi `Detail`.
   - Peminjaman: hanya melihat **miliknya sendiri**, dan tidak ada tombol `Hapus`.
8. Terakhir, buka langsung `https://<domain-frontend>/barang.html` lalu **refresh** → tunjukkan **tidak 404**.

> Kalimat penutup segmen: *"Semua data di atas berasal dari API Go Fiber v3 dengan MySQL, dan sekarang kita lihat kodenya satu per satu."*

---

## 4. Faiz — Autentikasi, Middleware, CRUD Kategori, Shared Fetch Helper

**Pembuka (ucapkan sekali):** *"Saya Muhammad Faiz Mubarok, NIM H1H024051. Saya mengerjakan autentikasi, middleware JWT dan role, CRUD Kategori, serta helper fetch bersama di frontend."*

Tunjukkan file dengan urutan berikut. Total ±60 detik.

### 4.1 `backend/internal/middleware/auth.go` — `WajibLogin`
**Yang ditunjukkan:** baris pembacaan header `Authorization` sampai `c.Next()`.

**Kalimat inti:** *"Middleware ini mengambil header `Authorization`, memastikan formatnya `Bearer <token>`, lalu memverifikasi tanda tangan JWT dengan metode HMAC dan secret dari environment. Kalau token tidak ada, salah tanda tangan, atau sudah kedaluwarsa, langsung dibalas 401. Kalau valid, `user_id` dan `role` disimpan di `c.Locals` supaya bisa dipakai handler dan middleware berikutnya."*

**Poin yang sering ditanya penguji:** kenapa `jwt.Parse` aman → karena algoritma dibatasi ke `SigningMethodHMAC`, sehingga token dengan `alg: none` ditolak.

### 4.2 `backend/internal/middleware/role.go` — `HanyaPeran`
**Kalimat inti:** *"Ini pengecekan role. Nilai `role` diambil dari `c.Locals` yang tadi diisi `WajibLogin`, lalu dicocokkan dengan daftar role yang diizinkan. Kalau tidak cocok, keluar 403 Forbidden."*

### 4.3 `backend/internal/router/kategori_routes.go` — urutan middleware
**Yang ditunjukkan:** baris `kategori.Post("/", middleware.HanyaPeran("admin"), h.Create)`.

**Kalimat inti:** *"Fiber v3 menjalankan handler sesuai urutan penulisan. Karena itu `HanyaPeran` saya letakkan sebelum handler, bukan sesudahnya — kalau dibalik, handler sudah merespons duluan dan pengecekan role tidak pernah jalan."*

### 4.4 `backend/internal/service/auth_service.go` — register & login
**Yang ditunjukkan:** dua hal saja.
1. `Register` → baris `Role: "member"`.
   **Kalimat inti:** *"Struct input register hanya punya nama, email, dan password. Role tidak pernah dibaca dari body, jadi kalau ada yang mengirim `role: admin`, nilainya diabaikan dan tetap tersimpan sebagai member."*
2. `Login` → `bcrypt.CompareHashAndPassword` lalu pembuatan `jwt.MapClaims` dengan `exp`.

   **Kalimat inti:** *"Password tidak pernah dibandingkan sebagai teks biasa, tapi lewat bcrypt. Token yang dikembalikan berisi `user_id`, `role`, waktu terbit, dan waktu kedaluwarsa dari `JWT_EXPIRE_HOURS`."*

### 4.5 `backend/internal/model/user.go`
**Kalimat inti:** *"Field `Password` diberi tag `json:\"-\"`, jadi hash password tidak akan pernah ikut terkirim di respons JSON apa pun — sudah saya buktikan dengan membuka `/api/auth/profile` di browser."*

### 4.6 `backend/internal/handler/auth_handler.go` — pola handler
**Kalimat inti:** *"Semua handler mengikuti pola yang sama: bind body, validasi, panggil service, lalu balas dengan format konsisten `{sukses, pesan, data}`. Kalau body bukan JSON valid, jawabannya 400. Kalau JSON-nya valid tapi nilainya salah, jawabannya 422 lewat validator."*

Tunjukkan juga `handleServiceError`: *"Error domain seperti 'email sudah terdaftar' dibawa oleh `AppError` sehingga status code-nya tepat. Error tak terduga dicatat di log server, tapi klien hanya menerima pesan generik — detail database tidak bocor."*

### 4.7 `backend/internal/handler/kategori_handler.go` + `service/kategori_service.go`
**Kalimat inti:** *"Ini CRUD penuh saya. Tiga aturan bisnisnya: nama kategori tidak boleh duplikat, pencarian nama memakai `LOWER()` supaya tidak peka huruf besar-kecil, dan kategori yang masih dipakai barang tidak boleh dihapus."*

Tunjukkan `Delete` → pemanggilan `CountBarangByKategori`.

### 4.8 `frontend/js/api.js` — `apiRequest`
**Kalimat inti:** *"Ini satu-satunya pintu keluar request dari frontend. Fungsinya menempelkan header `Authorization: Bearer <token>` otomatis. Kalau server membalas 401, token lokal dihapus dan pengguna diarahkan ke halaman login — kecuali untuk endpoint login dan register, supaya pesan 'email atau password salah' tetap tampil di form. Kalau 422, detail per field dari server dirangkai menjadi satu pesan yang bisa dibaca pengguna."*

Sebutkan juga helper tampilan yang ada di file ini (`formatTanggal`, `labelKondisi`, `labelStatus`) — *"supaya semua halaman menampilkan label yang seragam."*

### 4.9 `frontend/kategori.html` + `frontend/js/kategori.js`
**Kalimat inti:** *"Panel form hanya dibuka kalau `isAdmin()` benar, dan tombol Edit/Hapus juga hanya dirender untuk admin. Ini murni kenyamanan UI — pengaman sebenarnya tetap di backend, dan saya buktikan sekarang."*

### 4.10 Bukti `403` di terminal (opsional, sangat menguatkan)
```bash
# login sebagai member, ambil token, lalu coba tambah kategori
curl -i -X POST https://<domain-api>/api/kategori \
  -H "Authorization: Bearer <TOKEN_MEMBER>" \
  -H "Content-Type: application/json" \
  -d '{"nama":"Tes"}'
# → HTTP/1.1 403 Forbidden
```

**Penutup segmen:** *"Bagian saya memastikan semua endpoint selain status, login, dan register hanya bisa diakses dengan token valid, dan aksi tulis hanya untuk admin. Selanjutnya Dimas menjelaskan fondasi database dan CRUD Barang."*

---

## 5. Dimas — Environment, Database, Middleware Global, CRUD Barang

**Pembuka:** *"Saya Dimas Rafif Zaidam, NIM H1H024032. Saya mengerjakan koneksi database, seeder, middleware global, healthcheck, dan CRUD Barang dengan pagination serta pencarian."*

### 5.1 `backend/internal/config/env.go`
**Kalimat inti:** *"Semua konfigurasi dibaca dari environment. Port dibaca dari `PORT`, dan `APP_PORT` hanya fallback untuk lokal, karena platform hosting selalu menyediakan `PORT`. Sejak versi ini, `JWT_SECRET` wajib minimal 32 karakter dan nilai contoh di `.env.example` ditolak, jadi server tidak mungkin jalan dengan secret default."*

Tunjukkan juga pemecahan `CORS_ORIGINS` berdasarkan koma.

### 5.2 `backend/internal/config/database.go`
**Kalimat inti:** *"Di sini DSN MySQL dibentuk dari environment, lalu GORM membuka koneksi. `AutoMigrate` membuat atau menyesuaikan keempat tabel — users, kategori, barang, peminjaman — sehingga tidak perlu SQL manual. Setelah itu seeder memeriksa apakah akun demo sudah ada; kalau belum, dibuat dengan password yang di-hash bcrypt."*

### 5.3 `backend/cmd/api/main.go` — middleware global
**Kalimat inti:** *"Tiga middleware global dipasang berurutan: `recover` supaya panic tidak mematikan server, `logger` untuk mencatat setiap request, dan `cors` yang hanya mengizinkan origin frontend yang terdaftar. Selain itu ada `ErrorHandler` global: semua error diubah ke format JSON yang sama, dan detail internal tidak dikirim ke klien."*

### 5.4 `backend/internal/handler/status_handler.go`
**Kalimat inti:** *"`GET /api/status` melakukan `Ping()` ke database. Kalau database tidak merespons, jawabannya 500. Endpoint ini yang saya pakai sebagai healthcheck setelah deploy."*

### 5.5 `backend/internal/router/router.go`
**Kalimat inti:** *"Semua route dikumpulkan di grup `/api`, lalu tiap modul mendaftarkan routenya sendiri di file terpisah. Karena satu modul satu file, tiga orang bisa bekerja paralel tanpa konflik merge."*

### 5.6 `backend/internal/repository/barang_repo.go` — inti pagination & search
**Kalimat inti:** *"Query disusun bertahap. Pencarian memakai `LOWER(nama) LIKE` atau `LOWER(kode_barang) LIKE`, lalu filter kategori dan kondisi ditambahkan kalau ada. `Count` dijalankan dulu untuk mendapatkan total data, baru `Limit` dan `Offset` mengambil satu halaman. `Preload("Kategori")` membuat setiap barang membawa data kategorinya dalam satu query tambahan, bukan query per baris."*

**Kalau ditanya soal N+1:** *"`Preload` menghasilkan satu query terpisah dengan `IN (...)`, jadi tidak ada query per baris."*

### 5.7 `backend/internal/service/barang_service.go`
**Kalimat inti:** *"Service menormalkan `page` dan `limit` — limit dibatasi maksimal 100 — lalu menghitung `total_halaman` dengan pembulatan ke atas. Di sini juga ada aturan: kode barang disimpan huruf besar dan tidak boleh duplikat, kategori harus benar-benar ada, dan barang yang punya riwayat peminjaman tidak boleh dihapus."*

### 5.8 `backend/internal/handler/barang_handler.go`
**Kalimat inti:** *"Handler hanya memindahkan query param `page`, `limit`, `search`, `kategori_id`, dan `kondisi` menjadi satu struct filter, lalu menyerahkannya ke service. Jadi handler tetap tipis dan logikanya terpusat di service."*

### 5.9 `frontend/barang.html` + `frontend/js/barang.js`
**Kalimat inti:** *"Halaman ini memuat daftar barang, dropdown kategori untuk filter, kolom pencarian, dan tombol halaman. Setiap kali filter berubah, `currentPage` dikembalikan ke 1 supaya hasil pencarian tidak 'hilang' di halaman kosong. Label kondisi seperti `rusak_ringan` diterjemahkan menjadi 'Rusak Ringan' agar mudah dibaca."*

**Penutup segmen:** *"Jadi pagination, pencarian, dan relasi kategori semuanya ditangani di sisi server. Berikutnya Wendy menjelaskan modul peminjaman dan proses deployment."*

---

## 6. Wendy — CRUD Peminjaman, Stok, Kepemilikan Data, dan Deployment

**Pembuka:** *"Saya Wendy Virtus, NIM H1H0240xx. Saya mengerjakan CRUD Peminjaman beserta aturan stok dan kepemilikan data, lalu memimpin deployment ke aaPanel."*
*(Ganti NIM placeholder di atas dan di README sebelum pengumpulan.)*

### 6.1 `backend/internal/model/peminjaman.go`
**Kalimat inti:** *"Tabel peminjaman punya dua relasi: ke `User` lewat `user_id` sebagai peminjam, dan ke `Barang` lewat `barang_id`. Statusnya bertiga: dipinjam, dikembalikan, ditolak."*

### 6.2 `backend/internal/repository/peminjaman_repo.go`
**Yang ditunjukkan:** baris `Preload("User").Preload("Barang").Preload("Barang.Kategori")` dan `Where("user_id = ?")`.

**Kalimat inti:** *"Di sini aturan kepemilikan diterapkan di level query. Kalau pemanggilnya bukan admin, query otomatis difilter `user_id = ?` sehingga member secara struktural tidak mungkin menerima data orang lain — bukan hanya disembunyikan di tampilan."*

### 6.3 `backend/internal/service/peminjaman_service.go` — bagian paling penting
Tunjukkan tiga hal berurutan.

1. **`Create` di dalam `Transaction`**
   **Kalimat inti:** *"Pembuatan peminjaman berjalan dalam satu transaksi. Barang diambil dulu, stok diperiksa, stok dikurangi, baru data peminjaman dibuat. Kalau salah satu langkah gagal, semuanya dibatalkan sehingga stok tidak pernah berkurang tanpa catatan peminjaman."*
2. **Perpindahan status di `Update`**
   **Kalimat inti:** *"Status hanya boleh diubah admin. Saat status berpindah dari `dipinjam` ke `dikembalikan` atau `ditolak`, stok dikembalikan sejumlah unit yang dipinjam. Kalau peminjaman diaktifkan kembali, stok ditahan lagi dan diperiksa kecukupannya."*
3. **`validatePeminjamanUpdateAccess`**
   **Kalimat inti:** *"Ini pemisahan hak ubah. Admin bebas mengubah jumlah, tanggal, dan status. Member hanya boleh mengubah catatan pada peminjaman miliknya; kalau mencoba mengubah jumlah, tanggal, atau status, hasilnya 403."*

Tunjukkan juga `Delete`: *"Menghapus peminjaman yang masih berstatus `dipinjam` mengembalikan stoknya terlebih dahulu."*

### 6.4 `backend/internal/handler/peminjaman_handler.go`
**Kalimat inti:** *"`user_id` peminjam diambil dari `c.Locals` hasil verifikasi JWT, bukan dari body request. Jadi pengguna tidak bisa memalsukan peminjaman atas nama orang lain."*

### 6.5 `backend/internal/router/peminjaman_routes.go`
**Kalimat inti:** *"Route `PUT` sengaja tidak diberi `HanyaPeran`, karena member masih boleh mengubah catatannya sendiri. Penentuan haknya dilakukan di service lewat `validatePeminjamanUpdateAccess`. Hanya `DELETE` yang dikunci admin."*

### 6.6 `frontend/peminjaman.html` + `frontend/js/peminjaman.js`
**Kalimat inti:** *"Saat mengedit, member mendapat catatan penjelasan dan field jumlah serta tanggal dinonaktifkan; hanya catatan yang bisa diubah. Dropdown status hanya muncul untuk admin saat mengedit, karena peminjaman baru selalu dimulai dengan status `dipinjam` di server."*

### 6.7 Deployment aaPanel
Tunjukkan panel aaPanel sebentar, lalu file `backend/build-linux.sh` dan `frontend/js/config.js`.

**Kalimat inti:** *"Backend di-build menjadi binary Linux, dijalankan lewat Supervisor sebagai user `www`, dan diekspos melalui reverse proxy aaPanel dari `127.0.0.1:3000` ke subdomain HTTPS. Frontend di-upload sebagai situs statis di domain terpisah. Dua hal yang paling sering membuat deploy gagal: `productionApiBaseUrl` di `config.js` masih menunjuk domain contoh, dan `CORS_ORIGINS` di `.env` tidak sama persis dengan origin frontend — keduanya sudah saya samakan."*

**Penutup:** buka `https://<domain-api>/api/status` → tunjukkan `200` dan `{"sukses":true,...}`.

---

## 7. Penutup (4:30 – 5:00)

1. Tampilkan `git log --oneline --format="%h %an %s"` → pastikan **tiga nama anggota** muncul sebagai bukti commit masing-masing.
2. Tampilkan bagian **URL Deployment** pada `README.md`.
3. Ucapkan penutup singkat: *"Frontend di <URL frontend>, API di <URL API>, akun demo tersedia di README. Terima kasih."*

---

## 8. Larangan Saat Merekam

- Menampilkan isi `.env`, `JWT_SECRET`, password database, atau token JWT utuh di layar.
- Membuka file yang bukan tanggung jawab sendiri lalu mengklaim mengerjakannya.
- Demo dari `localhost` sebagai pengganti URL produksi.
- Merekam satu suara saja untuk seluruh video — soal mewajibkan tiap anggota tampil.

---

## 9. Checklist Akhir per Anggota

**Faiz**
- [ ] `middleware/auth.go`, `middleware/role.go`, `handler/auth_handler.go`, `handler/kategori_handler.go`, `js/api.js`, `kategori.html` dijelaskan.
- [ ] Bukti `403` untuk `member` ditampilkan (curl atau UI).
- [ ] Commit atas nama Muhammad Faiz Mubarok ada di `main`.

**Dimas**
- [ ] `config/database.go`, `config/env.go`, `cmd/api/main.go`, `router/router.go`, `repository/barang_repo.go`, `handler/barang_handler.go`, `barang.html`, `js/barang.js` dijelaskan.
- [ ] Pagination + search + filter didemokan hidup.
- [ ] Commit atas nama Dimas Rafif Zaidam ada di `main`.

**Wendy**
- [ ] `model/peminjaman.go`, `service/peminjaman_service.go`, `handler/peminjaman_handler.go`, `peminjaman.html`, `js/peminjaman.js`, konfigurasi deploy dijelaskan.
- [ ] Alur stok (turun saat pinjam, kembali saat dikembalikan) ditunjukkan.
- [ ] URL produksi dan `GET /api/status` diperlihatkan.
- [ ] NIM di `README.md` sudah diisi lengkap.
