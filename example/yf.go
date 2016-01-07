package main

// Yahoo Finance example. Shows how to pass a custom env.

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-yql"
)

func main() {
	db, _ := sql.Open("yql", "||store://datatables.org/alltableswithkeys")

	stmt, err := db.Query(
		"select * from yahoo.finance.historicaldata where symbol = ? and startDate = ? and endDate = ?",
		"YHOO",
		"2009-09-11",
		"2010-03-10")
	if err != nil {
		fmt.Println(err)
		return
	}
	for stmt.Next() {
		var data map[string]interface{}
		stmt.Scan(&data)
		fmt.Printf("%v %v %v %v %v\n", data["Date"], data["Open"], data["High"], data["Low"], data["Close"])
	}
}
