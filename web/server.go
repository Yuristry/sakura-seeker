package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Blog struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Author     string // 博客作者
	Title      string // 博客标题
	Content    string // 博客内容
	WordCount  int    // 字数统计
	CoverImage string // 封面链接
	Date       string // 博客日期
	URL        string // 博客详细页面链接
}

func createDB() *gorm.DB {
	dsn := "root:password@tcp(127.0.0.1:3306)/sakura?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}

func getBlogs(w http.ResponseWriter, r *http.Request) {
	db := createDB()
	var blogs []Blog
	db.Find(&blogs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func main() {
	http.HandleFunc("/api/blogs", getBlogs)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
