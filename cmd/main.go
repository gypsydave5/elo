package main

import (
	"fmt"

	"github.com/aceakash/elo"
)

func main() {

	

	table := elo.NewTable(32, 2000)
	table.Register("bruce")
	table.Register("clark")
	fmt.Print(table)
	table.AddResult("bruce", "clark")
	table.AddResult("bruce", "clark")
	table.AddResult("bruce", "clark")
	fmt.Print(table)
	store := elo.JsonFileTableStore{
		Filepath: "eloTable.json",
	}
	err := store.Save(table); if err != nil {
		fmt.Println("Error when saving to table")
		fmt.Print(err)
		panic(err)
	}
	table = elo.Table{}
	table, err = store.Load()
	if err != nil {
		fmt.Println("Error when saving to table")
		fmt.Print(err)
		panic(err)
	}
}
