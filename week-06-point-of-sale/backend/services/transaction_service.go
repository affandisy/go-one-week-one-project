package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/affandisy/pos-system/models"
	"github.com/affandisy/pos-system/repositories"
	"github.com/google/uuid"
)

type CartItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CheckoutRequest struct {
	Items         []CartItem `json:"items"`
	Discount      float64    `json:"discount"`
	PaymentMethod string     `json:"payment_method"` // "cash" atau "qris"
	CashGiven     float64    `json:"cash_given"`
}

type TransactionService interface {
	Checkout(cashierID string, req CheckoutRequest) (*models.Transaction, error)
}

type transactionService struct {
	trxRepo  repositories.TransactionRepository
	prodRepo repositories.ProductRepository
}

func NewTransactionService(trxRepo repositories.TransactionRepository, prodRepo repositories.ProductRepository) TransactionService {
	return &transactionService{trxRepo, prodRepo}
}

func (s *transactionService) Checkout(cashierID string, req CheckoutRequest) (*models.Transaction, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("keranjang kosong")
	}

	cashierUUID, err := uuid.Parse(cashierID)
	if err != nil {
		return nil, errors.New("ID Kasir tidak valid")
	}

	var totalAmount float64
	var details []models.TransactionDetail
	var movements []models.StockMovement

	// 1. Validasi Stok & Hitung Subtotal
	for _, item := range req.Items {
		product, err := s.prodRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("produk tidak ditemukan")
		}

		// Aturan Bisnis PRD: Stok tidak boleh negatif
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("stok %s tidak mencukupi (Sisa: %d)", product.Name, product.Stock)
		}

		subtotal := product.SellingPrice * float64(item.Quantity)
		totalAmount += subtotal

		// Siapkan Detail Transaksi
		details = append(details, models.TransactionDetail{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.SellingPrice,
			Subtotal:  subtotal,
		})

		// Siapkan Riwayat Stok (OUT)
		movements = append(movements, models.StockMovement{
			ProductID: product.ID,
			Type:      "OUT",
			Quantity:  item.Quantity,
			Note:      "Penjualan",
		})
	}

	// 2. Hitung Total Akhir & Kembalian
	finalAmount := totalAmount - req.Discount
	if finalAmount < 0 {
		finalAmount = 0
	}

	var changeAmount float64
	if req.PaymentMethod == "cash" {
		if req.CashGiven < finalAmount {
			return nil, errors.New("uang tunai tidak mencukupi")
		}
		changeAmount = req.CashGiven - finalAmount
	} else {
		// Jika QRIS/E-Wallet, uang yang diberikan dianggap pas
		req.CashGiven = finalAmount
		changeAmount = 0
	}

	// 3. Buat Nomor Struk (Contoh: TRX-20260422-153022)
	receiptNumber := fmt.Sprintf("TRX-%s", time.Now().Format("20060102-150405"))

	// 4. Rakit Objek Transaksi
	transaction := &models.Transaction{
		ReceiptNumber: receiptNumber,
		CashierID:     cashierUUID,
		TotalAmount:   totalAmount,
		Discount:      req.Discount,
		FinalAmount:   finalAmount,
		PaymentMethod: req.PaymentMethod,
		CashGiven:     req.CashGiven,
		ChangeAmount:  changeAmount,
		Details:       details,
	}

	// 5. Eksekusi penyimpanan ke Database (Atomic)
	if err := s.trxRepo.CreateTransaction(transaction, movements); err != nil {
		return nil, errors.New("gagal memproses transaksi")
	}

	return transaction, nil
}
