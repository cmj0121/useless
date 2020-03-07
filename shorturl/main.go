/* Copyright (C) 2020-2020 cmj. All right reserved. */
package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	driver, dsn := "sqlite3", "shorturl.sql"

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Catch panic")
			fmt.Println(r)
			os.Exit(1)
		}
	}()

	srv := MustNew(driver, dsn)
	srv.Create()
	srv.Get("demo")
	srv.Set(&Link{
		Key:        "demo",
		SourceLink: "http://example.com",
	})
}
