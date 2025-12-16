package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

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
type FreeGamesPromotionsElements struct {
	Title string                   `json:"title"`
	Price FreeGamesPromotionsPrice `json:"price"`
	// UpcomingPromotions FreeGamesUpcomingPromotion     `json:"promotions"`
	KeyImages     []FreeGamesPromotionsKeyImages `json:"keyImages"`
	OfferMappings []FreeGamesOfferMappings       `json:"offerMappings"`
}

type FreeGamesOfferMappings struct {
	PageSlug string `json:"pageSlug"`
}

type FreeGamesPromotionsKeyImages struct {
	Url string `json:"url"`
}

// type FreeGamesUpcomingPromotion struct {
// 	UpcomingPromotionalOffers []FreeGamesUpcommingPromotionalOffers `json:"upcomingPromotionalOffers"`
// }
// type FreeGamesUpcommingPromotionalOffers struct {
// 	PromotionalOffers []FreeGamesPromotionalOffers `json:"promotionalOffers"`
// }
// type FreeGamesPromotionalOffers struct {
// StartDate string `json:"startDate"`
// EndDate   string `json:"endDate"`
// DiscountSettings FreeGamesPromotionalOffersDiscountSettings `json:"discountSetting"`
// }
//	type FreeGamesPromotionalOffersDiscountSettings struct {
//		DiscountPercentage int `json:"discountPercentage"`
//	}

type FreeGamesPromotionsPrice struct {
	TotalPrice FreeGamesPromotionsTotalPrice `json:"totalPrice"`
}
type FreeGamesPromotionsTotalPrice struct {
	DiscountPrice int `json:"discountPrice"`
}

type GameInfo struct {
	Name    string
	Picture string
	Url     string
	// EndDate string
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
	var freeGames []GameInfo
	for _, element := range elements {
		if element.Price.TotalPrice.DiscountPrice == 0 {
			var gameInfo GameInfo
			gameInfo.Name = element.Title
			gameInfo.Picture = GetGamePicture(element)
			gameInfo.Url = GetGameUrl(element)
			// gameInfo.EndDate = element.UpcomingPromotions.UpcomingPromotionalOffers[0].PromotionalOffers[0].EndDate
			freeGames = append(freeGames, gameInfo)
			// fmt.Println("Found free game:", gameInfo.Url)
		}
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
