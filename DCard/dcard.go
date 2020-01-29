/* Copyright (C) 2020-2020 cmj. All right reserved. */
package dcard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

const (
	VERSION = "1.0.0"
	HOST    = "https://www.dcard.tw"
	API     = "_api"
)

var (
	default_agent = New()
)

type DCard struct {
	timeout time.Duration
	limit   int
	popular bool
}

func New() (out *DCard) {
	out = &DCard{
		timeout: time.Second * 40,
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

func (d *DCard) Boards() (out []DCardBoard) {
	data := d.fetch(HOST)
	re := regexp.MustCompile(`<script>([\s\S]+?)</script>`)

	type DCardConf struct {
		Forum struct {
			Stores []DCardBoard `json:"store"`
		} `json:"forums"`
	}
	conf := DCardConf{}

	for _, matched := range re.FindAllSubmatch(data, -1) {
		if bytes.Compare(matched[1][:14], []byte("window.$STATE=")) == 0 {
			if err := json.Unmarshal(matched[1][14:], &conf); err != nil {
				err = fmt.Errorf("Parse json failure - %s", err)
				panic(err)
			}

			out = conf.Forum.Stores
			break
		}
	}
	return
}

func (d *DCard) Posts(forum string, before int64) (out []DCardPost) {
	URL := fmt.Sprintf("%s/%s/forums/%s/posts?popular=%v&limit=%d", HOST, API, forum, d.popular, d.limit)
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

func (d *DCard) Comments(post DCardPost, after int) (out []DCardComment) {
	URL := fmt.Sprintf("%s/%s/posts/%v/comments", HOST, API, post.Id)
	if after >= 0 {
		URL = fmt.Sprintf("%s?after=%d", URL, after)
	}
	data := default_agent.fetch(URL)

	if err := json.Unmarshal(data, &out); err != nil {
		err = fmt.Errorf("Parse json failure - %s", err)
		panic(err)
	}

	return
}

func (d *DCard) AllComments(post DCardPost) (out []DCardComment) {
	after := 0
	for {
		tmp := d.Comments(post, after)
		out = append(out, tmp...)

		after = out[len(out)-1].Floor
		if after >= post.CommentCount || len(tmp) == 0 {
			break
		}
	}
	return
}
