package config

import (
	"log"

	"github.com/affandisy/school-system-sma/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=school_user password=school_password dbname=school_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke PostgreSQL: ", err)
	}

	log.Println("Berhasil terhubung ke PostgreSQL")
	DB = database

	err = DB.AutoMigrate(

		&models.User{},
		&models.AcademicYear{},
		&models.StudyProgram{},
		&models.SalaryComponent{},
		&models.GradeComponent{},

		&models.Subject{},
		&models.Class{},
		&models.TeacherSalaryConfig{},
		&models.Announcement{},
		&models.RecommendationCriteria{},

		&models.Student{},
		&models.TeacherAssignment{},
		&models.Attendance{},

		&models.Grade{},
		&models.SalaryPayslip{},
		&models.StudentRecommendation{},
		&models.ReportCard{},

		&models.AuditLog{},
	)

	if err != nil {
		log.Fatal("Gagal melakukan migrasi database: ", err)
	}

	log.Println("Migrasi Database Selesai dan Rapi!")
}
