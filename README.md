# Helpdesk Internal Application

Aplikasi sistem tiket dukungan (*helpdesk*) internal untuk menggantikan alur email dan spreadsheet dalam menangani permintaan bantuan karyawan. Dibangun dengan arsitektur bersih menggunakan **Go (Gin)** di sisi backend, **Vue 3 (Composition API + Vite)** di sisi frontend, dan **PostgreSQL** sebagai database relasional.

---

## Daftar Isi
1. [Tech Stack](#-tech-stack)
2. [Fitur Utama & Bonus Terimplementasi](#-fitur-utama--bonus-terimplementasi)
3. [Akun Seed Pengguna](#-akun-seed-pengguna)
4. [Cara Menjalankan Aplikasi](#-cara-menjalankan-aplikasi)
   - [Opsi 1: Menggunakan Docker Compose (Direkomendasikan)](#opsi-1-menggunakan-docker-compose-direkomendasikan)
   - [Opsi 2: Menjalankan Secara Manual (Lokal)](#opsi-2-menjalankan-secara-manual-lokal)
5. [Menjalankan Pengujian (Testing)](#-menjalankan-pengujian-testing)
6. [Daftar Endpoint API](#-daftar-endpoint-api)
7. [Dokumentasi API (Postman Collection)](#-dokumentasi-api-postman-collection)
8. [Catatan Arsitektur & Keputusan Desain](#-catatan-arsitektur--keputusan-desain)

---

## Tech Stack

- **Backend**: Go 1.24, Gin Web Framework, GORM, bcrypt, golang-jwt/jwt v5, godotenv.
- **Frontend**: Vue 3 (Composition API & `<script setup>`), Vite, TypeScript, Pinia, Vue Router, Tailwind CSS, Lucide Icons, Axios.
- **Database**: PostgreSQL 16.
- **Containerization**: Docker & Docker Compose (Multi-stage build).
- **CI/CD**: GitHub Actions Pipeline.

---

## Fitur Utama & Bonus Terimplementasi

### 1. Autentikasi & RBAC (Role-Based Access Control)
- Registrasi publik (khusus role `requester`).
- Login berbasis JWT token dengan password terenkripsi bcrypt.
- Tiga peran pengguna dengan hak akses tegas:
  - **Requester**: Membuat tiket, melihat tiket miliknya, berkomentar publik pada tiket miliknya, menutup tiket miliknya.
  - **Agent**: Melihat semua tiket, mengambil alih (*assign*) tiket, mengubah prioritas & status, menulis komentar publik & internal notes.
  - **Admin**: Memiliki seluruh akses agent ditambah manajemen pengguna (melihat daftar user dan mengubah role pengguna).

### 2. Alur Transisi Status Tiket (Workflow)
Backend menerapkan validasi transisi status tiket:

| Dari Status | Ke Status | Hak Akses Role |
| :--- | :--- | :--- |
| **Open** | In Progress | Agent, Admin |
| **In Progress** | Resolved | Agent, Admin |
| **Resolved** | Closed | Requester, Agent, Admin |
| **Resolved** | Open *(Reopen)* | Requester |
| **Open** | Closed *(Cancel)* | Requester, Admin |

*Transisi di luar tabel ini ditolak backend dengan status HTTP 403 Forbidden.*

### 3. Komentar & Catatan Internal (Internal Notes)
- Komentar publik dapat dibaca oleh pembuat tiket, agent, dan admin.
- Catatan internal (`is_internal = true`) hanya dapat ditulis dan dibaca oleh Agent dan Admin. Backend memfilter catatan internal sehingga tidak pernah terkirim ke Requester.

### 4. Fitur Bonus Opsional
- **Aturan SLA Sederhana**: Tiket dengan prioritas `urgent` yang belum ditangani (masih berstatus `open`) lebih dari 2 jam otomatis ditandai `is_overdue: true` dan menampilkan badge merah menyala `SLA Overdue`.
- **Halaman Statistik Admin**: Dashboard statistik operasional menampilkan metrik real-time: Total Tiket, Open, In Progress, Resolved, Closed, Total Pengguna, dan Tiket Overdue SLA.
- **CI Pipeline (GitHub Actions)**: Pipeline otomatis yang menguji kode backend (linting `go vet` dan integration tests) serta verifikasi tipe dan build frontend (`npm run build`).

---

## Akun Seed Pengguna

Aplikasi otomatis melakukan *seeding* akun saat backend pertama kali dijalankan:

| Role | Email | Password | Hak Akses Utama |
| :--- | :--- | :--- | :--- |
| **Admin** | `admin@helpdesk.com` | `$$admin$$` | Dashboard statistik, kelola user & role, semua aksi agent |
| **Agent** | `agent@helpdesk.com` | `$$agent$$` | Assign tiket, ubah status/prioritas, komentar internal notes |
| **Requester** | `user@helpdesk.com` | `$$user$$` | Buat tiket, pantau tiket sendiri, tutup tiket sendiri |

---

## Cara Menjalankan Aplikasi

### Opsi 1: Menggunakan Docker Compose (Direkomendasikan)
Cukup jalankan satu perintah di root folder repositori:

```bash
docker compose up --build
```

Setelah seluruh container menyala:
- **Frontend**: `http://localhost:3000` (Port 3000)
- **Backend API**: `http://localhost:8080`
- **Health Check**: `http://localhost:8080/api/health`

---

### Opsi 2: Menjalankan Secara Manual (Lokal)

#### 1. Persiapan Database
Pastikan PostgreSQL berjalan di lokal dan buat database:
```sql
CREATE DATABASE helpdesk_db;
```

#### 2. Backend (Go)
```bash
cd backend
# Salin .env.example menjadi .env lalu sesuaikan DB_PASSWORD
go mod download
go run .
```
*Backend berjalan di `http://localhost:8080`.*

#### 3. Frontend (Vue 3)
```bash
cd frontend
# Salin .env.example menjadi .env
npm install
npm run dev
```
*Frontend berjalan di `http://localhost:5173`.*

---

## Menjalankan Pengujian (Testing)

Pengujian backend mencakup integration test untuk Autentikasi, CRUD Tiket, Pembatasan Akses (RBAC), Transisi Status, Komentar, Catatan Internal, Admin Role Management, dan Validasi SLA:

1. Buat database pengujian di PostgreSQL:
   ```sql
   CREATE DATABASE helpdesk_test;
   ```
2. Jalankan perintah test di folder backend:
   ```bash
   cd backend
   go test -v ./tests/...
   ```

---

## Daftar Endpoint API

Semua endpoint berawalan prefix `/api`.

| Method | Endpoint | Keterangan | Hak Akses |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/health` | Cek status backend & database | Publik |
| `POST` | `/api/auth/register` | Pendaftaran akun baru (requester) | Publik |
| `POST` | `/api/auth/login` | Login akun & penerbitan token JWT | Publik |
| `GET` | `/api/me` | Profil pengguna saat ini | Authenticated |
| `GET` | `/api/tickets` | Daftar tiket (filter, search, pagination) | Authenticated (Requester hanya tiket miliknya) |
| `POST` | `/api/tickets` | Buat tiket baru | Authenticated |
| `GET` | `/api/tickets/:id` | Detail tiket beserta komentar | Pemilik, Agent, Admin |
| `PATCH`| `/api/tickets/:id` | Ubah tiket (status, prioritas, assignee) | Sesuai aturan role & workflow |
| `POST` | `/api/tickets/:id/comments` | Tambah komentar / internal notes | Pemilik, Agent, Admin |
| `GET` | `/api/tickets/:id/comments` | Riwayat komentar tiket | Pemilik, Agent, Admin |
| `GET` | `/api/users` | Daftar seluruh pengguna | Khusus Admin |
| `PATCH`| `/api/users/:id/role` | Ubah role pengguna | Khusus Admin |
| `GET` | `/api/users/stats` | Statistik operasional & SLA | Khusus Admin |

---

## Dokumentasi API (Postman Collection)

File Postman Collection v2.1 tersedia di root repositori:  
`helpdesk-api.postman_collection.json`

**Cara Pakai:**
1. Buka Postman -> Klik **Import** -> Pilih `helpdesk-api.postman_collection.json`.
2. Lakukan request `Login (Admin)` atau `Login (User)` -> Token JWT otomatis tersimpan di variabel environment collection untuk request berikutnya.

---

## Catatan Arsitektur & Keputusan Desain

*(Refleksi Teknis — Maksimal 1 Halaman)*

### 1. Keputusan Desain (Architectural Decisions)
- **Layered Clean Architecture pada Go Backend**: Kode backend dipisahkan secara ketat menjadi 4 lapisan: `handlers` (HTTP transport & response shaping), `services` (business logic, status transition rules, validation), `repositories` (abstraksi query GORM), dan `models` (entitas database). Pemisahan ini mempermudah penulisan integration test yang independen.
- **Workflow State Machine di Service Layer**: Validasi transisi status tiket sengaja ditempatkan murni di Go service layer (`statusTransitions` matrix map). Hal ini menjamin keamanan bisnis, sehingga request ilegal via cURL/Postman tetap ditolak dengan kode `403 Forbidden`.
- **Preload & Data Privacy Filtering**: Komentar internal difilter langsung di query level database (`is_internal = false` jika user adalah requester) demi privasi data dan efisiensi memori.
- **Vue 3 Composition API & Centralized State**: Menggunakan Pinia untuk sinkronisasi state autentikasi dengan `localStorage`. Axios response interceptor menangani status `401 Unauthorized` secara terpusat untuk membersihkan session dan redirect ke login.

### 2. Hal yang Sengaja Ditinggalkan (Trade-offs)
- **Soft Deletes**: Tidak mengaktifkan `gorm.DeletedAt` agar skema database tetap sederhana dan query performan tanpa overhead pengecekan timestamp deleted.
- **Monorepo Structure**: Menggunakan struktur monorepo satu repositori git dengan `docker-compose.yml` di root agar mempermudah reviewer mencoba aplikasi hanya dengan satu perintah.

### 3. Rencana Perbaikan (Future Improvements)
- **Real-time Notification**: Menambahkan WebSocket atau SSE agar pembaruan tiket dan komentar muncul seketika tanpa perlu me-refresh halaman.
- **File Attachment**: Menambahkan fitur upload file lampiran dokumen/gambar pendukung dengan validasi MIME-type dan penyimpanan objek (seperti MinIO/S3).
- **Automated Assignment**: Penugasan tiket otomatis menggunakan algoritma round-robin ke agent yang sedang aktif/tersedia.