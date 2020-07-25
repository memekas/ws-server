// package main

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/gorilla/websocket"
// )

// // func TestExample(t *testing.T) {
// // 	t.Run("PONGER", func(t *testing.T) {
// // 		s := httptest.NewServer(http.HandlerFunc(echo))
// // 		defer s.Close()

// // 		u := "ws" + strings.TrimPrefix(s.URL, "http")

// // 		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
// // 		if err != nil {
// // 			t.Error(err)
// // 		}
// // 		defer ws.Close()

// // 		for i := 0; i < 10; i++ {
// // 			if err := ws.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
// // 				t.Error(err)
// // 			}
// // 			_, p, err := ws.ReadMessage()
// // 			if err != nil {
// // 				t.Error(err)
// // 			}
// // 			if string(p) != "pong" {
// // 				t.Errorf("bad response: %s, wait: pong", p)
// // 			}
// // 			t.Logf("resp: %s", p)
// // 		}
// // 	})
// // }
