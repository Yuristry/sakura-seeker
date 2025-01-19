package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
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

func extractImageURL(styleAttr string) string {
	re := regexp.MustCompile(`url\((.*?)\)`) // 匹配 `url(...)`
	match := re.FindStringSubmatch(styleAttr)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func createDB() *gorm.DB {
	dsn := "root:password@tcp(127.0.0.1:3306)/sakura?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&Blog{})
	return db
}

func stopCrawler() {
	log.Println("Stopping crawler as no new blogs found.")
}

func tryAlert() {
	log.Println("New blogs found! Alerting...")
}

func main() {
	baseURL := "https://sakurazaka46.com/s/s46/diary/blog/list?ima=0000&page="
	var blogs []Blog

	c := colly.NewCollector(
		colly.AllowedDomains("sakurazaka46.com"),
	)
	db := createDB()

	c.OnHTML(".blog-top .com-blog-part .box", func(e *colly.HTMLElement) {
		blog := Blog{}
		blog.Date = strings.TrimSpace(e.ChildText(".date"))             // 日期
		blog.Author = strings.TrimSpace(e.ChildText(".name"))           // 博客作者
		blog.Title = strings.TrimSpace(e.ChildText(".title"))           // 标题
		blog.Content = strings.TrimSpace(e.ChildText(".lead"))          // 内容摘要
		blog.WordCount = len([]rune(blog.Content))                      // 计算字数
		blog.CoverImage = extractImageURL(e.ChildAttr(".img", "style")) // 封面链接
		blog.URL = e.ChildAttr("a", "href")                             // 博客链接
		if !strings.HasPrefix(blog.URL, "https://") {
			blog.URL = "https://sakurazaka46.com" + blog.URL
		}
		var existBlog Blog
		if result := db.Where("url = ?", blog.URL).First(&existBlog); result.Error == gorm.ErrRecordNotFound {
			blogs = append(blogs, blog)
		} else if result.Error == nil {
			fmt.Printf("Blog already exist: %s\n", blog.URL)
		} else {
			fmt.Println("error", result.Error)
			return
		}
	})
	maxPage := 1000
	for page := 0; page <= maxPage; page++ {
		fmt.Printf("Try to visit Page %d\n", page)
		oldNum := len(blogs)
		err := c.Visit(baseURL + strconv.Itoa(page))
		if err != nil {
			log.Fatal("Failed to scrape page:", err)
		}
		if len(blogs) == oldNum {
			stopCrawler()
			break
		}
	}

	fmt.Printf("Total Blogs: %d\n", len(blogs))

	if len(blogs) != 0 {
		tryAlert()
	}

	// 打印结果
	for _, blog := range blogs {
		fmt.Printf("Author: %s\n", blog.Author)
		fmt.Printf("Title: %s\n", blog.Title)
		fmt.Printf("Date: %s\n", blog.Date)
		fmt.Printf("Content: %s\n", blog.Content)
		fmt.Printf("Word Count: %d\n", blog.WordCount)
		fmt.Printf("Cover Image: %s\n", blog.CoverImage)
		fmt.Printf("URL: %s\n", blog.URL)
		fmt.Println("-------------------------------------------------")
		if err := db.Create(&blog).Error; err != nil {
			fmt.Println("error is: ", err)
		}
	}
}
