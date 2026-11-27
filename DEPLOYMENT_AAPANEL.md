# Panduan Deployment di aaPanel (VPS)

Panduan ini menyiapkan **Frontend** (`https://inventaris.domainkamu.com`), **Backend Go Fiber v3** (`https://api-inventaris.domainkamu.com`), dan **MySQL** di dalam satu server **aaPanel**.

---

## 1. Siapkan Database MySQL di aaPanel
1. Login ke dashboard **aaPanel** -> menu **Databases** -> **Add DB**.
2. Isi:
   - **Database name:** `siminventaris_db`
   - **Username:** `siminventaris_user`
   - **Password:** buat password kuat (simpan baik-baik)
   - **Access:** `Localhost (127.0.0.1)`
3. Klik **Submit**.

---

## 2. Build & Upload Backend Go Fiber v3
Karena `fiber/v3` versi terbaru memakai Go `1.26`, cara paling aman dan bebas error versi Go di VPS adalah compile binary statis Linux di laptop lalu upload ke aaPanel:

1. Di laptop (folder `backend/`):
   ```bash
   sh build-linux.sh
   ```
   Hasilnya ada di `backend/dist/siminventaris-api`.
   *(Jika VPS aaPanel Anda berarsitektur ARM64, ubah `GOARCH=arm64` di `build-linux.sh`.)*

2. Di aaPanel -> **Files**, buat folder `/www/wwwroot/api-inventaris`.
3. Upload file `siminventaris-api` ke `/www/wwwroot/api-inventaris/`, lalu ubah permission menjadi executable (`755` / `chmod +x siminventaris-api`).
4. Buat file `/www/wwwroot/api-inventaris/.env`:
   ```env
   PORT=3000
   APP_PORT=3000

   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=siminventaris_user
   DB_PASS=PASSWORD_DB_ANDA
   DB_NAME=siminventaris_db

   JWT_SECRET=ganti_dengan_string_acak_panjang_minimal_32_karakter
   JWT_EXPIRE_HOURS=24

   CORS_ORIGINS=https://inventaris.domainkamu.com
   ```
5. Jalankan backend agar selalu hidup:
   - **Opsi A (Go Project / Supervisor Manager di aaPanel):**
     - Buka **App Store** -> install **Supervisor Manager** (jika belum ada), lalu **Add Daemon**:
       - Name: `siminventaris-api`
       - Run User: `www` (atau `root`)
       - Run Dir: `/www/wwwroot/api-inventaris`
       - Start Command: `/www/wwwroot/api-inventaris/siminventaris-api`
     - Cek log: pastikan muncul `Server berjalan di port :3000` (AutoMigrate & seeder akun demo otomatis jalan).
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

---

## 3. Deploy Frontend (Static Site) di aaPanel
1. Edit `frontend/js/config.js` sebelum upload:
   ```js
   window.APP_CONFIG = {
     API_BASE_URL: isLocal
       ? "http://localhost:3000/api"
       : "https://api-inventaris.domainkamu.com/api",
   };
   ```
2. Di aaPanel -> **Website** -> **Add site**:
   - Domain: `inventaris.domainkamu.com`
   - Root Directory: `/www/wwwroot/inventaris`
   - PHP Version: **Pure static**
3. Upload seluruh isi folder `frontend/` (`index.html`, `login.html`, `barang.html`, folder `js/`, dll.) ke `/www/wwwroot/inventaris/`.
4. Masuk ke tab **SSL** -> **Let's Encrypt** -> **Apply**, lalu aktifkan **Force HTTPS**.

---

## 4. Verifikasi Wajib Sebelum Pengumpulan
- Buka `https://inventaris.domainkamu.com/login.html` di jendela **Incognito**.
- Login dengan akun demo (`admin@lab.test` / `Admin123!` dan `member@lab.test` / `Member123!`).
- Uji semua CRUD (`Kategori`, `Barang` dengan search/filter/pagination, `Peminjaman`).
- Refresh langsung halaman `https://inventaris.domainkamu.com/barang.html` dan `/peminjaman.html` — pastikan **tidak 404**.
- Pastikan saat login sebagai `member`, mencoba akses mutasi admin di backend mengembalikan **403**.
