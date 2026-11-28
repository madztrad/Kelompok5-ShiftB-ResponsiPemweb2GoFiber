// Helper Fetch & UI terpusat (dikelola bersama, dirintis Anggota 1)
function getToken() {
  return localStorage.getItem("token");
}

function getUser() {
  const raw = localStorage.getItem("user");
  try {
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

function isAdmin() {
  const user = getUser();
  return user && user.role === "admin";
}

function requireAuth() {
  if (!getToken()) {
    window.location.href = "login.html";
    return false;
  }
  return true;
}

function logoutLocal() {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
  // Hindari redirect berulang saat sudah berada di halaman login.
  if (!/login\.html$/.test(window.location.pathname)) {
    window.location.href = "login.html";
  }
}

// Tanggal ISO (YYYY-MM-DD) dari zona waktu perangkat, bukan UTC.
// toISOString() memakai UTC sehingga di WIB sebelum 07:00 menghasilkan tanggal kemarin.
function todayISO() {
  const now = new Date();
  const local = new Date(now.getTime() - now.getTimezoneOffset() * 60000);
  return local.toISOString().slice(0, 10);
}

const BULAN_SINGKAT = ["Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"];

// "2026-11-28" -> "28 Nov 2026"
function formatTanggal(value) {
  if (!value) return "-";
  const parts = String(value).split("-");
  if (parts.length !== 3) return String(value);
  const bulan = BULAN_SINGKAT[Number(parts[1]) - 1];
  if (!bulan || Number.isNaN(Number(parts[2]))) return String(value);
  return `${Number(parts[2])} ${bulan} ${parts[0]}`;
}

const LABEL_KONDISI = {
  baik: "Baik",
  rusak_ringan: "Rusak Ringan",
  rusak_berat: "Rusak Berat",
};

const LABEL_STATUS = {
  dipinjam: "Dipinjam",
  dikembalikan: "Dikembalikan",
  ditolak: "Ditolak",
};

function labelKondisi(value) {
  return LABEL_KONDISI[value] || value || "-";
}

function labelStatus(value) {
  return LABEL_STATUS[value] || value || "-";
}

function escapeHtml(value) {
  return String(value ?? "")
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

function showAlert(containerId, message, type = "error") {
  const el = document.getElementById(containerId);
  if (!el) return;
  const styles =
    type === "success"
      ? "bg-emerald-50 border-emerald-300 text-emerald-800"
      : "bg-rose-50 border-rose-300 text-rose-800";
  el.className = `mb-4 rounded-lg border px-4 py-3 text-sm ${styles}`;
  el.textContent = message;
  el.classList.remove("hidden");
}

function hideAlert(containerId) {
  const el = document.getElementById(containerId);
  if (el) el.classList.add("hidden");
}

function setButtonLoading(btn, loading, loadingText = "Memproses...") {
  if (!btn) return;
  if (loading) {
    btn.dataset.originalText = btn.textContent;
    btn.textContent = loadingText;
    btn.disabled = true;
    btn.classList.add("opacity-60", "cursor-not-allowed");
  } else {
    btn.textContent = btn.dataset.originalText || btn.textContent;
    btn.disabled = false;
    btn.classList.remove("opacity-60", "cursor-not-allowed");
  }
}

// Menyusun pesan yang bisa ditindaklanjuti saat fetch gagal total: bukan respons
// error dari server, melainkan koneksi putus, diblokir CORS, atau mixed content.
function pesanGagalKoneksi() {
  const base = window.APP_CONFIG.API_BASE_URL;
  const origin = window.location.origin;

  if (window.location.protocol === "file:") {
    return (
      "Halaman ini dibuka langsung dari file (file://), sehingga browser menolak memanggil API. " +
      "Jalankan lewat server statis (mis. python3 -m http.server 5500 di folder frontend/), " +
      "lalu buka http://localhost:5500/login.html"
    );
  }
  if (window.location.protocol === "https:" && base.startsWith("http://")) {
    return (
      `Halaman ini HTTPS tetapi API masih ${base}. Browser memblokir permintaan campuran (mixed content). ` +
      "Pakai URL API berawalan https://, mis. https://api-inventaris.domainmu.com/api"
    );
  }
  return (
    `Tidak dapat menghubungi API di ${base} (dipanggil dari ${origin}). ` +
    `Periksa berurutan: (1) buka ${base}/status di tab baru — bila gagal, daemon backend atau reverse proxy mati; ` +
    `(2) CORS_ORIGINS di backend/.env harus memuat persis ${origin} tanpa trailing slash, lalu restart backend; ` +
    "(3) pastikan URL di frontend/js/config.js tidak salah ketik."
  );
}

async function apiRequest(endpoint, options = {}) {
  if (window.APP_CONFIG.API_URL_IS_PLACEHOLDER) {
    throw new Error(
      "URL API masih domain contoh. Ubah productionApiBaseUrl di frontend/js/config.js menjadi URL backend Anda, lalu muat ulang halaman."
    );
  }

  const token = getToken();
  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  let res;
  try {
    res = await fetch(`${window.APP_CONFIG.API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });
  } catch {
    // fetch hanya melempar TypeError ("Failed to fetch") untuk masalah jaringan/CORS,
    // jadi penyebabnya diterjemahkan di sini.
    throw new Error(pesanGagalKoneksi());
  }

  if (res.status === 401 && !endpoint.startsWith("/auth/login") && !endpoint.startsWith("/auth/register")) {
    logoutLocal();
    throw new Error("Sesi berakhir, silakan login ulang.");
  }

  const data = await res.json().catch(() => ({
    sukses: false,
    pesan: "Respons server tidak valid",
  }));

  if (!res.ok) {
    if (res.status === 422 && data.data && typeof data.data === "object") {
      const details = Object.values(data.data).join(" | ");
      throw new Error(`${data.pesan}: ${details}`);
    }
    throw new Error(data.pesan || `Terjadi kesalahan (${res.status})`);
  }

  return data;
}

// Memeriksa keterjangkauan API lebih awal di halaman login/register, sehingga
// penyebabnya terlihat sebelum tombol Masuk ditekan (tanpa harus buka DevTools).
async function periksaKoneksiApi(containerId = "alert-box") {
  if (window.APP_CONFIG.API_URL_IS_PLACEHOLDER) return;
  try {
    const res = await fetch(`${window.APP_CONFIG.API_BASE_URL}/status`, {
      headers: { Accept: "application/json" },
    });
    if (!res.ok) {
      showAlert(
        containerId,
        `API di ${window.APP_CONFIG.API_BASE_URL} menjawab HTTP ${res.status}, bukan 200. ` +
          "Periksa log daemon backend dan target reverse proxy.",
        "error"
      );
    }
  } catch {
    showAlert(containerId, pesanGagalKoneksi(), "error");
  }
}
