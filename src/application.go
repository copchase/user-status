package main

import (
	"log"
	"net/http"

	"github.com/copchase/user-status/controller"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/userinfo", controller.GetUserInfo)

	log.Fatal(http.ListenAndServe(":5000", router))
}
