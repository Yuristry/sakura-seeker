package spider

import (
	"context"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)

type Spider struct{}

func NewSpider() *Spider {
	return &Spider{}
}

func (s *Spider) Run(ctx context.Context) error {
	c := colly.NewCollector()

	var blogTitles []string

	c.SetRequestTimeout(30 * time.Second)

	// 设置回调函数
	c.OnHTML("div.date-title", func(e *colly.HTMLElement) {
		title := e.ChildText(".title")
		blogTitles = append(blogTitles, title)
	})

	// 启动爬虫
	err := c.Visit("https://sakurazaka46.com/s/s46/diary/blog/list?ima=0000")
	if err != nil {
		return err
	}

	log.Println("Spider is running...")

	c.OnScraped(func(r *colly.Response) {
		// 爬取结束后，在这里可以对获取到的所有博客链接进行处理，比如打印或者存储等
		for _, link := range blogTitles {
			println(link)
		}
	})

	return nil
}
