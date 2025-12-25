package main

import (
	"database/sql"
	"testing"
	_ "github.com/lib/pq"
)

// Fungsi ini mengetes skenario transfer sukses
func TestExecTransfer_Success(t *testing.T) {
	// Setup koneksi database khusus testing
	db, _ := sql.Open("postgres", "user=postgres password=secret dbname=bank_db sslmode=disable")
	defer db.Close()

	req := TransferRequest{
		FromID:         1,
		ToID:           2,
		Amount:         100.00,
		IdempotencyKey: "test-key-success",
	}

	err := ExecTransfer(db, req)
	if err != nil {
		t.Errorf("Harusnya sukses, tapi dapet error: %v", err)
	}
}

// Fungsi ini mengetes jika saldo tidak cukup
func TestExecTransfer_InsufficientBalance(t *testing.T) {
	db, _ := sql.Open("postgres", "user=postgres password=secret dbname=bank_db sslmode=disable")
	defer db.Close()

	req := TransferRequest{
		FromID:         1,
		ToID:           2,
		Amount:         999999999.00, // Jumlah yang gak masuk akal
		IdempotencyKey: "test-key-fail",
	}

	err := ExecTransfer(db, req)
	if err == nil || err.Error() != "saldo tidak cukup" {
		t.Errorf("Harusnya error 'saldo tidak cukup', tapi malah: %v", err)
	}
}