<div align="center">

# SneakHub

### Marketplace Sneaker Multi-Vendor

API backend untuk platform jual-beli sepatu sneaker original dan bekas,
dibangun dengan **Go**, **Fiber v2**, **GORM**, dan **MySQL**.

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Fiber](https://img.shields.io/badge/Fiber-v2-000000?style=flat&logo=go)](https://gofiber.io/)
[![GORM](https://img.shields.io/badge/GORM-v1.31-00ADD8?style=flat)](https://gorm.io/)
[![MySQL](https://img.shields.io/badge/MySQL-8.x-4479A1?style=flat&logo=mysql)](https://www.mysql.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](#)

</div>

---

## Daftar Isi

- [Tentang](#tentang)
- [Fitur Utama](#fitur-utama)
- [Tech Stack](#tech-stack)
- [Persiapan](#persiapan)
- [Instalasi](#instalasi)
- [Environment Variables](#environment-variables)
- [Struktur Proyek](#struktur-proyek)
- [API Overview](#api-overview)
- [Role & Akses](#role--akses)
- [Alur Transaksi](#alur-transaksi)
- [Deployment](#deployment)
- [Postman Collection](#postman-collection)
- [License](#license)

---

## Tentang

**SneakHub** adalah platform marketplace sneaker multi-vendor yang memungkinkan
penjual dan pembeli melakukan transaksi sepatu sneaker original maupun bekas.
Platform ini mendukung fitur AI seperti smart filtering, rekomendasi personal,
pencarian berbasis gambar, dan prediksi harga pasar.

API ini dirancang untuk pasar Indonesia dengan semua field dan pesan dalam
Bahasa Indonesia.

---

## Fitur Utama

| Fitur | Deskripsi |
|---|---|
| **Smart Filter** | Pencarian produk cerdas dengan bobot prioritas (harga, kondisi, kepercayaan seller) |
| **Personalized Home** | Konten homepage yang dipersonalisasi berdasarkan preferensi ukuran dan brand favorit |
| **Search by Image** | Pencarian produk berdasarkan gambar menggunakan image embedding |
| **Price Prediction** | Prediksi harga pasar untuk seller berdasarkan brand, kondisi, dan usia produk |
| **Price & Restock Alert** | Notifikasi otomatis saat harga turun atau stok tersedia kembali |
| **Multi-Payment** | Dukungan Midtrans, Tripay, dan mode mock untuk testing |
| **Auto Shipping** | Integrasi Biteship untuk cek ongkir otomatis dan booking kurir |
| **Condition Score** | Sistem penilaian kondisi sepatu detail (upper, outsole, midsole, insole, dll) |
| **Seller Trust Score** | Skor kepercayaan seller berdasarkan penyelesaian order, rating, dan tingkat pembatalan |
| **Trending & Best Seller** | Daftar produk trending dan best seller mingguan |

---

## Tech Stack

| Komponen | Teknologi |
|---|---|
| **Bahasa** | Go 1.26 |
| **Web Framework** | [Fiber v2](https://gofiber.io/) |
| **ORM** | [GORM v1.31](https://gorm.io/) |
| **Database** | MySQL 8.x |
| **Autentikasi** | JWT (JSON Web Token) |
| **Password Hashing** | bcrypt |
| **Payment Gateway** | Midtrans / Tripay |
| **Shipping** | [Biteship](https://biteship.com/) |
| **Image Processing** | DCT-based perceptual hash (image embedding) |

---

## Persiapan

Sebelum menjalankan proyek ini, pastikan sudah terinstall:

- **Go** 1.26 atau lebih baru ([download](https://go.dev/dl/))
- **MySQL** 8.x atau MariaDB ([download](https://dev.mysql.com/downloads/mysql/))
- **Git** ([download](https://git-scm.com/))

Untuk fitur pembayaran, siapkan akun di salah satu payment gateway:
- **Midtrans**: [dashboard.midtrans.com](https://dashboard.midtrans.com/)
- **Tripay**: [tripay.co.id](https://tripay.co.id/)

Untuk fitur pengiriman otomatis:
- **Biteship**: [biteship.com](https://biteship.com/)

---

## Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/ehanz12/api-SneakHub.git
cd api-SneakHub
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Environment Variables

```bash
cp .env.example .env
```

Edit file `.env` sesuai konfigurasi kamu (lihat [Environment Variables](#environment-variables)).

### 4. Buat Database MySQL

```sql
CREATE DATABASE sneakhub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. Jalankan Server

```bash
go run main.go
```

Server akan berjalan di port yang sudah dikonfigurasi (default: `3000`).
Database akan otomatis di-migrasi dan di-seed saat pertama kali dijalankan.

---

## Environment Variables

| Variabel | Deskripsi | Default | Wajib |
|---|---|---|---|
| `DB_USER` | Username MySQL | `root` | Ya |
| `DB_PASSWORD` | Password MySQL | - | Ya |
| `DB_HOST` | Host MySQL | `127.0.0.1` | Ya |
| `DB_NAME` | Nama database | `sneakhub` | Ya |
| `DB_PORT` | Port MySQL | `3306` | Ya |
| `PORT` | Port server API | `3000` | Ya |
| `PUBLIC_URL` | URL publik server (untuk callback) | - | Ya (produksi) |
| `CORS_ORIGINS` | Origin frontend yang diizinkan | `http://localhost:3060` | Tidak |
| `JWT_SECRET` | Secret key untuk JWT | - | Ya |
| `PAYMENT_MODE` | Mode pembayaran: `mock`, `tripay`, `midtrans` | `midtrans` | Ya |
| `TRIPAY_API_KEY` | API Key Tripay | - | Opsional |
| `TRIPAY_PRIVATE_KEY` | Private Key Tripay | - | Opsional |
| `TRIPAY_MERCHANT_CODE` | Merchant Code Tripay | - | Opsional |
| `TRIPAY_IS_SANDBOX` | Mode sandbox Tripay | `true` | Opsional |
| `MIDTRANS_SERVER_KEY` | Server Key Midtrans | - | Opsional |
| `MIDTRANS_CLIENT_KEY` | Client Key Midtrans | - | Opsional |
| `MIDTRANS_IS_SANDBOX` | Mode sandbox Midtrans | `true` | Opsional |
| `BITESHIP_API_KEY` | API Key Biteship (ongkir otomatis) | - | Opsional |
| `SHIPPING_POLL_MINUTES` | Interval polling tracking (menit) | `30` | Tidak |
| `AUTO_COMPLETE_DAYS` | Hari auto-complete pesanan | `7` | Tidak |

> **Catatan:**
> - Jika `BITESHIP_API_KEY` kosong atau data toko belum lengkap, ongkir menggunakan flat rate **Rp15.000**.
> - Untuk `PAYMENT_MODE=mock`, gunakan endpoint `/api/mock-pay` dan `/api/payments/mock-settle` untuk testing tanpa gateway asli.

---

## Struktur Proyek

```
api-SneakHub/
|
|-- main.go                          # Entry point aplikasi
|-- go.mod                           # Module dependencies
|-- .env.example                     # Template environment variables
|-- API_DOCUMENTATION.docs           # Dokumentasi API lengkap
|-- SneakHub.postman_collection.json # Postman collection
|
|-- config/                          # Konfigurasi
|   `-- config.go                    # Struct config + LoadEnv()
|
|-- database/                        # Koneksi & seeding database
|   |-- database.go                  # Koneksi MySQL via GORM
|   `-- seeders/
|       |-- brand.go                 # Seed 15 brand sneaker
|       `-- category.go              # Seed 8 kategori produk
|
|-- models/                          # Model data GORM (25 model)
|   |-- user.go                      # User (customer/seller/admin)
|   |-- seller.go                    # Profil toko
|   |-- product.go                   # Produk
|   |-- order.go                     # Pesanan
|   |-- payment.go                   # Transaksi pembayaran
|   `-- ...                          # Model lainnya
|
|-- handlers/                        # HTTP request handlers (18 file)
|   |-- user.go                      # Register, login, profil
|   |-- seller.go                    # Manajemen toko
|   |-- product.go                   # CRUD produk
|   |-- checkout.go                  # Alur checkout
|   |-- discovery.go                 # Smart filter, rekomendasi
|   `-- ...                          # Handler lainnya
|
|-- services/                        # Business logic (24 file)
|   |-- discovery.go                 # Engine AI/rekomendasi
|   |-- evaluation.go                # Penilaian kondisi
|   |-- price.go                     # Prediksi harga
|   |-- shipping.go                  # Kalkulasi ongkir
|   |-- biteship.go                  # Integrasi Biteship API
|   |-- payment_gateway.go           # Payment provider (Tripay/Mock)
|   |-- payment_midtrans.go          # Midtrans Snap
|   `-- ...                          # Service lainnya
|
|-- middleware/
|   `-- jwt.go                       # Middleware JWT + role-based access
|
|-- requests/                        # Request DTOs/validation (15 file)
|-- responses/                       # Response DTOs (15 file)
|-- mappers/                         # Model-to-response mappers (14 file)
|-- routes/                          # Definisi rute (13 file)
|   `-- route.go                     # Master route + rate limiter
|
|-- utils/                           # Utility functions
|   |-- jwt.go                       # JWT generation & verification
|   |-- validation.go                # Validasi ukuran sepatu (35-48)
|   |-- number.go                    # Validasi nomor HP
|   `-- embedding.go                 # Image embedding (DCT-based)
|
`-- uploads/                         # Upload gambar produk (gitignored)
```

---

## API Overview

Base URL: `http://localhost:3000/api`

> Dokumentasi API lengkap ada di [`API_DOCUMENTATION.docs`](API_DOCUMENTATION.docs).

### Autentikasi

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| POST | `/api/auth/register` | Public | Registrasi customer baru |
| POST | `/api/auth/login` | Public | Login semua role |

### User

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| GET | `/api/users/me` | Semua | Lihat profil sendiri |
| PUT | `/api/users/me` | Semua | Update profil sendiri |

### Seller (Toko)

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| POST | `/api/users/me/seller-activation` | Semua | Aktivasi akun sebagai seller |
| GET | `/api/seller/me` | Seller | Profil toko sendiri |
| PUT | `/api/seller/me` | Seller | Update profil toko |
| GET | `/api/seller/products` | Seller | Produk milik toko |
| GET | `/api/seller/orders` | Seller | Pesanan masuk |
| POST | `/api/seller/orders/:order_id/ship` | Seller | Kirim pesanan (resi otomatis) |
| GET | `/api/seller/dashboard` | Seller | Statistik toko |

### Admin

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| GET | `/api/admin/users` | Admin | Daftar semua user |
| PATCH | `/api/admin/users/:user_id/status` | Admin | Ubah status akun |
| PATCH | `/api/admin/users/:user_id/role` | Admin | Ubah role user |
| GET | `/api/admin/sellers` | Admin | Daftar registrasi seller |
| PATCH | `/api/admin/sellers/:seller_id/verification` | Admin | Verifikasi/tolak seller |
| GET | `/api/admin/products` | Admin | Daftar semua produk |
| PATCH | `/api/admin/products/:product_id/status` | Admin | Ubah status publikasi |
| GET | `/api/admin/orders` | Admin | Daftar semua pesanan |
| GET | `/api/admin/reports` | Admin | Laporan platform |

### Produk

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| GET | `/api/products` | Semua | List & cari produk |
| GET | `/api/products/:product_id` | Semua | Detail produk |
| POST | `/api/products` | Seller/Admin | Buat produk |
| PUT | `/api/products/:product_id` | Seller/Admin | Update produk |
| DELETE | `/api/products/:product_id` | Seller/Admin | Hapus produk |
| POST | `/api/products/:product_id/images` | Seller/Admin | Upload gambar |
| GET | `/api/products/:product_id/images` | Semua | List gambar |
| DELETE | `/api/products/:product_id/images/:image_id` | Seller/Admin | Hapus gambar |
| POST | `/api/products/search-by-image` | Semua | Cari produk via gambar (AI) |
| GET | `/api/products/:product_id/price-insight` | Semua | Insight harga vs pasar |
| POST | `/api/products/:product_id/condition-score` | Seller/Admin | Set skor kondisi |
| GET | `/api/products/:product_id/condition-score` | Semua | Lihat skor kondisi |

### Discovery & AI

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| POST | `/api/discovery/smart-filter` | Semua | Smart filter produk |
| GET | `/api/home/personalized` | Semua | Homepage personalisasi |
| GET | `/api/recommendation/cocok-untuk-kamu` | Semua | Rekomendasi personal |
| GET | `/api/trending` | Semua | Produk trending |
| GET | `/api/best-seller/weekly` | Semua | Best seller mingguan |
| POST | `/api/price/prediction` | Seller | Prediksi harga |

### Cart

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| GET | `/api/carts` | Customer | Lihat keranjang |
| POST | `/api/carts/items` | Customer | Tambah item |
| PUT | `/api/carts/items/:cart_item_id` | Customer | Update jumlah |
| DELETE | `/api/carts/items/:cart_item_id` | Customer | Hapus item |

### Address

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| GET | `/api/addresses` | Customer | Daftar alamat |
| GET | `/api/addresses/:address_id` | Customer | Detail alamat |
| POST | `/api/addresses` | Customer | Buat alamat |
| PUT | `/api/addresses/:address_id` | Customer | Update alamat |
| DELETE | `/api/addresses/:address_id` | Customer | Hapus alamat |

### Wishlist

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| GET | `/api/wishlist` | Customer | Daftar wishlist |
| POST | `/api/wishlist` | Customer | Tambah ke wishlist |
| DELETE | `/api/wishlist/:product_id` | Customer | Hapus dari wishlist |
| POST | `/api/wishlist/:product_id/price-alert` | Customer | Aktifkan price alert |
| DELETE | `/api/wishlist/:product_id/price-alert` | Customer | Matikan price alert |
| POST | `/api/wishlist/:product_id/restock-alert` | Customer | Aktifkan restock alert |
| DELETE | `/api/wishlist/:product_id/restock-alert` | Customer | Matikan restock alert |

### Checkout & Payment

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| POST | `/api/checkout` | Customer | Checkout keranjang |
| POST | `/api/checkout/shipping-rates` | Customer | Hitung ongkir per toko |
| POST | `/api/payments/notification` | Public | Callback Tripay |
| GET | `/api/mock-pay/:order_id` | Public | Halaman mock payment |
| POST | `/api/payments/mock-settle` | Public | Settle mock payment |
| POST | `/api/payments/midtrans/notification` | Public | Callback Midtrans |

### Orders

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| GET | `/api/orders` | Semua | Daftar pesanan (role-aware) |
| GET | `/api/orders/:order_id` | Semua | Detail pesanan |
| POST | `/api/orders` | Semua | Buat pesanan manual |
| PUT | `/api/orders/:order_id` | Semua | Update status pesanan |
| DELETE | `/api/orders/:order_id` | Semua | Hapus pesanan |
| POST | `/api/orders/:order_id/review` | Customer | Review produk |
| POST | `/api/orders/:order_id/confirm` | Customer | Konfirmasi pesanan selesai |

### Notification

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| GET | `/api/notifications` | Semua | Daftar notifikasi |
| PATCH | `/api/notifications/:notification_id/read` | Semua | Tandai sudah dibaca |

### Category & Brand

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| GET | `/api/category` | Semua | Daftar kategori |
| POST | `/api/category` | Admin | Buat kategori |
| PUT | `/api/category/:category_id` | Admin | Update kategori |
| DELETE | `/api/category/:category_id` | Admin | Hapus kategori |
| GET | `/api/brand` | Semua | Daftar brand |
| POST | `/api/brand` | Admin | Buat brand |
| PUT | `/api/brand/:brand_id` | Admin | Update brand |
| DELETE | `/api/brand/:brand_id` | Admin | Hapus brand |

---

## Role & Akses

| Role | Akses |
|---|---|
| **Customer** | Belanja: wishlist, cart, alamat, checkout, orders, reviews, notifikasi |
| **Seller** | Kelola toko: produk, orders, dashboard, prediksi harga |
| **Admin** | Kelola user/produk/order, moderasi, laporan, kategori & brand |

### Rute yang Diizinkan per Role

| Route Group | Role yang Diizinkan |
|---|---|
| `auth/*`, `payments/*`, `mock-pay/*` | **Public** |
| `users/me`, `discovery/*`, `home/*`, `recommendation/*`, `trending`, `best-seller/*`, `products` (GET), `category` (GET), `notifications/*`, `orders` (GET/POST/PUT/DELETE) | **Semua role** |
| `wishlist/*`, `carts/*`, `addresses/*`, `checkout`, `orders/:id/review` | **Customer** |
| `seller/*` (dashboard, products, orders), `price/prediction` | **Seller** |
| `admin/*`, `category` (POST/PUT/DELETE) | **Admin** |
| `products` (POST/PUT/DELETE), `products/:id/images` (POST/DELETE), `products/:id/condition-score` (POST) | **Seller & Admin** |

---

## Alur Transaksi

```
Register/Login
    |
    v
Aktivasi Seller ----> Buat Produk ----> Upload Gambar
    |                                       |
    v                                       v
Customer: Tambah Cart <----- Browse Produk
    |
    v
Buat Alamat
    |
    v
Checkout ----> Pilih Kurir ----> Pilih Pembayaran
    |
    v
Bayar (Midtrans/Tripay/Mock)
    |
    v
Seller: Kirim Pesanan (resi otomatis Biteship)
    |
    v
Customer: Konfirmasi Terima (atau auto-complete 7 hari)
    |
    v
Review Produk
```

### Status Pesanan

```
pending --> diproses --> dikirim --> selesai
   |           |
   +-----> dibatalkan (stok dikembalikan)
```

---

## Deployment

### Environment Produksi

1. Buat file `.env.production` dengan konfigurasi produksi
2. Set `PAYMENT_MODE=midtrans` atau `tripay` (bukan `mock`)
3. Set `MIDTRANS_IS_SANDBOX=false` atau `TRIPAY_IS_SANDBOX=false`
4. Isi `PUBLIC_URL` dengan URL server produksi (contoh: `https://api-sneakhub.up.railway.app`)
5. Isi `CORS_ORIGINS` dengan URL frontend produksi
6. Gunakan `JWT_SECRET` yang kuat dan unik

### Railway / VPS

```bash
# Build
go build -o sneakhub main.go

# Jalankan
./sneakhub
```

> **Catatan:** Folder `uploads/` perlu persistent storage untuk menyimpan gambar produk.

---

## Postman Collection

File collection Postman sudah tersedia di repository:

- `SneakHub.postman_collection.json` - Versi 1
- `SneakHub.postman_collection_v2.json` - Versi 2 (recommended)

Import ke Postman untuk testing semua endpoint.

---

## License

MIT License

Copyright (c) 2026 ehanz12

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
