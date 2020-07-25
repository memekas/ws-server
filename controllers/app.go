package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/memekas/ws-server/models"
	"github.com/memekas/ws-server/utils"
	"github.com/sirupsen/logrus"
)

var channels = make(map[uint]chan []byte)
var usersConnectCount uint32
var mu = &sync.Mutex{}

// Echo -
func Echo(log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error(err)
			return
		}
		i := 0
		for {
			timer := time.NewTimer(time.Second)
			<-timer.C
			if i == 5 {
				if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye")); err != nil {
					log.Error(err)
					return
				}
				a := c.RemoteAddr()
				log.Infof("Close connection %s", a.String())
				return
			}
			if err := c.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(i))); err != nil {
				log.Error(err)
				return
			}
			i++
		}
	})
}

// Root - redirect to login or home page
func Root(log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		_, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
	})
}

// NotificationSub - subscribe user to notifications
func NotificationSub(log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		cookie, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		// upgrade connection
		upgrader := websocket.Upgrader{}
		con, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error(err)
			return
		}

		// add user count
		atomic.AddUint32(&usersConnectCount, 1)
		defer atomic.AddUint32(&usersConnectCount, ^uint32(0))

		// get user id from cookie
		tk := &models.Token{}
		tk.Decrypt(cookie.Value)

		// create new channel for messages. Delete when connection refused
		in := make(chan []byte)

		mu.Lock()
		channels[tk.UserID] = in
		mu.Unlock()

		defer func() {
			mu.Lock()
			delete(channels, tk.UserID)
			mu.Unlock()
		}()

		// read from channel and send to user
		for {
			if err := con.WriteMessage(websocket.TextMessage, <-in); err != nil {
				log.Error(err)
				return
			}
		}

	})
}

// NotificationSendToUser - send notification to user
func NotificationSendToUser(log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		_, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		not := &models.Notification{}
		err = json.NewDecoder(r.Body).Decode(not)
		if err != nil {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, "Invalid request"))
			return
		}
		log.Info(*not)

		out, ok := channels[not.ToUser]
		if !ok {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, "toUser not found"))
			return
		}

		out <- []byte(not.Msg)
	})
}

// NotificationSendToAll - send notification to all users
func NotificationSendToAll(log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		_, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		not := &models.Notification{}
		err = json.NewDecoder(r.Body).Decode(not)
		if err != nil {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, "Invalid request"))
			return
		}
		log.Info(*not)

		mu.Lock()
		for _, val := range channels {
			val <- []byte(not.Msg)
		}
		mu.Unlock()
	})
}
