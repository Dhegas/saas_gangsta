package main

import (
	"fmt"
	"log"

	"github.com/dhegas/saas_gangsta/internal/database"

	"github.com/joho/godotenv"
)

// Contoh model tabel sederhana
type UserTest struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string `gorm:"unique"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	database.Connect()

	db := database.DB

	// 1. Membuat tabel di Supabase (secara otomatis)
	log.Println("Migrasi tabel UserTest...")
	db.AutoMigrate(&UserTest{})

	// 2. Insert atau Get data test
	newUser := UserTest{Name: "Admin Gangsta", Email: "admin@saas-gangsta.com"}

	result := db.Where(UserTest{Email: newUser.Email}).FirstOrCreate(&newUser)
	if result.Error != nil {
		log.Println("Gagal menyimpan data:", result.Error)
	} else {
		log.Println("Data berhasil disimpan / ditemukan dengan ID:", newUser.ID)
	}

	// 3. Menampilkan isi tabel
	var users []UserTest
	db.Find(&users)
	fmt.Println("=== DAFTAR USER DI SUPABASE ===")
	for _, u := range users {
		fmt.Printf("- %d | %s | %s\n", u.ID, u.Name, u.Email)
	}
	fmt.Println("===============================")

	// router...
}
