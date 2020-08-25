package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var domain map[string]interface{}
var (
	black = Readfile("wordlist/Black_Url.list")
	con = Readfile("wordlist/Black_Con.list")
)

func (self self)Start() {
	domain = make(map[string]interface{})

	var wg sync.WaitGroup
	wg.Add(4)
	go self.BaiduDomain(&wg)
	go self.CershAPI(&wg)
	go self.XimcxAPI(&wg)
	go self.spyseAPI(&wg)
	go self.Subdomain(&wg)
	wg.Wait()

	for d := range domain {
		d := d
		wg.Add(1)
		go func() {
			defer wg.Done()
			b := func() bool{
				for _, burl := range black {
					if strings.Contains(d, burl) {
						return false
					}
				}
				return true
			}()
			if b {
				resp ,err := http.Get("http://" + d)
				if err != nil {
					resp, err = http.Get("https://" + d)
					if err == nil {
						body, _ := ioutil.ReadAll(resp.Body)
						x := func() bool{
							for _, c := range con {
								if strings.Contains(string(body), c) {
									return false
								}
							}
							return true
						}()
						if x {
							title := GetTitle(body)
							fmt.Println("http://" + d, title)
						}
					}
				}else {
					body, _ := ioutil.ReadAll(resp.Body)
					x := func() bool{
						for _, c := range con {
							if strings.Contains(string(body), c) {
								return false
							}
						}
						return true
					}()
					if x {
						title := GetTitle(body)
						fmt.Println("https://" + d, title)
					}
				}
			}
		}()
	}
	wg.Wait()

}

