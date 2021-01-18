package main

import (
	"database/sql"
	"fmt"
)

func main() {

	sql.Register("driver", &PgnDriver{})
	fmt.Println(sql.Drivers())

	db, err := sql.Open("driver", "games.pgn")

	if err != nil {
		panic(err)
	}

	query := "SELECT White,Black,Result,PlyCount FROM games WHERE Result=?"

	rows, err := db.Query(query, "1-0")

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var white, black, result, plycount string
		err := rows.Scan(&white, &black, &result, &plycount)

		if err != nil {
			panic(err)
		}

		fmt.Println("white: "+white,
			"black: "+black,
			"result: "+result,
			"plycount: "+plycount)
	}
}

