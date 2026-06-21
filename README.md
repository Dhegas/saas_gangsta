# 🚀 SaaS Gangsta Backend API
> **Enterprise-Grade Platform SaaS POS & Self-Ordering UMKM Kuliner**

SaaS Gangsta adalah layanan backend API skala industri yang dirancang khusus untuk mendigitalisasi operasional UMKM kuliner di Indonesia. Dengan arsitektur **multi-tenant**, setiap partner mendapatkan ekosistem data yang terisolasi, aman, dan performan.

---

## 📖 Terminologi Domain
Agar tidak ada ambiguitas dalam arsitektur sistem, berikut adalah standar istilah yang digunakan:
* **Partner** = Pelaku usaha / Penjual yang mendaftar ke platform. (Satu partner bisa memiliki banyak tenant).
* **Tenant** = Unit bisnis / Toko fisik cabang yang beroperasi. (Data menu, pesanan, dan meja terikat pada entitas ini).
* **Customer** = Pengguna akhir (pembeli) yang menggunakan layanan self-ordering.
* **Admin** = Super user pengelola platform SaaS.

---

## 🏗 Architecture & Design Patterns

Sistem ini dibangun menggunakan standar **Clean Architecture** yang dikombinasikan dengan pendekatan **Modular Monolith (Domain-Oriented)**.

