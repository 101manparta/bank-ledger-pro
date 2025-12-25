package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Driver Postgres
)

// Struktur data
type TransferRequest struct {
	FromID         int
	ToID           int
	Amount         float64
	IdempotencyKey string
}

func main() {
	// 1. Koneksi Database
	db, err := sql.Open("postgres", "user=postgres password=secret dbname=bank_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Simulasi Request dari User
	req := TransferRequest{
		FromID:         1,
		ToID:           2,
		Amount:         500.00,
		IdempotencyKey: "unique-uuid-dari-frontend-123",
	}

	// Jalankan Transfer
	err = ExecTransfer(db, req)
	if err != nil {
		fmt.Printf("Gagal Transfer: %v\n", err)
	} else {
		fmt.Println("Transfer Berhasil & Aman!")
	}
}

func ExecTransfer(db *sql.DB, req TransferRequest) error {
	ctx := context.Background()

	// 2. Mulai Database Transaction (ACID dimulai di sini)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Fungsi helper untuk Rollback otomatis jika ada error
	defer tx.Rollback()

	// 3. CEK IDEMPOTENCY
	// Masalah: User klik 2x. Solusi: Cek apakah key ini sudah pernah sukses?
	var exists int
	err = tx.QueryRowContext(ctx, "SELECT count(1) FROM transfers WHERE idempotency_key = $1", req.IdempotencyKey).Scan(&exists)
	if exists > 0 {
		return errors.New("transaksi ini sudah pernah diproses (idempotent)")
	}

	// 4. KUNCI & AMBIL SALDO PENGIRIM (SELECT FOR UPDATE)
	// Kita kunci baris ini agar tidak bisa diubah transaksi lain sampai tx ini selesai.
	var fromBalance float64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", req.FromID).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("akun pengirim tidak ditemukan: %v", err)
	}

	if fromBalance < req.Amount {
		return errors.New("saldo tidak cukup")
	}

	// 5. EKSEKUSI PERPINDAHAN UANG
	// Kurangi saldo pengirim
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", req.Amount, req.FromID)
	if err != nil {
		return err
	}

	// Tambah saldo penerima
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", req.Amount, req.ToID)
	if err != nil {
		return err
	}

	// 6. CATAT RIWAYAT & LOCK IDEMPOTENCY
	_, err = tx.ExecContext(ctx, "INSERT INTO transfers (from_account_id, to_account_id, amount, idempotency_key) VALUES ($1, $2, $3, $4)",
		req.FromID, req.ToID, req.Amount, req.IdempotencyKey)
	if err != nil {
		return err
	}

	// 7. COMMIT (Uang resmi berpindah)
	return tx.Commit()
}