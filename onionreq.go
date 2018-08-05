package requests

import (
	"bytes"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func fatalf(fmtStr string, args interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args)
	os.Exit(-1)
}

// ReqHandler unsurprisingly handles requests
type ReqHandler struct {
	Header http.Header
	client http.Client
}

// Constructor
func NewHandler(headers http.Header) (*ReqHandler, error) {
	tbProxyURL, err := url.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		return nil, err
	}
	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{Dial: tbDialer.Dial}
	client := http.Client{Transport: tr}
	r := &ReqHandler{Header: headers, client: client}
	return r, nil
}

// Do runs the requests and returns the response object
func (r *ReqHandler) Do(method, url string, data *string, header *http.Header) (*http.Response, error) {
	var buf io.Reader
	if data != nil {
		buf = bytes.NewBufferString(*data)
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header = *header
	} else {
		req.Header = r.Header
	}
	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *ReqHandler) Request(method, url string, data *string, header *http.Header) ([]byte, error) {
	res, err := r.Do(method, url, data, header)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return content, nil
}
