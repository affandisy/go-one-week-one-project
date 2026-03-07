package main

import (
	"fmt"
	"log"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/config"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/brianvoe/gofakeit/v6"
)

func SeederDatabase() {
	// 1. Inisialisasi Database
	config.ConnectDatabase()
	db := config.DB

	// 2. Konfigurasi Seeder
	totalData := 100000
	batchSize := 5000 // GORM akan memasukkan 5000 data sekaligus dalam 1 query INSERT
	var products []models.Product

	fmt.Printf("🚀 Memulai proses seeding %d data produk...\n", totalData)
	startTime := time.Now()

	// Inisialisasi seed untuk gofakeit agar datanya selalu acak
	gofakeit.Seed(0)

	// 3. Proses Looping & Chunking
	for i := 1; i <= totalData; i++ {
		// Buat data produk palsu yang realistis
		products = append(products, models.Product{
			SKU:          gofakeit.Regex("^[A-Z]{3}-[0-9]{4}-[A-Z]{2}$"), // Contoh: JKT-1234-AB
			Name:         gofakeit.ProductName(),
			Description:  gofakeit.Sentence(10),
			CurrentStock: gofakeit.Number(0, 1000), // Angka acak 0 - 1000
		})

		// 4. Eksekusi Batch Insert saat keranjang (slice) sudah penuh
		if len(products) == batchSize {
			// CreateInBatches adalah fitur GORM untuk bulk insert yang sangat cepat
			if err := db.CreateInBatches(&products, batchSize).Error; err != nil {
				log.Fatal("❌ Gagal melakukan batch insert:", err)
			}

			fmt.Printf("✅ Berhasil menyimpan %d baris data...\n", i)

			// Kosongkan keranjang (slice) untuk batch berikutnya agar RAM tidak penuh
			products = nil
		}
	}

	// 5. Simpan sisa data jika ada (misal total data tidak habis dibagi batchSize)
	if len(products) > 0 {
		if err := db.CreateInBatches(&products, len(products)).Error; err != nil {
			log.Fatal("❌ Gagal melakukan sisa batch insert:", err)
		}
	}

	fmt.Printf("🎉 Seeding selesai! Waktu eksekusi: %v\n", time.Since(startTime))
}
