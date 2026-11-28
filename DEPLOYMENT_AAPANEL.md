# Panduan Deployment di aaPanel (VPS)

Panduan ini menyiapkan **Frontend** (`https://inventaris.domainkamu.com`), **Backend Go Fiber v3** (`https://api-inventaris.domainkamu.com`), dan **MySQL** di dalam satu server **aaPanel**.
Ganti kedua domain contoh tersebut dengan domain milik kelompok. Deployment belum siap diuji sebelum DNS, URL di konfigurasi frontend, dan origin CORS memakai domain aktual yang sama.

---

## 0. Prasyarat
1. Arahkan record DNS `A` untuk domain frontend dan subdomain API ke alamat IP VPS.
2. Pastikan port 80 dan 443 dapat diakses dari internet. Port aplikasi Go (contoh `3000`) dan MySQL tidak perlu dibuka ke publik.
3. Siapkan aaPanel, Go minimal sesuai directive `go` di `backend/go.mod` pada mesin build, serta akses ke database MySQL aaPanel.

## 1. Siapkan Database MySQL di aaPanel
1. Login ke dashboard **aaPanel** -> menu **Databases** -> **Add DB**.
2. Isi:
   - **Database name:** `siminventaris_db`
   - **Username:** `siminventaris_user`
   - **Password:** buat password kuat, simpan baik-baik, dan jangan commit ke Git
   - **Access:** `Localhost (127.0.0.1)`
3. Klik **Submit**.

---

## 2. Build & Upload Backend Go Fiber v3
Cara sederhana adalah compile binary Linux di mesin lokal lalu upload ke aaPanel. Gunakan Go yang memenuhi versi di `backend/go.mod`:

1. Di laptop (folder `backend/`):
   ```bash
   sh build-linux.sh
   ```
   Hasilnya ada di `backend/dist/siminventaris-api`.
   Script menghasilkan binary Linux `amd64`. Jika VPS menggunakan ARM64, ubah target `GOARCH` di `build-linux.sh` sebelum build.

> **Jika repo di-clone langsung di server:** folder proyek Go adalah `<repo>/backend` (di situlah `go.mod` berada), bukan root monorepo. Paket `main` ada di `backend/cmd/api`, sehingga build harus dijalankan **di dalam** `backend/` dengan target eksplisit:
> ```bash
> cd <repo>/backend
> go build -o siminventaris-api ./cmd/api
> ```
> Menjalankan `go build` polos di root monorepo akan gagal karena tidak ada file Go di sana. Perhatikan juga `go.mod` meminta Go **1.26**; bila Go di server lebih lama, build akan menolak dengan pesan `go.mod requires go >= 1.26.0` — dalam kasus itu tetap pakai binary hasil build dari laptop. Yang penting: **binary dan `.env` harus berada di `<repo>/backend/`**, karena `.env` dibaca dari direktori kerja proses.

**Pembaruan berikutnya: tidak perlu upload binary lagi.** Setelah repo ada di server dan Go tersedia, cukup jalankan dari folder repo:

```bash
sh backend/deploy-server.sh
```

Script itu menarik `git pull --ff-only`, mem-build `./cmd/api` menjadi `backend/siminventaris-api`, lalu me-restart daemon lewat `supervisorctl` bila tersedia (kalau tidak, restart dari aaPanel). `.env` sengaja tidak disentuh karena tidak disimpan di Git.

2. Di aaPanel -> **Files**, buat folder `/www/wwwroot/api-inventaris`.
3. Upload file `siminventaris-api` ke `/www/wwwroot/api-inventaris/`, lalu ubah permission menjadi executable (`755` / `chmod +x siminventaris-api`).
4. Buat file `/www/wwwroot/api-inventaris/.env` di server. File ini hanya berisi rahasia produksi dan **tidak boleh dimasukkan ke repositori**:
   ```env
   PORT=3000

   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=siminventaris_user
   DB_PASS=PASSWORD_DB_ANDA
   DB_NAME=siminventaris_db

   JWT_SECRET=GANTI_DENGAN_SECRET_ACAK_MINIMAL_32_KARAKTER
   JWT_EXPIRE_HOURS=24

   CORS_ORIGINS=https://inventaris.domainkamu.com
   ```
   Buat secret acak, misalnya dengan `openssl rand -base64 48`, dan jangan memakai nilai contoh di atas. Gunakan URL frontend tanpa trailing slash untuk `CORS_ORIGINS`. Hosting sebaiknya menyediakan `PORT`; aplikasi memakai `APP_PORT` hanya sebagai fallback lokal.
