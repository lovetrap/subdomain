package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

//Get return body
func Get(client *http.Client, url string) (body []byte, err error) {
	rsp, err := client.Get(url)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%v", rsp.StatusCode)
		return
	}
	return ioutil.ReadAll(rsp.Body)
}

func NewClient(proxyAddr string) *http.Client {
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return nil
	}
	netTransport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		Dial: func(network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, time.Second*time.Duration(10))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

/* 获取API返回内容 */
func GetApi(api string,accept string,Type string)(respose string,error error){

	if accept == ""{
		accept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	}

	if Type == "" {
		Type = "text/html;charset=utf-8"
	}

	client := &http.Client{Timeout:time.Second * 5}
	req, _ := http.NewRequest("GET",api,nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36")
	req.Header.Add("Accept-Language","zh-CN,zh;q=0.9")
	req.Header.Add("Accept",accept)
	req.Header.Add("Content-Type", Type)
	req.Header.Add("Referer", "https://www.google.com/")
	req.Header.Add("Connection","keep-alive")

	rep, err := client.Do(req)
	if err != nil{
		error = err
	}else {
		body, _ := ioutil.ReadAll(rep.Body)
		respose = string(body)
		defer rep.Body.Close()
	}

	return respose,err
}

/* 读取字典文件函数 */
func Readfile(filename string) (str []string)  {
	fi, _ := os.Open(filename)
	defer fi.Close()
	b, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil
	}
	length := bytes.Count(b, []byte("\n"))
	br := bufio.NewReader(fi)

	for i := 0; i < length; i++ {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		str = append(str, string(a))
	}
	return str
}