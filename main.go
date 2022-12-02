package crawler

import (
	"github.com/gocolly/colly/v2"
)

// Head 请求头结构体
type Head struct {
	Name  string
	Value string
}

// Crawler 爬虫包结构体
type Crawler struct {
	Url           string               // 爬虫地址
	Urls          []string             // 爬虫地址
	HeaderInfo    []Head               // 爬虫地址请求头
	Data          []string             // 爬虫地址请求参数
	WorkerMaxNum  int                  // 最大工作人数
	Results       [][]byte             // 批量请求结果获取
	DOM           []*colly.HTMLElement // DOM截取结果
	FailUrl       []string             // 失败结果集合
	requestMethod string               // 并发请求方式
	workerNum     int                  // 工作人数
	workPassage   chan int             // 工作通道
	workDone      chan bool            // 工作完成
	status        int                  // 爬虫方式【只影响并发】 0 获取内容 1 dom
	dos           string               // dom并发分配dom字段
}

// Request 单条请求
func (c *Crawler) Request(method string) (successNum int, errorNum int) {
	client, err := c.client(method, c.Url)
	if err != nil {
		c.FailUrl = append(c.FailUrl, c.Url)
		return
	}
	c.Results = append(c.Results, client)

	return len(c.Results), len(c.FailUrl)
}

// Requests 多条请求
func (c *Crawler) Requests(method string) (successNum int, errorNum int) {
	// 初始化
	c.initialize()

	c.requestMethod = method
	c.workers(1)
	go c.allocateAcquisition(0, true)
	c.wait()

	return len(c.Results), len(c.FailUrl)
}

// RequestOnHtml 单条请求选取DOM
func (c *Crawler) RequestOnHtml(dom string) (successNum int, errorNum int) {
	// 初始化
	c.initialize()

	backDom, err := c.domCrawl(dom, c.Url)
	if err != nil {
		c.FailUrl = append(c.FailUrl, c.Url)
	}

	c.DOM = append(c.DOM, backDom)
	return len(c.DOM), len(c.FailUrl)
}

// RequestsOnHtml 多条请求选取DOM
func (c *Crawler) RequestsOnHtml(dom string) (successNum int, errorNum int) {
	// 初始化
	c.initialize()

	c.status = 1
	c.dos = dom
	c.workers(1)
	go c.allocateAcquisition(0, true)
	c.wait()

	return len(c.Results), len(c.FailUrl)
}
