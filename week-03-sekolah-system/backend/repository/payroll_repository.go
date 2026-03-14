package repository

import (
	"github.com/affandisy/school-system-sma/models"
	"gorm.io/gorm"
)

type PayrollRepository interface {
	GetTeacherConfig(teacherID string) (*models.TeacherSalaryConfig, error)
	SavePayslip(payslip *models.SalaryPayslip) error
	CheckPayslipExists(teacherID string, month, year int) bool
}

type payrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) PayrollRepository {
	return &payrollRepository{db}
}

func (r *payrollRepository) GetTeacherConfig(teacherID string) (*models.TeacherSalaryConfig, error) {
	var config models.TeacherSalaryConfig
	err := r.db.Where("teacher_id = ?", teacherID).First(&config).Error
	return &config, err
}

func (r *payrollRepository) SavePayslip(payslip *models.SalaryPayslip) error {
	return r.db.Create(payslip).Error
}

func (r *payrollRepository) CheckPayslipExists(teacherID string, month, year int) bool {
	var count int64
	r.db.Model(&models.SalaryPayslip{}).
		Where("teacher_id = ? AND month = ? AND year = ?", teacherID, month, year).
		Count(&count)
	return count > 0
}
