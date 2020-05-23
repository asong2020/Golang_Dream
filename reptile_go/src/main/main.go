package main

import (
	"fmt"
	"github.com/jackdanger/collectlinks"
	"io/ioutil"
	"net/http"
	"src/src/reptile"
)

func demo1()  {
	resp,err := http.Get("http://www.baidu.com")
	if err != nil{
		fmt.Println("http get err",err)
		return
	}
	body,err:= ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("read error",err)
		return
	}
	fmt.Println(string(body))
}

func download(url string,queue chan string)  {
	client := &http.Client{}
	req,_:= http.NewRequest("GET",url,nil)
	req.Header.Set("User-Agent","Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resq,err:=client.Do(req)
	if err!=nil{
		fmt.Println("http get error",err)
		return
	}
	defer resq.Body.Close()
	fmt.Println(resq.Body)
	links := collectlinks.All(resq.Body)
	for _,link := range links{
		fmt.Println("parse url",link)
		go func() {
			queue <- link
		}()
	}

}

func main()  {
/*	fmt.Println("Hello, world")
	url := "http://www.baidu.com/"
	queue := make(chan string) //声明一个通道
	go func() {
		queue <- url
	}()
	for ur := range queue{
		download(ur,queue)
	}*/
	reptile.Get_Book("https://www.xsbiquge.com/91_91600/")
/*	s := `<div id="list"><dl><dt>正文</dt><dd><a href="test">绯红</a></dd></dl></div>`
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href"  && a.Val == "test"{
					for c:=n.FirstChild;c!=nil;c=c.NextSibling{
						fmt.Println(c.Data)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)*/
}
