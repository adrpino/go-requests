package requests

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func fatalf(fmtStr string, args interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args)
	os.Exit(-1)
}

func OnionRequest(u string) (*http.Response, error) {
	// Create a transport that uses Tor Browser's SocksPort.  If
	// talking to a system tor, this may be an AF_UNIX socket, or
	// 127.0.0.1:9050 instead.
	tbProxyURL, err := url.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		fatalf("Failed to parse proxy URL: %v\n", err)
	}

	// Get a proxy Dialer that will create the connection on our
	// behalf via the SOCKS5 proxy.  Specify the authentication
	// and re-create the dialer/transport/client if tor's
	// IsolateSOCKSAuth is needed.
	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		fatalf("Failed to obtain proxy dialer: %v\n", err)
	}

	// Make a http.Transport that uses the proxy dialer, and a
	// http.Client that uses the transport.
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client := &http.Client{Transport: tbTransport}

	// Example: Fetch something.  Real code will probably want to use
	// client.Do() so they can change the User-Agent.
	resp, err := client.Get(u)
	return resp, err
}

// ReqHandler unsurprisingly handles requests
type ReqHandler struct {
	client    *http.Client
	transport *http.Transport
	headers   http.Header
}

// Constructor
func NewHandler() ReqHandler {
	userAgent := []string{"Mozilla/5.0 (Windows NT 6.1; rv:45.0) Gecko/20100101 Firefox/45.0"}
	encoding := []string{"gzip, deflate"}
	headers := http.Header{
		"User-Agent":      userAgent,
		"Accept-Encoding": encoding,
	}
	tbProxyURL, err := url.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		fatalf("Failed to parse proxy URL: %v\n", err)
	}
	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		fatalf("Failed to obtain proxy dialer: %v\n", err)
	}
	tr := &http.Transport{Dial: tbDialer.Dial}
	client := &http.Client{Transport: tr}
	r := ReqHandler{transport: tr, headers: headers, client: client}
	return r
}

// Runs a request and returns the resulting bytes
// TODO separate in case the the response object is required
func (r *ReqHandler) Request(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = r.headers
	res, err := r.client.Do(req)
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
