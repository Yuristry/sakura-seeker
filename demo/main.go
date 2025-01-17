package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Blog struct {
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
		return match[1] // 提取 URL
	}
	return ""
}

func main() {
	baseURL := "https://sakurazaka46.com/s/s46/diary/blog/list?ima=0000&page="
	var blogs []Blog

	// 创建 Colly 收集器
	c := colly.NewCollector(
		colly.AllowedDomains("sakurazaka46.com"),
	)

	// 解析博客列表
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

		blogs = append(blogs, blog)
	})
	maxPage := 1
	// 启动爬取
	for page := 0; page <= maxPage; page++ {
		err := c.Visit(baseURL + strconv.Itoa(page))
		if err != nil {
			log.Fatal("Failed to scrape page:", err)
		}
	}

	fmt.Printf("Total Blogs: %d\n", len(blogs))

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
	}
}
