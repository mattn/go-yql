YQL Go wrapper
=============================


Simple Yql Queries:

	db, _ := sql.Open("yql", "")

	stmt, err := db.Query(
		"select * from rss where url = ?",
		"http://blog.golang.org/feeds/posts/default?alt=rss")
	
	if err != nil {
		return
	}
	for stmt.Next() {
		var data map[string]interface{}
		stmt.Scan(&data)
		fmt.Printf("%v\n", data["link"])
		fmt.Printf("  %v\n\n", data["title"])
	}


Private Yql Queries:

	db, _ := sql.Open("yql", *conKey+"|"+*conSecret)

	stmt, err := db.Query(
		"select * from contentanalysis.analyze where url=?",
		"http://www.espn.com")

	if err != nil {
		return
	}
	for stmt.Next() {
		var data interface{}
		stmt.Scan(&data)
		fmt.Printf("%v\n", data)
	}