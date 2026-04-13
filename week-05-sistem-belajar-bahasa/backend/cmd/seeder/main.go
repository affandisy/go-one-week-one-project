package main

import (
	"encoding/json"
	"log"

	"github.com/affandi/belajar-bahasa/config"
	"github.com/affandi/belajar-bahasa/models"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Memulai proses Data Seeding...")

	if err := godotenv.Load(); err != nil {
		log.Println("Catatan: File .env tidak ditemukan, menggunakan variabel OS.")
	}

	config.ConnectDatabase()
	db := config.DB

	var count int64
	db.Model(&models.Module{}).Count(&count)
	if count > 0 {
		log.Println("Peringatan: Database sudah berisi data Modul. Proses seeding dibatalkan agar data aman.")
		return
	}

	modules := []models.Module{
		{Title: "Alphabet", Description: "Pengenalan huruf dan bunyi dasar.", LevelOrder: 1, IsLocked: false},
		{Title: "Vocabulary", Description: "Kosakata benda, hewan, dan kata kerja umum.", LevelOrder: 2, IsLocked: true},
		{Title: "Sentences", Description: "Membangun struktur kalimat sederhana (S-P-O).", LevelOrder: 3, IsLocked: true},
	}

	for i := range modules {
		if err := db.Create(&modules[i]).Error; err != nil {
			log.Fatalf("Gagal membuat modul %s: %v", modules[i].Title, err)
		}
	}
	log.Println("✅ Berhasil menyisipkan 3 Modul Level.")

	// Helper untuk mengubah array string menjadi JSON string
	makeOptions := func(opts []string) string {
		bytes, _ := json.Marshal(opts)
		return string(bytes)
	}

	mod1ID := modules[0].ID
	materialsMod1 := []models.Material{
		// Flashcards (Learn Mode)
		{ModuleID: mod1ID, ContentType: "learn_card", Question: "A", CorrectAnswer: "Dibaca: Ei (Apple)", DisplayOrder: 1},
		{ModuleID: mod1ID, ContentType: "learn_card", Question: "B", CorrectAnswer: "Dibaca: Bi (Banana)", DisplayOrder: 2},
		{ModuleID: mod1ID, ContentType: "learn_card", Question: "C", CorrectAnswer: "Dibaca: Si (Cat)", DisplayOrder: 3},

		// Quiz (Kuis)
		{
			ModuleID: mod1ID, ContentType: "quiz_mcq",
			Question:      "Huruf apa yang berada di awal kata 'Apple'?",
			CorrectAnswer: "A",
			Options:       makeOptions([]string{"A", "B", "C", "D"}),
			DisplayOrder:  4,
		},
		{
			ModuleID: mod1ID, ContentType: "quiz_mcq",
			Question:      "Pilih cara membaca huruf C dalam bahasa Inggris:",
			CorrectAnswer: "Si",
			Options:       makeOptions([]string{"Ci", "Si", "Ke", "Se"}),
			DisplayOrder:  5,
		},
	}

	// --- MATERI UNTUK LEVEL 2: VOCABULARY ---
	mod2ID := modules[1].ID
	materialsMod2 := []models.Material{
		// Flashcards
		{ModuleID: mod2ID, ContentType: "learn_card", Question: "Dog", CorrectAnswer: "Anjing", DisplayOrder: 1},
		{ModuleID: mod2ID, ContentType: "learn_card", Question: "Cat", CorrectAnswer: "Kucing", DisplayOrder: 2},
		{ModuleID: mod2ID, ContentType: "learn_card", Question: "Eat", CorrectAnswer: "Makan", DisplayOrder: 3},

		// Quiz
		{
			ModuleID: mod2ID, ContentType: "quiz_mcq",
			Question:      "Apa bahasa Inggris dari 'Kucing'?",
			CorrectAnswer: "Cat",
			Options:       makeOptions([]string{"Dog", "Cat", "Eat", "Bird"}),
			DisplayOrder:  4,
		},
		{
			ModuleID: mod2ID, ContentType: "quiz_mcq",
			Question:      "Kata kerja untuk 'Makan' adalah?",
			CorrectAnswer: "Eat",
			Options:       makeOptions([]string{"Run", "Sleep", "Eat", "Drink"}),
			DisplayOrder:  5,
		},
	}

	// Eksekusi insert massal ke database
	db.Create(&materialsMod1)
	db.Create(&materialsMod2)

	log.Println("✅ Berhasil menyisipkan Materi dan Soal Kuis awal.")
	log.Println("🎉 Proses Data Seeding selesai secara keseluruhan!")
}
