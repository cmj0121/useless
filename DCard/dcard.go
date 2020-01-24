/* Copyright (C) 2020-2020 cmj. All right reserved. */
package dcard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	VERSION = "1.0.0"
	API     = "https://www.dcard.tw/_api"
)

type DCard struct {
	timeout time.Duration
	limit   int
	popular bool
}

func New() (out *DCard) {
	out = &DCard{
		timeout: time.Second * 4,
		limit:   20,
		popular: false,
	}
	return
}

func (d *DCard) Limit(limit int) (out *DCard) {
	out = d
	out.limit = limit
	return
}

func (d *DCard) Popular(popular bool) (out *DCard) {
	out = d
	out.popular = popular
	return
}

func (d *DCard) fetch(uri string) []byte {
	cli := &http.Client{
		Timeout: d.timeout,
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		err = fmt.Errorf("DCard fetch failure - %s", err)
		panic(err)
	}

	res, err := cli.Do(req)
	if err != nil {
		err = fmt.Errorf("DCard fetch failure - %s", err)
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("DCard fetch failure - %s", err)
		panic(err)
	}

	return body
}

func (d *DCard) Posts(forum string, before int64) (out []DCardPost) {
	URL := fmt.Sprintf("%s/forums/%s/posts?popular=%v&limit=%d", API, forum, d.popular, d.limit)
	if before > 0 {
		URL = fmt.Sprintf("%s&before=%d", URL, before)
	}
	data := d.fetch(URL)

	if err := json.Unmarshal(data, &out); err != nil {
		err = fmt.Errorf("Parse json failure - %s", err)
		panic(err)
	}

	return
}
