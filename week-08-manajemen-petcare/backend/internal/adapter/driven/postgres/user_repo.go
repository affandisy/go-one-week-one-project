package postgres

import (
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"gorm.io/gorm"
)

type UserModel struct {
	ID           string `gorm:"type:uuid;primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null"`
	CreatedAt    time.Time
}

func (UserModel) TableName() string { return "users" }

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	db.AutoMigrate(&UserModel{})
	return &userRepository{db}
}

func (r *userRepository) Save(user *domain.User) error {
	model := UserModel{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
	}
	return r.db.Create(&model).Error
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var model UserModel
	if err := r.db.Where("username = ?", username).First(&model).Error; err != nil {
		return nil, err
	}
	return &domain.User{
		ID:           model.ID,
		Username:     model.Username,
		PasswordHash: model.PasswordHash,
		Role:         model.Role,
	}, nil
}
