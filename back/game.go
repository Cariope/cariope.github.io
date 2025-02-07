package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	playersMutex sync.RWMutex
	players      = []*Player{}
	gameStarted  = false
	verbs        []string
	objects      []string

	missions        = []*Mission{}
	missionIndex    = 0
	givenMissions   = []*Mission{}
	spaceShipHealth = 100

	config = Config{}
)

func gameLoop() {
	configure()

	for {
		if len(players) == config.NbPlayers {
			gameStarted = true
			sendMessageToAllPlayers(WebsocketMessage{Type: "start_game"})
			time.Sleep(3 * time.Second)

			countDecrease := time.Now()

			for spaceShipHealth > 0 && gameStarted {
				checkExpired()
				fillMissions()

				time.Sleep(200 * time.Millisecond)

				if countDecrease.Add(30 * time.Second).Before(time.Now()) {
					countDecrease = time.Now()
					config.DefaultTimeout = max(5, config.DefaultTimeout-config.DecreaseTimeout)
				}

			}
			sendMessageToAllPlayers(WebsocketMessage{Type: "spaceship_explosion"})
			configure()

		}
		time.Sleep(1 * time.Second)
	}
}

func configure() {

	gameStarted = false
	players = []*Player{}
	missions = []*Mission{}
	missionIndex = 0
	givenMissions = []*Mission{}
	spaceShipHealth = 100

	config = Config{
		NbPlayers:          10,
		NbTargetMissions:   4,
		NbSolvableMissions: 8,
		FailHealth:         -10,
		SuccessHealth:      5,
		DefaultTimeout:     20,
		DecreaseTimeout:    1,
	}

	verbsFile, err := loadFile("./verbs.txt")
	if err != nil {
		log.Fatal(err)
	}
	verbs = verbsFile
	objectsFile, err := loadFile("./objects.txt")
	if err != nil {
		log.Fatal(err)
	}
	objects = objectsFile

	// Generate missions
	generatedMissions := map[string]bool{}

	for i := 0; i < 100; i++ {
		// random value between 0 and len(verbs)
		randVerb := rand.Intn(len(verbs))
		verb := verbs[randVerb]
		randObject := rand.Intn(len(objects))
		object := objects[randObject]
		if generatedMissions[verb+object] {
			i--
			continue
		}
		mission := &Mission{
			Id:     fmt.Sprintf("%d_%d", randVerb, randObject),
			Verb:   verb,
			Action: object,
		}
		missions = append(missions, mission)
	}
}

func fillMissions() {
	// Fill missions for each player
	playersMutex.Lock()
	defer playersMutex.Unlock()
	for _, player := range players {
		toFill := config.NbSolvableMissions - len(player.SolvableMissions)

		for i := 0; i < toFill; i++ {
			player.SolvableMissions = append(player.SolvableMissions, missions[missionIndex])
			givenMissions = append(givenMissions, missions[missionIndex])
			player.Write(WebsocketMessage{Type: "new_solvable_mission", Mission: missions[missionIndex]})
			missionIndex++
		}
	}
	shuffleList(givenMissions)

	givenMissionIndex := 0
	for _, player := range players {
		toFill := config.NbTargetMissions - len(player.TargetMissions)
		for i := 0; i < toFill; i++ {
			givenMissions[givenMissionIndex].TimeoutSeconds = config.DefaultTimeout
			givenMissions[givenMissionIndex].Timeout = time.Now().Add(time.Duration(config.DefaultTimeout) * time.Second)
			player.TargetMissions = append(player.TargetMissions, givenMissions[givenMissionIndex])
			player.Write(WebsocketMessage{Type: "new_target_mission", Mission: givenMissions[givenMissionIndex]})
			givenMissionIndex++

		}
	}
	givenMissions = givenMissions[givenMissionIndex:]

}

func checkExpired() {
	for _, player := range players {
		for _, mission := range player.TargetMissions {

			if mission.Timeout.Before(time.Now()) {

				player.TargetMissions = removeMissionFromSlice(mission.Id, player.TargetMissions)
				player.Write(WebsocketMessage{
					Type:    "mission_target_failed",
					Mission: mission,
				})
				givenMissions = append(givenMissions, mission)
				updateHealth(config.FailHealth)

			}
		}
	}
}

func solveMission(player *Player, missionID string) {

	// Player should have this mission
	foundSolvable := false
	for _, m := range player.SolvableMissions {
		if m.Id == missionID {
			foundSolvable = true
		}
	}

	if foundSolvable {

		//  This mission should be in TargetMission of a player
		for _, p := range players {
			for _, m := range p.TargetMissions {
				if m.Id == missionID {
					updateHealth(config.SuccessHealth)
					p.Write(WebsocketMessage{
						Type:    "mission_target_solved",
						Mission: m,
					})

					player.Write(WebsocketMessage{
						Type:    "mission_solvable_solved",
						Mission: m,
					})

					player.SolvableMissions = removeMissionFromSlice(m.Id, player.SolvableMissions)
					p.TargetMissions = removeMissionFromSlice(m.Id, p.TargetMissions)
					return
				}
			}
		}
		updateHealth(config.FailHealth)

		player.Write(WebsocketMessage{
			Type: "mission_solvable_failed",
			Mission: &Mission{
				Id: missionID,
			},
		})

	} else {
		println("Player ask to solve a mission he doesn't have")
	}

}

func updateHealth(value int) {
	spaceShipHealth = min(100, max(0, spaceShipHealth+value))
	sendMessageToAllPlayers(WebsocketMessage{
		Type:   "update_health",
		Health: &spaceShipHealth,
	})
}
