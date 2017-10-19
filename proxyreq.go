package requests

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "golang.org/x/net/proxy"
	"math/rand"
	"net/http"
	"net/url"
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
	SetPool() error
}

type ProxyList struct {
	Pool         []*Proxy
	Source       string
	client       *http.Client
	currentProxy *Proxy
	currentInd   int
}

func (pl *ProxyList) Init() error {
	pl.Source = "https://proxy-list.org/english/index.php"
	err := pl.SetPool()
	if err != nil {
		return err
	}
	err = pl.SetProxy()
	if err != nil {
		return err
	}
	return nil
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
			if pr == nil {
				continue
			}
			ps := string(pr)
			numDots := len(strings.Split(ps, "."))
			if numDots == 4 {
				pool = append(pool, NewProxy(ps))
				continue
			}
		}

	})
	pl.Pool = pool
	return nil
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
	pr := pl.Pool[ind]
	proxyStr := fmt.Sprintf("http://%s", pr)
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		fmt.Println("Error parsing proxy")
		return err
	}
	if err != nil {
		fmt.Println("error setting proxy")
		return err
	}
	Transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: Transport}
	pl.client = client
	//	pl.currentProxy = NewProxy(proxyStr)
	pl.currentProxy = pl.Pool[ind]
	pl.currentInd = ind
	return nil
}

// Returns currently set proxy
func (pl *ProxyList) CurrentProxy() (*Proxy, error) {
	pr := pl.currentProxy
	fmt.Println(pr)
	if pr.p == "" {
		err := errors.New("currentProxy is unset")
		return nil, err
	}
	return pr, nil
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

func (pl *ProxyList) Get(url string) (*http.Response, error) {
	res, err := pl.client.Get(url)
	//	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return res, nil
}
