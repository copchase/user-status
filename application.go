package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/userinfo", GetUserInfo)

	http.ListenAndServe(":5000", router)
}
