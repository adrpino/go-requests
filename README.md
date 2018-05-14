# go-requests

Simple library to do HTTP requests routed through public proxies and Tor.

## Installation

```bash
go get github.com/adrpino/go-requests
```


### Example usage
If you have Tor running:

```go
package main

import (
	"fmt"
	req "github.com/adrpino/go-requests"
	"io/ioutil"
)

func main() {
	res, _ := req.OnionRequest("https://icanhazip.com")
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	// Panicking is for cowards
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

}
```
