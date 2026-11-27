// Dikelola oleh Anggota 2: CRUD Barang + Preload Kategori + Search/Filter/Pagination
let currentPage = 1;
const pageLimit = 5;

document.addEventListener("DOMContentLoaded", () => {
  if (!requireAuth()) return;
  renderNavbar("barang");

  if (isAdmin()) {
    document.getElementById("admin-form-panel")?.classList.remove("hidden");
  }

  document.getElementById("barang-form")?.addEventListener("submit", handleSaveBarang);
  document.getElementById("btn-reset")?.addEventListener("click", resetBarangForm);
  document.getElementById("filter-form")?.addEventListener("submit", (e) => {
    e.preventDefault();
    currentPage = 1;
    loadBarang();
  });
  document.getElementById("btn-prev")?.addEventListener("click", () => {
    if (currentPage > 1) {
      currentPage--;
      loadBarang();
    }
  });
  document.getElementById("btn-next")?.addEventListener("click", () => {
    currentPage++;
    loadBarang();
  });

  loadKategoriOptions();
  loadBarang();
});

async function loadKategoriOptions() {
  try {
    const res = await apiRequest("/kategori");
    const items = res.data || [];
    const options = items
      .map((k) => `<option value="${k.id}">${escapeHtml(k.nama)}</option>`)
      .join("");

    const formSelect = document.getElementById("kategori_id");
    if (formSelect) {
      formSelect.innerHTML = `<option value="">-- Pilih Kategori --</option>${options}`;
    }
    const filterSelect = document.getElementById("filter-kategori");
    if (filterSelect) {
      filterSelect.innerHTML = `<option value="">Semua Kategori</option>${options}`;
    }
  } catch (err) {
    showAlert("alert-box", `Gagal memuat opsi kategori: ${err.message}`, "error");
  }
}

async function loadBarang() {
  const tbody = document.getElementById("barang-tbody");
  tbody.innerHTML = `<tr><td colspan="7" class="px-4 py-6 text-center text-slate-500">Memuat data barang...</td></tr>`;

  const params = new URLSearchParams({
    page: String(currentPage),
    limit: String(pageLimit),
  });
  const search = document.getElementById("search")?.value.trim();
  const kategoriId = document.getElementById("filter-kategori")?.value;
  const kondisi = document.getElementById("filter-kondisi")?.value;
  if (search) params.set("search", search);
  if (kategoriId) params.set("kategori_id", kategoriId);
  if (kondisi) params.set("kondisi", kondisi);

  try {
    const res = await apiRequest(`/barang?${params.toString()}`);
    const items = res.data?.items || [];
    const meta = res.data?.meta || { page: 1, total_halaman: 1, total_data: 0 };

    document.getElementById("page-info").textContent =
      `Halaman ${meta.page} dari ${meta.total_halaman} (Total: ${meta.total_data} barang)`;
    document.getElementById("btn-prev").disabled = meta.page <= 1;
    document.getElementById("btn-next").disabled = meta.page >= meta.total_halaman;

    if (items.length === 0) {
      tbody.innerHTML = `<tr><td colspan="7" class="px-4 py-6 text-center text-slate-500">Data barang tidak ditemukan.</td></tr>`;
      return;
    }

    tbody.innerHTML = items
      .map(
        (b) => `
      <tr class="border-t border-slate-200 hover:bg-slate-50">
        <td class="px-4 py-3 font-mono text-xs font-semibold">${escapeHtml(b.kode_barang)}</td>
        <td class="px-4 py-3 font-semibold">${escapeHtml(b.nama)}</td>
        <td class="px-4 py-3">${escapeHtml(b.kategori?.nama || "-")}</td>
        <td class="px-4 py-3">${escapeHtml(b.stok)}</td>
        <td class="px-4 py-3"><span class="rounded-full bg-slate-100 px-2 py-0.5 text-xs">${escapeHtml(b.kondisi)}</span></td>
        <td class="px-4 py-3 text-right space-x-2">
          <button onclick="detailBarang(${b.id})" class="rounded border border-slate-300 px-2.5 py-1 text-xs hover:bg-slate-100">Detail</button>
          ${
            isAdmin()
              ? `<button onclick="editBarang(${b.id})" class="rounded border border-amber-300 bg-amber-50 px-2.5 py-1 text-xs text-amber-800 hover:bg-amber-100">Edit</button>
                 <button onclick="deleteBarang(${b.id})" class="rounded border border-rose-300 bg-rose-50 px-2.5 py-1 text-xs text-rose-700 hover:bg-rose-100">Hapus</button>`
              : ""
          }
        </td>
      </tr>`
      )
      .join("");
  } catch (err) {
    showAlert("alert-box", err.message, "error");
    tbody.innerHTML = `<tr><td colspan="7" class="px-4 py-6 text-center text-rose-600">Gagal memuat data barang.</td></tr>`;
  }
}

async function handleSaveBarang(e) {
  e.preventDefault();
  hideAlert("alert-box");
  const btn = document.getElementById("btn-save");
  const id = document.getElementById("barang-id").value;
  setButtonLoading(btn, true, "Menyimpan...");

  try {
    const payload = {
      kode_barang: document.getElementById("kode_barang").value.trim(),
      nama: document.getElementById("nama").value.trim(),
      stok: Number(document.getElementById("stok").value),
      kondisi: document.getElementById("kondisi").value,
      kategori_id: Number(document.getElementById("kategori_id").value),
    };
    const res = await apiRequest(id ? `/barang/${id}` : "/barang", {
      method: id ? "PUT" : "POST",
      body: JSON.stringify(payload),
    });
    showAlert("alert-box", res.pesan, "success");
    resetBarangForm();
    loadBarang();
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  } finally {
    setButtonLoading(btn, false);
  }
}

async function detailBarang(id) {
  hideAlert("alert-box");
  try {
    const res = await apiRequest(`/barang/${id}`);
    const b = res.data;
    document.getElementById("detail-box").classList.remove("hidden");
    document.getElementById("detail-content").innerHTML = `
      <strong>[${escapeHtml(b.kode_barang)}] ${escapeHtml(b.nama)}</strong> —
      Kategori: ${escapeHtml(b.kategori?.nama || "-")},
      Stok: ${escapeHtml(b.stok)},
      Kondisi: ${escapeHtml(b.kondisi)}
    `;
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  }
}

async function editBarang(id) {
  hideAlert("alert-box");
  try {
    const res = await apiRequest(`/barang/${id}`);
    const b = res.data;
    document.getElementById("barang-id").value = b.id;
    document.getElementById("kode_barang").value = b.kode_barang;
    document.getElementById("nama").value = b.nama;
    document.getElementById("stok").value = b.stok;
    document.getElementById("kondisi").value = b.kondisi;
    document.getElementById("kategori_id").value = b.kategori_id;
    document.getElementById("form-title").textContent = `Edit Barang #${b.id}`;
    window.scrollTo({ top: 0, behavior: "smooth" });
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  }
}

async function deleteBarang(id) {
  if (!confirm("Yakin ingin menghapus barang ini?")) return;
  hideAlert("alert-box");
  try {
    const res = await apiRequest(`/barang/${id}`, { method: "DELETE" });
    showAlert("alert-box", res.pesan, "success");
    loadBarang();
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  }
}

function resetBarangForm() {
  document.getElementById("barang-id").value = "";
  document.getElementById("barang-form").reset();
  document.getElementById("form-title").textContent = "Tambah Barang";
}
