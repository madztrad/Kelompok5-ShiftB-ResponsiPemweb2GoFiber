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

  const res = await fetch(`${window.APP_CONFIG.API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

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
