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

        player := game.AddPlayer()

        for {
            messageType, p, err := conn.ReadMessage()
            var f interface{}
            if err != nil {
                game.DeleteUnit(player)
                return
            }
            err = json.Unmarshal(p, &f)
            if err != nil {
                return
            }

            response := parseCommand(game, f.(map[string]interface{}), player)

            if response != nil {
                if err = conn.WriteMessage(messageType, response); err != nil {
                    game.DeleteUnit(player)
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

func parseCommand(game *obj.Game, message map[string]interface{}, player *obj.Unit) []byte {
    switch command := message["get"]; command {
        case "units":
            return game.UnitsToJSON()
        case "boom":
            coords := message["args"].(map[string]interface {})
            game.MakeBoom(float32(coords["x"].(float64)), float32(coords["y"].(float64)))
        case "move":
            player.SetPlayerMoveSpeed(message["args"].(map[string]interface {}))
        case "direction":
            coords := message["args"].(map[string]interface {})
            player.DX, player.DY = float32(coords["x"].(float64)), float32(coords["y"].(float64))
        case "fire":
            player.F = message["args"].(bool)
    }
    return nil
}