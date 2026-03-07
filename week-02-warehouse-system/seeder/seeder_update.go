package main

import (
	"fmt"
	"log"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/config"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	config.ConnectDatabase()
	db := config.DB

	fmt.Println("Memulai proses Seeding WMS...")
	startTime := time.Now()
	gofakeit.Seed(0)

	fmt.Println("Seeding Users")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	passString := string(hashedPassword)

	users := []models.User{
		{Name: "Super Admin", Email: "admin@wms.com", Password: passString, Role: "admin"},
		{Name: "Budi Manager", Email: "manager@wms.com", Password: passString, Role: "manager"},
		{Name: "Siti Operator", Email: "operator1@wms.com", Password: passString, Role: "operator"},
		{Name: "Agus Operator", Email: "operator2@wms.com", Password: passString, Role: "operator"},
		{Name: "Dewi Operator", Email: "operator3@wms.com", Password: passString, Role: "operator"},
	}
	db.Create(&users)

	fmt.Println("Seeding 100 Warehouse")
	var warehouses []models.Warehouse
	for i := 1; i <= 100; i++ {
		warehouses = append(warehouses, models.Warehouse{
			Code:    fmt.Sprintf("WH-%03d", i),
			Name:    fmt.Sprintf("Gudang %s %s", gofakeit.City(), gofakeit.CompanySuffix()),
			Address: gofakeit.Address().Address,
		})
	}
	db.CreateInBatches(&warehouses, 50)

	fmt.Println("Seeding 100 Locations")
	var locations []models.Location
	for i := 1; i <= 100; i++ {
		code := fmt.Sprintf("%s-%s-%02d", gofakeit.RandomString([]string{"RAK", "BLOK", "ZONE"}), gofakeit.LetterN(1), gofakeit.Number(1, 50))
		locations = append(locations, models.Location{
			Code:        code,
			Name:        "Area Penyimpanan " + code,
			Description: "Kapasitas " + gofakeit.RandomString([]string{"Besar", "Sedang", "Kecil"}),
		})
	}
	db.CreateInBatches(&locations, 50)

	fmt.Println("⏳ Seeding 10 Partners...")
	var partners []models.Partner
	for i := 1; i <= 10; i++ {
		pType := "SUPPLIER"
		if i > 5 {
			pType = "CUSTOMER" // 5 Supplier, 5 Customer
		}
		partners = append(partners, models.Partner{
			Name:    gofakeit.Company(),
			Type:    pType,
			Phone:   gofakeit.Phone(),
			Address: gofakeit.Address().Address,
		})
	}
	db.Create(&partners)

	fmt.Println("⏳ Seeding 1,000 Products...")
	var products []models.Product
	for i := 1; i <= 1000; i++ {
		locID := uint(gofakeit.Number(1, 100))
		products = append(products, models.Product{
			SKU:          fmt.Sprintf("SKU-%s-%04d", gofakeit.LetterN(3), i),
			Name:         gofakeit.ProductName(),
			Description:  gofakeit.ProductDescription(),
			Category:     gofakeit.ProductCategory(),
			Unit:         gofakeit.RandomString([]string{"pcs", "box", "kg"}),
			Price:        gofakeit.Price(10000, 5000000),
			MinStock:     gofakeit.Number(10, 50),
			MaxStock:     gofakeit.Number(500, 1000),
			CurrentStock: gofakeit.Number(50, 800), // Stok global
			LocationID:   &locID,
			IsActive:     true,
		})
	}
	db.CreateInBatches(&products, 250)

	fmt.Println("Seeding Warehouse Stocks, Batches & Adjustments...")
	var whStocks []models.WarehouseStock
	var batches []models.ProductBatch
	var adjustments []models.StockAdjustment

	for i := 1; i <= 300; i++ {
		prodID := uint(gofakeit.Number(1, 1000))
		whID := uint(gofakeit.Number(1, 1000))

		whStocks = append(whStocks, models.WarehouseStock{
			WarehouseID: whID,
			ProductID:   prodID,
			Stock:       gofakeit.Number(10, 200),
		})

		batches = append(batches, models.ProductBatch{
			ProductID:  prodID,
			BatchNo:    fmt.Sprintf("BATCH-%s-%03d", gofakeit.LetterN(2), i),
			ExpiryDate: time.Now().AddDate(gofakeit.Number(0, 2), gofakeit.Number(1, 11), 0),
			Stock:      gofakeit.Number(10, 100),
		})

		adjustments = append(adjustments, models.StockAdjustment{
			ProductID: prodID,
			UserID:    uint(gofakeit.Number(1, 5)),
			Qty:       gofakeit.RandomInt([]int{-5, -2, -1, 3, 5, 10}),
			Reason:    gofakeit.Sentence(5),
			Status:    "approved",
		})
	}
	db.CreateInBatches(&whStocks, 100)
	db.CreateInBatches(&batches, 100)
	db.CreateInBatches(&adjustments, 100)

	fmt.Println("Seeding 36.000 Transaction")

	totalMonths := 36
	txPerMonth := 1000
	var transactions []models.Transaction

	startDate := time.Now().AddDate(-3, 0, 0)

	for m := 0; m < totalMonths; m++ {
		currentMonthDate := startDate.AddDate(0, m, 0)
		for t := 1; t <= txPerMonth; t++ {
			txType := gofakeit.RandomString([]string{"INBOUND", "OUTBOUND", "TRANSFER"})
			status := gofakeit.RandomString([]string{"approved", "draft"})

			daysInMonth := 28
			randomDay := gofakeit.Number(0, daysInMonth-1)
			txDate := currentMonthDate.AddDate(0, 0, randomDay)

			creatorID := uint(gofakeit.Number(3, 5))
			var approverID *uint
			if status == "approved" {
				mngrID := uint(2)
				approverID = &mngrID
			}

			whID := uint(gofakeit.Number(1, 100))
			var targetWhID *uint
			var partnerID *uint

			if txType == "TRANSFER" {
				tID := uint(gofakeit.Number(1, 100))
				if tID == whID {
					tID = (whID % 99) + 1
				} // Pastikan tidak sama
				targetWhID = &tID
			} else if txType == "INBOUND" {
				pID := uint(gofakeit.Number(1, 5)) // 1-5 adalah Supplier
				partnerID = &pID
			} else {
				pID := uint(gofakeit.Number(6, 10)) // 6-10 adalah Customer
				partnerID = &pID
			}

			numItems := gofakeit.Number(1, 3)
			var items []models.TransactionItem
			for k := 0; k < numItems; k++ {
				qty := gofakeit.Number(5, 50)
				price := float64(gofakeit.Number(10000, 500000))
				items = append(items, models.TransactionItem{
					ProductID: uint(gofakeit.Number(1, 1000)),
					Quantity:  qty,
					UnitPrice: price,
					SubTotal:  float64(qty) * price,
				})
			}

			tx := models.Transaction{
				ReferenceNo:       fmt.Sprintf("TRX-%s-%d-%05d", txType[:2], currentMonthDate.Unix(), t),
				TransactionDate:   txDate,
				Type:              txType,
				Status:            status,
				WarehouseID:       whID,
				TargetWarehouseID: targetWhID,
				PartnerID:         partnerID,
				CreatedByID:       creatorID,
				ApprovedByID:      approverID,
				Notes:             gofakeit.Sentence(4),
				Items:             items,
			}
			transactions = append(transactions, tx)
		}

		if err := db.CreateInBatches(&transactions, 200).Error; err != nil {
			log.Fatalf("❌ Gagal seeding transaksi pada bulan ke-%d: %v", m+1, err)
		}

		transactions = nil
		fmt.Printf("... Bulan %d/36 selesai\n", m+1)
	}

	fmt.Printf("🎉 MEGA-SEEDING SELESAI! Waktu eksekusi total: %v\n", time.Since(startTime))
	fmt.Println("Coba jalankan API Export Report Anda sekarang dan rasakan sensasi mengelola jutaan baris data!")
}
