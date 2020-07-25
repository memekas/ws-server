package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var email = flag.String("user", "qwertyaq@yandex.ru", "")
var password = flag.String("pass", "secret", "")

var log = logrus.New()

func auth() (*http.Response, error) {
	// User auth
	user := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    *email,
		Password: *password,
	}

	b, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://"+*addr+"/user/login", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func main() {
	flag.Parse()

	// Auth user and get cookies
	resp, err := auth()
	if err != nil {
		log.Error(err)
		return
	}
	cArr := resp.Cookies()

	// create ws connection with cookies
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/notification/subscribe"}

	rawConn, err := net.Dial("tcp", u.Host)
	if err != nil {
		log.Error(err)
		return
	}
	wsHeaders := http.Header{
		"Cookie": {cArr[0].String()},
	}

	log.Printf("connecting to %s", u.String())
	ws, resp, err := websocket.NewClient(rawConn, &u, wsHeaders, 1024, 1024)
	if err != nil {
		log.Error(err)
		return
	}
	defer ws.Close()

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("recv: %s", p)
	}
}
