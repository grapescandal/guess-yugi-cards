package logic

import (
	"fmt"
	"strings"
)

var answer string
var isStart bool

func StartGame() {
	isStart = true
	answer = "Rat Bat"
}

func GetHint() string {
	hint := ""
	for _, a := range answer {
		if string(a) != " " {
			hint += "-"
		} else {
			hint += " "
		}
	}
	return hint
}

func Answer(message string) (bool, string) {
	answerLower := strings.ToLower(answer)
	if len(message) != len(answer) {
		return false, fmt.Sprintf("Please be sure your answer length is %v", len(answer))
	}
	if message == answerLower {
		return true, answer
	} else {
		hint := ""
		for _, m := range message {
			wrongCounter := 0

			for i, a := range answerLower {
				if m == a {
					hint += string(answer[i])
					break
				} else {
					wrongCounter += 1
				}
			}

			if wrongCounter == len(message) {
				hint += "-"
			}
		}
		return false, hint
	}
}
