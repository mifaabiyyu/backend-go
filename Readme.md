# 🛡️ Backend API - Golang Clean Architecture

Backend API ini dibangun menggunakan bahasa Go dengan pendekatan Clean Architecture, JWT Authentication, Role & Permission, Redis Cache, dan fitur Rate Limiting.

## 🚀 Fitur Utama

- ✅ **Authentication**: Login dengan JWT dan Basic Auth.
- 🧠 **Role & Permission**: Otentikasi dan otorisasi berdasarkan role user.
- 🔐 **JWT Middleware**: Validasi token dengan claims custom.
- ⚡ **Redis Caching**: Cache user data untuk mempercepat respon.
- 🛡️ **Rate Limiter**: Mencegah abuse dari IP yang sama.
- 🧱 **PostgreSQL + sqlc**: Query builder yang efisien dan type-safe.
- 🧪 **Testing Friendly**: Struktur yang mudah di-test dan scalable.

## 🪰 Teknologi

| Tools  | Keterangan                     |
| ------ | ------------------------------ |
| Go     | Bahasa utama backend           |
| Chi    | Lightweight router             |
| sqlc   | Query builder untuk PostgreSQL |
| Redis  | Cache untuk user session       |
| JWT v5 | Token-based authentication     |

## 🧽 Struktur Proyek

```
.
├── cmd/                    # Entry-point aplikasi
|   ├── api/
|   ├── migration/
├── internal/
│   ├── auth/               # JWT & Authenticator
│   ├── db/
│   │   └── generated/      # Hasil generate sqlc
│   ├── store/              # Redis layer
│   └── service/            # Business logic (opsional)
├── api/                    # Middleware & routing
├── utils/                  # Helper, response writer, logger
└── main.go
```

## 🔐 Authentication Flow

1. **Login** → Mendapatkan JWT Token
2. **Request selanjutnya** → Kirim token via `Authorization: Bearer <token>`
3. **Middleware**:
   - Validasi token JWT
   - Ambil `user_id` & `role_id` dari claims
   - Ambil data user dari Redis / DB
   - Cek permission berdasarkan role

## 📦 Instalasi & Jalankan

```bash
git clone https://github.com/namauser/backend-go-app.git
cd backend-go-app

go mod tidy
go run main.go
```

### .env Contoh

```
JWT_SECRET=mysecret
JWT_ISS=backend-apps
JWT_AUD=rahasia

REDIS_HOST=localhost:6379
REDIS_ENABLED=true
```

## 🧪 Contoh API

### 🔑 Login

```
POST /login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "your-password"
}
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVC..."
}
```

### 🛡️ Protected Endpoint

```
GET /dashboard
Authorization: Bearer <token>
```

## ✨ Kontribusi

Feel free to fork dan pull request! Jangan lupa bintangin ⭐ repo ini jika membantu kamu!

## 👨‍💼 Author

**Mifa Abiyyu**\
📧 [mifaabiyyu@gmail.com](mailto:mifaabiyyu@gmail.com)\
👥 (+62) 8151-5141-186

---

> Backend ini cocok dijadikan fondasi untuk proyek SaaS, aplikasi internal perusahaan, hingga MVP startup 🚀
