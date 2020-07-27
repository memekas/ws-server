package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/memekas/ws-server/pkg/db"
	"github.com/memekas/ws-server/pkg/rabbit"
	"github.com/memekas/ws-server/pkg/ws_server"
	"github.com/sirupsen/logrus"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "http service address")
	log := logrus.New()

	err := godotenv.Load()
	if err != nil {
		log.Error(err)
		return
	}

	flag.Parse()
	log.Info(*addr)

	db := &db.DB{}
	if err := db.Init(); err != nil {
		log.Error(err)
		return
	}
	defer db.Close()

	rabbit := &rabbit.RabbitMQ{}
	if err := rabbit.Init(); err != nil {
		log.Error(err)
		return
	}
	defer rabbit.Close()

	go ws_server.Sender(log, rabbit)

	router := mux.NewRouter()

	router.Handle("/user/new", ws_server.RegUser(db, log)).Methods("POST")
	router.Handle("/user/login", ws_server.LoginUser(db, log)).Methods("POST")

	router.Handle("/notification/subscribe", ws_server.NotificationSub(log))
	router.Handle("/notification/send", ws_server.NotificationSend(log, rabbit)).Methods("POST")

	router.Handle("/notification/users", ws_server.GetUsersOnline(log))

	err = http.ListenAndServe(*addr, router)
	if err != nil {
		log.Error(err)
	}
}
