// Dikelola oleh Anggota 3: CRUD Peminjaman (/api/peminjaman) + Preload Barang & User
document.addEventListener("DOMContentLoaded", () => {
  if (!requireAuth()) return;
  renderNavbar("peminjaman");

  if (isAdmin()) {
    document.getElementById("status-group")?.classList.remove("hidden");
  }

  // Default tanggal pinjam hari ini
  const tglInput = document.getElementById("tanggal_pinjam");
  if (tglInput && !tglInput.value) {
    tglInput.value = new Date().toISOString().slice(0, 10);
  }

  setPeminjamanEditMode(false);
  document.getElementById("peminjaman-form")?.addEventListener("submit", handleSavePeminjaman);
  document.getElementById("btn-reset")?.addEventListener("click", resetPeminjamanForm);

  loadBarangOptions();
  loadPeminjaman();
});

async function loadBarangOptions() {
  try {
    const res = await apiRequest("/barang?limit=100");
    const items = res.data?.items || [];
    const select = document.getElementById("barang_id");
    if (!select) return;
    select.innerHTML =
      `<option value="">-- Pilih Barang --</option>` +
      items
        .map(
          (b) =>
            `<option value="${b.id}">${escapeHtml(b.kode_barang)} - ${escapeHtml(b.nama)} (Stok: ${escapeHtml(b.stok)})</option>`
        )
        .join("");
  } catch (err) {
    showAlert("alert-box", `Gagal memuat daftar barang: ${err.message}`, "error");
  }
}

async function loadPeminjaman() {
  const tbody = document.getElementById("peminjaman-tbody");
  tbody.innerHTML = `<tr><td colspan="7" class="px-4 py-6 text-center text-slate-500">Memuat data peminjaman...</td></tr>`;

  try {
    const res = await apiRequest("/peminjaman");
    const items = res.data || [];
    if (items.length === 0) {
      tbody.innerHTML = `<tr><td colspan="7" class="px-4 py-6 text-center text-slate-500">Belum ada data peminjaman.</td></tr>`;
      return;
    }

    tbody.innerHTML = items
      .map(
        (p) => `
      <tr class="border-t border-slate-200 hover:bg-slate-50">
        <td class="px-4 py-3 font-medium">${escapeHtml(p.id)}</td>
        <td class="px-4 py-3">${escapeHtml(p.user?.nama || "-")}</td>
        <td class="px-4 py-3 font-semibold">${escapeHtml(p.barang?.nama || "-")}</td>
        <td class="px-4 py-3">${escapeHtml(p.jumlah)}</td>
        <td class="px-4 py-3 text-xs">${escapeHtml(p.tanggal_pinjam)} ${p.tanggal_kembali ? `s/d ${escapeHtml(p.tanggal_kembali)}` : ""}</td>
        <td class="px-4 py-3"><span class="rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-medium">${escapeHtml(p.status)}</span></td>
        <td class="px-4 py-3 text-right space-x-2">
          <button onclick="detailPeminjaman(${p.id})" class="rounded border border-slate-300 px-2.5 py-1 text-xs hover:bg-slate-100">Detail</button>
          <button onclick="editPeminjaman(${p.id})" class="rounded border border-amber-300 bg-amber-50 px-2.5 py-1 text-xs text-amber-800 hover:bg-amber-100">Edit</button>
          ${
            isAdmin()
              ? `<button onclick="deletePeminjaman(${p.id})" class="rounded border border-rose-300 bg-rose-50 px-2.5 py-1 text-xs text-rose-700 hover:bg-rose-100">Hapus</button>`
              : ""
          }
        </td>
      </tr>`
      )
      .join("");
  } catch (err) {
    showAlert("alert-box", err.message, "error");
    tbody.innerHTML = `<tr><td colspan="7" class="px-4 py-6 text-center text-rose-600">Gagal memuat data peminjaman.</td></tr>`;
  }
}

