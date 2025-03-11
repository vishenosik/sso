package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
)

type Fetcher interface {
	Fetch() (string, error)
}

type FetcherFunc func() (string, error)

func (f FetcherFunc) Fetch() (string, error) {
	return f()
}

// customDriver implements the `sql.Driver` interface.
type customDriver struct {
	// base is the base database driver
	base driver.Driver

	// fetcherFn is the function that will be used to fetch the database credential
	fetcherFn FetcherFunc
}

func (d *customDriver) Open(_ string) (driver.Conn, error) {
	// fetch the database credential
	dsn, err := d.fetcherFn()
	if err != nil {
		return nil, err
	}

	// open the database connection using the fetched credential and base driver
	return d.base.Open(dsn)
}

func OpenWithRotator(name string, base driver.Driver, fetcher Fetcher) (*sql.DB, error) {
	sql.Register(name, &customDriver{
		base:      base,
		fetcherFn: fetcher.Fetch,
	})

	// you don't need to fill the dsn, it will be fetched from the fetcher.
	return sql.Open(name, "")
}

var counter int

func simpleFetcher() (string, error) {
	log.Println("fetcher called")

	// add your custom logic
	// e.g. fetching from vault / config / etc.

	counter++
	return fmt.Sprintf("file:foobar-%d.sqlite", counter), nil
}

func main() {
	// you can adjust the `&sqlite3.SQLiteDriver{}` accordingly (e.g. &pq.Driver{}, etc.)
	db, err := OpenWithRotator("foobar", &sqlite3.SQLiteDriver{}, FetcherFunc(simpleFetcher))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(2 * time.Second)

	for range time.Tick(time.Second) {
		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
	}
}
