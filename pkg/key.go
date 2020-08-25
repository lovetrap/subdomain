package pkg

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
	"tool/cralw/pkg/config"
)

//type censys struct {
//	Status  string `json:"status"`
//	Results []struct {
//		ParsedNames     []string `json:"parsed.names"`
//		ParsedSubjectDn string   `json:"parsed.subject_dn"`
//	} `json:"results"`
//}

//Censys CensysAPI
//func (self self) Censys(wg *sync.WaitGroup) {
//	defer wg.Done()
//	api := "https://www.censys.io/api/v1/search/certificates"
//	page := 1
//	for {
//		data := fmt.Sprintf(`{
//	  "query":"parsed.names: %s",
//	  "page":%d,
//	  "fields":["parsed.subject_dn", "parsed.names"],
//	  "flatten":true
//	}`,self.Host,page)
//		client := pkg.NewClient(config.C.Proxy)
//		req, _ := http.NewRequest("POST", api, strings.NewReader(data))
//		req.SetBasicAuth(config.C.API.CensysAPIID, config.C.API.CensysAPISecret)
//		rep, err := client.Do(req)
//
//		if err != nil {
//			fmt.Printf("[-] %s censysAPI请求失败！！\n", self.Host)
//		} else {
//			body, _ := ioutil.ReadAll(rep.Body)
//			res := &censys{}
//			_ = json.Unmarshal(body,&res)
//			if res.Status != "ok" {
//				break
//			}else {
//				page++
//			}
//		}
//		time.Sleep(time.Second * 3)
//	}
//
//}

type Shodan struct {
	Domain string   `json:"domain"`
	Tags   []string `json:"tags"`
	Data   []struct {
		Subdomain string    `json:"subdomain"`
		Type      string    `json:"type"`
		Value     string    `json:"value"`
		LastSeen  time.Time `json:"last_seen"`
	} `json:"data"`
	Subdomains []string `json:"subdomains"`
	More       bool     `json:"more"`
}

//Shodan ShodanAPI
func (self self) Shodan(wg *sync.WaitGroup) {
	defer wg.Done()
	api := "https://api.shodan.io/dns/domain"
	//page := 1
	time.Sleep(time.Second * 1)
	URL := fmt.Sprintf("%s/%s?key=%s", api, self.Host, config.C.API.ShodanKey)
	client := NewClient(config.C.Proxy)
	body, err := Get(client, URL)
	if err != nil {
		return
		//fmt.Printf("[-] %s ShodanAPI请求失败！！\n", self.Host)
	} else {
		res := &Shodan{}
		_ = json.Unmarshal(body,&res)
		for _, k := range  res.Subdomains {
			link := k + "." + self.Host
			domain[link] = struct {}{}
		}

	}

}

//secAPI secAPI
func (self self) secAPI(wg *sync.WaitGroup) {
	defer wg.Done()
	api := fmt.Sprintf("https://api.securitytrails.com/v1/domain/%s/subdomains?apikey=%s", self.Host, config.C.API.Seckey)
	body, err := GetApi(api, "", "")
	if err != nil {
		return
		//fmt.Printf("[-] %s 请求securitytrails_API失败！！\n", self.Host)
	} else {
		res := struct {
			Subdomains []string `json:"subdomains"`
		}{}
		_ = json.Unmarshal([]byte(body), &res)
		for _, k := range res.Subdomains {

			link := k + "." + self.Host
			domain[link] = struct {}{}
		}

	}
}

type virustotal struct {
	Subdomains []string `json:"subdomains"`
}

//virustotalAPI virustotalAPI
func (self self) virustotalAPI(wg *sync.WaitGroup) {
	defer wg.Done()
	api := fmt.Sprintf("https://www.virustotal.com/vtapi/v2/domain/report?apikey=%s&domain=%s", config.C.API.VirustotalKey, self.Host)

	body, err := GetApi(api, "", "")
	if err != nil {
		return
		//fmt.Printf("[-] %s VirustotalAPI请求失败！！\n", self.Host)
	} else {
		res := &virustotal{}
		_ = json.Unmarshal([]byte(body), &res)
		for _, link := range res.Subdomains {
			domain[link] = struct {}{}
		}
	}
}



type spyse struct {
	Records []struct {
		Domain string `json:"domain"`
	} `json:"records"`
}

//spyseAPI spyseAPI
func (self self) spyseAPI(wg *sync.WaitGroup) {
	defer wg.Done()
	page := 1
	res := &spyse{}
	for {
		api := fmt.Sprintf("https://api.spyse.com/v1/subdomains?domain=%s&api_tpken=%s&page=%d", self.Host, config.C.API.SpyseToken, page)
		body, err := GetApi(api, "", "")
		if err != nil {
			return
			//fmt.Printf("[-] %s SpyseAPI请求失败！！\n", self.Host)
			break
		} else {
			if strings.Contains(body, self.Host) {
				page++
				_ = json.Unmarshal([]byte(body), &res)
				for _, k := range res.Records {
					domain[k.Domain] = struct {}{}

				}
			} else {
				break
			}
		}
		time.Sleep(time.Second * 3)
	}

}