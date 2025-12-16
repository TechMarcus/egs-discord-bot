package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// General structure
type FreeGamesPromotions struct {
	Data FreeGamesPromotionsData `json:"data"`
}
type FreeGamesPromotionsData struct {
	Catalog FreeGamesPromotionsCatalog `json:"Catalog"`
}
type FreeGamesPromotionsCatalog struct {
	SearchStore FreeGamesPromotionsSearchStore `json:"searchStore"`
}
type FreeGamesPromotionsSearchStore struct {
	Elements []FreeGamesPromotionsElements `json:"elements"`
}

// Game elements

type FreeGamesPromotionsElements struct {
	Title         string                         `json:"title"`
	Price         FreeGamesPromotionsPrice       `json:"price"`
	AllPromotions FreeGamesAllPromotions         `json:"promotions"`
	KeyImages     []FreeGamesPromotionsKeyImages `json:"keyImages"`
	OfferMappings []FreeGamesOfferMappings       `json:"offerMappings"`
}

// All promotions

type FreeGamesAllPromotions struct {
	CurrentPromotions  []FreeGameCurrentPromotion  `json:"promotionalOffers"`
	UpcomingPromotions []FreeGameUpcomingPromotion `json:"upcomingPromotionalOffers"`
}

// Current promotions

type FreeGameCurrentPromotion struct {
	GamesOffer []FreeGamesOffers `json:"promotionalOffers"`
}
type FreeGamesOffers struct {
	EndDate   string `json:"endDate"`
	StartDate string `json:"startDate"`
}

// Upcoming promotions

type FreeGameUpcomingPromotion struct {
	UpcommingGameOffer []UpcommingFreeGamesOffer `json:"promotionalOffers"`
}

type UpcommingFreeGamesOffer struct {
	EndDate   string `json:"endDate"`
	StartDate string `json:"startDate"`
}

// Url mapping

type FreeGamesOfferMappings struct {
	PageSlug string `json:"pageSlug"`
}

type FreeGamesPromotionsKeyImages struct {
	Url string `json:"url"`
}

// Price info

type FreeGamesPromotionsPrice struct {
	TotalPrice FreeGamesPromotionsTotalPrice `json:"totalPrice"`
}
type FreeGamesPromotionsTotalPrice struct {
	DiscountPrice int `json:"discountPrice"`
}

// Game info struct

type GameInfo struct {
	Name      string
	Picture   string
	Url       string
	EndDate   string
	StartDate string
}

func CheckFreeGame() ([]GameInfo, error) {
	var listurl = "https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=en-US&country=UA&allowCountries=UA"
	resp, err := http.Get(listurl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var freeGamesPromotions FreeGamesPromotions
	err = json.Unmarshal([]byte(b), &freeGamesPromotions)
	if err != nil {
		log.Fatal(err)
	}

	elements := freeGamesPromotions.Data.Catalog.SearchStore.Elements
	freeGames := make([]GameInfo, 0, len(elements))
	for _, element := range elements {
		if element.Price.TotalPrice.DiscountPrice != 0 {
			continue
		}
		gameInfo := GameInfo{
			Name:    element.Title,
			Picture: GetGamePicture(element),
			Url:     GetGameUrl(element),
		}
		promos := element.AllPromotions
		if len(promos.CurrentPromotions) > 0 &&
			len(promos.CurrentPromotions[0].GamesOffer) > 0 {

			offer := promos.CurrentPromotions[0].GamesOffer[0]
			gameInfo.StartDate = offer.StartDate
			gameInfo.EndDate = offer.EndDate

		} else if len(promos.UpcomingPromotions) > 0 &&
			len(promos.UpcomingPromotions[0].UpcommingGameOffer) > 0 {

			offer := promos.UpcomingPromotions[0].UpcommingGameOffer[0]
			gameInfo.StartDate = offer.StartDate
			gameInfo.EndDate = offer.EndDate
		}
		gameInfo.StartDate = ConvertTimeFormat(gameInfo.StartDate)
		gameInfo.EndDate = ConvertTimeFormat(gameInfo.EndDate)
		freeGames = append(freeGames, gameInfo)
	}
	return freeGames, nil
}

func GetGamePicture(element FreeGamesPromotionsElements) string {
	gamePicture := element.KeyImages[0].Url
	url := gamePicture
	return url
}

func GetGameUrl(element FreeGamesPromotionsElements) string {
	if len(element.OfferMappings) == 0 {
		return "Not found"
	}
	gameUrl := element.OfferMappings[0].PageSlug
	return "https://www.epicgames.com/store/en-US/p/" + gameUrl
}

func GameInfoToJson(games []GameInfo, jsonfile string) (error, bool) {
	result, error := json.Marshal(games)
	if error != nil {
		return error, false
	}
	file, err := os.OpenFile(jsonfile, os.O_RDONLY, 0755)
	if err != nil {
		return err, false
	}
	defer file.Close()

	// overwriting of file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		firstLine := scanner.Text()
		if string(result) == firstLine {
			fmt.Println("No new games")
			return nil, false
		}
	}
	file, error = os.OpenFile(jsonfile, os.O_RDWR|os.O_TRUNC, 0755)
	if error != nil {
		return error, false
	}
	defer file.Close()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(result)
	if err != nil {
		return err, false
	}

	return nil, true
}

func ConvertTimeFormat(givenTime string) string {
	t, err := time.Parse(time.RFC3339Nano, givenTime)
	if err != nil {
		panic(err)
	}
	loc, err := time.LoadLocation("Europe/Kyiv")
	if err != nil {
		panic(err)
	}
	uaTime := t.In(loc)

	return uaTime.Format("2006-01-02 15:04")
}