async function handleSavePeminjaman(e) {
  e.preventDefault();
  hideAlert("alert-box");
  const btn = document.getElementById("btn-save");
  const id = document.getElementById("peminjaman-id").value;
  setButtonLoading(btn, true, "Menyimpan...");

  try {
    let payload;
    if (id && !isAdmin()) {
      payload = {
        catatan: document.getElementById("catatan").value.trim(),
      };
    } else if (id) {
      payload = {
        jumlah: Number(document.getElementById("jumlah").value),
        tanggal_pinjam: document.getElementById("tanggal_pinjam").value,
        tanggal_kembali: document.getElementById("tanggal_kembali").value,
        catatan: document.getElementById("catatan").value.trim(),
        status: document.getElementById("status").value,
      };
    } else {
      payload = {
        barang_id: Number(document.getElementById("barang_id").value),
        jumlah: Number(document.getElementById("jumlah").value),
        tanggal_pinjam: document.getElementById("tanggal_pinjam").value,
        tanggal_kembali: document.getElementById("tanggal_kembali").value,
        catatan: document.getElementById("catatan").value.trim(),
      };
    }

    const res = await apiRequest(id ? `/peminjaman/${id}` : "/peminjaman", {
      method: id ? "PUT" : "POST",
      body: JSON.stringify(payload),
    });
    showAlert("alert-box", res.pesan, "success");
    resetPeminjamanForm();
    loadBarangOptions();
    loadPeminjaman();
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  } finally {
    setButtonLoading(btn, false);
  }
}

async function detailPeminjaman(id) {
  hideAlert("alert-box");
  try {
    const res = await apiRequest(`/peminjaman/${id}`);
    const p = res.data;
    document.getElementById("detail-box").classList.remove("hidden");
    document.getElementById("detail-content").innerHTML = `
      <strong>Peminjaman #${escapeHtml(p.id)}</strong> —
      Peminjam: ${escapeHtml(p.user?.nama || "-")},
      Barang: ${escapeHtml(p.barang?.nama || "-")} (${escapeHtml(p.jumlah)} unit),
      Status: <strong>${escapeHtml(p.status)}</strong>,
      Catatan: ${escapeHtml(p.catatan || "-")}
    `;
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  }
}

async function editPeminjaman(id) {
  hideAlert("alert-box");
  try {
    const res = await apiRequest(`/peminjaman/${id}`);
    const p = res.data;
    document.getElementById("peminjaman-id").value = p.id;
    document.getElementById("barang_id").value = p.barang_id;
    document.getElementById("barang_id").disabled = true;
    document.getElementById("jumlah").value = p.jumlah;
    document.getElementById("tanggal_pinjam").value = p.tanggal_pinjam;
    document.getElementById("tanggal_kembali").value = p.tanggal_kembali || "";
    document.getElementById("catatan").value = p.catatan || "";
    if (isAdmin()) {
      document.getElementById("status").value = p.status;
    }
    setPeminjamanEditMode(true);
    document.getElementById("form-title").textContent = `Edit Peminjaman #${p.id}`;
    window.scrollTo({ top: 0, behavior: "smooth" });
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  }
}

async function deletePeminjaman(id) {
  if (!confirm("Yakin ingin menghapus data peminjaman ini?")) return;
  hideAlert("alert-box");
  try {
    const res = await apiRequest(`/peminjaman/${id}`, { method: "DELETE" });
    showAlert("alert-box", res.pesan, "success");
    loadBarangOptions();
    loadPeminjaman();
  } catch (err) {
    showAlert("alert-box", err.message, "error");
  }
}

function resetPeminjamanForm() {
  document.getElementById("peminjaman-id").value = "";
  setPeminjamanEditMode(false);
  document.getElementById("peminjaman-form").reset();
  document.getElementById("tanggal_pinjam").value = new Date().toISOString().slice(0, 10);
  document.getElementById("form-title").textContent = "Ajukan Peminjaman";
}

function setPeminjamanEditMode(isEditing) {
  const memberEdit = isEditing && !isAdmin();
  document.getElementById("barang_id").disabled = isEditing;
  document.getElementById("jumlah").disabled = memberEdit;
  document.getElementById("tanggal_pinjam").disabled = memberEdit;
  document.getElementById("tanggal_kembali").disabled = memberEdit;
  document.getElementById("member-edit-note")?.classList.toggle("hidden", !memberEdit);
}
