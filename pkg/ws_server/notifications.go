package ws_server

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/memekas/ws-server/pkg/auth"
	"github.com/memekas/ws-server/pkg/rabbit"
	"github.com/memekas/ws-server/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var channels = make(map[uint]chan []byte)
var usersConnectCount uint32
var mu = &sync.Mutex{}

func NotificationSub(log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		cookie, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			utils.Respond(w, http.StatusUnauthorized, utils.Message(false, "Auth cookie not found"))
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
		tk := &auth.Token{}
		if err := tk.Decrypt(cookie.Value); err != nil {
			utils.Respond(w, http.StatusUnauthorized, utils.Message(false, "Failed to decrypt cookie"))
			return
		}

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

// NotificationSend - send notification. If toUser == 0 send to all
func NotificationSend(log *logrus.Logger, rabbit *rabbit.RabbitMQ) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		_, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			utils.Respond(w, http.StatusUnauthorized, utils.Message(false, "Auth cookie not found"))
			return
		}

		not := &Notification{}
		err = json.NewDecoder(r.Body).Decode(not)
		if err != nil {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, "Invalid request"))
			log.Error(err)
			return
		}

		ch, err := rabbit.Get().Channel()
		if err != nil {
			utils.Respond(w, http.StatusInternalServerError, utils.Message(false, "Failed to open a rabbit channel"))
			log.Error(err)
			return
		}
		defer ch.Close()

		err = ch.ExchangeDeclare(
			os.Getenv("RABBITMQ_EXCHANGE_NAME_NOTIFICATIONS"), //name
			"direct", //type
			true,     //durable
			false,    // auto-delete
			false,    //internal
			false,    //no-wait
			nil,      //arguments
		)
		if err != nil {
			utils.Respond(w, http.StatusInternalServerError, utils.Message(false, "Failed to declare rabbit exchange"))
			log.Error(err)
			return
		}

		bNot, err := json.Marshal(not)
		if err != nil {
			utils.Respond(w, http.StatusInternalServerError, utils.Message(false, "Failed Marshal"))
			log.Error(err)
			return
		}

		err = ch.Publish(
			os.Getenv("RABBITMQ_EXCHANGE_NAME_NOTIFICATIONS"),     // exchange
			os.Getenv("RABBITMQ_EXCHANGE_ROUT_KEY_NOTIFICATIONS"), // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        bNot,
			})
		if err != nil {
			utils.Respond(w, http.StatusInternalServerError, utils.Message(false, "Failed to publish rabbit message"))
			log.Error(err)
			return
		}

	})
}

// Sender - Worker that sends notifications to users
func Sender(log *logrus.Logger, rabbit *rabbit.RabbitMQ) {
	ch, err := rabbit.Get().Channel()
	if err != nil {
		log.Error(err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		os.Getenv("RABBITMQ_EXCHANGE_NAME_NOTIFICATIONS"), //name
		"direct", //type
		true,     //durable
		false,    // auto-delete
		false,    //internal
		false,    //no-wait
		nil,      //arguments
	)
	if err != nil {
		log.Error(err)
		return
	}

	q, err := ch.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE_NAME_NOTIFICATIONS"), // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Error(err)
		return
	}

	err = ch.QueueBind(
		q.Name, // queue
		os.Getenv("RABBITMQ_EXCHANGE_ROUT_KEY_NOTIFICATIONS"), // routing key
		os.Getenv("RABBITMQ_EXCHANGE_NAME_NOTIFICATIONS"),     // exchange
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Error(err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		log.Error(err)
		return
	}

	for msg := range msgs {
		sendNotification(log, msg.Body)
	}
}

func sendNotification(log *logrus.Logger, body []byte) {
	not := &Notification{}
	err := json.Unmarshal(body, not)
	if err != nil {
		log.Error(err)
		return
	}

	mu.Lock()
	if not.ToUser == 0 {
		for _, value := range channels {
			value <- []byte(not.Msg)
		}
	} else {
		channels[not.ToUser] <- []byte(not.Msg)
	}
	mu.Unlock()
}
