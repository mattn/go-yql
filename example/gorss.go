package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-yql"
)

func main() {
	db, _ := sql.Open("yql", "")

	stmt, err := db.Query(
		"select * from rss where url = ?",
		"http://blog.golang.org/feeds/posts/default?alt=rss")
	if err != nil {
		fmt.Println(err)
		return
	}
	for stmt.Next() {
		var data map[string]interface{}
		stmt.Scan(&data)
		fmt.Printf("%v\n", data["link"])
		fmt.Printf("  %v\n\n", data["title"])
	}
}
