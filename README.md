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
   - [Health Check](#91-health-check)
   - [Auth](#92-auth)
   - [Customer](#93-customer)
   - [Merchant](#94-merchant)
   - [Admin](#95-admin)
10. [CORS Configuration](#10-cors-configuration)
11. [Response Schema Convention](#11-response-schema-convention)
12. [Error Handling Convention](#12-error-handling-convention)
13. [Docker & Container Setup](#13-docker--container-setup)
14. [Clean Code Guidelines](#14-clean-code-guidelines)
15. [Module Structure per Feature](#15-module-structure-per-feature)
16. [Development Roadmap](#16-development-roadmap)
17. [Definition of Done](#17-definition-of-done)
18. [Technical Risks](#18-technical-risks)

---

## 1. Project Overview

`saas_gangsta` adalah layanan backend untuk platform SaaS yang membantu **UMKM kuliner Indonesia** mengelola operasional toko secara digital.

### Tujuan Bisnis

Platform ini menggabungkan kebutuhan utama operasional toko makanan:

| Modul | Keterangan |
|---|---|
| Digital Menu | Pelanggan scan QR → lihat menu digital |
| Self Ordering | Pelanggan order dari meja tanpa panggil pelayan |
| POS / Kasir | Merchant input order manual via kasir digital |
| Manajemen Meja | Monitor kondisi dan status meja real-time |
| Laporan Penjualan | Rekap harian/mingguan/bulanan per merchant |
| Membership / Subscription | Admin kelola paket langganan SaaS |

### Tiga Role Utama

| Role | Deskripsi |
|---|---|
| **Customer** | Pelanggan toko. Scan QR, lihat menu, order, bayar, cek status, review |
| **Merchant** | Pemilik toko. Kelola menu, terima order, POS, laporan, profil toko |
| **Admin** | Pengelola platform SaaS. Kelola tenant, membership, billing, user global |

### Model Bisnis

Sistem menggunakan model **SaaS multi-tenant**, di mana setiap merchant (tenant) memiliki data yang terisolasi satu sama lain. Merchant berlangganan bulanan untuk mengakses platform.

---

## 2. Tech Stack

| Komponen | Teknologi |
|---|---|
| Language | Go 1.23+ |
| HTTP Framework | Gin (`github.com/gin-gonic/gin`) |
| Database | PostgreSQL via Supabase (managed) |
| ORM | GORM (`gorm.io/gorm`) + `pgx` |
| API Gateway | Nginx (reverse proxy & rate limiter) |
| Authentication | JWT (access token + refresh token) |
| Container | Docker + Docker Compose |
| Environment Config | `.env` + `godotenv` |
| Validation | `go-playground/validator/v10` |
| Logging | `log/slog` (standard library Go 1.21+) |
| UUID | `github.com/google/uuid` |
| Schema Convention | **camelCase** untuk semua JSON response & request |

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
                        │  │ CORS     │  │ /merchant        │  │
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

### Prinsip Arsitektur Wajib

- **Multi-tenant first**: Semua tabel bisnis wajib punya `tenantId`. Semua query merchant/customer harus difilter `tenantId`.
- **Role-separated routes**: Route dipisah per role (`/customer`, `/merchant`, `/admin`).
- **Tenant context dari JWT**: `tenantId` diambil dari JWT claim, bukan dari input user.
- **DB transaction untuk flow kritis**: Create order, payment, void transaksi wajib dalam DB transaction.
- **Idempotency key**: Digunakan untuk endpoint create order dan payment.
- **Standard response format**: Semua endpoint menggunakan envelope response yang konsisten.

---

## 4. Folder Structure

```
saas_gangsta/
├── cmd/
│   └── api/
│       └── main.go                  # Entry point aplikasi
│
├── internal/
│   ├── bootstrap/
│   │   └── app.go                   # Registrasi semua dependency (router, db, middleware)
│   │
│   ├── common/
│   │   ├── config/
│   │   │   └── config.go            # Load & parse env config
│   │   ├── middleware/
│   │   │   ├── auth.go              # JWT auth middleware
│   │   │   ├── role_guard.go        # Role-based access control
│   │   │   ├── tenant_guard.go      # Tenant context injector
│   │   │   ├── cors.go              # CORS middleware
│   │   │   ├── logger.go            # Request logging middleware
│   │   │   └── recovery.go          # Panic recovery middleware
│   │   ├── response/
│   │   │   └── response.go          # Standard response envelope helper
│   │   ├── errors/
│   │   │   ├── app_error.go         # Custom error types
│   │   │   └── error_handler.go     # Global error mapping
│   │   ├── auth/
│   │   │   └── jwt.go               # JWT generate & parse helper
│   │   └── tenant/
│   │       └── context.go           # Tenant context helper (get tenantId dari ctx)
│   │
│   ├── database/
│   │   └── db.go                    # GORM database connection
│   │
│   └── modules/
│       ├── auth/
│       │   ├── delivery/http/
│       │   │   └── auth_handler.go
│       │   ├── usecase/
│       │   │   └── auth_usecase.go
│       │   ├── domain/
│       │   │   └── auth_domain.go
│       │   ├── repository/
│       │   │   └── auth_repository.go
│       │   └── dto/
│       │       ├── login_request.go
│       │       └── auth_response.go
│       │
│       ├── customer_menu/
│       │   ├── delivery/http/
│       │   ├── usecase/
│       │   ├── domain/
│       │   ├── repository/
│       │   └── dto/
│       │
│       ├── customer_order/
│       │   ├── delivery/http/
│       │   ├── usecase/
│       │   ├── domain/
│       │   ├── repository/
│       │   └── dto/
│       │
│       ├── payment/
│       │   ├── delivery/http/
│       │   ├── usecase/
│       │   ├── domain/
│       │   ├── repository/
│       │   └── dto/
│       │
│       ├── merchant_menu/
│       │   ├── delivery/http/
│       │   ├── usecase/
│       │   ├── domain/
│       │   ├── repository/
│       │   └── dto/
│       │
│       ├── merchant_pos/
│       │   ├── delivery/http/
│       │   ├── usecase/
│       │   ├── domain/
│       │   ├── repository/
│       │   └── dto/
│       │
│       ├── merchant_table/
│       │   ├── delivery/http/
│       │   ├── usecase/
│       │   ├── domain/
│       │   ├── repository/
│       │   └── dto/
│       │
│       ├── merchant_report/
│       │   ├── delivery/http/
│       │   ├── usecase/
│       │   ├── domain/
│       │   ├── repository/
│       │   └── dto/
│       │
│       ├── admin_tenant/
│       │   ├── delivery/http/
│       │   ├── usecase/
│       │   ├── domain/
│       │   ├── repository/
│       │   └── dto/
│       │
│       └── admin_subscription/
│           ├── delivery/http/
│           ├── usecase/
│           ├── domain/
│           ├── repository/
│           └── dto/
│
├── pkg/
│   ├── validator/
│   │   └── validator.go             # Custom validator helper
│   ├── pagination/
│   │   └── pagination.go            # Pagination helper
│   └── logger/
│       └── logger.go                # Logger instance (wraps slog)
│
├── migrations/
│   └── *.sql                        # SQL migration files (ordered by timestamp)
│
├── deployments/
│   ├── docker-compose.yml           # Docker Compose (api + postgres local dev)
│   ├── Dockerfile                   # Multi-stage Dockerfile
│   └── nginx/
│       └── nginx.conf               # Nginx config untuk API Gateway
│
├── docs/
│   └── openapi.yaml                 # OpenAPI/Swagger documentation
│
├── scripts/
│   └── seed.sql                     # Seed data untuk development
│
├── .env                             # Environment variables (jangan di-commit)
├── .env.example                     # Template environment variables
├── go.mod
├── go.sum
└── README.md
```

### Struktur Per Module (Wajib Konsisten)

Setiap module di dalam `internal/modules/<module_name>/` **wajib** mengikuti struktur berikut:

```
internal/modules/<module_name>/
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

Buat file `.env` di root folder `saas_gangsta/`:

```env
# Application
APP_ENV=development
APP_PORT=8080
APP_NAME=saas_gangsta

# Database (Supabase PostgreSQL)
DATABASE_URL=postgresql://postgres:<password>@<supabase-host>:5432/postgres?sslmode=require

# JWT
JWT_SECRET=your-very-strong-jwt-secret-key-minimum-32-chars
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=7d

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com

# Supabase (jika butuh Supabase client langsung)
SUPABASE_URL=https://<project-ref>.supabase.co
SUPABASE_SERVICE_ROLE_KEY=your-supabase-service-role-key
```

> ⚠️ **PENTING**: Jangan pernah commit `.env` ke repository. Gunakan `.env.example` sebagai template.
>
> ⚠️ `SUPABASE_SERVICE_ROLE_KEY` hanya boleh ada di backend. **Jangan pernah expose ke frontend.**

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

| Fungsi | Keterangan |
|---|---|
| Reverse Proxy | Forward request dari port 80/443 ke Go server :8080 |
| Rate Limiting | Batasi jumlah request per IP/client |
| SSL Termination | Handle HTTPS (jika digunakan) |
| Request Buffering | Buffer upload sebelum diteruskan ke backend |
| Health Check | Probe endpoint `/health` |

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

Contoh: `POST /api/v1/auth/login`, `GET /api/v1/merchant/menus`

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

### JWT Payload (Claims)

```json
{
  "sub":      "user-uuid",
  "role":     "merchant",
  "tenantId": "tenant-uuid",
  "iat":      1700000000,
  "exp":      1700000900
}
```

### Role Values

| Role | Keterangan |
|---|---|
| `customer` | Pelanggan toko |
| `merchant` | Pemilik / pengelola toko |
| `admin` | Admin platform SaaS |

### Middleware Stack per Route Group

```
Public Routes (/auth, /health):
  → CORS → Logger → Recovery

Customer Routes (/api/v1/customer/**):
  → CORS → Logger → Recovery → JWTAuth → RoleGuard("customer") → TenantGuard

Merchant Routes (/api/v1/merchant/**):
  → CORS → Logger → Recovery → JWTAuth → RoleGuard("merchant") → TenantGuard

Admin Routes (/api/v1/admin/**):
  → CORS → Logger → Recovery → JWTAuth → RoleGuard("admin")
```

---

## 8. Database Schema (PostgreSQL Supabase)

### Prinsip Schema

- Semua tabel menggunakan `UUID` sebagai primary key
- Kolom `tenantId` wajib ada di semua tabel bisnis (kecuali tabel platform-level)
- Timestamp: `createdAt`, `updatedAt`, dan `deletedAt` (soft delete)
- Naming convention kolom: **snake_case** di database, **camelCase** di JSON response/request

### Entitas Inti

```sql
-- Tabel tenant (data per merchant/toko)
tenants
  id UUID PK
  name VARCHAR
  slug VARCHAR UNIQUE
  status VARCHAR  -- active | inactive | suspended
  created_at TIMESTAMP
  updated_at TIMESTAMP

-- User (seluruh platform)
users
  id UUID PK
  tenant_id UUID FK → tenants.id (nullable untuk admin platform)
  email VARCHAR UNIQUE
  password_hash VARCHAR
  role VARCHAR  -- customer | merchant | admin
  is_active BOOLEAN
  created_at TIMESTAMP
  updated_at TIMESTAMP

-- Profil merchant (data toko)
merchant_profiles
  id UUID PK
  tenant_id UUID FK → tenants.id
  store_name VARCHAR
  address TEXT
  phone VARCHAR
  logo_url VARCHAR
  opening_hours JSONB
  created_at TIMESTAMP
  updated_at TIMESTAMP

-- Kategori menu
categories
  id UUID PK
  tenant_id UUID FK
  name VARCHAR
  description TEXT
  sort_order INTEGER
  is_active BOOLEAN
  created_at TIMESTAMP
  updated_at TIMESTAMP

-- Menu item
menus
  id UUID PK
  tenant_id UUID FK
  category_id UUID FK → categories.id
  name VARCHAR
  description TEXT
  price NUMERIC(12,2)
  image_url VARCHAR
  is_available BOOLEAN
  created_at TIMESTAMP
  updated_at TIMESTAMP

-- Meja
tables
  id UUID PK
  tenant_id UUID FK
  table_number VARCHAR
  capacity INTEGER
  status VARCHAR  -- empty | occupied | reserved
  created_at TIMESTAMP
  updated_at TIMESTAMP

-- Order
orders
  id UUID PK
  tenant_id UUID FK
  table_id UUID FK → tables.id
  user_id UUID FK → users.id (nullable untuk guest order)
  idempotency_key VARCHAR UNIQUE
  status VARCHAR  -- pending | accepted | cooking | ready | done | canceled
  subtotal NUMERIC(12,2)
  tax NUMERIC(12,2)
  total NUMERIC(12,2)
  notes TEXT
  order_source VARCHAR  -- self_order | pos
  created_at TIMESTAMP
  updated_at TIMESTAMP

-- Order items
order_items
  id UUID PK
  order_id UUID FK → orders.id
  menu_id UUID FK → menus.id
  quantity INTEGER
  unit_price NUMERIC(12,2)
  subtotal NUMERIC(12,2)
  notes TEXT

-- Payment
payments
  id UUID PK
  order_id UUID FK → orders.id
  tenant_id UUID FK
  idempotency_key VARCHAR UNIQUE
  method VARCHAR  -- cash | qris | transfer
  status VARCHAR  -- pending | paid | failed | refunded
  amount NUMERIC(12,2)
  paid_at TIMESTAMP
  created_at TIMESTAMP
  updated_at TIMESTAMP

-- Subscription plans (admin level)
subscription_plans
  id UUID PK
  name VARCHAR
  description TEXT
  price NUMERIC(12,2)
  billing_cycle VARCHAR  -- monthly | yearly
  features JSONB
  is_active BOOLEAN
  created_at TIMESTAMP
  updated_at TIMESTAMP

-- Subscriptions per tenant
subscriptions
  id UUID PK
  tenant_id UUID FK
  plan_id UUID FK → subscription_plans.id
  status VARCHAR  -- active | expired | canceled | trial
  started_at TIMESTAMP
  expires_at TIMESTAMP
  created_at TIMESTAMP
  updated_at TIMESTAMP

-- Audit log
audit_logs
  id UUID PK
  tenant_id UUID FK
  user_id UUID FK
  action VARCHAR  -- e.g. MENU_DELETED, ORDER_CANCELED, VOID_TRANSACTION
  entity_type VARCHAR
  entity_id UUID
  metadata JSONB
  created_at TIMESTAMP
```

---

## 9. API Endpoints

Semua endpoint menggunakan prefix `/api/v1/`.

---

### 9.1 Health Check

| Method | Path | Auth | Deskripsi |
|---|---|---|---|
| GET | `/health` | Public | Cek status server hidup |
| GET | `/ready` | Public | Cek readiness (DB connection check) |

**Response `GET /health`:**
```json
{
  "status": "ok",
  "service": "saas_gangsta",
  "timestamp": "2026-04-15T22:30:00Z"
}
```

---

### 9.2 Auth

| Method | Path | Auth | Deskripsi |
|---|---|---|---|
| POST | `/api/v1/auth/login` | Public | Login user (semua role) |
| POST | `/api/v1/auth/refresh` | Public | Refresh access token |
| POST | `/api/v1/auth/logout` | JWT | Logout & invalidate refresh token |
| GET | `/api/v1/auth/me` | JWT | Get current user info |

**Request `POST /api/v1/auth/login`:**
```json
{
  "email": "merchant@example.com",
  "password": "secret123"
}
```

**Response `POST /api/v1/auth/login`:**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "accessToken": "eyJhbG...",
    "refreshToken": "eyJhbG...",
    "user": {
      "id": "uuid",
      "email": "merchant@example.com",
      "role": "merchant",
      "tenantId": "tenant-uuid"
    }
  }
}
```

---

### 9.3 Customer

> **Auth required**: Bearer token (role: `customer`)
> **TenantId**: otomatis dari JWT, tidak perlu dikirim di request

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/customer/menus` | Lihat daftar menu (berdasarkan tenantId dari JWT) |
| GET | `/api/v1/customer/menus/:id` | Detail menu |
| GET | `/api/v1/customer/categories` | Daftar kategori |
| POST | `/api/v1/customer/orders` | Buat order baru (self-order) |
| GET | `/api/v1/customer/orders/:id` | Detail order |
| GET | `/api/v1/customer/orders/:id/status` | Cek status order real-time |
| POST | `/api/v1/customer/payments` | Buat payment untuk order |
| GET | `/api/v1/customer/payments/:id/status` | Cek status payment |
| GET | `/api/v1/customer/transactions` | Riwayat transaksi customer |
| GET | `/api/v1/customer/transactions/:id` | Detail transaksi |
| POST | `/api/v1/customer/reviews` | Submit review/rating |

**Request `POST /api/v1/customer/orders`:**
```json
{
  "tableId": "table-uuid",
  "idempotencyKey": "client-generated-uuid",
  "items": [
    {
      "menuId": "menu-uuid",
      "quantity": 2,
      "notes": "tanpa sambal"
    }
  ],
  "notes": "pesanan meja 3"
}
```

---

### 9.4 Merchant

> **Auth required**: Bearer token (role: `merchant`)
> **TenantId**: otomatis dari JWT claim

#### 9.4.1 Menu Management

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/merchant/categories` | Daftar kategori milik tenant |
| POST | `/api/v1/merchant/categories` | Buat kategori baru |
| PATCH | `/api/v1/merchant/categories/:id` | Update kategori |
| DELETE | `/api/v1/merchant/categories/:id` | Hapus kategori |
| GET | `/api/v1/merchant/menus` | Daftar menu (support: `?page`, `?limit`, `?search`, `?categoryId`) |
| POST | `/api/v1/merchant/menus` | Buat menu baru |
| PATCH | `/api/v1/merchant/menus/:id` | Update menu |
| DELETE | `/api/v1/merchant/menus/:id` | Hapus menu |
| PATCH | `/api/v1/merchant/menus/:id/availability` | Toggle status tersedia/tidak |

#### 9.4.2 Table Management

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/merchant/tables` | Daftar meja |
| POST | `/api/v1/merchant/tables` | Tambah meja |
| PATCH | `/api/v1/merchant/tables/:id` | Update info meja |
| DELETE | `/api/v1/merchant/tables/:id` | Hapus meja |

#### 9.4.3 Order Board (POS & Incoming Orders)

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/merchant/orders` | Daftar order aktif |
| GET | `/api/v1/merchant/orders/board` | Order board per meja (untuk tampilan kasir) |
| GET | `/api/v1/merchant/orders/:id` | Detail order |
| PATCH | `/api/v1/merchant/orders/:id/status` | Update status order |
| POST | `/api/v1/merchant/pos/orders` | Buat order manual lewat POS |

#### 9.4.4 Transaction

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/merchant/transactions` | Daftar transaksi (support filter tanggal, status) |
| GET | `/api/v1/merchant/transactions/:id` | Detail transaksi |

#### 9.4.5 Report

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/merchant/reports/daily` | Laporan harian (`?date=YYYY-MM-DD`) |
| GET | `/api/v1/merchant/reports/weekly` | Laporan mingguan (`?week=YYYY-WNN`) |
| GET | `/api/v1/merchant/reports/monthly` | Laporan bulanan (`?month=YYYY-MM`) |
| GET | `/api/v1/merchant/reports/summary` | Summary dashboard merchant |

#### 9.4.6 Merchant Profile

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/merchant/profile` | Get profil toko |
| PUT | `/api/v1/merchant/profile` | Update profil toko |

---

### 9.5 Admin

> **Auth required**: Bearer token (role: `admin`)
> Admin TIDAK terikat tenant → dapat akses semua data platform

#### 9.5.1 Tenant Management

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/admin/tenants` | Daftar semua tenant |
| GET | `/api/v1/admin/tenants/:id` | Detail tenant |
| POST | `/api/v1/admin/tenants` | Registrasi tenant baru |
| PATCH | `/api/v1/admin/tenants/:id` | Update info tenant |
| PATCH | `/api/v1/admin/tenants/:id/status` | Aktifkan / nonaktifkan tenant |

#### 9.5.2 Subscription Management

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/admin/subscription-plans` | Daftar paket berlangganan |
| POST | `/api/v1/admin/subscription-plans` | Buat paket baru |
| PUT | `/api/v1/admin/subscription-plans/:id` | Update paket |
| DELETE | `/api/v1/admin/subscription-plans/:id` | Hapus paket |
| GET | `/api/v1/admin/subscriptions` | Monitor status berlangganan semua tenant |
| PATCH | `/api/v1/admin/tenants/:id/subscription` | Update subscription tenant |

#### 9.5.3 Admin Dashboard

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/admin/dashboard` | Overview platform (jumlah tenant aktif, total transaksi, dll) |
| GET | `/api/v1/admin/users` | Daftar user platform |
| PATCH | `/api/v1/admin/users/:id/status` | Aktifkan / nonaktifkan user |

---

## 10. CORS Configuration

CORS dikonfigurasi di middleware `internal/common/middleware/cors.go`.

```go
// Contoh implementasi CORS middleware
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        allowedOrigins := cfg.CORSAllowedOrigins // dari env

        c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ","))
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers",
            "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Idempotency-Key",
        )
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
        c.Writer.Header().Set("Access-Control-Max-Age", "86400")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

**Aturan CORS:**

| Setting | Value |
|---|---|
| Allowed Origins | Dari env `CORS_ALLOWED_ORIGINS` (comma-separated) |
| Allowed Methods | GET, POST, PUT, PATCH, DELETE, OPTIONS |
| Allowed Headers | Authorization, Content-Type, X-Idempotency-Key |
| Credentials | true |
| Max Age | 86400 seconds (24 jam) |

---

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

| Status | Kasus |
|---|---|
| `200 OK` | GET, PATCH, DELETE berhasil |
| `201 Created` | POST create resource berhasil |
| `400 Bad Request` | Validation error / payload tidak valid |
| `401 Unauthorized` | Token tidak ada atau expired |
| `403 Forbidden` | Role tidak punya akses ke resource ini |
| `404 Not Found` | Resource tidak ditemukan |
| `409 Conflict` | Duplicate resource (idempotency key conflict) |
| `500 Internal Server Error` | Error server yang tidak terduga |

---

## 12. Error Handling Convention

Semua error harus dipetakan ke response envelope yang konsisten.

### Error Code List

| Code | Status | Keterangan |
|---|---|---|
| `VALIDATION_ERROR` | 400 | Input tidak valid |
| `UNAUTHORIZED` | 401 | Token tidak valid / tidak ada |
| `FORBIDDEN` | 403 | Role tidak memiliki izin |
| `NOT_FOUND` | 404 | Resource tidak ditemukan |
| `CONFLICT` | 409 | Duplicate / idempotency conflict |
| `TENANT_NOT_FOUND` | 404 | Tenant tidak ditemukan |
| `TENANT_INACTIVE` | 403 | Tenant tidak aktif / suspended |
| `ORDER_ALREADY_PAID` | 409 | Order sudah dibayar |
| `INTERNAL_ERROR` | 500 | Error internal server |

---

## 13. Docker & Container Setup

### `deployments/Dockerfile`

```dockerfile
# Stage 1: Builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api/main.go

# Stage 2: Runtime
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/bin/api .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./api"]
```

### `deployments/docker-compose.yml`

```yaml
version: "3.9"

services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - api

  api:
    build:
      context: ..
      dockerfile: deployments/Dockerfile
    env_file:
      - ../.env
    ports:
      - "8080:8080"
    restart: unless-stopped
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: saas_gangsta_dev
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

### Menjalankan dengan Docker

```bash
# Jalankan semua service
docker compose -f deployments/docker-compose.yml up --build

# Jalankan di background
docker compose -f deployments/docker-compose.yml up -d --build

# Stop semua service
docker compose -f deployments/docker-compose.yml down
```

### Menjalankan tanpa Docker (Development)

```bash
# Pastikan .env sudah dikonfigurasi
cd saas_gangsta/

# Install dependencies
go mod tidy

# Jalankan server development
go run ./cmd/api/main.go
```

---

## 14. Clean Code Guidelines

### Prinsip Wajib

1. **Single Responsibility**: Setiap file/struct punya satu tanggung jawab utama
2. **Dependency Injection**: Semua dependency diinjeksikan via constructor, bukan diakses global
3. **Interface-driven**: Usecase dan Repository menggunakan interface sehingga mudah di-mock
4. **No business logic di handler**: Handler hanya bind request → call usecase → return response
5. **No DB query di usecase**: Query hanya boleh ada di repository layer
6. **camelCase JSON**: Semua field JSON request/response menggunakan camelCase (`json:"camelCase"`)
7. **snake_case kolom DB**: Kolom database menggunakan snake_case, mapping via GORM tag

### Naming Convention

| Konteks | Convention | Contoh |
|---|---|---|
| Package | lowercase | `auth`, `middleware`, `config` |
| File Go | snake_case | `auth_handler.go`, `menu_usecase.go` |
| Struct/Interface | PascalCase | `MenuUsecase`, `OrderRepository` |
| Function/Method | camelCase (atau PascalCase jika exported) | `getMenuById`, `CreateOrder` |
| Variable lokal | camelCase | `userId`, `tenantId`, `menuList` |
| Constant | PascalCase atau SCREAMING_SNAKE_CASE | `StatusPending`, `MAX_PAGE_SIZE` |
| JSON field | camelCase | `"tenantId"`, `"createdAt"`, `"isAvailable"` |
| DB column | snake_case | `tenant_id`, `created_at`, `is_available` |
| Route path | kebab-case | `/merchant/menu-items`, `/admin/subscription-plans` |

### Interface Pattern (Repository & Usecase)

```go
// domain/menu_domain.go
type Menu struct {
    ID          string
    TenantID    string
    CategoryID  string
    Name        string
    Price       float64
    IsAvailable bool
    CreatedAt   time.Time
}

// repository interface
type MenuRepository interface {
    FindAll(ctx context.Context, tenantID string, filter MenuFilter) ([]Menu, int64, error)
    FindByID(ctx context.Context, tenantID string, id string) (*Menu, error)
    Create(ctx context.Context, menu *Menu) error
    Update(ctx context.Context, menu *Menu) error
    Delete(ctx context.Context, tenantID string, id string) error
}

// usecase interface
type MenuUsecase interface {
    GetMenuList(ctx context.Context, tenantID string, filter MenuFilter) ([]MenuResponse, PaginationMeta, error)
    GetMenuByID(ctx context.Context, tenantID string, id string) (*MenuResponse, error)
    CreateMenu(ctx context.Context, tenantID string, req CreateMenuRequest) (*MenuResponse, error)
    UpdateMenu(ctx context.Context, tenantID string, id string, req UpdateMenuRequest) (*MenuResponse, error)
    DeleteMenu(ctx context.Context, tenantID string, id string) error
}
```

### Context Usage

Semua function yang mengakses DB atau melakukan IO **wajib** menerima `context.Context` sebagai parameter pertama:

```go
func (r *menuRepository) FindAll(ctx context.Context, tenantID string, filter MenuFilter) ([]Menu, int64, error) {
    // ...
}
```

### Tenant Guard

TenantId selalu diambil dari JWT claim yang sudah divalidasi middleware, **bukan dari parameter request user**:

```go
// internal/common/tenant/context.go
func GetTenantID(c *gin.Context) (string, error) {
    tenantID, exists := c.Get("tenantId")
    if !exists {
        return "", errors.New("tenantId not found in context")
    }
    return tenantID.(string), nil
}
```

---

## 15. Module Structure per Feature

### Contoh Detail: Modul `merchant_menu`

```
internal/modules/merchant_menu/
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

## 16. Development Roadmap

### Fase 0 — Foundation (Week 1)

- [ ] Bootstrap project structure sesuai folder layout ini
- [ ] Setup config env loader (`internal/common/config/config.go`)
- [ ] Setup GORM connection ke Supabase (`internal/database/db.go`)
- [ ] Setup base middleware: logger, recovery, CORS
- [ ] Setup Gin router dengan prefix `/api/v1/`
- [ ] Implement `GET /health` dan `GET /ready`
- [ ] Setup Docker + docker-compose (nginx + api)
- [ ] Setup migration awal (`migrations/001_init_schema.sql`)

### Fase 1 — Auth & Multi-Tenant Core (Week 2)

- [ ] Implement JWT generate & parse (`internal/common/auth/jwt.go`)
- [ ] Implement `POST /auth/login`, `POST /auth/refresh`, `POST /auth/logout`
- [ ] Implement JWT auth middleware
- [ ] Implement role guard middleware (`RoleGuard("merchant")`)
- [ ] Implement tenant guard middleware (inject `tenantId` ke context)
- [ ] Unit test module auth (coverage ≥ 70%)

### Fase 2 — Merchant Menu Management (Week 3)

- [ ] CRUD kategori menu
- [ ] CRUD menu item
- [ ] Toggle status tersedia/tidak tersedia
- [ ] Pagination & search untuk list menu
- [ ] Filter `tenantId` konsisten di semua query
- [ ] Integration test CRUD menu

### Fase 3 — Customer Order Flow (Week 4-5)

- [ ] Endpoint customer lihat menu berdasarkan tenant
- [ ] Cart preview (POST /customer/carts/preview)
- [ ] Create order (dengan idempotency key)
- [ ] Order status lifecycle (state machine)
- [ ] Endpoint cek status order

### Fase 4 — POS & Transaction (Week 6)

- [ ] Merchant create order manual via POS
- [ ] Sinkron format order dari self-order dan POS
- [ ] Nomor struk unik per tenant
- [ ] Endpoint daftar & detail transaksi merchant
- [ ] Audit log untuk cancel/void transaksi

### Fase 5 — Payment (Week 7)

- [ ] Endpoint create payment intent
- [ ] Payment status lifecycle
- [ ] Idempotency untuk payment endpoint

### Fase 6 — Table Management & Real-time Board (Week 8)

- [ ] CRUD meja merchant
- [ ] Order board per meja
- [ ] Update status meja otomatis berdasarkan status order

### Fase 7 — Report & Admin Subscription (Week 9-10)

- [ ] Laporan harian/mingguan/bulanan per merchant
- [ ] Admin: kelola tenant, paket subscription, monitoring membership
- [ ] Admin dashboard overview

---

## 17. Definition of Done

Setiap endpoint/modul dinyatakan **Done** jika memenuhi semua kriteria berikut:

- [ ] Endpoint sesuai spec di README ini (method, path, payload, response)
- [ ] Semua field response menggunakan **camelCase**
- [ ] Semua query merchant/customer difilter oleh `tenantId`
- [ ] Input validation menggunakan `validator.v10` disertai pesan error yang deskriptif
- [ ] Middleware auth, role guard, dan tenant guard aktif di route yang sesuai
- [ ] Standard response envelope digunakan konsisten
- [ ] Logging minimal: request info + error yang relevan
- [ ] Tidak ada hardcode config (port, secret, dsb) — semua dari env
- [ ] Unit test untuk usecase (happy path + minimal 1 error path)

---

## 18. Technical Risks

| Risiko | Mitigasi |
|---|---|
| Kebocoran data antar tenant | Tenant guard di middleware + wajib filter `tenantId` di setiap query |
| Race condition pada update status order | DB transaction + row-level locking pada flow order kritis |
| Double payment akibat webhook retry | Idempotency key untuk semua endpoint create order & payment |
| Query report lambat saat data membesar | Index pada `tenant_id`, `created_at`, `status` sejak awal |
| JWT secret bocor | Simpan di env, rotasi berkala, gunakan secret length minimum 32 char |
| Salah konfigurasi RLS Supabase | Review RLS policy per role/tenant dengan test otomatis |

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

---

> **Last Updated**: April 2026
> **Module**: `github.com/dhegas/saas_gangsta`
> **Stack**: Go 1.23+ · Gin · GORM · PostgreSQL (Supabase) · Docker · Nginx
