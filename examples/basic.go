package main

import (
	"fmt"
	req "github.com/adrpino/go-requests"
	"io/ioutil"
)

func main() {
	h := req.NewHandler()
	res, err := h.Request("https://icanhazheaders.com")
	fmt.Println(string(res))

}
