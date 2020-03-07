/* Copyright (C) 2020-2020 cmj. All right reserved. */
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Link struct {
	Key        string `json:"key"`
	SourceLink string `json:"source"`
}

type ShortURL struct {
	*sql.DB
	*log.Logger

	Table    string
	MaxRetry int

	driver string
}

func MustNew(driver, dsn string) (surl *ShortURL) {
	var db *sql.DB
	var err error

	if db, err = sql.Open(driver, dsn); err != nil {
		panic(err)
		return
	}
	if err = db.Ping(); err != nil {
		panic(err)
		return
	}

	surl = &ShortURL{
		DB:       db,
		Logger:   log.New(os.Stderr, "", log.Lshortfile|log.LUTC),
		Table:    "short_url",
		MaxRetry: 5,

		driver: driver,
	}

	return
}

func (surl *ShortURL) Create() {
	stmt := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id	        STRING PRIMARY KEY,
			link        STRING
		)
	`, surl.Table)

	if _, err := surl.Exec(stmt); err != nil {
		surl.Panicf("Cannot create table: %s\n----\n%s\n", err, stmt)
		return
	}
}

func (surl *ShortURL) Get(in string) (out *Link) {
	var stmt string

	switch surl.driver {
	case "sqlite3":
		stmt = fmt.Sprintf("SELECT link FROM %s WHERE id = ?", surl.Table)
		break
	default:
		surl.Panicf("SQL driver not implemented: %s", surl.driver)
		break
	}

	rows, err := surl.Query(stmt, in)
	if err != nil {
		/* Cannot query SQL */
		surl.Panicf("Cannot query %s: %s", stmt, err)
	}
	defer rows.Close()

	for rows.Next() {
		var link string

		if err := rows.Scan(&link); err != nil {
			/* Cannot get the row from SQL */
			surl.Panicf("Cannot get row: %s", err)
		}
		out = &Link{
			Key:        in,
			SourceLink: link,
		}

		break
	}

	return
}

func (surl *ShortURL) Set(in *Link) {
	var stmt string

	switch surl.driver {
	case "sqlite3":
		stmt = fmt.Sprintf("INSERT OR REPLACE INTO %s (id, link) VALUES (?, ?)", surl.Table)
		break
	default:
		surl.Panicf("SQL driver not implemented: %s", surl.driver)
		break
	}

	if _, err := surl.Exec(stmt, in.Key, in.SourceLink); err != nil {
		/* execute SQL failure */
		surl.Panicf("Execute %s: %s", stmt, err)
	}
	return
}
