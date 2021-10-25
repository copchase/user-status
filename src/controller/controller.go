package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/copchase/user-status/irc"
	"github.com/julienschmidt/httprouter"
)

var db = irc.CreateTwitchDatabase()

func GetUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	channel := strings.ToLower(r.URL.Query().Get("channel"))
	user := strings.ToLower(r.URL.Query().Get("user"))
	if channel == "" {
		http.Error(w, "missing channel", http.StatusBadRequest)
	} else if user == "" {
		http.Error(w, "missing user", http.StatusBadRequest)
	}

	userInfo := db.ReadUserInfo(channel, user)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(userInfo)
	if err != nil {
		http.Error(w, "error marshaling UserInfo", http.StatusInternalServerError)
	}
}
