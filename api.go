package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/memekas/ws-server/controllers"
	"github.com/memekas/ws-server/models"
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

	db := &models.DB{}
	if err := db.Init(); err != nil {
		log.Error(err)
		return
	}
	defer db.Close()

	rabbit := &models.RabbitMQ{}
	if err := rabbit.Init(); err != nil {
		log.Error(err)
		return
	}
	defer rabbit.Close()

	go controllers.Sender(log, rabbit)

	router := mux.NewRouter()

	router.Handle("/user/new", controllers.RegUser(db, log)).Methods("POST")

	router.Handle("/user/login", controllers.LoginUser(db, log)).Methods("POST")

	router.Handle("/notification/subscribe", controllers.NotificationSub(log))
	router.Handle("/notification/send", controllers.NotificationSend(log, rabbit)).Methods("POST")

	err = http.ListenAndServe(*addr, router)
	if err != nil {
		log.Error(err)
	}
}
