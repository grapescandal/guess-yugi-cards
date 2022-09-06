package api

import (
	"encoding/json"
	"fmt"
	"guess-yugioh-cards-bot/model"
	"io"
	"log"
	"net/http"
	"os"
)

const url = "https://db.ygoprodeck.com/api/v7/cardinfo.php?"

func GetCardsData() (*model.CardResponse, error) {
	resp := new(model.CardResponse)
	body, err := http.Get(url + "&startdate=01/01/1999&enddate=12/30/2004&dateregion=ocg_date")
	if err != nil {
		fmt.Println(err)
	}
	defer body.Body.Close()

	if body.StatusCode == 200 {
		err = json.NewDecoder(body.Body).Decode(&resp)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
	}

	return resp, err
}

func GetCardImage(url string) (*os.File, error) {
	// don't worry about errors
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create("card.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Download card image Success!")
	return file, err
}
