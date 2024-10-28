package main

import (
	"candlelight-api/Accounts"
	"candlelight-api/CreationStudio"
	"candlelight-api/Lobby"
	"os"

	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"gopkg.in/natefinch/lumberjack.v2"
)

const SERVER_VERSION = "v0.2.28 - Oct 26, 2024"

func main() {
	startServer()
}

func startServer() {
	//Log file & Server startup
	logName := "./logs/server.log"
	log.SetPrefix("CANDLELIGHT-API: ")

	compressLogs := false
	if os.Getenv("CANDLELIGHT_COMPRESS_LOGS") == "true" {
		compressLogs = true
	}

	log.SetOutput(&lumberjack.Logger{
		Filename: logName,
		MaxSize:  1,
		MaxAge:   7,
		Compress: compressLogs,
	})

	log.Println("Starting HTTP listener...")

	//Start the server at localhost:10000 & register all paths
	mux := http.NewServeMux()
	registerPathHandlers(mux)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	handler := c.Handler(mux)

	http.ListenAndServe(":10000", handler)
	// err = http.ListenAndServeTLS(":10000", "./candlelight-api/server.crt", "./candlelight-api/server.key", handler)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
func registerPathHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", heartbeat)
	mux.HandleFunc("/dummy", CreationStudio.GenerateJSON)
	mux.HandleFunc("/version", version)

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

func version(w http.ResponseWriter, r *http.Request) {
	log.Printf("==Version==: Returning Server Version (%s)", SERVER_VERSION)
	fmt.Fprint(w, SERVER_VERSION)
}
