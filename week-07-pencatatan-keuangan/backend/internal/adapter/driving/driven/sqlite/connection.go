package sqlite

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabaseConnection menginisialisasi SQLite dengan optimasi tingkat tinggi
func NewDatabaseConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		// Matikan log bawaan di produksi agar tidak memenuhi I/O disk
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// === OPTIMASI PRAGMA SQLITE ===

	// 1. Write-Ahead Logging: Jauh lebih cepat dari jurnal Rollback bawaan.
	// Memungkinkan pembacaan (SELECT) dan penulisan (INSERT) terjadi bersamaan.
	db.Exec("PRAGMA journal_mode = WAL;")

	// 2. NORMAL synchronous: Aman untuk WAL mode, mempercepat INSERT hingga 10x lipat.
	db.Exec("PRAGMA synchronous = NORMAL;")

	// 3. Menyimpan file temporary di memori (RAM) alih-alih di disk.
	db.Exec("PRAGMA temp_store = MEMORY;")

	// 4. Memastikan Foreign Key constraint aktif (SQLite defaultnya mati).
	db.Exec("PRAGMA foreign_keys = ON;")

	// 5. Menaikkan cache size (sekitar 20MB) agar pembacaan laporan langsung dari RAM.
	db.Exec("PRAGMA cache_size = -20000;")

	log.Println("✅ SQLite dioptimalkan dengan mode WAL dan In-Memory Temp Store")
	return db, nil
}