### Architecture Diagram
[ERD SaaS Gangsta](https://viewer.diagrams.net/?tags=%7B%7D&lightbox=1&highlight=0000ff&edit=_blank&layers=1&nav=1&title=diagram%20gangsta.drawio&dark=auto#Uhttps%3A%2F%2Fdrive.google.com%2Fuc%3Fid%3D14pQkdyTJauy4w241hsFAc8Z1uFC_wkXd%26export%3Ddownload)

### Layer Responsibility
- **Delivery**: Menangani binding request JSON, validasi input (`validator/v10`), dan parsing context JWT.
- **Usecase**: Jantung bisnis aplikasi. Menjalankan logika bisnis, caching, dan orkestrasi repository.
- **Domain**: Definisi entitas bisnis (Struct) dan kontrak interface.
- **Repository**: Abstraksi data akses menggunakan GORM untuk PostgreSQL.

### Teknologi Pendukung
- **Database**: PostgreSQL yang di-hosting di **Supabase**.
- **Cache**: In-memory cache lokal (**LocalCache**) untuk mempercepat query list data dinamis dan mengurangi load database (tidak menggunakan Redis).
- **Payment Gateway**: Integrasi pembayaran otomatis menggunakan **Midtrans**.

### Current Internal Structure
- **Routing Bootstrap**: dipisah ke `internal/bootstrap/routes.go`, `internal/bootstrap/customer_routes.go`, `internal/bootstrap/partner_routes.go`, dan `internal/bootstrap/admin_routes.go`.
- **Cross-cutting Infra**: `internal/config`, `internal/middleware`, `internal/infrastructure/database`.
- **Business Domains**: `internal/domains/{menu,order,payment,report,subscription,table,tenant,user}`.
- **Shared Utilities**: `internal/common/errors` dan `internal/common/response`.

---

## 🔑 Akun Testing (Untuk Dosen / Penilaian)
Untuk mempermudah pengujian fitur secara mandiri oleh Dosen, berikut adalah data akun yang siap digunakan:

* **👑 Admin Account** (Super user pengelola platform SaaS)
  * **Email**: `admin@gmail.com`
  * **Password**: `ADMIN123`

* **🏪 Partner Account** (Mitra / Pelaku usaha UMKM)
  * **Email**: `purnomo@gmail.com`
  * **Password**: `password`

* **📱 Customer Account** (Pelanggan untuk self-ordering)
  * **Email**: `popo@gmail.com`
  * **Password**: `password`

---

## 👥 Roles & Access Control
| Role | Deskripsi | Hak Akses Utama |
|---|---|---|
| 👑 **Admin** | Pengelola Platform | Kelola Tenant, Billing, Monitoring Global |
| 🏪 **Partner** | Pelaku Usaha / Mitra | Kelola Tenant, Menu, POS, Meja, Laporan Penjualan |
| 📱 **Customer** | Pelanggan | Scan QR, Self-Ordering, Cek Status Pesanan |

---

## 🛠 Developer Onboarding Guide

Selamat bergabung di tim pengembang! Ikuti panduan ini untuk mulai berkontribusi.

### 1. Local Environment Setup
1. **Clone & Install**:
   ```bash
   git clone <repo_url>
   cd saas_gangsta
   go mod tidy
   ```
2. **Setup Environment**:
   * File `.env` berisi key sensitif (seperti kredensial database Supabase dan token Midtrans) sehingga tidak di-commit ke dalam Git.
   * File `.env` lengkap yang valid telah dilampirkan pada form pengumpulan tugas di LMS (dalam bentuk teks). Silakan salin isi konfigurasi tersebut dan buat file baru bernama `.env` di root direktori project ini.
   * Secara default, port aplikasi local diatur pada port `8080` (diakses melalui `http://localhost:8080`).
3. **Run Application**:
   ```bash
   go run cmd/api/main.go
   ```

### 2. Workflow Pengembangan Domain
Setiap penambahan fitur baru harus mengikuti struktur domain berikut di `internal/domains/<nama_domain>/`:
- **dto/**: Request & Response object untuk mapping data luar.
- **delivery/http/**: Handler Gin untuk menangani request masuk.
- **usecase/**: Logika bisnis murni.
- **domain/**: Abstraksi interface dan model entitas.
- **repository/**: Query database (GORM).

### 3. Konvensi Kode (Coding Standards)
- **JSON Naming**: Selalu gunakan `camelCase` (e.g., `storeName`).
- **DB Naming**: Selalu gunakan `snake_case` (e.g., `store_name`).
- **Multi-tenancy**: Semua query bisnis **wajib** menyertakan filter `tenant_id` yang diambil dari context JWT (jangan dari input user).
- **Responses**: Gunakan helper di `internal/common/response` untuk memastikan respons seragam:
  ```go
  response.Success(c, http.StatusOK, "Berhasil", data)
  ```

---

## 📋 Operational Guide

### Database & Migrasi
- **Update Schema**: Tambahkan file `.sql` baru di folder `supabase/migrations/` dengan format timestamp berurutan.
- **Seed Data**: Gunakan `scripts/seed.sql` untuk mengisi data awal di lingkungan development.

### API Documentation (Swagger)
Kami menggunakan Swaggo. Jika ada perubahan pada komentar di handler, update dokumentasi dengan:
```bash
swag init -g cmd/api/main.go
```
Akses di: `http://localhost:8080/swagger/index.html`

---

## ✅ Definition of Done (DoD)
Sebelum submit Pull Request (PR), pastikan:
- [ ] Endpoint sesuai spek (Method, Path, Payload).
- [ ] Validasi input berfungsi (menggunakan `validator.v10`).
- [ ] Filter `tenantId` sudah diterapkan di layer Repository.
- [ ] Response menggunakan format standard envelope.
- [ ] Error ditangani dengan tepat dan tidak mengekspos detail teknis ke user luar.
- [ ] Dokumentasi Swagger sudah di-update (`swag init`).

---

## 👥 Pengembang & Git Flow

Repository ini menggunakan beberapa branch dengan fungsi masing-masing untuk menjaga kerapian dan integritas kode:

* **Branch Fitur (Developer Pengembang)**:
  * `dev-dekgus`, `dev-renata`, `dev-dhegas`: Merupakan branch kerja mandiri bagi masing-masing pengembang untuk membangun dan menyempurnakan fitur secara paralel.
* **Branch Pengujian (Testing)**:
  * `dev`: Digunakan untuk menggabungkan seluruh fitur baru dan melakukan pengujian fungsional (testing) sebelum siap digabungkan ke kode stabil.
* **Branch Utama (Stable/Fix)**:
  * `main`: Menyimpan kode final yang telah teruji, stabil, bersih dari bug, dan siap untuk dirilis secara resmi.
* **Branch Deployment (Hosting)**:
  * `production`: Merupakan branch produksi yang dihubungkan langsung ke **Railway** untuk hosting aplikasi secara online.
