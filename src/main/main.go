package main

import (
    "fmt"
    "net/http"
    "html/template"
    "obj"
    "ws"
    "os"
)

var homeTempl = template.Must(template.ParseFiles("templates/home.html"))
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}

var game *obj.Game;
func main() {
    game = obj.NewGame()
    go game.Start()

    fmt.Println("start")
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    http.HandleFunc("/", serveHome)
    http.HandleFunc("/ws", ws.HandlerFactory(game))

    bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_DIY_IP"), os.Getenv("OPENSHIFT_DIY_PORT"))
	fmt.Printf("listening on %s...", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
    fmt.Println("stop")
}