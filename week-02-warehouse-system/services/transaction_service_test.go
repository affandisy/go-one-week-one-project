package services_test

import (
	"errors"
	"testing"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/mocks"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessTransaction_WarehouseInvalid(t *testing.T) {
	mockWhRepo := new(mocks.WarehouseRepository)

	txService := services.NewTransactionService(nil, nil, nil, nil, mockWhRepo, nil)

	mockWhRepo.On("FindByID", uint(99)).Return((*models.Warehouse)(nil), errors.New("record not found"))

	req := services.CreateTransactionRequest{
		WarehouseID: 99,
		Type:        "TRANSFER",
		Items: []services.TransactionItemDTO{
			{ProductID: 1, Quantity: 10},
		},
	}

	result, err := txService.ProcessTransaction(req, 1)

	assert.Error(t, err)
	assert.Equal(t, "Gudang asal tidak valid", err.Error())
	assert.Nil(t, result)

	mockWhRepo.AssertExpectations(t)
}

func TestProcessTransaction_InboundSuccess(t *testing.T) {
	mockTxRepo := new(mocks.TransactionRepository)
	mockWhRepo := new(mocks.WarehouseRepository)
	mockPartnerRepo := new(mocks.PartnerRepository)

	txService := services.NewTransactionService(mockTxRepo, nil, mockPartnerRepo, nil, mockWhRepo, nil)

	validWarehouse := &models.Warehouse{ID: 1, Code: "JKT-01"}
	validPartner := &models.Partner{ID: 1, Type: "SUPPLIER"}

	mockWhRepo.On("FindByID", uint(1)).Return(validWarehouse, nil)
	mockPartnerRepo.On("FindByID", uint(1)).Return(validPartner, nil)

	mockTxRepo.On("ExecuteTransaction", mock.AnythingOfType("*models.Transaction"), mock.Anything).Return(nil)

	partnerID := uint(1)
	req := services.CreateTransactionRequest{
		ReferenceNo: "PO-001",
		WarehouseID: 1,
		PartnerID:   &partnerID,
		Type:        "INBOUND",
		Items: []services.TransactionItemDTO{
			{ProductID: 1, Quantity: 50, UnitPrice: 10000},
		},
	}

	result, err := txService.ProcessTransaction(req, 100)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "draft", result.Status)
	assert.Equal(t, "INBOUND", result.Type)
	assert.Len(t, result.Items, 1)

	mockWhRepo.AssertExpectations(t)
	mockPartnerRepo.AssertExpectations(t)
	mockTxRepo.AssertExpectations(t)
}
