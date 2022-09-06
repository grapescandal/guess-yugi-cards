package logic

import (
	"fmt"
	"guess-yugioh-cards-bot/api"
	"guess-yugioh-cards-bot/model"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/oliamb/cutter"
)

var answer model.Answer
var isStart bool
var openPieces []int
var turn int = 0
var maxTurn int = 0
var currentScore int = 10
var maxScore int = 10

func InitGame() {
	isStart = false
	openPieces = []int{}
	currentScore = maxScore
}

func StartGame() {
	if !isStart {
		isStart = true
		card := GetRandomCard()
		GetCardImage(card.CardImages[0].ImageURL)
		answer = model.Answer{
			Name:     card.Name,
			CardInfo: card,
		}
		turn = 0
	}
}

func GetTurn() int {
	return turn
}

func SetMaxTurn(number int) {
	maxTurn = number
}

func NextTurn() {
	turn += 1
	if turn > maxTurn {
		turn = 0
	}
}

func GetRandomCard() model.CardInfo {
	cardsData, err := api.GetCardsData()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomNumber := r1.Intn(len(cardsData.Data))
	card := cardsData.Data[randomNumber]
	fmt.Printf("Pick card Name: %v\n", card.Name)
	return card
}

func GetCardImage(url string) {
	_, err := api.GetCardImage(url)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	_, err = os.Open("card.jpg")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func ReadCardImage() *os.File {
	file, err := os.Open("card.jpg")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return file
}

func GetPieceCardImage(index int) (*os.File, error) {

	isAlreadyOpen := false
	for _, i := range openPieces {
		if index == i {
			isAlreadyOpen = true
			break
		}
	}

	if isAlreadyOpen {
		err := fmt.Errorf("%v is already open", index)
		return nil, err
	}

	file, err := os.Open("card.jpg")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Failed to decode: %v", err)
	}

	col := 3
	indexY := 0
	indexX := 0
	actualIndex := index - 1
	if index > 3 {
		indexY = actualIndex / col
		indexX = actualIndex % col
	} else {
		indexX = actualIndex
	}
	width := 108
	height := 108
	x := 50 + (108 * (indexX))
	y := 110 + (108 * (indexY))
	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:  width,
		Height: height,
		Anchor: image.Point{x, y},
		Mode:   cutter.TopLeft,
	})
	fmt.Printf("x: %v, y: %v\n", x, y)
	if err != nil {
		fmt.Printf("Error croppedImg: %v", err)
	}

	f, err := os.Create("piece.jpg")
	if err != nil {
		fmt.Printf("Failed to create: %v", err)
	}
	defer f.Close()

	err = jpeg.Encode(f, croppedImg, nil)
	if err != nil {
		fmt.Printf("Failed to encode: %v", err)
	}

	finalFile, err := os.Open("piece.jpg")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	openPieces = append(openPieces, index)
	return finalFile, nil
}

func DecreaseScore() {
	currentScore -= 1
}

func GetHint() string {
	hint := ""
	for _, a := range answer.Name {
		isAlphabets := isAlphabets(a)
		if isAlphabets {
			hint += "-"
		} else {
			hint += string(a)
		}
	}

	return hint
}

func isAlphabets(c rune) bool {
	return unicode.IsLetter(c)
}

func Answer(message string) (bool, bool, string, *model.Answer) {
	answerLower := strings.ToLower(answer.Name)
	if len(message) != len(answer.Name) {
		return false, false, fmt.Sprintf("Please be sure your answer length is %v", len(answer.Name)), nil
	}
	if message == answerLower {
		return true, true, "", &answer
	} else {
		hint := ""
		for i, a := range answerLower {

			if string(message[i]) == string(a) {
				hint += string(answer.Name[i])
			} else {
				hint += "-"
			}
		}
		return false, true, hint, nil
	}
}
