package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	// 创建一个新的 Collector
	c := colly.NewCollector()

	// 处理错误
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Error:", err)
	})

	var blogTitles []string

	// 设置请求超时
	c.SetRequestTimeout(30 * time.Second)

	// 等待 class 为 "other-project-text" 的元素出现
	c.OnHTML("div.date-title", func(e *colly.HTMLElement) {
		title := e.ChildText(".title")
		blogTitles = append(blogTitles, title)
	})

	c.OnScraped(func(r *colly.Response) {
		// 爬取结束后，在这里可以对获取到的所有博客链接进行处理，比如打印或者存储等
		for _, link := range blogTitles {
			println(link)
		}
	})

	// 访问目标页面
	c.Visit("https://sakurazaka46.com/s/s46/diary/blog/list?ima=0000")
}
