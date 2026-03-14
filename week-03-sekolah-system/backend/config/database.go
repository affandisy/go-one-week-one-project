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
		log.Fatal("❌ Gagal terhubung ke PostgreSQL: ", err)
	}

	log.Println("✅ Berhasil terhubung ke PostgreSQL")
	DB = database

	// URUTAN MIGRASI SANGAT PENTING!
	// Tabel yang tidak punya Foreign Key harus dibuat pertama kali.
	// Tabel yang bergantung pada tabel lain dibuat belakangan.
	err = DB.AutoMigrate(
		// 1. Independent Tables (Master Data Dasar)
		&models.User{},
		&models.AcademicYear{},
		&models.StudyProgram{},
		&models.SalaryComponent{},
		&models.GradeComponent{},

		// 2. Tables with 1st Level Dependencies
		&models.Subject{},             // Bergantung pada StudyProgram
		&models.Class{},               // Bergantung pada StudyProgram, AcademicYear, User
		&models.TeacherSalaryConfig{}, // Bergantung pada User

		// 3. Tables with 2nd Level Dependencies
		&models.Student{},           // Bergantung pada User, Class, AcademicYear
		&models.TeacherAssignment{}, // Bergantung pada User, Class, Subject, AcademicYear
		&models.Attendance{},        // Bergantung pada User

		// 4. Transactional & Heavy Relational Tables
		&models.Grade{},         // Bergantung pada Student, Subject, GradeComponent, User
		&models.SalaryPayslip{}, // Bergantung pada User
	)

	if err != nil {
		log.Fatal("❌ Gagal melakukan migrasi database: ", err)
	}

	log.Println("✅ Migrasi Database Selesai dan Rapi!")
}
