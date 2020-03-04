/* Copyright (C) 2020-2020 cmj. All right reserved. */
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type ETF struct {
	*TWSE

	ListingDate string
	ID          string
	Name        string
	Issuer      string
	Link        string
}

func MustNewETF(id string) (etf *ETF) {
	etf = &ETF{
		TWSE: MustNew(),
		ID:   id,
	}
	etf.Reload()
	return
}

func (etf *ETF) String() (out string) {
	out = fmt.Sprintf("%-8s %s", etf.ID, etf.Name)
	return
}

func (etf *ETF) Reload() {
	URL := fmt.Sprintf("https://www.twse.com.tw/zh/ETF/fund/%s", etf.ID)
	resp, err := etf.Get(URL)
	if err != nil {
		etf.Panicf("Cannot get %s: %s", URL, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		etf.Panicf("Cannot read HTML from %s: %s", URL, err)
		return
	}

	html := string(body)
	etf.ListingDate = etf.parseHTML(`<td>上市日期</td>\s+<td>(.*?)</td>`, html)
	etf.Name = etf.parseHTML(`<td>ETF簡稱</td>\s+<td>(.*?)</td>`, html)
	etf.Issuer = etf.parseHTML(`<td>基金經理公司</td>\s+<td>(.*?)</td>`, html)
	etf.Link = etf.parseHTML(`<a href="(.*?)" target="_blank">申購買回清單PCF`, html)
}

func (etf *ETF) Constituent() (out []string) {
	resp, err := etf.Get(etf.Link)
	if err != nil {
		etf.Panicf("Cannot get %s: %s", etf.Link, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		etf.Panicf("Cannot read HTML from %s: %s", etf.Link, err)
		return
	}

	html := string(body)

	etf.Printf("%s", html)
	return
}

func (etf *ETF) parseHTML(pattern, html string) (out string) {
	re := regexp.MustCompile(pattern)
	ret := re.FindStringSubmatch(html)

	if len(ret) < 1 {
		etf.Panicf("Cannot find pattern %s", pattern)
		return
	}

	out = ret[1]
	return
}

type TWSE struct {
	*http.Client
	*log.Logger
}

func MustNew() (twse *TWSE) {
	twse = &TWSE{
		Client: &http.Client{},
		Logger: log.New(os.Stderr, "", log.Lshortfile|log.LUTC),
	}
	return
}

func (twse *TWSE) Types() (out []string) {
	out = []string{
		"domestic",
		"foreign",
		"offshore",
		"LI", "lever", "inverse",
		"futures",
		"LIfuteres",
	}

	return
}

func (twse *TWSE) List(filter string) (list []*ETF) {
	URL := ""
	offset := 0

	switch filter {
	case "domestic":
		URL = "https://www.twse.com.tw/zh/page/ETF/domestic_etf.html"
		break
	case "foreign":
		URL = "https://www.twse.com.tw/zh/page/ETF/foreign_etf.html"
		break
	case "offshore":
		URL = "https://www.twse.com.tw/zh/page/ETF/offshore_etf.html"
		break
	case "LI", "lever", "inverse":
		URL = "https://www.twse.com.tw/zh/page/ETF/LI_etf.html"
		break
	case "futures":
		URL = "https://www.twse.com.tw/zh/page/ETF/vanillafutures_etf.html"
		break
	case "LIfuteres":
		URL = "https://www.twse.com.tw/zh/page/ETF/LIfutures_etf.html"
		break
	default:
        URL = "https://www.twse.com.tw/zh/page/ETF/list.html"
		offset = 1
        break
	}

	resp, err := twse.Get(URL)
	if err != nil {
		twse.Panicf("Cannot get %s: %s", URL, err)
		return
	}

	html, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		twse.Panicf("Cannot read HTML from %s: %s", URL, err)
		return
	}

	re_table := regexp.MustCompile(`<tr>([\s\S]+?)</tr>`)
	re_column := regexp.MustCompile(`<td>(?:<.*?>)?(.*?)<`)
	for _, table := range re_table.FindAllSubmatch(html, -1) {
		columns := re_column.FindAllSubmatch(table[1], -1)

		if len(columns) < 2 {
			continue
		}

		list = append(list, &ETF{
			TWSE:   twse,
			ID:     string(columns[offset+0][1]),
			Name:   string(columns[offset+1][1]),
		})
	}
	return
}

func usage() {
	fmt.Printf("Usage: %s ACTION [OPTION]\n", os.Args[0])
	fmt.Printf("\n")
	fmt.Printf("    help    Show this message\n")
	fmt.Printf("    type    List all type of ETF\n")
	fmt.Printf("    list    Get the ETF ID and Name\n")
	os.Exit(1)
}

func main() {
	var action string
	var option string

	if action = ""; len(os.Args) > 1 {
		action = os.Args[1]
	}
	if option = ""; len(os.Args) > 2 {
		option = os.Args[2]
	}

	switch action {
	case "", "help":
		usage()
	case "type":
		twse := MustNew()
		fmt.Printf("%s\n", strings.Join(twse.Types(), ", "))
		break
	case "list":
		twse := MustNew()
		for _, etf := range twse.List(option) {
			fmt.Printf("%v\n", etf)
		}
		break
	default:
		fmt.Printf("Unknown action %s\n", action)
		break
	}

	return
}
