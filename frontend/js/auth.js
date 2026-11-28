// Dikelola oleh Anggota 1: Login, Register, Profile
document.addEventListener("DOMContentLoaded", () => {
  const loginForm = document.getElementById("login-form");
  if (loginForm) {
    if (getToken()) {
      window.location.href = "index.html";
      return;
    }
    periksaKoneksiApi();
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      hideAlert("alert-box");
      const btn = document.getElementById("btn-submit");
      setButtonLoading(btn, true, "Memeriksa...");
      try {
        const res = await apiRequest("/auth/login", {
          method: "POST",
          body: JSON.stringify({
            email: document.getElementById("email").value.trim(),
            password: document.getElementById("password").value,
          }),
        });
        localStorage.setItem("token", res.data.token);
        localStorage.setItem("user", JSON.stringify(res.data.user));
        window.location.href = "index.html";
      } catch (err) {
        showAlert("alert-box", err.message, "error");
      } finally {
        setButtonLoading(btn, false);
      }
    });
  }

  const registerForm = document.getElementById("register-form");
  if (registerForm) {
    periksaKoneksiApi();
    registerForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      hideAlert("alert-box");
      const btn = document.getElementById("btn-submit");
      setButtonLoading(btn, true, "Mendaftar...");
      try {
        await apiRequest("/auth/register", {
          method: "POST",
          body: JSON.stringify({
            nama: document.getElementById("nama").value.trim(),
            email: document.getElementById("email").value.trim(),
            password: document.getElementById("password").value,
          }),
        });
        showAlert("alert-box", "Registrasi berhasil! Mengarahkan ke halaman login...", "success");
        setTimeout(() => {
          window.location.href = "login.html";
        }, 1000);
      } catch (err) {
        showAlert("alert-box", err.message, "error");
      } finally {
        setButtonLoading(btn, false);
      }
    });
  }

  const profileCard = document.getElementById("profile-card");
  if (profileCard && requireAuth()) {
    loadProfile(profileCard);
  }
});

async function loadProfile(container) {
  try {
    const res = await apiRequest("/auth/profile");
    const u = res.data;
    localStorage.setItem("user", JSON.stringify(u));
    container.innerHTML = `
      <dl class="divide-y divide-slate-200 text-sm">
        <div class="py-3 flex justify-between"><dt class="text-slate-500">ID</dt><dd class="font-medium">${escapeHtml(u.id)}</dd></div>
        <div class="py-3 flex justify-between"><dt class="text-slate-500">Nama</dt><dd class="font-medium">${escapeHtml(u.nama)}</dd></div>
        <div class="py-3 flex justify-between"><dt class="text-slate-500">Email</dt><dd class="font-medium">${escapeHtml(u.email)}</dd></div>
        <div class="py-3 flex justify-between"><dt class="text-slate-500">Role</dt><dd class="font-semibold uppercase">${escapeHtml(u.role)}</dd></div>
      </dl>
    `;
  } catch (err) {
    showAlert("alert-box", err.message, "error");
    container.textContent = "Gagal memuat profil.";
  }
}
