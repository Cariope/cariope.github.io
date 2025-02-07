package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	Ok                 = "ok"
	NotEnoughPlayers   = "not_enough_players"
	MissionNotSolvable = "mission_not_solvable"

	StartGame    = "start_game"
	SolveMission = "solve_mission"
)

func handleMessage(ws *websocket.Conn, msg WebsocketMessage) {

	log.Printf("Received message: %s - %+v", msg.Type, msg.Mission.Id)
	player := getPlayerByClientId(ws.RemoteAddr().String())

	switch msg.Type {
	case SolveMission:
		solveMission(player, msg.Mission.Id)
	}
}

func computeMessage(msg string, resp string) []byte {
	return []byte(fmt.Sprintf("%s:%s", msg, resp))
}
