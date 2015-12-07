package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/redis.v3"
)

const (
	R_PP_CARD         = "pp_card_"
	R_PP_DATE         = "date"
	R_PP_LINK         = "href"
	R_PP_IMG_URL      = "img_url"
	R_PP_TITLE        = "title"
	R_PP_SUMMARY      = "summary"
	R_PP_CONTENT      = "content"
	R_PP_CONTENT_IMGS = "content_imgs"
)

var (
	cli *redis.Client
)

func main() {
	var err error
	cli = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if ping := cli.Ping(); ping.Err() != nil {
		log.Fatal(err)
	}

	log.Println("start...")
	for i := 0; i < 13; i++ {
		DownloadCards(i)
	}
}

func DownloadCards(page int) {
	url := `http://www.thepaper.cn/load_index.jsp?nodeids=25990&topCids=&pageidx=`
	doc, err := goquery.NewDocument(fmt.Sprint(url, page))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("downlaod over.")
	doc.Find(".news_li").Each(func(i int, s1 *goquery.Selection) { // 1
		log.Println("Article: ", i)

		var key string
		date := s1.Find(".pdtt_trbs span").Text()
		if date != "" {
			fmt.Println("date: ", FormatDate(date))
			if FormatDate(date).Before(time.Date(2015, 0, 0, 0, 0, 0, 0, time.Local)) {
				return
			}
		}

		// card link
		link, ok := s1.Find("h2 a").Attr("href")
		if ok {
			log.Println("link href: ", link)
			key = R_PP_CARD + link
			cli.HSet(key, R_PP_LINK, link)
			cli.HSet(key, R_PP_DATE, FormatDate(date).Format("2006-01-2"))
		} else {
			return
		}

		// img url
		imgUrl, ok := s1.Find(".news_tu img").Attr("src")
		if ok {
			log.Println("img url: ", imgUrl)
			cli.HSet(key, R_PP_IMG_URL, imgUrl)
		}

		// cart title
		cardTitle := s1.Find("h2 a").Text()
		if cardTitle != "" {
			log.Println("card title: ", cardTitle)
			cli.HSet(key, R_PP_TITLE, cardTitle)
		}

		// summary
		summary := s1.Find("p").Text()
		if summary != "" {
			log.Println("summary: ", summary)
			cli.HSet(key, R_PP_SUMMARY, summary)
		}
		DownloadContent(key, link)
	})
}

func FormatDate(str string) time.Time {
	if len(str) >= 10 {
		if t, err := time.Parse("2006-01-2", str[:10]); err == nil {
			return t
		}
	}
	if h := strings.LastIndex(str, "小时前"); h > 0 {
		if hour, err := strconv.Atoi(str[:h]); err == nil {
			return time.Now().Add(time.Duration(hour*-1) * time.Hour)
		}
	}
	if d := strings.LastIndex(str, "天前"); d > 0 {
		if day, err := strconv.Atoi(str[:d]); err == nil {
			return time.Now().AddDate(0, 0, -1*day)
		}
	}
	log.Println("can not decode date str: ", str)
	return time.Time{}
}

func DownloadContent(key, href string) {
	url := "http://www.thepaper.cn/" + href
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("downlaod content over.")
	doc.Find(".news_txt").Each(func(i int, s1 *goquery.Selection) {
		txt := s1.Text()
		imgs := make([]string, 0)
		s1.Find("img").Each(func(j int, s2 *goquery.Selection) {
			if src, ok := s2.Attr("src"); ok {
				imgs = append(imgs, src)
			}
		})
		cli.HSet(key, R_PP_CONTENT, txt)
		cli.HSet(key, R_PP_CONTENT_IMGS, strings.Join(imgs, ","))
		log.Println("get content over...")
	})
}
