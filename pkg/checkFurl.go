package pkg

import (
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"regexp"
)

func GetContext(URL string)([]byte){

	rep	,err := http.Get(URL)
	if err != nil {
		return nil
	}
	body ,err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil
	}
	return body
}

func GetTitle(body []byte)string{
	Title := regexp.MustCompile(`<title>(.*?)</title>`)
	title := Title.FindAllStringSubmatch(string(body),-1)
	if len(title) < 1 {
		return ""
	}
	return title[0][1]

}

func checkFURL(Schema string, Url string)(v bool){
	v = true
	a,b := GetContext(Schema + "://bjbjbjbjbjbjbjbjbj." + Url), GetContext(Schema + "://shshshshshshshshsh." + Url)
	p,d := 0,0
	var min int
	if a != nil && b != nil {
		if len(a) >= len(b) {
			min = len(b) - 1
		}else {
			min = len(a) - 1
		}
		if GetTitle(a) == GetTitle(b) {
			for range a {
				if p <= min {
					if a[p] == b[p]{
						d++
					}
				}
				p++
			}
			s, _ := decimal.NewFromFloat(float64(d)).Div(decimal.NewFromFloat(float64(len(a)))).Float64()
			if s * 100 > 75 {
				v = false
			}
		}
	}
	return
}

//检测是否泛解析
func CheckFURL(Url string)(v bool){
	v = false
	if checkFURL("http", Url) && checkFURL("https", Url) {
		v = true
	}
	return
}
