package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/memekas/ws-server/models"
	"github.com/memekas/ws-server/utils"
	"github.com/sirupsen/logrus"
)

// RegUser - register new user in *DB
func RegUser(con *models.DB, log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		user := &models.Account{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, "Invalid request"))
			return
		}

		if err := user.Create(con); err != nil {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, err.Error()))
			return
		}

		resp := utils.Message(true, "User registered")
		resp["user"] = user
		utils.Respond(w, http.StatusOK, resp)
	})
}

// LoginUser - login user
func LoginUser(con *models.DB, log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoHandleFunc(log, r)

		user := &models.Account{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			utils.Respond(w, http.StatusBadRequest, utils.Message(false, "Invalid request"))
			return
		}

		if err := user.Login(con); err != nil {
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
