package logic

import (
	"fmt"
	"guess-yugioh-cards-bot/helper"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const PREFIX = ".yugi"

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := m.Content
	if len(content) <= len(PREFIX) {
		return
	}
	if content[:len(PREFIX)] != PREFIX {
		return
	}
	content = content[len(PREFIX):]
	if len(content) < 1 {
		return
	}
	args := strings.Fields(content)
	command := strings.ToLower(args[0])

	if command == "create" {
		lobby := CreateLobby(m.ChannelID)
		message := fmt.Sprintf("Lobby: %s was created", lobby.Id)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	} else if command == "lobby" {
		lobby := GetLobby(m.ChannelID)
		message := ""
		if lobby == nil {
			message += "Lobby not found"
		} else {
			message += fmt.Sprintf("Lobby:%s \n", lobby.Id)
		}

		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	} else if command == "join" {
		playerName := helper.FilterInput(m.Content, PREFIX+" "+"join")
		player := CreatePlayer(playerName, m.Author.ID)
		lobbyID := JoinLobby(m.ChannelID, player)

		message := fmt.Sprintf("Player: %v has joined to lobby: %v", playerName, lobbyID)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	} else if command == "start" {
		message := ""
		lobby := GetLobby(m.ChannelID)
		if lobby == nil {
			message += "Please create lobby first"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		} else if lobby != nil && len(lobby.Player) == 0 {
			message += "Please join lobby first"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		StartGame()
		hint := GetHint()
		message += "Game Started!\n"

		message += fmt.Sprintf("Answer is %s", hint)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	} else if command == "answer" {
		message := ""

		if !isStart {
			message += "Game is not start yet"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		answer := strings.ToLower(helper.FilterInput(m.Content, PREFIX+" "+"answer"))
		player := GetPlayer(m.ChannelID, m.Author.ID)
		players := GetPlayers(m.ChannelID)
		result, answer := Answer(answer)
		if result {
			player.Score += 10
			message += fmt.Sprintf("Player: %v win, Answer is %v \n", player.Name, answer)
			message += "----------Scoreboard----------\n"
			sort.SliceStable(players, func(i, j int) bool {
				return players[i].Score > players[j].Score
			})
			for _, player := range players {
				message += fmt.Sprintf("%v : %v", player.Name, player.Score)
			}
		} else {
			message += fmt.Sprintf("Try again, Answer is %v", answer)
		}
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
