package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

func main() {
	port := "5000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	ServeBack(port)
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func ServeBack(port string) {
	r := mux.NewRouter()

	// Add a log to each entring call
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request: %s %s", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/", handleConnections)
	r.HandleFunc("/test", testHandler).Methods("GET")
	r.HandleFunc("/reset", reset).Methods("GET")
	r.HandleFunc("/set_config", setConfig).Methods("POST")

	// Wrap the router with CORS middleware
	handler := cors.Default().Handler(r)

	go gameLoop()

	log.Println("Server started on :", port)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clientID := ws.RemoteAddr().String()
	log.Printf("Client connected with id: %s", clientID)

	if gameStarted {
		ws.WriteMessage(websocket.TextMessage, []byte("game_already_started"))
	} else {
		player := &Player{
			ClientID: clientID,
			Write: func(msg WebsocketMessage) error {
				if msg.Mission != nil {
					println("write message", msg.Type, msg.Mission.Id)
				} else if msg.Health != nil {
					println("write message", msg.Type, *msg.Health)
				} else if msg.NbPlayers != nil {
					println("write message", msg.Type, *msg.NbPlayers)
				} else {
					println("write message", msg.Type)
				}
				msgBytes, _ := json.Marshal(msg)
				return ws.WriteMessage(websocket.TextMessage, []byte(msgBytes))
			},
		}
		addPlayer(player)

		player.Write(WebsocketMessage{
			Type:   "config",
			Config: &config,
		})
	}

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Client disconnected with id: %s", clientID)
			removePlayer(clientID)
			break
		}
		msgDecoded := WebsocketMessage{}
		json.Unmarshal(msg, &msgDecoded)

		handleMessage(ws, msgDecoded)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func setConfig(w http.ResponseWriter, r *http.Request) {

	var newConfig Config
	err := json.NewDecoder(r.Body).Decode(&newConfig)
	if err != nil {
		log.Printf("Error decoding config: %v", err)
		http.Error(w, "Invalid config", http.StatusBadRequest)
		return
	}

	config = newConfig

	sendMessageToAllPlayers(WebsocketMessage{
		Type:   "config",
		Config: &config,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func reset(w http.ResponseWriter, r *http.Request) {
	gameStarted = false
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
