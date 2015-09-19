package ws

import (
    "fmt"
    "github.com/gorilla/websocket"
    "net/http"
    "encoding/json"
    "obj"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func HandlerFactory(game *obj.Game) func(http.ResponseWriter, *http.Request) {
    res := func (w http.ResponseWriter, r *http.Request) {
        fmt.Println("connection")
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            fmt.Println(err)
            return
        }
        for {
            messageType, p, err := conn.ReadMessage()
            var f interface{}
            if err != nil {
                return
            }
            err = json.Unmarshal(p, &f)
            if err != nil {
                return
            }

            response := commands(game, f.(map[string]interface{}))

            if response != nil {
                if err = conn.WriteMessage(messageType, response); err != nil {
                    return
                }
            }
        }
    }
    return res
}

type response struct {
    Get string
    Data string
}

func commands(game *obj.Game, message map[string]interface{}) []byte {
    switch command := message["get"]; command {
        case "units":
            resp := map[string]interface{}{
                "get": "units",
                "units": game.World.Units,
            }
            b, err := json.Marshal(resp)
            if err != nil {
                return nil
            }
            return b
        case "boom":
            coords := message["args"].(map[string]interface {})
            game.MakeBoom(float32(coords["x"].(float64)), float32(coords["y"].(float64)))
    }
    return nil
}