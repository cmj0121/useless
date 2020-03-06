/* Copyright (C) 2020-2020 cmj. All right reserved. */
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type ETF struct {
	*TWSE `yaml:"-"`

	ListingDate string `yaml:"listing_date"`
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Issuer      string `yaml:"issuer"`
	Link        string `yaml:"link"`
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

type TWSECache struct {
	List map[string][]*ETF	`yaml:"etf"`
}

func LoadCache(file string) (out *TWSECache, err error) {
	var data []byte

	out = &TWSECache {
		List: make(map[string][]*ETF, 0),
	}

	if data, err = ioutil.ReadFile(file); err != nil {
		/* Cannot read file */
		return
	}
	err = yaml.Unmarshal(data, &out)
	return
}

func (cache *TWSECache) StoreCache(file string) (err error) {
	var data []byte

	if data, err = yaml.Marshal(&cache); err != nil {
		/* Cannot marshal the structure to yaml */
		return
	}

	err = ioutil.WriteFile(file, data, 0600)
	return
}

type TWSE struct {
	*http.Client
	*log.Logger
	*TWSECache

	CacheFile string
	Force bool
}

func MustNew() (twse *TWSE) {
	twse = &TWSE{
		Client: &http.Client{},
		Logger: log.New(os.Stderr, "", log.Lshortfile|log.LUTC),

		Force: false,
	}

	usr, err := user.Current()
	if err != nil {
		twse.Panicf("Cannot get current user")
		return
	}
	/* always load the configuration */
	twse.CacheFile = fmt.Sprintf("%s/.twse", usr.HomeDir)
	twse.TWSECache, _ = LoadCache(twse.CacheFile)
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

	if data, ok := twse.TWSECache.List[filter]; ok && ! twse.Force && len(data) != 0 {
		list = data
		return
	}

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

	twse.TWSECache.List[filter] = list
	twse.TWSECache.StoreCache(twse.CacheFile)
	return
}

func (twse *TWSE) Parse(name string, args... string) {
	actions := []string{}

	for idx := 0; idx < len(args); idx ++ {
		switch args[idx] {
		case "-f":
			if idx ++; idx >= len(args) {
				twse.Usage(name, "-f option need file path")
			}
		case "-F", "--force":
			twse.Force = true
			break
		default:
			actions = append(actions, args[idx])
		}
	}

	if len(actions) == 0 {
		twse.Usage(name)
		return
	}

	switch actions[0] {
		case "type":
			result := twse.Types()
			fmt.Printf("%s\n", strings.Join(result, " "))
			break
		case "list":
			var filter string

			if filter = ""; len(actions) > 1 {
				filter = actions[1]
			}

			for _, etf := range twse.List(filter) {
				fmt.Printf("%s\n", etf)
			}
			break
		default:
			twse.Usage(name)
			break
	}
}

func (twse *TWSE) Usage(name string, err ...interface{}) {
	/* disable the debug log format */
	twse.SetFlags(0)
	if len(err) >0 {
		for _, e := range err {
			twse.Printf("error: %v\n", e)
		}
	}

	twse.Printf("Usage: %s [OPTION] ACTION\n", name)
	twse.Printf("\n")
	twse.Printf("ACTION\n")
	twse.Printf("    help           Show this message\n")
	twse.Printf("    type           List all type of ETF\n")
	twse.Printf("    list           Get the ETF ID and Name\n")
	twse.Printf("OPTION\n")
	twse.Printf("    -f FILE        The cache file (default: ~/.twetf)\n")
	twse.Printf("    -F, --force    Force fetch\n")
	os.Exit(1)
}

func main() {
	twse := MustNew()
	defer func() {
		if r := recover(); r != nil {
			twse.Usage(os.Args[0], r)
		}
	}()

	twse.Parse(os.Args[0], os.Args[1:]...)
}
