// Dikelola oleh Anggota 1: CRUD Kategori (/api/kategori)
document.addEventListener("DOMContentLoaded", () => {
  if (!requireAuth()) return;
  renderNavbar("kategori");

  const adminPanel = document.getElementById("admin-form-panel");
  if (isAdmin() && adminPanel) {
    adminPanel.classList.remove("hidden");
  }

  document.getElementById("kategori-form")?.addEventListener("submit", handleSaveKategori);
  document.getElementById("btn-reset")?.addEventListener("click", resetKategoriForm);
  loadKategori();
});

async function loadKategori() {
  const tbody = document.getElementById("kategori-tbody");
  tbody.innerHTML = `<tr><td colspan="4" class="px-4 py-6 text-center text-slate-500">Memuat data kategori...</td></tr>`;
  try {
    const res = await apiRequest("/kategori");
    const items = res.data || [];
    if (items.length === 0) {
      tbody.innerHTML = `<tr><td colspan="4" class="px-4 py-6 text-center text-slate-500">Belum ada data kategori.</td></tr>`;
      return;
    }

    tbody.innerHTML = items
      .map(
        (k) => `
      <tr class="border-t border-slate-200 hover:bg-slate-50">
        <td class="px-4 py-3 font-medium">${escapeHtml(k.id)}</td>
        <td class="px-4 py-3 font-semibold">${escapeHtml(k.nama)}</td>
        <td class="px-4 py-3 text-slate-600">${escapeHtml(k.deskripsi || "-")}</td>
        <td class="px-4 py-3 text-right space-x-2">
          <button onclick="detailKategori(${k.id})" class="rounded border border-slate-300 px-2.5 py-1 text-xs hover:bg-slate-100">Detail</button>
          ${
            isAdmin()
              ? `<button onclick="editKategori(${k.id})" class="rounded border border-amber-300 bg-amber-50 px-2.5 py-1 text-xs text-amber-800 hover:bg-amber-100">Edit</button>
                 <button onclick="deleteKategori(${k.id})" class="rounded border border-rose-300 bg-rose-50 px-2.5 py-1 text-xs text-rose-700 hover:bg-rose-100">Hapus</button>`
              : ""
          }
        </td>
      </tr>`
      )
      .join("");
  } catch (err) {
    showAlert("alert-box", err.message, "error");
    tbody.innerHTML = `<tr><td colspan="4" class="px-4 py-6 text-center text-rose-600">Gagal memuat data.</td></tr>`;
  }
}

async function handleSaveKategori(e) {
  e.preventDefault();
  if (!isAdmin()) {
    showAlert("alert-box", "Akses ditolak: hanya admin yang dapat menyimpan kategori", "error");
    return;
  }
  hideAlert("alert-box");
  const btn = document.getElementById("btn-save");
  const id = document.getElementById("kategori-id").value;
  setButtonLoading(btn, true, "Menyimpan...");

  try {
    const payload = {
      nama: document.getElementById("nama").value.trim(),
      deskripsi: document.getElementById("deskripsi").value.trim(),
    };
    const res = await apiRequest(id ? `/kategori/${id}` : "/kategori", {
      method: id ? "PUT" : "POST",
      body: JSON.stringify(payload),
    });
    showAlert("alert-box", res.pesan, "success");
    resetKategoriForm();
    loadKategori();
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  } finally {
    setButtonLoading(btn, false);
  }
}

async function detailKategori(id) {
  hideAlert("alert-box");
  try {
    const res = await apiRequest(`/kategori/${id}`);
    const k = res.data;
    document.getElementById("detail-box").classList.remove("hidden");
    document.getElementById("detail-content").innerHTML = `
      <strong>#${escapeHtml(k.id)} — ${escapeHtml(k.nama)}</strong><br/>
      <span class="text-slate-600">${escapeHtml(k.deskripsi || "Tidak ada deskripsi.")}</span>
    `;
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  }
}

async function editKategori(id) {
  if (!isAdmin()) {
    showAlert("alert-box", "Akses ditolak: hanya admin yang dapat mengubah kategori", "error");
    return;
  }
  hideAlert("alert-box");
  try {
    const res = await apiRequest(`/kategori/${id}`);
    const k = res.data;
    document.getElementById("kategori-id").value = k.id;
    document.getElementById("nama").value = k.nama;
    document.getElementById("deskripsi").value = k.deskripsi || "";
    document.getElementById("form-title").textContent = `Edit Kategori #${k.id}`;
    window.scrollTo({ top: 0, behavior: "smooth" });
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  }
}

async function deleteKategori(id) {
  if (!isAdmin()) {
    showAlert("alert-box", "Akses ditolak: hanya admin yang dapat menghapus kategori", "error");
    return;
  }
  if (!confirm("Yakin ingin menghapus kategori ini?")) return;
  hideAlert("alert-box");
  try {
    const res = await apiRequest(`/kategori/${id}`, { method: "DELETE" });
    showAlert("alert-box", res.pesan, "success");
    loadKategori();
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  }
}

function resetKategoriForm() {
  document.getElementById("kategori-id").value = "";
  document.getElementById("kategori-form").reset();
  document.getElementById("form-title").textContent = "Tambah Kategori";
}