5. Jalankan backend agar selalu hidup:
   - **Opsi A (Go Project / Supervisor Manager di aaPanel):**
     - Buka **App Store** -> install **Supervisor Manager** (jika belum ada), lalu **Add Daemon**:
       - Name: `siminventaris-api`
       - Run User: `www` (hindari menjalankan service sebagai `root`)
       - Run Dir: `/www/wwwroot/api-inventaris`
       - Start Command: `/www/wwwroot/api-inventaris/siminventaris-api`
     - Cek log: pastikan muncul `Server berjalan di port :3000`. Saat startup, aplikasi menjalankan AutoMigrate dan membuat akun demo jika belum tersedia.
   - **Opsi B (Website -> Go Project):**
     - Pilih path `/www/wwwroot/api-inventaris`, executable `/www/wwwroot/api-inventaris/siminventaris-api`, port `3000`.

6. Pasang Domain & HTTPS untuk Backend:
   - Buka **Website** -> **Add site** (atau langsung dari Go Project domain binding) dengan domain `api-inventaris.domainkamu.com`, PHP Version: **Pure static**.
   - Masuk ke **Settings** site -> **Reverse Proxy** -> **Add Reverse Proxy**:
     - Proxy name: `go_api`
     - Target URL: `http://127.0.0.1:3000`
     - Sent Domain: `$host`
   - Masuk ke tab **SSL** -> **Let's Encrypt** -> centang domain -> **Apply**, lalu aktifkan **Force HTTPS**.
   - Tes di browser: buka `https://api-inventaris.domainkamu.com/api/status` -> harus menjawab `200 OK` dengan `{"sukses":true,...}`.
   - Pastikan reverse proxy meneruskan semua path `/api/...` tanpa menghapus prefix `/api`.

---

## 3. Deploy Frontend (Static Site) di aaPanel
1. Edit `productionApiBaseUrl` di `frontend/js/config.js` **sebelum upload**. Nilainya harus URL API HTTPS yang sama dengan domain yang dipakai di CORS:
   ```js
   const productionApiBaseUrl = "https://api-inventaris.domainkamu.com/api";
   ```
2. Di aaPanel -> **Website** -> **Add site**:
   - Domain: `inventaris.domainkamu.com`
   - Root Directory: `/www/wwwroot/inventaris`
   - PHP Version: **Pure static**
3. Upload seluruh isi folder `frontend/` (`index.html`, `login.html`, `barang.html`, folder `js/`, dll.) ke `/www/wwwroot/inventaris/`.
4. Masuk ke tab **SSL** -> **Let's Encrypt** -> **Apply**, lalu aktifkan **Force HTTPS**.
5. Pastikan halaman seperti `login.html`, `barang.html`, `kategori.html`, dan `peminjaman.html` tersedia pada document root supaya bisa dibuka langsung atau di-refresh tanpa 404.

---

## 4. Verifikasi Wajib Sebelum Pengumpulan
- Buka `https://inventaris.domainkamu.com/login.html` di jendela **Incognito** dan pastikan tidak ada mixed-content error di DevTools.
- Cek `https://api-inventaris.domainkamu.com/api/status` mengembalikan HTTP 200 dari internet.
- Login dengan akun demo (`admin@lab.test` / `Admin123!` dan `member@lab.test` / `Member123!`).
- Uji registrasi publik, login, ketiga CRUD, validasi input, dan relasi Preload dari halaman frontend.
- Sebagai `member`, pastikan daftar/detail hanya menampilkan data milik sendiri, perubahan catatan sendiri berhasil, perubahan status/jumlah ditolak, dan operasi admin menerima **403**.
- Sebagai `admin`, pastikan bisa melihat seluruh peminjaman, mengubah status, menghapus data, dan stok barang konsisten setelah peminjaman dibuat/dikembalikan/dihapus.
- Refresh langsung halaman `login.html`, `barang.html`, `kategori.html`, dan `peminjaman.html` — pastikan **tidak 404**.
- Pastikan `frontend/js/config.js` menunjuk ke API produksi dan `CORS_ORIGINS` memuat origin frontend produksi yang tepat.


---

## 5. Kendala Umum Saat Deploy

