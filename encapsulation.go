package crawler

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

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

		fmt.Println(c.Urls[i])
		client, err := c.client(c.requestMethod, c.Urls[i])
		if err != nil {
			return
		}
		c.BackResults = append(c.BackResults, client)
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
	for _, k := range c.HeaderInfo {
		req.Header.Add(k.Name, k.Value)
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
