package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Actor struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Films     []Film    `json:"films" gorm:"many2many:actor_films;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Film struct {
	ID         uint       `json:"id"`
	Title      string     `json:"title"`
	LanguageID uint       `json:"-"`
	Language   Language   `json:"language"`
	Categories []Category `json:"categories" gorm:"many2many:film_categories;"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Language struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func initialData(db *gorm.DB) {
	tx := db.Begin()
	defer tx.Commit()

	lang := Language{
		Name: "English",
	}
	tx.Create(&lang)
	cate := Category{
		Name: "Documentary",
	}
	tx.Create(&cate)
	film := Film{
		Title:      "Young Language",
		Language:   lang,
		Categories: []Category{cate},
	}
	tx.Create(&film)
	actor := Actor{
		FirstName: "Ed",
		LastName:  "Chase",
		Films:     []Film{film},
	}
	tx.Create(&actor)
}

func main() {
	connURL := os.Getenv("POSTGRESQL_URL")
	db, err := gorm.Open("postgres", connURL)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	initialData(db)

	var dest []Actor
	db.Preload("Films.Language").Preload("Films.Categories").Find(&dest)
	json.NewEncoder(os.Stdout).Encode(&dest)
}
