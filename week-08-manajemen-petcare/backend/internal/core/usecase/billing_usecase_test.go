package usecase_test

import (
	"errors"
	"testing"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/usecase"
)

// --- 1. MOCK REPOSITORY ---
// Ini adalah "Database Palsu" yang mengimplementasikan port.BillingRepository
type mockBillingRepo struct {
	savedInvoice *domain.Invoice // Menyimpan hasil akhir untuk divalidasi
	shouldFail   bool            // Jika kita ingin mensimulasikan database error
}

func (m *mockBillingRepo) SaveInvoice(invoice *domain.Invoice) error {
	if m.shouldFail {
		return errors.New("simulasi kegagalan database")
	}
	m.savedInvoice = invoice
	return nil
}

// --- 2. SKENARIO PENGUJIAN ---

func TestGenerateInvoice_Success(t *testing.T) {
	// Persiapan (Arrange)
	mockRepo := &mockBillingRepo{}
	uc := usecase.NewBillingUseCase(mockRepo)
	ownerID := "owner-123"
	items := []domain.InvoiceItem{
		{PetID: "pet-A", ServiceID: "srv-1", Price: 150000}, // Grooming Kucing A
		{PetID: "pet-B", ServiceID: "srv-2", Price: 50000},  // Konsultasi Kucing B
	}

	// Eksekusi (Act)
	invoice, err := uc.GenerateInvoice(ownerID, items)

	// Validasi (Assert)
	if err != nil {
		t.Fatalf("Diharapkan sukses, tetapi mendapat error: %v", err)
	}

	if invoice == nil {
		t.Fatal("Diharapkan objek invoice dikembalikan, tetapi mendapat nil")
	}

	// Memverifikasi total harga (150k + 50k = 200k)
	expectedTotal := 200000.0
	if invoice.TotalAmount != expectedTotal {
		t.Errorf("Diharapkan total %v, mendapat %v", expectedTotal, invoice.TotalAmount)
	}

	// Memverifikasi detail bahwa data benar-benar diteruskan ke Mock Repository
	if mockRepo.savedInvoice.OwnerID != ownerID {
		t.Errorf("Diharapkan OwnerID %s, mendapat %s", ownerID, mockRepo.savedInvoice.OwnerID)
	}

	if len(mockRepo.savedInvoice.Items) != 2 {
		t.Errorf("Diharapkan 2 item, mendapat %v", len(mockRepo.savedInvoice.Items))
	}
}

func TestGenerateInvoice_Fail_EmptyItems(t *testing.T) {
	// Persiapan
	mockRepo := &mockBillingRepo{}
	uc := usecase.NewBillingUseCase(mockRepo)

	// Eksekusi dengan array kosong
	invoice, err := uc.GenerateInvoice("owner-123", []domain.InvoiceItem{})

	// Validasi
	if err == nil {
		t.Error("Diharapkan menghasilkan error karena item kosong, tetapi sukses")
	}
	if invoice != nil {
		t.Error("Diharapkan invoice bernilai nil jika gagal")
	}
	if err.Error() != "invoice harus memiliki minimal satu layanan" {
		t.Errorf("Pesan error tidak sesuai: %v", err.Error())
	}
}

func TestGenerateInvoice_Fail_MissingPetID(t *testing.T) {
	// Persiapan
	mockRepo := &mockBillingRepo{}
	uc := usecase.NewBillingUseCase(mockRepo)
	items := []domain.InvoiceItem{
		{PetID: "pet-A", ServiceID: "srv-1", Price: 150000},
		{PetID: "", ServiceID: "srv-2", Price: 50000}, // PetID sengaja dikosongkan
	}

	// Eksekusi
	invoice, err := uc.GenerateInvoice("owner-123", items)

	// Validasi
	if err == nil {
		t.Error("Diharapkan menghasilkan error karena ada PetID yang kosong, tetapi sukses")
	}
	if invoice != nil {
		t.Error("Diharapkan invoice bernilai nil jika gagal")
	}
	if err.Error() != "setiap layanan wajib diikat pada ID hewan peliharaan spesifik" {
		t.Errorf("Pesan error tidak sesuai: %v", err.Error())
	}
}
