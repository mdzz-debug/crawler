## Crawler

#### golang 封装爬虫依赖

1. 可对json、接口进行单独获取以及并发获取
2. 可对网站单独或者批量获取DOM元素内的文本内容



单独获取接口 / JSON

```go
var cra crawler.Crawler
cra.Url = "爬取地址"
cra.Request("get") // 设置请求方法get/post...

for _,v := range cra.Results{
// v 即为结果 []byte
// 用json.Unmarshal()自由处理
}
```

并发获取接口 / JSON

```go
var cra crawler.Crawler

cra.WorkerMaxNum = 128 // 协程池中的最大协程数,默认值为128，可忽略
cra.Urls = [...]
cra.Requests("get") // 设置请求方法get/post...

for _,v := range cra.Results{
// v 即为结果 []byte
// 用json.Unmarshal()自由处理
}
```

单独获取DOM

```go
var cra crawler.Crawler
cra.Url = "爬取地址"
cra.RequestOnHtml("div[class=...]")

for _,v := range cra.Results{
v.ChildTexts("a") // 获取节点下标签a文本
v.ChildAttrs("a", "href") // 获取节点下a标签的href
...  // 具体方法和colly无差别
}
```

批量获取DOM

```go
var cra crawler.Crawler

cra.WorkerMaxNum = 128 // 协程池中的最大协程数,默认值为128，可忽略
cra.Urls = [...]
cra.RequestsOnHtml("div[class=...]")

for _,v := range cra.DOM{
v.ChildTexts("a") // 获取节点下标签a文本
v.ChildAttrs("a", "href") // 获取节点下a标签的href
...  // 具体方法和colly无差别
}
```

