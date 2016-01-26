package main

import (
	"fmt"
	"html/template"
	"net/http"
	"obj"
	"ws"
)

var homeTempl = template.Must(template.ParseFiles("templates/home.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}

var game *obj.Game

func main() {
	game = obj.NewGame()
	go game.Start()

	fmt.Println("start")
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", ws.HandlerFactory(game))

	err := http.ListenAndServe(":7101", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("stop")
}
