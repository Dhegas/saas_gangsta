# Railway Deployment Setup

Panduan ini menyiapkan deployment backend ke Railway untuk:
- Service
- Environment variables
- Health check
- Domain

## 1. Create Service

1. Masuk ke Railway Dashboard.
2. Pilih New Project.
3. Pilih Deploy from GitHub Repo dan pilih repository ini.
4. Railway akan membaca [railway.toml](../railway.toml) dan menggunakan Dockerfile di [deployments/Dockerfile](../deployments/Dockerfile).

## 2. Environment Variables

Tambahkan variables berikut di Railway Service:

- APP_ENV=production
- APP_NAME=saas_gangsta
- DATABASE_URL=<Supabase pooled or direct connection string>
- REDIS_URL=<Railway Redis URL>
- JWT_SECRET=<strong secret min 32 chars>
- JWT_ACCESS_TOKEN_EXPIRY=15m
- JWT_REFRESH_TOKEN_EXPIRY=168h
- CORS_ALLOWED_ORIGINS=https://${{RAILWAY_PUBLIC_DOMAIN}}

Optional:
- APP_DOMAIN=https://${{RAILWAY_PUBLIC_DOMAIN}}

Catatan:
- APP_PORT tidak perlu diisi di Railway karena aplikasi otomatis membaca PORT dari Railway.
- Untuk local dev, tetap gunakan APP_PORT=8080.

## 3. Health Check

Health check sudah dikonfigurasi di [railway.toml](../railway.toml):
- Path: /health

Endpoint tersedia:
- /health
- /ready

Gunakan /health untuk platform-level health check agar service tetap stabil walau dependency eksternal sedang bermasalah sementara.

## 4. Domain Setup

1. Buka service di Railway.
2. Masuk tab Settings -> Networking.
3. Klik Generate Domain untuk domain default Railway.
4. Jika pakai custom domain, klik Custom Domain lalu tambahkan domain Anda.
5. Setelah domain aktif, update CORS_ALLOWED_ORIGINS agar mencakup domain frontend dan domain API.

## 5. Verify Deployment

Setelah deploy sukses:

1. Hit endpoint health:
   - https://<your-domain>/health
2. Hit readiness:
   - https://<your-domain>/ready
3. Pastikan log service menunjukkan startup pada port dari variable PORT Railway.
