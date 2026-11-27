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
  if (typeof window !== "undefined" && typeof window.location !== "undefined" && window.location.pathname !== "/login.html") {
    window.location.href = "login.html";
  }
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
