package requests

import (
	"encoding/base64"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	_ "net/url"
	"strings"
)

type Proxy struct {
	p string
}

// Constructor
func NewProxy(ps string) *Proxy {
	return &Proxy{ps}
}

func (p *Proxy) String() string {
	return string(p.p)
}

type ProxyPool interface {
	SetProxy()
	SetPool()
	Proxies()
	CurrentProxy()
	DropProxy()
	NumProxies()
}

type ProxyList struct {
	Pool         []*Proxy
	Source       string
	currentProxy *Proxy
	currentInd   int
}

func (pl *ProxyList) Init() error {
	pl.Source = "https://proxy-list.org/english/index.php"
	err := pl.SetPool()
	pl.SetProxy()
	return err
}

func (pl *ProxyList) SetPool() error {
	res, err := OnionRequest(pl.Source)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		panic(err)
	}
	var pool = make([]*Proxy, 10)
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		elem := s.Find("li").Find("script").Text()
		encoded := strings.Split(elem, "'")
		for _, enc := range encoded {
			pr, err := base64.StdEncoding.DecodeString(enc)
			if err != nil {
				continue
			}
			//			fmt.Println(string(pr))
			ps := string(pr)
			pool = append(pool, NewProxy(ps))
		}

	})
	pl.Pool = pool
	return err
}

func (pl *ProxyList) NumProxies() int {
	return len(pl.Pool)
}

func (pl *ProxyList) SetProxy() error {
	n := pl.NumProxies()
	if n == 0 {
		err := errors.New("Empty pool of proxies")
		return err
	}
	ind := rand.Intn(n)
	pl.currentProxy = pl.Pool[ind]
	pl.currentInd = ind
	return nil
}

// Returns currently set proxy
func (pl *ProxyList) CurrentProxy() (*Proxy, error) {
	proxy := pl.currentProxy
	if proxy.p == "" {
		err := errors.New("currentProxy is unset")
		return nil, err
	}
	return proxy, nil
}

func (pl *ProxyList) DeleteProxy() error {
	if pl.currentProxy.p == "" {
		err := errors.New("currentProxy is unset")
		return err
	}
	ind := pl.currentInd
	pl.Pool = append(pl.Pool[:ind], pl.Pool[ind+1:]...)
	return nil
}
