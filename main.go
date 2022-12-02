package crawler

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"math/rand"
)

// Head 请求头结构体
type Head struct {
	Name  string
	Value string
}

// Crawler 爬虫包结构体
type Crawler struct {
	Url           string    // 爬虫地址
	Urls          []string  // 爬虫地址
	HeaderInfo    []Head    // 爬虫地址请求头
	Data          []string  // 爬虫地址请求参数
	WorkerMaxNum  int       // 最大工作人数
	BackResults   [][]byte  // 批量请求结果获取
	BackErrorUrl  []string  // 批量请求结果获取
	requestMethod string    // 并发请求方式
	workerNum     int       // 工作人数
	workPassage   chan int  // 工作通道
	workDone      chan bool // 工作完成
}

// Requests 多条请求
func (c *Crawler) Requests(method string) (successNum int, errorNum int) {
	if c.WorkerMaxNum == 0 {
		c.WorkerMaxNum = 128
	}
	c.workPassage = make(chan int)
	c.workDone = make(chan bool)
	c.requestMethod = method
	c.workers(1)
	go c.allocateAcquisition(0, true)
	c.wait()

	return len(c.BackResults), len(c.BackErrorUrl)
}

// Request 单条请求
func (c *Crawler) Request(method string) (successNum int, errorNum int) {
	client, err := c.client(method, c.Url)
	if err != nil {
		c.BackErrorUrl = append(c.BackErrorUrl, c.Url)
		return
	}
	c.BackResults = append(c.BackResults, client)

	return len(c.BackResults), len(c.BackErrorUrl)
}

func (c *Crawler) RequestOnHtml(dom string) (*colly.HTMLElement, error) {
	var ele *colly.HTMLElement
	b := colly.NewCollector(
		colly.UserAgent(userAgent[rand.Intn(len(userAgent))]),
	)
	b.OnHTML(dom, func(element *colly.HTMLElement) {
		ele = element
	})
	err := b.Visit(c.Url)
	if err != nil {
		fmt.Println(err)
		return ele, err
	}
	return ele, nil
}
