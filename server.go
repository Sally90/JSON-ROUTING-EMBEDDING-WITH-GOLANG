package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonContentType = "application/json"

type Player struct {
	Name string
	Wins int
}

// PlayerStore stores score information about players.
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

// store and http.Handler are assgined to server
// http.Handler / Router should only be set up once including the Handle functions, so that all requests use the same http.Handler / router
type PlayerServer struct {
	store        PlayerStore
	http.Handler // Interface Embedding: PlayerServer struct has now all methods that http.Handler has which is ServeHTTP
}

//store has to be of type PlayerStore interface meaning that store has to implement the GetPlayerScore and RecordWin methods
//in test file, we pass a pointer to NewPlayerServer - why is this allowed?
//because in test file, the pointer implements the store interface (the two methods GetPlayerScore and RecordWin)
func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer) //new function returns a pointer to the given type
	p.store = store
	router := http.NewServeMux()
	// Handle (!) is a function that tells which Handler function should be used for which path
	// Handler (r!) is a function that takes in http.ResponseWriter and *http.Request and tells what to do given an incoming request
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	p.Handler = router // router implements the Handler interface (because ServeMus has ServeHTTP method) and can thus be assigned as http.Handler to server
	return p
}

//convention in Go: receiver is named by first letter of type
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
