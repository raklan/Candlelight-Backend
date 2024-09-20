package main

import (
	"candlelight-api/Accounts"
	"candlelight-api/CreationStudio"
	"candlelight-api/Lobby"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
)

func main() {
	startServer()
}

func startServer() {
	//Go does the strangest datetime string formatting I've ever seen. You give it a specific date/time (Specifically Jan 2, 2006 3:04:05 PM GMT-7)
	//in the format you want, and it'll match whatever the object is into that format
	logName := fmt.Sprintf("./logs/%v.log", time.Now().Format("2006-01-02_15-04-05"))

	//Log file & Server startup
	log.SetPrefix("CANDLELIGHT-API: ")
	logfile, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	log.Println("Starting HTTP listener...")

	//Start the server at localhost:10000 & register all paths
	mux := http.NewServeMux()
	registerPathHandlers(mux)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	handler := c.Handler(mux)

	http.ListenAndServe(":10000", handler)
}
func registerPathHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", heartbeat)
	mux.HandleFunc("/dummy", CreationStudio.GenerateJSON)

	//Gamedef-related requests
	mux.HandleFunc("/studio", CreationStudio.Studio)
	mux.HandleFunc("/allGames", CreationStudio.GetAllGames)

	//Lobby-related requests
	mux.HandleFunc("/joinLobby", Lobby.HandleJoinLobby)
	mux.HandleFunc("/hostLobby", Lobby.HostLobby)
	mux.HandleFunc("/rejoinLobby", Lobby.HandleRejoinLobby)

	//Account-related Requests
	mux.HandleFunc("/createAccount", Accounts.CreateAccount)
	mux.HandleFunc("/login", Accounts.Login)
	//mux.HandleFunc("/changePassword", Accounts.ChangePassword)
}

// Simple heartbeat endpoint to test if the server is up and running
func heartbeat(w http.ResponseWriter, r *http.Request) {
	log.Println("==Heartbeat==: Returning dummy response...")
	fmt.Fprintf(w, "Buh-dump, buh-dump")
}
