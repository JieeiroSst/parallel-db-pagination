package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	Name       string `json:"name`
	Desciption string `json:"desciption"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func insert(db *gorm.DB) error {
	for i := 0; i < 1000; i++ {
		book := Book{
			Name:       randSeq(4),
			Desciption: randSeq(5),
		}

		db.Create(&book)
	}
	return nil
}

func total(db *gorm.DB) int {
	var total int64
	db.Model(&Book{}).Count(&total)
	return int(total)
}

func fetchData(db *gorm.DB, page, limit int) ([]Book, error) {
	res := make([]Book, 0)
	if page < 1 {
		page = 1
	}

	err := db.Offset((page - 1) * limit).Limit(limit).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func main() {
	dsn := "host=localhost user=hoadev password=hoadev123 dbname=hoadev_db port=5433 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&Book{})

	// insert(db)

	total := total(db)
	fmt.Println("===", total)

	limit := 100

	var wg sync.WaitGroup
	ch := make(chan []Book)
	totalPages := math.Ceil(float64(total) / float64(limit))
	totalPagess := int(totalPages)

	for i := 1; i <= totalPagess; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			data, err := fetchData(db, page, limit)
			if err != nil {
				return
			}
			ch <- data
		}(i)
	}

	var allData []Book
	for i := 0; i < totalPagess; i++ {
		result := <-ch
		allData = append(allData, result...)
	}

	fmt.Println("Total Pages:", totalPagess)
	aa, _ := json.Marshal(&allData)
	fmt.Println("=====", string(aa))
	fmt.Println("===== len ", len(allData))

	wg.Wait()
}
