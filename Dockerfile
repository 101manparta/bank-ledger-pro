# STAGE 1: Membangun (Build) aplikasi
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
# Compile kode Go menjadi file executable bernama 'main'
RUN go build -o main .

# STAGE 2: Menjalankan (Run) aplikasi
FROM alpine:latest
WORKDIR /root/
# Hanya ambil file hasil compile dari stage 1 (biar ringan)
COPY --from=builder /app/main .
# Jalankan aplikasi
CMD ["./main"]