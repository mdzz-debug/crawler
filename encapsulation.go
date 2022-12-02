package crawler

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

var (
	userAgent = []string{
		"Mozilla/5.0 (compatible; U; ABrowse 0.6; Syllable) AppleWebKit/420+ (KHTML, like Gecko)",
		"Mozilla/5.0 (compatible; U; ABrowse 0.6;  Syllable) AppleWebKit/420+ (KHTML, like Gecko)",
		"Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; Acoo Browser 1.98.744; .NET CLR 3.5.30729)",
		"Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; Acoo Browser 1.98.744; .NET CLR   3.5.30729)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; Acoo Browser; GTB5; Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1) ; InfoPath.1; .NET CLR 3.5.30729; .NET CLR 3.0.30618)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; SV1; Acoo Browser; .NET CLR 2.0.50727; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729; Avant Browser)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; SLCC1;   .NET CLR 2.0.50727; Media Center PC 5.0; .NET CLR 3.0.04506)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; GTB5; Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1) ; Maxthon; InfoPath.1; .NET CLR 3.5.30729; .NET CLR 3.0.30618)",
		"Mozilla/4.0 (compatible; Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; Acoo Browser 1.98.744; .NET CLR 3.5.30729); Windows NT 5.1; Trident/4.0)",
		"Mozilla/4.0 (compatible; Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; GTB6; Acoo Browser; .NET CLR 1.1.4322; .NET CLR 2.0.50727); Windows NT 5.1; Trident/4.0; Maxthon; .NET CLR 2.0.50727; .NET CLR 1.1.4322; InfoPath.2)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; Acoo Browser; GTB6; Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1) ; InfoPath.1; .NET CLR 3.5.30729; .NET CLR 3.0.30618)",
	}
)

// 初始化
func (c *Crawler) initialize() {
	if c.WorkerMaxNum == 0 {
		c.WorkerMaxNum = 128
	}
	c.workPassage = make(chan int)
	c.workDone = make(chan bool)
	c.DOM = []*colly.HTMLElement{}
	c.Results = [][]byte{}
	c.FailUrl = []string{}
}

// 工位增减
func (c *Crawler) workers(num int) {
	if c.workerNum+num < 0 {
		return
	}
	c.workerNum = c.workerNum + num
}

// 分配处理
func (c *Crawler) allocateAcquisition(i int, concurrent bool) {
	if i < len(c.Urls) {
		if c.workerNum < c.WorkerMaxNum {
			c.workPassage <- i + 1
		} else {
			c.allocateAcquisition(i+1, false)
		}

		if c.status == 1 {
			backDom, err := c.domCrawl(c.dos, c.Urls[i])
			if err != nil {
				c.FailUrl = append(c.FailUrl, c.Urls[i])
				return
			}
			c.DOM = append(c.DOM, backDom)
		} else {
			client, err := c.client(c.requestMethod, c.Urls[i])
			if err != nil {
				c.FailUrl = append(c.FailUrl, c.Urls[i])
				return
			}
			c.Results = append(c.Results, client)
		}
	}

	if concurrent {
		c.workDone <- true
	}
}

// 等待程序
func (c *Crawler) wait() {
	for {
		select {
		case key := <-c.workPassage:
			c.workers(1)
			go c.allocateAcquisition(key, true)
		case <-c.workDone:
			c.workers(-1)
			if c.workerNum == 0 {
				return
			}
		}
	}
}

// GET请求方法
func (c *Crawler) client(method string, url string) ([]byte, error) {
	client := http.Client{}
	var req *http.Request
	var resp *http.Response
	var err error

	switch strings.ToUpper(method) {
	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("request method error")
	}

	agent := false
	for _, k := range c.HeaderInfo {
		if strings.ToLower(k.Name) == "user-agent" {
			agent = true
		}
		req.Header.Add(k.Name, k.Value)
	}
	if !agent {
		req.Header.Add("user-agent", userAgent[rand.Intn(len(userAgent))])
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	backData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return backData, nil
}

// DOM爬取
func (c *Crawler) domCrawl(dom string, url string) (*colly.HTMLElement, error) {
	var ele *colly.HTMLElement
	b := colly.NewCollector(
		colly.UserAgent(userAgent[rand.Intn(len(userAgent))]),
	)
	b.OnHTML(dom, func(element *colly.HTMLElement) {
		ele = element
	})
	err := b.Visit(url)
	if err != nil {
		return nil, err
	}
	return ele, nil
}
