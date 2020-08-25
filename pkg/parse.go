package pkg

import (
	"net/url"
	"regexp"
	"strings"
)

type self struct {
	Scheme string
	Host string
}

func Parse(domain string) *self{
	if regexp.MustCompile(`^http`).FindString(domain) == "" {
		domain = "http://" + domain
	}
	URL, _ := url.Parse(domain)
	if URL.Scheme == "" {
		URL.Scheme = "http"
	}
	host := strings.Split(URL.Host, ".")
	p := host[len(host) - 2:]
	temp := p[0] + "." + p[1]
	if temp == "com.cn" {
		p = host[len(host) - 3:]
		temp = p[0] + "." + p[1] + "." + p[2]
	}
	return &self{URL.Scheme, temp}
}