| Gejala | Penyebab | Solusi |
| :--- | :--- | :--- |
| Log daemon berhenti dengan `JWT_SECRET wajib diisi string acak minimal 32 karakter (nilai contoh/bawaan ditolak)` | `.env` belum dibuat di folder binary, atau `JWT_SECRET` masih nilai contoh dari `.env.example` | Buat `.env` di **Run Dir** yang sama dengan binary, isi dengan hasil `openssl rand -base64 48`, lalu restart daemon |
| Log: `Gagal terhubung ke database: ... Access denied for user 'siminventaris_user'@'localhost'` | MySQL di server memakai `skip_name_resolve`, sehingga koneksi TCP dari `127.0.0.1` tidak dianggap sebagai `localhost` | Buat user/grant yang sama untuk host `127.0.0.1`, atau sesuaikan `DB_HOST` dengan host user yang terdaftar di aaPanel |
| Console browser: `TypeError: Failed to fetch` | Permintaan tidak pernah sampai ke backend: `CORS_ORIGINS` tidak memuat origin frontend, halaman HTTPS memanggil API `http://` (mixed content), atau daemon/reverse proxy mati | Ikuti urutan diagnosis di bagian **5.1** di bawah |
| Frontend memanggil `api-inventaris.domainkamu.com` padahal domain sudah diganti | `productionApiBaseUrl` di `frontend/js/config.js` belum diubah **sebelum** upload | Perbaiki `config.js`, upload ulang, lalu muat ulang halaman dengan hard reload |
| `https://api-inventaris.domainkamu.com/api/status` membalas **502** | Daemon backend mati, atau reverse proxy menunjuk port yang salah | Cek status Supervisor/Go Project, pastikan target proxy `http://127.0.0.1:3000`, dan baca log daemon |
| Refresh `barang.html` menghasilkan **404** | Isi folder `frontend/` ter-upload ke subfolder, bukan ke document root site | Upload **isi** folder `frontend/` (bukan foldernya) langsung ke root site |
| Bug lama muncul lagi saat demo/video | Browser menyimpan `js/*.js` versi lama | Hard reload (`Cmd`/`Ctrl` + `Shift` + `R`) atau gunakan jendela Incognito |

Setelah semua rampung, perbarui bagian **URL Deployment** dan **Tautan Video** pada `README.md` — dua bagian itu masih placeholder dan termasuk yang diperiksa asisten.

---

## 5.1 Diagnosis `Failed to fetch`

Halaman login dan register kini memeriksa keterjangkauan API saat dimuat, lalu menampilkan penyebabnya: URL API yang dipanggil, origin halaman, dan langkah periksa. Bila banner itu muncul, atau muncul `TypeError: Failed to fetch` saat menekan tombol, ikuti urutan ini:

1. **Buka `https://api-inventaris.domainmu.com/api/status` di tab baru.**
   - Koneksi ditolak / timeout -> daemon backend mati, atau reverse proxy salah target. Cek status Supervisor/Go Project dan baca log daemon.
   - `502 Bad Gateway` -> reverse proxy hidup tetapi tidak menemukan backend. Pastikan targetnya `http://127.0.0.1:3000` dan backend benar-benar listen di port itu.
   - `200` dengan `{"sukses":true,...}` -> API sehat, lanjut ke langkah 2.
2. **Coba login sekali lagi, lalu lihat log daemon backend.** Bila muncul baris
   `CORS: origin "https://inventaris.domainmu.com" tidak ada di CORS_ORIGINS [...]`,
   artinya origin frontend belum terdaftar. Tambahkan origin tersebut **persis** (tanpa trailing slash) ke `CORS_ORIGINS` di `backend/.env`, lalu restart backend.
3. **DevTools -> tab Network -> klik request yang gagal:**
   - `(blocked:mixed-content)` -> halaman HTTPS memanggil API `http://`. Ganti `productionApiBaseUrl` menjadi URL API `https://`.
   - `CORS error` atau respons tanpa header `Access-Control-Allow-Origin` -> ulangi langkah 2.
   - `net::ERR_NAME_NOT_RESOLVED` -> DNS subdomain API belum mengarah ke VPS.
4. **Pastikan `frontend/js/config.js` menunjuk domain API yang benar**, lalu hard reload (`Cmd`/`Ctrl` + `Shift` + `R`) karena berkas `.js` sering tersimpan di cache browser.

