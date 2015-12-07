package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
)

var (
	url = `http://www.thepaper.cn/load_index.jsp?nodeids=25990&topCids=&pageidx=`
)

func DownloadPage() (str string, err error) {
	req, err := http.Get(fmt.Sprint(url, 1))
	if err != nil {
		fmt.Println("download error, ", err)
		return
	}
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("get body error, ", err)
		return
	}

	return string(bs), nil
}
