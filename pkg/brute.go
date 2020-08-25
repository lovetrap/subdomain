package pkg

import (
	"fmt"
	"sync"
)

var (
	DomainDict = Readfile("wordlist/SubDomainDict.list")
	//NextDict = Readfile("wordlist/NextSubDomainDict.list")
)

func (self self)Subdomain(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, i := range DomainDict {
		link := fmt.Sprintf("%s.%s" ,i ,self.Host)
		domain[link] = struct {}{}
	}
}
