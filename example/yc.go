package main

import (
	"database/sql"
	"flag"
	"fmt"
	oauth "github.com/akrennmair/goauth"
	_ "github.com/mattn/go-yql"
	"io/ioutil"
	"log"
)

var (
	conKey    *string = flag.String("key", "", "Consumer Key")
	conSecret *string = flag.String("secret", "", "Consumer Secret")
)

func main() {
	flag.Parse()
	log.SetFlags(log.Ltime | log.Lshortfile)

	log.Printf("key=%s   secret=%s", *conKey, *conSecret)

	db, _ := sql.Open("yql", *conKey+"|"+*conSecret)

	stmt, err := db.Query(
		"select * from contentanalysis.analyze where url=?",
		"http://www.nydailynews.com/entertainment/bar-refaeli-hottest-moments-gallery-1.1061268?old=%2Fgallery.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	for stmt.Next() {
		var data interface{}
		stmt.Scan(&data)
		fmt.Printf("%v\n", data)
	}
}
