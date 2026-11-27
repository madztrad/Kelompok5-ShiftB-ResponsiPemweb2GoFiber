// Setel URL produksi ini ke domain HTTPS backend aaPanel sebelum mengunggah frontend.
const isLocal =
  window.location.hostname === "localhost" ||
  window.location.hostname === "127.0.0.1";
const productionApiBaseUrl = "https://api-inventaris.domainkamu.com/api";

window.APP_CONFIG = {
  API_BASE_URL: isLocal
    ? "http://localhost:3000/api"
    : productionApiBaseUrl,
};
