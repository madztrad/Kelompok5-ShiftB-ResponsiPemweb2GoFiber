function renderNavbar(activePage) {
  const nav = document.getElementById("navbar");
  if (!nav) return;

  const user = getUser();
  const links = [
    { href: "index.html", label: "Dashboard", key: "dashboard" },
    { href: "kategori.html", label: "Kategori", key: "kategori" },
    { href: "barang.html", label: "Barang", key: "barang" },
    { href: "peminjaman.html", label: "Peminjaman", key: "peminjaman" },
    { href: "profile.html", label: "Profil", key: "profile" },
  ];

  const menuHtml = links
    .map((item) => {
      const isActive = item.key === activePage;
      const cls = isActive
        ? "bg-slate-900 text-white"
        : "text-slate-600 hover:bg-slate-100 hover:text-slate-900";
      return `<a href="${item.href}" class="rounded-lg px-3 py-2 text-sm font-medium ${cls}">${item.label}</a>`;
    })
    .join("");

  nav.innerHTML = `
    <div class="border-b border-slate-200 bg-white">
      <div class="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
        <div class="flex flex-wrap items-center gap-4">
          <a href="index.html" class="text-lg font-bold text-slate-900">SimInventaris</a>
          <div class="flex flex-wrap items-center gap-1">${menuHtml}</div>
        </div>
        <div class="flex items-center gap-3">
          ${
            user
              ? `<span class="whitespace-nowrap text-xs font-semibold uppercase tracking-wide rounded-full bg-slate-100 px-2.5 py-1 text-slate-700">${escapeHtml(user.nama)} (${escapeHtml(user.role)})</span>`
              : ""
          }
          <button id="btn-logout" class="rounded-lg border border-slate-300 px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-100">
            Logout
          </button>
        </div>
      </div>
    </div>
  `;

  const btnLogout = document.getElementById("btn-logout");
  if (btnLogout) {
    btnLogout.addEventListener("click", async () => {
      try {
        await apiRequest("/auth/logout", { method: "POST" });
      } catch {
        // Tetap hapus token lokal meskipun endpoint logout gagal/belum aktif
      } finally {
        logoutLocal();
      }
    });
  }
}
