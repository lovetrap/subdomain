package pkg

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

//BaiduDomain 百度CE API
func (self self) BaiduDomain(wg *sync.WaitGroup) {
	defer wg.Done()
	api := "http://ce.baidu.com/index/getRelatedSites?site_address=" + self.Host
	rep, err := http.Get(api)

	if err != nil {
		return
		//fmt.Printf("[-] %s BaiduAPI请求超时！！\n", self.Host)
	} else {
		body, _ := ioutil.ReadAll(rep.Body)

		res := struct{
			Data []struct {
				Domain string `json:"domain"`
			} `json:"data"`
		}{}

		_ = json.Unmarshal(body, &res)
		for _, i := range res.Data {
			domain[i.Domain] = struct {}{}
		}

	}
}

//CershAPI CershAPI
func (self self) CershAPI(wg *sync.WaitGroup) {
	defer wg.Done()
	api := "https://crt.sh/?q=" + self.Host
	DNS, _ := regexp.Compile(`DNS:(.*?)<BR>`)
	urls := make(map[string]struct{})

	client := &http.Client{Timeout:time.Second * 5}
	c := colly.NewCollector()
	c.SetClient(client)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		HrefURL := e.Request.AbsoluteURL(e.Attr("href"))
		if !strings.Contains(HrefURL, "caid") && !strings.Contains(HrefURL, self.Host) {
			if HrefURL != "https://sectigo.com/" && HrefURL != "https://github.com/crtsh" && HrefURL != "https://crt.sh/" {
				if _, ok := urls[HrefURL]; !ok {
					if HrefURL != "" {
						urls[HrefURL] = struct{}{}
					}
				}
			}
		}
	})

	_ = c.Visit(api)

	if len(urls) == 0 {
		//fmt.Printf("[-] %s CershAPI请求失败！！\n", Url.URL)
		return
	} else {
		for v := range urls {
			body, err := Get(client,v)
			if err != nil {
				return
				//fmt.Printf("[-] %s CershAPI请求Domain失败！！\n", self.Host)
				break
			} else {
				for _, link := range DNS.FindAllString(string(body), -1) {
					link = strings.Replace(link, "DNS:", "", -1)
					link = strings.Replace(link, "<BR>", "", -1)
					domain[link] = struct {}{}
				}
			}
			time.Sleep( time.Second * 2 )
		}
	}

}

//XimcxAPI XimcxAPI
func (self self) XimcxAPI(wg *sync.WaitGroup) {
	defer wg.Done()
	api := "http://sbd.ximcx.cn/DomainServlet"
	rep, err := http.PostForm(api, url.Values{
		"domain": {self.Host},
	})
	if err != nil {
		return
		//fmt.Printf("[-] %s ximcxAPI请求失败！！\n", self.Host)
	} else {
		body, _ := ioutil.ReadAll(rep.Body)
		defer rep.Body.Close()

		res := struct{
			Data []struct {
				Domain string `json:"domain"`
			} `json:"data"`
		}{}

		_ = json.Unmarshal(body, &res)
		for _, k := range res.Data {
			link := k.Domain
			domain[link] = struct {}{}

		}

	}

}

//EntrusAPI EntrusAPI
//func (Url *URL) EntrusAPI(wg *sync.WaitGroup) {
//	defer wg.Done()
//	api := fmt.Sprintf("https://ctsearch.entrust.com/api/v1/certificates?fields=subjectDN&domain=%s&includeExpired=true", Url.URL)
//	re, _ := regexp.Compile(`(\w+\.)+\w+`)
//	client := pkg.NewClient(config.C.Proxy)
//	body, err := function.Get(client,api)
//	if err != nil {
//		fmt.Printf("[-] %s EntrusAPI请求失败！！\n", Url.URL)
//	} else {
//
//		var res []struct {
//			SubjectDN string `json:"subjectDN"`
//		}
//
//		_ = json.Unmarshal(body, &res)
//		for _, k := range res {
//
//			link := strings.Replace(strings.Replace(re.FindAllString(k.SubjectDN, -1)[0], "www.itrus.com.cn", "", -1), "www.verisign.com", "", -1)
//			if function.CheckArray(Url.Domain, link) {
//				Url.Domain = append(Url.Domain, link)
//			}
//		}
//
//	}
//}

//CertspotterAPI certspotterAPI
//func (Url *URL) CertspotterAPI(wg *sync.WaitGroup) {
//	defer wg.Done()
//	api := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names", Url.URL)
//	client := pkg.NewClient(config.C.Proxy)
//	body, err := function.Get(client, api)
//	if err != nil {
//		fmt.Printf("[-] %s CertspotterAPI请求失败！！\n", Url.URL)
//	} else {
//
//		var res []struct{
//			DNSNames []string `json:"dns_names"`
//		}
//
//		_ = json.Unmarshal(body, &res)
//		for _, k := range res {
//			for _, link := range k.DNSNames {
//				if function.CheckArray(Url.Domain, link) {
//					Url.Domain = append(Url.Domain, link)
//				}
//			}
//		}
//	}
//}

////HackertAPI HackertAPI
//func (self self) HackertAPI(wg *sync.WaitGroup) {
//	defer wg.Done()
//	api := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", self.Host)
//	body, err := pkg.GetApi(api, "", "")
//	if err != nil {
//		fmt.Printf("[-] %s HackertAPI请求失败！！\n", self.Host)
//	} else {
//		for _, k := range strings.Split(body, "\n") {
//			link := strings.Split(k, ",")[0]
//			domain[link] = struct {}{}
//
//		}
//
//	}
//}