// Konfigurasi alamat API untuk seluruh halaman frontend.
// Urutan penentuan (yang pertama cocok dipakai):
// 1. Parameter ?api=<url> pada URL halaman (untuk uji cepat; disimpan di localStorage,
//    hapus dengan ?api=reset).
// 2. localStorage.apiBaseUrl.
// 3. localhost / 127.0.0.1 -> http://localhost:3000/api.
// 4. productionApiBaseUrl di bawah, bila sudah diisi dengan domain API asli.
// 5. Cadangan: origin halaman + /api, yaitu ketika API dipasang sebagai reverse proxy
//    pada domain frontend yang sama, sehingga tidak ada CORS sama sekali.
const productionApiBaseUrl = "https://api-inventaris.domainkamu.com/api";

function apiBaseUrlDariOverride() {
  const nilai = new URLSearchParams(window.location.search).get("api");
  if (!nilai) return localStorage.getItem("apiBaseUrl");
  if (nilai === "reset") {
    localStorage.removeItem("apiBaseUrl");
    return null;
  }
  const bersih = nilai.replace(/\/+$/, "");
  localStorage.setItem("apiBaseUrl", bersih);
  return bersih;
}

function tentukanKonfigurasiApi() {
  const override = apiBaseUrlDariOverride();
  if (override) return { base: override, mode: "override" };

  if (["localhost", "127.0.0.1"].includes(window.location.hostname)) {
    return { base: "http://localhost:3000/api", mode: "lokal" };
  }
  if (!productionApiBaseUrl.includes("domainkamu.com")) {
    return { base: productionApiBaseUrl, mode: "domain-api" };
  }
  return { base: `${window.location.origin}/api`, mode: "same-origin" };
}

const konfigurasiApi = tentukanKonfigurasiApi();

window.APP_CONFIG = {
  API_BASE_URL: konfigurasiApi.base,
  // "override" | "lokal" | "domain-api" | "same-origin" — ditampilkan di diagnosa.html
  // dan dipakai api.js untuk menyusun pesan galat yang sesuai konteks.
  API_MODE: konfigurasiApi.mode,
};
