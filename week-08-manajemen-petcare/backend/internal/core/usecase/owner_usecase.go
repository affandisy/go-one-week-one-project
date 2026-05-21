package usecase

import (
	"errors"
	"strings"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/google/uuid"
)

type ownerUseCase struct {
	repo port.OwnerRepository
}

func NewOwnerUseCase(repo port.OwnerRepository) port.OwnerUseCase {
	return &ownerUseCase{repo}
}

func (uc *ownerUseCase) RegisterOwner(name, phone string) (*domain.Owner, error) {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return nil, errors.New("nama pemilik wajib diisi")
	}

	owner := &domain.Owner{
		ID:    uuid.NewString(),
		Name:  cleanName,
		Phone: strings.TrimSpace(phone),
	}

	if err := uc.repo.Save(owner); err != nil {
		return nil, err
	}
	return owner, nil
}

func (uc *ownerUseCase) ListOwners() ([]domain.Owner, error) {
	return uc.repo.FindAll()
}
