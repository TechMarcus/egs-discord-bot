package main

import (
	"log"
	"time"
)

func main() {
	jsonfile := "games_info.json"
	freeGames, err := CheckFreeGame()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = GameInfoToJson(freeGames, jsonfile)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = SendGames(jsonfile)
	if err != nil {
		log.Fatal(err)
		return
	}
	time.Sleep(100000)
}
