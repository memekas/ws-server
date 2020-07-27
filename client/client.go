package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var email = flag.String("user", "qwertyaq@yandex.ru", "")
var password = flag.String("pass", "secret", "")

var ws *websocket.Conn
var exit = make(chan os.Signal, 1)

var log = logrus.New()

func main() {
	flag.Parse()
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	go exitf()

	// Auth user and get cookies
	log.Info("Reg user: " + *email)
	resp, err := reg()
	if err != nil {
		log.Error(err)
		return
	}
	if resp.StatusCode == 400 {
		log.Info("User " + *email + " already exist")
	}

	// Auth user and get cookies
	log.Info("Auth user: " + *email)
	resp, err = auth()
	if err != nil {
		log.Error(err)
		return
	}
	if resp.StatusCode == 400 {
		log.Info("Fail to auth")
		return
	}

	cArr := resp.Cookies()

	forever := make(chan interface{})

	go reciever(cArr[0].String())
	go reader(cArr[0].String())

	<-forever
}

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

func reg() (*http.Response, error) {
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

	req, err := http.NewRequest("POST", "http://"+*addr+"/user/new", bytes.NewBuffer(b))
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

func reciever(cookie string) {
	// create ws connection with cookies
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/notification/subscribe"}

	rawConn, err := net.Dial("tcp", u.Host)
	if err != nil {
		log.Error(err)
		return
	}
	wsHeaders := http.Header{
		"Cookie": {cookie},
	}

	log.Printf("connecting to %s", u.String())
	ws, _, err = websocket.NewClient(rawConn, &u, wsHeaders, 1024, 1024)
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

func reader(cookie string) {
	help()

	for {
		fmt.Println("Enter command: ")
		r := bufio.NewReader(os.Stdin)
		text, err := r.ReadString('\n')
		if err != nil {
			log.Error(err)
			continue
		}
		text = strings.TrimRight(text, "\n")

		args := strings.Split(text, " ")

		switch args[0] {
		case "send":
			send(cookie, args)
		case "users":
			users()
		case "exit":
			ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			ws.Close()
			os.Exit(0)
		default:
			help()
		}
	}
}

func help() {
	fmt.Println("usage:\n")
	fmt.Println("send [ID] [MSG]")
	fmt.Println("send MSG to user with ID. If ID == 0, send to all online users\n")

	fmt.Println("users")
	fmt.Println("show the count of online users\n")

	fmt.Println("exit")
}

func send(cookie string, args []string) {
	if len(args) < 3 {
		help()
		return
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		log.Error(err)
		return
	}

	not := struct {
		ToUser uint   `json:"toUser"`
		Msg    string `json:"msg"`
	}{
		ToUser: uint(id),
		Msg:    args[2],
	}

	b, err := json.Marshal(not)
	if err != nil {
		log.Error(err)
		return
	}

	req, err := http.NewRequest("POST", "http://"+*addr+"/notification/send", bytes.NewBuffer(b))
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "session", Value: cookie})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Something go wrong")
		return
	}
	fmt.Println("send OK!")
}

func users() {
	req, err := http.NewRequest("GET", "http://"+*addr+"/notification/users", nil)
	if err != nil {
		log.Error(err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Something go wrong")
		return
	}

	count := struct {
		Count uint32 `json:"count"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&count)
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Online users count: %d\n", count.Count)
}

func exitf() {
	<-exit
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	ws.Close()
	os.Exit(0)
}
