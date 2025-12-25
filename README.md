Bank Ledger Engine (Golang)
Pendahuluan
Proyek ini adalah implementasi sistem Ledger (pencatatan transaksi) keuangan sederhana namun memiliki standar keamanan industri perbankan. Fokus utama proyek ini bukan pada fitur yang banyak, melainkan pada Integritas Data dan Resistensi terhadap kegagalan sistem.

Sebagai Backend Engineer, saya membangun sistem ini untuk menyelesaikan masalah klasik di sistem finansial: Double Spending, Race Condition, dan Ketidakkonsistenan data saat terjadi gangguan jaringan.

Masalah Nyata & Solusi Engineering
Dalam membangun sistem ini, saya menerapkan beberapa logika krusial untuk menangani skenario dunia nyata:

1. Penanganan Race Condition (Concurrency)
Masalah: Bagaimana jika seorang nasabah melakukan dua transfer di milidetik yang sama sedangkan saldonya hanya cukup untuk satu transaksi? Solusi: Saya menerapkan Row-Level Locking menggunakan perintah SELECT FOR UPDATE pada PostgreSQL. Hal ini memaksa database untuk mengantrekan transaksi yang mencoba mengakses baris (akun) yang sama, sehingga saldo tidak akan pernah menjadi negatif atau salah hitung.

2. Idempotency (Mencegah Transaksi Ganda)
Masalah: User seringkali menekan tombol "Bayar" berkali-kali karena koneksi internet yang lambat. Solusi: Saya menambahkan Idempotency Key pada setiap transaksi. Sistem akan mengecek apakah kunci unik tersebut sudah pernah sukses diproses sebelumnya. Jika ya, sistem akan menolak pemrosesan ulang tanpa memotong saldo nasabah untuk kedua kalinya.

3. Integritas Data (Atomic Transactions)
Masalah: Kegagalan sistem di tengah proses (Misal: Saldo pengirim sudah berkurang, tapi server mati sebelum menambah saldo penerima). Solusi: Seluruh proses transfer dibungkus dalam Database Transaction (ACID). Jika salah satu langkah gagal, sistem akan melakukan Rollback otomatis ke kondisi semula, sehingga tidak ada uang yang "lenyap di tengah jalan".

Teknologi yang Digunakan
Golang: Dipilih karena performa konkurensi yang tinggi dan sistem pengetikan data yang ketat (strong typing) untuk meminimalisir bug saat runtime.

PostgreSQL: Digunakan karena kepatuhannya terhadap prinsip ACID yang sangat vital bagi data finansial.

Docker & Docker Compose: Untuk standarisasi lingkungan pengembangan, memastikan sistem berjalan sama persis di laptop developer maupun di server produksi.

Unit Testing: Menggunakan package testing bawaan Go untuk memvalidasi logika bisnis secara otomatis sebelum dideploy.

Struktur Project
main.go: Inti dari logika bisnis dan koneksi database.

main_test.go: Skenario pengujian otomatis untuk memastikan keamanan sistem.

docker-compose.yaml: Orkestrasi infrastruktur (Database Postgres).

Dockerfile: Containerization aplikasi untuk kemudahan deployment.

Cara Menjalankan
Pastikan Docker Desktop sudah berjalan.

Jalankan infrastruktur:

Bash

docker-compose up -d
Jalankan pengujian otomatis:

Bash

go test -v
Analisis & Pengembangan Selanjutnya
Sistem ini dirancang untuk dapat dikembangkan ke arah Microservices. Langkah selanjutnya yang direncanakan adalah:

Menambahkan lapisan REST API menggunakan Gin atau Echo.

Implementasi Centralized Logging untuk audit trail yang lebih mendalam.

Menambahkan Circuit Breaker untuk menangani kegagalan koneksi database secara lebih elegan.

Tips buat kamu (The Human Touch):
Saat nanti ada yang tanya di interview tentang project ini, kamu bisa bilang gini:

"Saya sengaja pakai NUMERIC di Postgres, bukan FLOAT, karena saya sadar presisi satu perak pun sangat berharga di perbankan. Saya juga lebih memilih SELECT FOR UPDATE daripada optimasi di level aplikasi karena saya ingin database yang menjadi 'Single Source of Truth' untuk validasi saldo."
