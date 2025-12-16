package main

import (
	"fmt"
	"log"
)

func main() {
	jsonfile := "games_info.json"
	freeGames, err := CheckFreeGame()
	if err != nil {
		log.Fatal(err)
		return
	}
	err, isUpdated := GameInfoToJson(freeGames, jsonfile)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(isUpdated)

	if !isUpdated {
		log.Printf("No new free games found.")
		return
	}
	err = SendGames(jsonfile)
	if err != nil {
		log.Fatal(err)
		return
	}

}
