package reptile

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

var find bool   // 寻找文章标志位
var listFind bool // 寻找文章列表标志位
var buf bytes.Buffer
var num int = 10  //爬取章节数量

const server  = "https://www.xsbiquge.com"

// server 根url target 目标url
func Get_Book(target string)  error{
	resp,err := http.Get(target)
	if err!=nil{
		fmt.Println("get err http",err)
		return err
	}
	defer resp.Body.Close()
	doc,err := html.Parse(resp.Body)
	if err != nil{
		fmt.Println("html parse err",err)
		return err
	}
	parseList(doc)
	return nil
}
//找到文章列表
func parseList(n *html.Node)  {
 	if n.Type == html.ElementNode && n.Data == "div"{
 		for _,a := range n.Attr{
 			if a.Key == "id" && a.Val == "list"{
 				listFind = true
				parseTitle(n)
 				break
			}
		}
	}
	if !listFind{
		for c := n.FirstChild;c!=nil;c=c.NextSibling{
			parseList(c)
		}
	}
}
//获取文章头部
func parseTitle(n *html.Node)  {
	if n.Type == html.ElementNode && n.Data == "a"{
		for _, a := range n.Attr {
			if a.Key == "href"{
				//获取文章title
				for c:=n.FirstChild;c!=nil;c=c.NextSibling{
					buf.WriteString(c.Data+ "\n")
				}
				url := a.Val
				target := server + url // 得到 文章url
				everyChapter(target)
				num--
			}
		}
	}
	if num <= 0{
		return
	}else {
		for c := n.FirstChild;c!=nil;c=c.NextSibling{
			parseTitle(c)
		}
	}
}
//获取每个章节
func everyChapter(target string)  {
	fmt.Println(target)
	resp,err := http.Get(target)
	if err!=nil{
		fmt.Println("get err http",err)
	}
	defer resp.Body.Close()
	doc,err := html.Parse(resp.Body)
	find = false
	parse(doc)
	text,err := os.Create("三国之他们非要打种地的我.txt")
	if err!=nil{
		fmt.Println("get create file err",err)
	}
	file := strings.NewReader(buf.String())
	file.WriteTo(text)
}


//解析文章
func parse(n *html.Node)  {
	if n.Type == html.ElementNode && n.Data == "div"{
		for _,a := range n.Attr{
			if a.Key == "id" && a.Val == "content" {
				find = true
				parseTxt(&buf,n)
				break
			}
		}
	}
	if !find{
		for c := n.FirstChild;c!=nil;c=c.NextSibling{
			parse(c)
		}
	}
}
//提取文字
func parseTxt(buf *bytes.Buffer,n *html.Node)  {
	for c:=n.FirstChild;c!=nil;c=c.NextSibling{
		if c.Data != "br"{
			buf.WriteString(c.Data+"\n")
		}
	}
}
