package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-yql"
	"log"
)

var (
	conKey    *string = flag.String("key", "", "Consumer Key")
	conSecret *string = flag.String("secret", "", "Consumer Secret")
)

func main() {
	flag.Parse()
	db, err := sql.Open("yql", *conKey+"|"+*conSecret)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, arg := range flag.Args() {
		stmt, err := db.Query("select * from contentanalysis.analyze where url=?", arg)
		if err != nil {
			log.Fatal(err)
			return
		}
		for stmt.Next() {
			var data interface{}
			err = stmt.Scan(&data)
			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Printf("%v\n", data)
		}
	}
}
