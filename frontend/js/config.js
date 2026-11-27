// Otomatis pakai localhost saat development, dan URL produksi saat di-deploy di aaPanel.
// GANTI URL produksi di bawah sesuai domain/subdomain backend Anda di aaPanel!
const isLocal =
  window.location.hostname === "localhost" ||
  window.location.hostname === "127.0.0.1";

window.APP_CONFIG = {
  API_BASE_URL: isLocal
    ? "http://localhost:3000/api"
    : "https://api-inventaris.domainkamu.com/api",
};
