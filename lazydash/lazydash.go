/* Copyright (C) 2020-2020 cmj. All right reserved. */
package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

type TableSpec struct {
	/* unique SQL SPEC name */
	Name string `yaml:"name"`
	/* SQL related table */
	Table string `yaml:"table"`
	/* Return fields */
	Field []struct {
		Name   string `yaml:"name"`
		Column string `yaml:"column"`
	}
}

type LazyDashConfig struct {
	/* project name */
	Name     string `yaml:"name"`
	Template string `yaml:"template"`
	/* Service bind on*/
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	/* max query size per each query */
	PageSize int `yaml:"page-size"`
	/* Database configuration */
	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	}
}

type LazyDash struct {
	/* API version */
	Version int `yaml:"version"`
	/* global configuration */
	Config LazyDashConfig `yaml: config`
	/* Table specification */
	Spec map[string]TableSpec `yaml:"spec"`

	/* private method */
	db *sql.DB
}

func New(config string) (out *LazyDash, err error) {
	fd, e := ioutil.ReadFile(config)
	if e != nil {
		err = fmt.Errorf("Cannot open file - %s", e)
		return
	}

	out = &LazyDash{
		Version: 1,
		Config: LazyDashConfig{
			Name:     "Lazy Dashboard",
			Template: "tmpl/lazydash.htm",
			Host:     "",
			Port:     9999,
			PageSize: 40,
		},
	}
	if e := yaml.Unmarshal(fd, &out); e != nil {
		err = fmt.Errorf("Cannot unmarshal config - %s", e)
		return
	}

	if out.Config.Database.Username == "" && out.Config.Database.Password == "" {
		uri := fmt.Sprintf("%s/%s", out.Config.Database.Host, out.Config.Database.Database)
		out.db, err = sql.Open(out.Config.Database.Driver, uri)
	} else {
		uri := fmt.Sprintf(
			"%s:%s@%s/%s",
			out.Config.Database.Username, out.Config.Database.Password,
			out.Config.Database.Host, out.Config.Database.Database,
		)
		out.db, err = sql.Open(out.Config.Database.Driver, uri)
	}
	return
}

func (in *LazyDash) Close() {
	/* close all resource */
	in.db.Close()
}

func (in *LazyDash) CheckEverything() (err error) {
	if err = in.db.Ping(); err != nil {
		/* Cannot connect to database */
		return
	}

	return
}

func (in *LazyDash) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var buff bytes.Buffer

	if strings.HasPrefix(r.URL.Path, "/static/") {
		fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
		fs.ServeHTTP(w, r)
	} else {
		t := template.Must(template.ParseFiles(in.Config.Template))
		t.Execute(&buff, in)
		w.Write(buff.Bytes())
	}
}

func (in *LazyDash) Run() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", in.Config.Host, in.Config.Port), in))
}

func (in *LazyDash) Name() (out string) {
	if out = in.Config.Name; out == "" {
		out = "Lazy Dashboard"
	}
	return
}

func main() {
	lazydash, err := New("config.yaml")
	if err != nil {
		fmt.Printf("New LazyDash failure - %s\n", err)
		return
	}
	defer lazydash.Close()

	if err := lazydash.CheckEverything(); err != nil {
		fmt.Printf("Check LazyDash failure - %s\n", err)
		return
	}

	lazydash.Run()
	return
}
