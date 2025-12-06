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
	Title              string                         `json:"title"`
	Price              FreeGamesPromotionsPrice       `json:"price"`
	UpcomingPromotions FreeGamesUpcomingPromotion     `json:"promotions"`
	KeyImages          []FreeGamesPromotionsKeyImages `json:"keyImages"`
	OfferMappings      []FreeGamesOfferMappings       `json:"offerMappings"`
}

type FreeGamesOfferMappings struct {
	PageSlug string `json:"pageSlug"`
}

type FreeGamesPromotionsKeyImages struct {
	Url string `json:"url"`
}
type FreeGamesUpcomingPromotion struct {
	UpcomingPromotionalOffers []FreeGamesUpcommingPromotionalOffers `json:"upcomingPromotionalOffers"`
}
type FreeGamesUpcommingPromotionalOffers struct {
	PromotionalOffers []FreeGamesPromotionalOffers `json:"promotionalOffers"`
}
type FreeGamesPromotionalOffers struct {
	StartDate        string                                     `json:"startDate"`
	DiscountSettings FreeGamesPromotionalOffersDiscountSettings `json:"discountSetting"`
}

type FreeGamesPromotionalOffersDiscountSettings struct {
	DiscountPercentage int `json:"discountPercentage"`
}
type FreeGamesPromotionsPrice struct {
	TotalPrice FreeGamesPromotionsTotalPrice `json:"totalPrice"`
}

type FreeGamesPromotionsTotalPrice struct {
	OriginalPrice int `json:"originalPrice"`
	DiscountPrice int `json:"discountPrice"`
}

type GameInfo struct {
	Name    string
	Picture string
	Url     string
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
	if err != nil {
		return nil, err
	}

	var freeGames []GameInfo
	for _, element := range elements {
		if element.Price.TotalPrice.DiscountPrice == 0 {
			var gameInfo GameInfo
			gameInfo.Name = element.Title
			gameInfo.Picture = GetGamePicture(element)
			gameInfo.Url = GetGameUrl(element)
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

func GameInfoToJson(games []GameInfo, jsonfile string) error {
	result, error := json.Marshal(games)
	if error != nil {
		return error
	}
	file, err := os.OpenFile(jsonfile, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		firstLine := scanner.Text()
		if string(result) == firstLine {
			fmt.Println("No new games")
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(result)
	if err != nil {
		return err
	}

	return nil
}
