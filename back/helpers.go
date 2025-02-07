package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func shuffleList[T any](list []T) {
	for i := range list {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}
}

func addPlayer(player *Player) {
	playersMutex.Lock()
	defer playersMutex.Unlock()
	players = append(players, player)
	updatePlayerCount()
}

func updatePlayerCount() {
	len := len(players)
	sendMessageToAllPlayers(WebsocketMessage{
		Type:      "new_player",
		NbPlayers: &len,
	})
}

func removePlayer(clientID string) {
	playersMutex.Lock()
	defer playersMutex.Unlock()
	for i, player := range players {
		if player.ClientID == clientID {
			players = append(players[:i], players[i+1:]...)
			break
		}
	}
	updatePlayerCount()
}

func removeMissionFromSlice(missionID string, missions []*Mission) []*Mission {
	newSlice := []*Mission{}
	for i, mission := range missions {
		if mission.Id != missionID {
			newSlice = append(newSlice, missions[i])
		}
	}
	return newSlice
}

func getPlayerByClientId(clientID string) *Player {
	playersMutex.RLock()
	defer playersMutex.RUnlock()
	for _, player := range players {
		if player.ClientID == clientID {
			return player
		}
	}
	return nil
}

func loadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var inputs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		inputs = append(inputs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return inputs, nil
}

func sendMessageToAllPlayers(message WebsocketMessage) {
	for _, player := range players {
		fmt.Printf("message: %+v\n", message)
		player.Write(message)
	}
}

func logMissions(missions []*Mission) string {
	str := ""
	for _, mission := range missions {
		str += fmt.Sprintf("%v ", mission.Id)
	}
	return str
}
