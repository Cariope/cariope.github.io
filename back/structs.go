package main

import (
	"time"
)

type Config struct {
	NbPlayers          int `json:"nb_players"`
	NbTargetMissions   int `json:"nb_target_missions"`
	NbSolvableMissions int `json:"nb_solvable_missions"`
	FailHealth         int `json:"fail_health"`
	SuccessHealth      int `json:"success_health"`
	DefaultTimeout     int `json:"default_timeout"`
	DecreaseTimeout    int `json:"decrease_timeout"`
}

type Mission struct {
	Id             string `json:"id"`
	Verb           string `json:"verb"`
	Action         string `json:"action"`
	TimeoutSeconds int    `json:"timeout_seconds"`

	Timeout time.Time
	solved  *bool
}

type WebsocketMessage struct {
	Type      string   `json:"type"`
	Mission   *Mission `json:"mission"`
	Health    *int     `json:"health"`
	Config    *Config  `json:"config"`
	NbPlayers *int     `json:"nb_players"`
}

type Player struct {
	ClientID string
	Write    func(WebsocketMessage) error
	Name     string `json:"name"`
	MissionStatus
}

type MissionStatus struct {
	TargetMissions   []*Mission `json:"target_mission"`
	SolvableMissions []*Mission `json:"solvable_mission"`
}
