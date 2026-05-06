# saas_gangsta — Backend API

> **Platform SaaS POS & Self-Order untuk UMKM Kuliner Indonesia**
> Backend API Service · Go + Gin · PostgreSQL (Supabase) · Docker · Nginx API Gateway

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Tech Stack](#2-tech-stack)
3. [Architecture Overview](#3-architecture-overview)
4. [Folder Structure](#4-folder-structure)
5. [Environment Variables](#5-environment-variables)
6. [API Gateway (Nginx)](#6-api-gateway-nginx)
7. [Authentication & Authorization (JWT)](#7-authentication--authorization-jwt)
8. [Database Schema (PostgreSQL Supabase)](#8-database-schema-postgresql-supabase)
9. [API Endpoints](#9-api-endpoints)

- [Public & Health](#91-public--health)
- [Auth](#92-auth)
- [Users (Partner/Admin)](#93-users-partneradmin)
- [Tenant Context & Partner](#94-tenant-context--partner)
- [Catalog (Categories & Menus)](#95-catalog-categories--menus)
- [Orders & Customers](#96-orders--customers)
- [Dining Tables](#97-dining-tables)
- [Reports](#98-reports)

10. [CORS Configuration](#10-cors-configuration)
11. [Response Schema Convention](#11-response-schema-convention)
12. [Error Handling Convention](#12-error-handling-convention)
13. [Docker & Container Setup](#13-docker--container-setup)
14. [Clean Code Guidelines](#14-clean-code-guidelines)
15. [Domain Structure per Feature](#15-domain-structure-per-feature)
16. [Development Roadmap](#16-development-roadmap)
17. [Definition of Done](#17-definition-of-done)
18. [Technical Risks](#18-technical-risks)

---

## 1. Project Overview

`saas_gangsta` adalah layanan backend untuk platform SaaS yang membantu **UMKM kuliner Indonesia** mengelola operasional toko secara digital.

### Tujuan Bisnis

Platform ini menggabungkan kebutuhan utama operasional toko makanan:

| Modul                     | Keterangan                                      |
| ------------------------- | ----------------------------------------------- |
| Digital Menu              | Pelanggan scan QR → lihat menu digital          |
| Self Ordering             | Pelanggan order dari meja tanpa panggil pelayan |
| POS / Kasir               | Partner input order manual via kasir digital    |
| Manajemen Meja            | Monitor kondisi dan status meja real-time       |
| Laporan Penjualan         | Rekap harian/mingguan/bulanan per partner       |
| Membership / Subscription | Admin kelola paket langganan SaaS               |

### Tiga Role Utama

| Role         | Deskripsi                                                                |
| ------------ | ------------------------------------------------------------------------ |
| **Customer** | Pelanggan toko. Scan QR, lihat menu, order, bayar, cek status, review    |
| **Partner**  | Pemilik toko. Kelola menu, terima order, POS, laporan, profil toko       |
| **Admin**    | Pengelola platform SaaS. Kelola tenant, membership, billing, user global |

### Model Bisnis

Sistem menggunakan model **SaaS multi-tenant**, di mana setiap partner (tenant) memiliki data yang terisolasi satu sama lain. Partner berlangganan bulanan untuk mengakses platform.

---

## 2. Tech Stack

| Komponen           | Teknologi                                         |
| ------------------ | ------------------------------------------------- |
| Language           | Go 1.23+                                          |
| HTTP Framework     | Gin (`github.com/gin-gonic/gin`)                  |
| Database           | PostgreSQL via Supabase (managed)                 |
| ORM                | GORM (`gorm.io/gorm`) + `pgx`                     |
| API Gateway        | Nginx (reverse proxy & rate limiter)              |
| Authentication     | JWT (access token + refresh token)                |
| Container          | Docker + Docker Compose                           |
| Environment Config | `.env` + `godotenv`                               |
| Validation         | `go-playground/validator/v10`                     |
| Logging            | `log/slog` (standard library Go 1.21+)            |
| UUID               | `github.com/google/uuid`                          |
| Schema Convention  | **camelCase** untuk semua JSON response & request |

---

## 3. Architecture Overview

Proyek menggunakan pendekatan **Clean Architecture + Modular Monolith** yang cocok untuk MVP dan tetap scalable untuk production.

```
                        ┌──────────────────────────────────────┐
                        │           NGINX (API Gateway)         │
                        │  - Reverse Proxy                      │
                        │  - Rate Limiting                      │
                        │  - SSL Termination                    │
                        └────────────────┬─────────────────────┘
                                         │
                        ┌────────────────▼─────────────────────┐
                        │         Go + Gin HTTP Server          │
                        │                                       │
                        │  ┌──────────┐  ┌──────────────────┐  │
                        │  │ Middleware│  │  Route Groups    │  │
                        │  │ JWT Auth │  │ /customer        │  │
                        │  │ CORS     │  │ /partner        │  │
                        │  │ Logging  │  │ /admin           │  │
                        │  │ Recovery │  │ /auth            │  │
                        │  └──────────┘  └──────────────────┘  │
                        │                                       │
                        │  ┌────────────────────────────────┐  │
                        │  │        Internal Modules         │  │
                        │  │  delivery → usecase → domain   │  │
                        │  │         → repository           │  │
                        │  └────────────────────────────────┘  │
                        └────────────────┬─────────────────────┘
                                         │
                        ┌────────────────▼─────────────────────┐
                        │    PostgreSQL (Supabase Managed)       │
                        │    - Multi-tenant data model          │
                        │    - Row Level Security (RLS)         │
                        └──────────────────────────────────────┘
```

---

## 4. Folder Structure

### Struktur Per Domain (Wajib Konsisten)

Setiap domain di dalam `internal/domains/<domain_name>/` **wajib** mengikuti struktur berikut:

```
internal/domains/<domain_name>/
├── delivery/
│   └── http/
│       └── <module>_handler.go       # HTTP handler: bind request, call usecase, return response
├── usecase/
│   └── <module>_usecase.go           # Business logic, orchestrasi repository
├── domain/
│   └── <module>_domain.go            # Struct entity / domain model (bukan GORM model)
├── repository/
│   └── <module>_repository.go        # DB query: GORM query + SQL
└── dto/
    ├── <module>_request.go           # Request DTO (payload dari client)
    └── <module>_response.go          # Response DTO (data ke client)
```

**Alur data yang wajib diikuti:**

```
HTTP Request
    ↓
Delivery/Handler     (bind & validate request DTO, parse JWT context)
    ↓
Usecase              (business rule, orchestrate repository calls)
    ↓
Repository           (GORM query / raw SQL ke database)
    ↓
Database (Supabase PostgreSQL)
    ↓
Repository           (return domain/entity)
    ↓
Usecase              (map ke response DTO)
    ↓
Delivery/Handler     (return standard response)
    ↓
HTTP Response
```

---

## 5. Environment Variables

### .env.example

```env
APP_ENV=development
APP_PORT=8080
APP_NAME=saas_gangsta

DATABASE_URL=postgresql://postgres:password@host:5432/postgres?sslmode=require

JWT_SECRET=change-me-to-a-strong-secret
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=7d

CORS_ALLOWED_ORIGINS=http://localhost:3000

SUPABASE_URL=https://<project-ref>.supabase.co
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key
```

---

## 6. API Gateway (Nginx)

Nginx digunakan sebagai **API Gateway** di depan Go server. Semua request dari client masuk melalui Nginx terlebih dahulu.

### Fungsi Nginx

| Fungsi            | Keterangan                                          |
| ----------------- | --------------------------------------------------- |
| Reverse Proxy     | Forward request dari port 80/443 ke Go server :8080 |
| Rate Limiting     | Batasi jumlah request per IP/client                 |
| SSL Termination   | Handle HTTPS (jika digunakan)                       |
| Request Buffering | Buffer upload sebelum diteruskan ke backend         |
| Health Check      | Probe endpoint `/health`                            |

### Konfigurasi (`deployments/nginx/nginx.conf`)

```nginx
worker_processes auto;

events {
    worker_connections 1024;
}

http {
    # Rate limiting zone
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=30r/m;

    upstream backend {
        server api:8080;
    }

    server {
        listen 80;
        server_name _;

        location /health {
            proxy_pass http://backend;
        }

        location /api/v1/ {
            limit_req zone=api_limit burst=10 nodelay;

            proxy_pass         http://backend;
            proxy_http_version 1.1;
            proxy_set_header   Host              $host;
            proxy_set_header   X-Real-IP         $remote_addr;
            proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
            proxy_set_header   X-Forwarded-Proto $scheme;

            # Timeout
            proxy_connect_timeout 10s;
            proxy_read_timeout    30s;
            proxy_send_timeout    10s;
        }
    }
}
```

### URL Prefix

Semua endpoint backend wajib menggunakan prefix:

```
/api/v1/
```

Contoh: `POST /api/v1/auth/login`, `GET /api/v1/partner/menus`

---

## 7. Authentication & Authorization (JWT)

### Flow JWT

```
1. Client POST /api/v1/auth/login    → { accessToken, refreshToken }
2. Client request endpoint           → Header: Authorization: Bearer <accessToken>
3. Middleware JWT                    → Parse & validate token
4. Middleware inject claim ke ctx    → { userId, role, tenantId }
5. Handler ambil dari ctx            → gunakan untuk query
```

### Role Values

| Role       | Keterangan                             |
| ---------- | -------------------------------------- |
| `CUSTOMER` | Pelanggan toko (default saat register) |
| `PARTNER`  | Pemilik / pengelola toko               |
| `ADMIN`    | Admin platform SaaS                    |

Catatan: endpoint [customer context](#94-tenant-context--partner) masih memakai `RoleGuard("BASIC")` sebagai legacy alias; saat ini berlaku sebagai role customer, namun perlu diseragamkan ke `CUSTOMER`.

### Middleware Stack per Route Group

```
Public:
  → CORS → Logger → Recovery
  - /swagger/*any, /health, /ready
  - /api/v1/health, /api/v1/ready
  - /api/v1/auth/register, /api/v1/auth/login, /api/v1/auth/refresh

Auth-only:
  → CORS → Logger → Recovery → JWTAuth
  - /api/v1/auth/logout, /api/v1/auth/me

Shared path (role-based access):
  /api/v1/menus
    - GET: JWTAuth → RoleGuards("CUSTOMER", "PARTNER", "ADMIN")
    - POST/PUT/DELETE/PATCH: JWTAuth → RoleGuard("PARTNER") → TenantGuard
  /api/v1/orders
    - POST, GET /:id, /:id/customer: JWTAuth → RoleGuards("CUSTOMER", "PARTNER", "ADMIN")
    - GET, PATCH /:id/status, DELETE /:id: JWTAuth → RoleGuard("PARTNER") → TenantGuard

Partner/Tenant scoped:
  → CORS → Logger → Recovery → JWTAuth
  - /api/v1/categories, /api/v1/dining-tables, /api/v1/reports: RoleGuard("PARTNER") → TenantGuard
  - /api/v1/users: RoleGuards("PARTNER", "ADMIN") → TenantGuard
  - /api/v1/partner/*: RoleGuard("PARTNER") (+ TenantGuard khusus /partner/me)
  - /api/v1/customer/me: RoleGuard("BASIC") → TenantGuard
```

---

## 8. Database Schema (PostgreSQL Supabase)

### 8.1 Migrations (Current Tables)

Sumber: seluruh file SQL di folder `supabase/migrations`.

**tenants**

| Column       | Type         | Notes                         |
| ------------ | ------------ | ----------------------------- |
| id           | UUID         | PK, default gen_random_uuid() |
| user_id      | UUID         | FK -> users.id, nullable      |
| name         | VARCHAR(150) | not null                      |
| slug         | VARCHAR(80)  | not null, unique              |
| status       | VARCHAR(20)  | default 'active'              |
| description  | TEXT         | nullable                      |
| address      | TEXT         | nullable                      |
| phone_number | VARCHAR(20)  | nullable                      |
| open_hours   | VARCHAR(100) | nullable                      |
| logo_url     | TEXT         | nullable                      |
| banner_url   | TEXT         | nullable                      |
| created_at   | TIMESTAMPTZ  | default now()                 |
| updated_at   | TIMESTAMPTZ  | default now()                 |
| deleted_at   | TIMESTAMPTZ  | nullable                      |

**users**

| Column        | Type         | Notes                                   |
| ------------- | ------------ | --------------------------------------- |
| id            | UUID         | PK, default gen_random_uuid()           |
| email         | VARCHAR(255) | not null, unique                        |
| password_hash | TEXT         | not null                                |
| full_name     | VARCHAR(150) | not null                                |
| role          | VARCHAR(20)  | CHECK in ('CUSTOMER','PARTNER','ADMIN') |
| is_active     | BOOLEAN      | default true                            |
| last_login_at | TIMESTAMPTZ  | nullable                                |
| created_at    | TIMESTAMPTZ  | default now()                           |
| updated_at    | TIMESTAMPTZ  | default now()                           |
| deleted_at    | TIMESTAMPTZ  | nullable                                |

**categories**

| Column      | Type         | Notes                         |
| ----------- | ------------ | ----------------------------- |
| id          | UUID         | PK, default gen_random_uuid() |
| tenant_id   | UUID         | FK -> tenants.id              |
| name        | VARCHAR(120) | not null                      |
| description | TEXT         | nullable                      |
| sort_order  | INTEGER      | default 0                     |
| is_active   | BOOLEAN      | default true                  |
| created_at  | TIMESTAMPTZ  | default now()                 |
| updated_at  | TIMESTAMPTZ  | default now()                 |
| deleted_at  | TIMESTAMPTZ  | nullable                      |

**menus**

| Column       | Type          | Notes                         |
| ------------ | ------------- | ----------------------------- |
| id           | UUID          | PK, default gen_random_uuid() |
| tenant_id    | UUID          | FK -> tenants.id              |
| category_id  | UUID          | FK -> categories.id, nullable |
| name         | VARCHAR(180)  | not null                      |
| description  | TEXT          | nullable                      |
| price        | NUMERIC(12,2) | not null, >= 0                |
| image_url    | TEXT          | nullable                      |
| is_available | BOOLEAN       | default true                  |
| created_at   | TIMESTAMPTZ   | default now()                 |
| updated_at   | TIMESTAMPTZ   | default now()                 |
| deleted_at   | TIMESTAMPTZ   | nullable                      |

**dining_tables**

| Column     | Type        | Notes                         |
| ---------- | ----------- | ----------------------------- |
| id         | UUID        | PK, default gen_random_uuid() |
| tenant_id  | UUID        | FK -> tenants.id              |
| table_name | VARCHAR(50) | not null                      |
| created_at | TIMESTAMPTZ | default now()                 |
| updated_at | TIMESTAMPTZ | default now()                 |
| deleted_at | TIMESTAMPTZ | nullable                      |

**orders**

| Column           | Type          | Notes                            |
| ---------------- | ------------- | -------------------------------- |
| id               | UUID          | PK, default gen_random_uuid()    |
| tenant_id        | UUID          | FK -> tenants.id                 |
| user_id          | UUID          | FK -> users.id, nullable         |
| dining_tables_id | UUID          | FK -> dining_tables.id, nullable |
| status           | VARCHAR(20)   | not null                         |
| total_price      | NUMERIC(12,2) | not null, >= 0                   |
| created_at       | TIMESTAMPTZ   | default now()                    |
| updated_at       | TIMESTAMPTZ   | default now()                    |
| deleted_at       | TIMESTAMPTZ   | nullable                         |

**order_items**

| Column     | Type          | Notes                         |
| ---------- | ------------- | ----------------------------- |
| id         | UUID          | PK, default gen_random_uuid() |
| order_id   | UUID          | FK -> orders.id               |
| menu_id    | UUID          | FK -> menus.id                |
| menu_name  | VARCHAR(180)  | not null                      |
| quantity   | INTEGER       | not null, > 0                 |
| unit_price | NUMERIC(12,2) | not null, >= 0                |
| subtotal   | NUMERIC(12,2) | not null, >= 0                |
| notes      | TEXT          | nullable                      |
| deleted_at | TIMESTAMPTZ   | nullable                      |

**customers**

| Column       | Type         | Notes                         |
| ------------ | ------------ | ----------------------------- |
| id           | UUID         | PK, default gen_random_uuid() |
| order_id     | UUID         | FK -> orders.id               |
| tenant_id    | UUID         | FK -> tenants.id              |
| full_name    | VARCHAR(150) | not null                      |
| phone_number | VARCHAR(20)  | nullable                      |
| created_at   | TIMESTAMPTZ  | default now()                 |
| deleted_at   | TIMESTAMPTZ  | nullable                      |

### 8.2 Legacy / Removed by Migrations

- `tenant_profiles` dibuat pada awalnya, lalu dipindahkan ke `tenants` dan dihapus pada migrasi berikutnya.

## 9. API Endpoints

### 9.1 Public & Health

| Method | Endpoint       | Fungsi                       | Role   |
| ------ | -------------- | ---------------------------- | ------ |
| GET    | /swagger/\*any | Swagger UI                   | Public |
| GET    | /health        | Liveness check               | Public |
| GET    | /ready         | Readiness check (DB + Redis) | Public |
| GET    | /api/v1/health | API health                   | Public |
| GET    | /api/v1/ready  | API readiness                | Public |

### 9.2 Auth

| Method | Endpoint              | Fungsi                                  | Role                     |
| ------ | --------------------- | --------------------------------------- | ------------------------ |
| POST   | /api/v1/auth/register | Register user (default role `CUSTOMER`) | Public                   |
| POST   | /api/v1/auth/login    | Login semua role                        | Public                   |
| POST   | /api/v1/auth/refresh  | Refresh token                           | Public                   |
| POST   | /api/v1/auth/logout   | Logout user                             | CUSTOMER, PARTNER, ADMIN |
| GET    | /api/v1/auth/me       | Ambil profil user login                 | CUSTOMER, PARTNER, ADMIN |

### 9.3 Users (Partner/Admin)

| Method | Endpoint                        | Fungsi                 | Role           |
| ------ | ------------------------------- | ---------------------- | -------------- |
| GET    | /api/v1/users                   | List user dalam tenant | PARTNER, ADMIN |
| GET    | /api/v1/users/:id               | Detail user            | PARTNER, ADMIN |
| PUT    | /api/v1/users/:id               | Update user            | PARTNER, ADMIN |
| DELETE | /api/v1/users/:id               | Soft delete user       | PARTNER, ADMIN |
| PATCH  | /api/v1/users/:id/toggle-active | Toggle aktif user      | PARTNER, ADMIN |

### 9.4 Tenant Context & Partner

| Method | Endpoint                | Fungsi                             | Role                 |
| ------ | ----------------------- | ---------------------------------- | -------------------- |
| GET    | /api/v1/customer/me     | Validasi context customer + tenant | BASIC (legacy alias) |
| POST   | /api/v1/partner/tenants | Buat tenant untuk partner          | PARTNER              |
| GET    | /api/v1/partner/tenants | List tenant milik partner          | PARTNER              |
| GET    | /api/v1/partner/me      | Validasi context partner + tenant  | PARTNER              |

### 9.5 Catalog (Categories & Menus)

**Categories (Partner-only + TenantGuard)**

| Method | Endpoint                             | Fungsi                | Role    |
| ------ | ------------------------------------ | --------------------- | ------- |
| POST   | /api/v1/categories                   | Create category       | PARTNER |
| GET    | /api/v1/categories                   | List categories       | PARTNER |
| GET    | /api/v1/categories/:id               | Detail category       | PARTNER |
| PUT    | /api/v1/categories/:id               | Update category       | PARTNER |
| DELETE | /api/v1/categories/:id               | Soft delete category  | PARTNER |
| PATCH  | /api/v1/categories/:id/toggle-active | Toggle aktif category | PARTNER |
| PATCH  | /api/v1/categories/reorder           | Reorder category      | PARTNER |

**Menus (shared path dengan role berbeda)**

| Method | Endpoint                           | Fungsi                   | Role                     |
| ------ | ---------------------------------- | ------------------------ | ------------------------ |
| GET    | /api/v1/menus                      | List menu                | CUSTOMER, PARTNER, ADMIN |
| GET    | /api/v1/menus/:id                  | Detail menu              | CUSTOMER, PARTNER, ADMIN |
| POST   | /api/v1/menus                      | Create menu              | PARTNER                  |
| PUT    | /api/v1/menus/:id                  | Update menu              | PARTNER                  |
| DELETE | /api/v1/menus/:id                  | Soft delete menu         | PARTNER                  |
| PATCH  | /api/v1/menus/:id/toggle-available | Toggle ketersediaan menu | PARTNER                  |

Catatan: endpoint GET menus tidak memakai `TenantGuard`; tenant biasanya diambil dari query param (mis. `tenantId`).

### 9.6 Orders & Customers

| Method | Endpoint                    | Fungsi                      | Role                     |
| ------ | --------------------------- | --------------------------- | ------------------------ |
| POST   | /api/v1/orders              | Create order                | CUSTOMER, PARTNER, ADMIN |
| GET    | /api/v1/orders/:id          | Detail order                | CUSTOMER, PARTNER, ADMIN |
| POST   | /api/v1/orders/:id/customer | Create customer untuk order | CUSTOMER, PARTNER, ADMIN |
| GET    | /api/v1/orders/:id/customer | Detail customer order       | CUSTOMER, PARTNER, ADMIN |
| PUT    | /api/v1/orders/:id/customer | Update customer order       | CUSTOMER, PARTNER, ADMIN |
| GET    | /api/v1/orders              | List order (partner view)   | PARTNER                  |
| PATCH  | /api/v1/orders/:id/status   | Update status order         | PARTNER                  |
| DELETE | /api/v1/orders/:id          | Soft delete order           | PARTNER                  |

Catatan: endpoint order untuk customer tidak memakai `TenantGuard`; tenant biasanya diambil dari query param saat order.

### 9.7 Dining Tables

| Method | Endpoint                         | Fungsi           | Role    |
| ------ | -------------------------------- | ---------------- | ------- |
| POST   | /api/v1/dining-tables            | Create meja      | PARTNER |
| GET    | /api/v1/dining-tables            | List meja        | PARTNER |
| GET    | /api/v1/dining-tables/:id        | Detail meja      | PARTNER |
| GET    | /api/v1/dining-tables/:id/status | Status meja      | PARTNER |
| PUT    | /api/v1/dining-tables/:id        | Update meja      | PARTNER |
| DELETE | /api/v1/dining-tables/:id        | Soft delete meja | PARTNER |

### 9.8 Reports

| Method | Endpoint                        | Fungsi         | Role    |
| ------ | ------------------------------- | -------------- | ------- |
| GET    | /api/v1/reports/revenue         | Rekap revenue  | PARTNER |
| GET    | /api/v1/reports/top-menus       | Top menu       | PARTNER |
| GET    | /api/v1/reports/orders-by-table | Order per meja | PARTNER |
| GET    | /api/v1/reports/daily-summary   | Summary harian | PARTNER |

## 11. Response Schema Convention

### Standard Response Envelope

Semua endpoint **wajib** menggunakan format response yang sama:

**Success:**

```json
{
  "success": true,
  "message": "Data berhasil diambil",
  "data": { ... }
}
```

**Success dengan Pagination:**

```json
{
  "success": true,
  "message": "Data berhasil diambil",
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8
  }
}
```

**Error:**

```json
{
  "success": false,
  "message": "Pesan error yang deskriptif",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [ ... ]
  }
}
```

### Schema Convention: camelCase

- Semua field di JSON request dan response menggunakan **camelCase**
- Contoh: `tenantId`, `createdAt`, `isAvailable`, `orderSource`, `unitPrice`
- Kolom database menggunakan `snake_case` (dikonversi di layer DTO)

### HTTP Status Code Convention

| Status                      | Kasus                                         |
| --------------------------- | --------------------------------------------- |
| `200 OK`                    | GET, PATCH, DELETE berhasil                   |
| `201 Created`               | POST create resource berhasil                 |
| `400 Bad Request`           | Validation error / payload tidak valid        |
| `401 Unauthorized`          | Token tidak ada atau expired                  |
| `403 Forbidden`             | Role tidak punya akses ke resource ini        |
| `404 Not Found`             | Resource tidak ditemukan                      |
| `409 Conflict`              | Duplicate resource (idempotency key conflict) |
| `500 Internal Server Error` | Error server yang tidak terduga               |

---

## 12. Error Handling Convention

Semua error harus dipetakan ke response envelope yang konsisten.

### Error Code List

| Code                 | Status | Keterangan                       |
| -------------------- | ------ | -------------------------------- |
| `VALIDATION_ERROR`   | 400    | Input tidak valid                |
| `UNAUTHORIZED`       | 401    | Token tidak valid / tidak ada    |
| `FORBIDDEN`          | 403    | Role tidak memiliki izin         |
| `NOT_FOUND`          | 404    | Resource tidak ditemukan         |
| `CONFLICT`           | 409    | Duplicate / idempotency conflict |
| `TENANT_NOT_FOUND`   | 404    | Tenant tidak ditemukan           |
| `TENANT_INACTIVE`    | 403    | Tenant tidak aktif / suspended   |
| `ORDER_ALREADY_PAID` | 409    | Order sudah dibayar              |
| `INTERNAL_ERROR`     | 500    | Error internal server            |

---

## 15. Domain Structure per Feature

### Contoh Detail: Domain `menu`

```
internal/domains/menu/
├── delivery/http/
│   └── menu_handler.go
│       - GetMenuListHandler(c *gin.Context)
│       - GetMenuDetailHandler(c *gin.Context)
│       - CreateMenuHandler(c *gin.Context)
│       - UpdateMenuHandler(c *gin.Context)
│       - DeleteMenuHandler(c *gin.Context)
│       - ToggleAvailabilityHandler(c *gin.Context)
│
├── usecase/
│   └── menu_usecase.go
│       - interface MenuUsecase
│       - struct menuUsecase implements MenuUsecase
│       - GetMenuList(ctx, tenantID, filter)
│       - CreateMenu(ctx, tenantID, req)
│       - UpdateMenu(ctx, tenantID, id, req)
│       - DeleteMenu(ctx, tenantID, id)
│
├── domain/
│   └── menu_domain.go
│       - struct Menu {}
│       - struct MenuFilter {}
│
├── repository/
│   └── menu_repository.go
│       - interface MenuRepository
│       - struct menuRepository implements MenuRepository
│       - FindAll(ctx, tenantID, filter)
│       - FindByID(ctx, tenantID, id)
│       - Create(ctx, menu)
│       - Update(ctx, menu)
│       - Delete(ctx, tenantID, id)
│
└── dto/
    ├── menu_request.go
    │   - struct CreateMenuRequest { ... json:"camelCase" }
    │   - struct UpdateMenuRequest { ... json:"camelCase" }
    │   - struct MenuFilter { ... }
    │
    └── menu_response.go
        - struct MenuResponse { ... json:"camelCase" }
```

---

## Quick Start

```bash
# 1. Clone repository
git clone <repo-url>
cd saas_gangsta

# 2. Copy dan isi environment variables
cp .env.example .env
# Edit .env sesuai konfigurasi Supabase dan JWT kamu

# 3. Install dependencies
go mod tidy

# 4. Jalankan development server
go run ./cmd/api/main.go

# ATAU jalankan dengan Docker
docker compose -f deployments/docker-compose.yml up --build
```

Server akan berjalan di `http://localhost:8080`
Melalui Nginx (Docker): `http://localhost:80`

Health check: `GET http://localhost:8080/health`

Railway deployment guide: [docs/railway.md](docs/railway.md)

---

> **Last Updated**: April 2026
> **Module**: `github.com/dhegas/saas_gangsta`
> **Stack**: Go 1.23+ · Gin · GORM · PostgreSQL (Supabase) · Docker · Nginx
