package ws_server

import (
	"encoding/json"
	"net/http"

	"github.com/memekas/ws-server/pkg/db"
	"github.com/memekas/ws-server/pkg/utils"
	"github.com/sirupsen/logrus"
)

func RegUser(con *db.DB, log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		user := &db.Account{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, "Invalid request"))
			return
		}

		if err := con.CreateUser(user); err != nil {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, err.Error()))
			return
		}

		resp := utils.Message(true, "User registered")
		resp["user"] = user
		utils.Respond(w, http.StatusOK, resp)
	})
}

func LoginUser(con *db.DB, log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		user := &db.Account{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, "Invalid request"))
			return
		}

		if err := con.LoginUser(user); err != nil {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, "User or password are incorrect"))
			return
		}

		// create cookie
		cookie := http.Cookie{
			Name:  "session",
			Value: user.Token,
		}
		http.SetCookie(w, &cookie)

		resp := utils.Message(true, "User logged in")
		resp["user"] = user
		utils.Respond(w, http.StatusOK, resp)
	})
}
